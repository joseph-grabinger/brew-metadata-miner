package parser

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"main/config"
	"main/parser/delegate"
)

type parser struct {
	config *config.Config

	// A map of formulas, where the key is the name of the formula.
	formulas map[string]*formula
}

// NewParser creates a new parser.
func NewParser(config *config.Config) *parser {
	return &parser{
		config:   config,
		formulas: make(map[string]*formula),
	}
}

// Parse parses the core repository and extracts the formulas.
func (p *parser) Parse() error {
	return p.readFormulas()
}

func (p *parser) Analyze() {
	valid := make([]*formula, 0)
	noRepo := make([]*formula, 0)

	for _, value := range p.formulas {
		if value.repoURL != "" {
			valid = append(valid, value)
		} else {
			noRepo = append(noRepo, value)
		}
	}

	fmt.Println("Total number of formulas:", len(p.formulas))
	fmt.Println("Number of valid formulas:", len(valid))
	fmt.Println("Number of formulas without a repository:", len(noRepo))
	fmt.Println("Formulas without a repository:", noRepo)
}

// ReadFormaulas reads all formulas from the core repository into the formulas map.
func (p *parser) readFormulas() error {
	// Match parent directories of the fomula files.
	matches, err := filepath.Glob(p.config.CoreRepo.Dir + "/Formula/*")
	if err != nil {
		return err
	}

	for _, match := range matches {
		// Walk through the directory to read each file.
		if err := filepath.Walk(match, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error accessing path %s: %v\n", path, err)
				return err
			}
			// Skip directories.
			if info.IsDir() {
				return nil
			}

			log.Println("Reading file: ", path)

			file, err := os.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			// Parse Formula from file.
			sourceFormula, err := parseFromFile(file)
			if err != nil {
				log.Printf("Error parsing file %s: %v\n", path, err)
				return err
			}

			formula := fromSourceFormula(sourceFormula)

			// Add the formula to the formulas map.
			p.formulas[formula.name] = formula

			//log.Println("Successfully parsed formula:", formula)

			return nil
		}); err != nil {
			log.Printf("Error walking directory: %v\n", err)
		}
	}
	return nil
}

// parseFromFile parses a formula from a file into a Formula struct.
func parseFromFile(file *os.File) (*sourceFormula, error) {
	scanner := bufio.NewScanner(file)
	formulaParser := &delegate.FormulaParser{Scanner: scanner}

	base := filepath.Base(file.Name())
	name := strings.TrimSuffix(base, ".rb")

	formula := &sourceFormula{name: name}

	homepage, err := formulaParser.ParseField(homepagePattern, "homepage")
	if err != nil {
		return nil, err
	}
	formula.homepage = homepage

	fields := []delegate.ParseStrategy{
		delegate.NewSLM("url", urlPattern, *formulaParser),
		delegate.NewSLM("mirror", mirrorPattern, *formulaParser),
		delegate.NewMLM("license", licensePattern, *formulaParser, isBeginLicenseSequence, hasUnopenedBrackets, cleanLicenseSequence),
		//delegate.NewSLMM("head", headURLPattern, *formulaParser, []string{headVCSPattern}),
		delegate.NewMLM("head", headURLPattern, *formulaParser, isBeginHeadSequence, isEndHeadSequence, cleanHeadSequence),
		delegate.NewMLM("dependency", dependencyPattern, *formulaParser, isBeginDependencySequence, isEndDependencySequence, cleanDependencySequence),
	}

	results, err := formulaParser.ParseFields(fields)
	if err != nil {
		log.Panicln("Error parsing fields:", err)
		return nil, err
	}

	formula.url = results["url"].(string)
	if results["mirror"] != nil {
		formula.mirror = results["mirror"].(string)
	}
	if results["license"] != nil {
		formula.license = results["license"].(string)
	}

	// Set the license to "pseudo" if it is empty.
	if formula.license == "" {
		formula.license = "pseudo"
	}

	if results["head"] != nil {
		if formulaHead, ok := results["head"].(*head); ok {
			formula.head = formulaHead
			// if len(formulaHead) > 1 {
			// 	formula.head = &head{url: formulaHead[0], vcs: formulaHead[1]}
			// } else {
			// 	formula.head = &head{url: formulaHead[0]}
			// }
		} else {
			headURL := results["head"].(string)
			formula.head = &head{url: headURL}
		}
		log.Println("Head:", formula.head)
	}

	dependencies := make([]*dependency, 0)
	if results["dependency"] != nil {
		for _, dep := range results["dependency"].([][]string) {
			dependency := &dependency{name: dep[0]}
			if len(dep) > 1 {
				dependency.depType = dep[1]
			}
			dependencies = append(dependencies, dependency)
		}
	}
	formula.dependencies = dependencies

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return formula, nil
}

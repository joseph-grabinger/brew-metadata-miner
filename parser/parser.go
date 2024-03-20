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
	"main/parser/types"
)

type parser struct {
	config *config.Config

	// A map of formulas, where the key is the name of the formula.
	formulas map[string]*types.Formula
}

// NewParser creates a new parser.
func NewParser(config *config.Config) *parser {
	return &parser{
		config:   config,
		formulas: make(map[string]*types.Formula),
	}
}

// Parse parses the core repository and extracts the formulas.
func (p *parser) Parse() error {
	return p.readFormulas()
}

func (p *parser) Analyze() {
	valid := make([]*types.Formula, 0)
	noRepo := make([]*types.Formula, 0)

	for _, value := range p.formulas {
		if value.RepoURL != "" {
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

			formula := types.FromSourceFormula(sourceFormula)

			// Add the formula to the formulas map.
			p.formulas[formula.Name] = formula

			//log.Println("Successfully parsed formula:", formula)

			return nil
		}); err != nil {
			log.Printf("Error walking directory: %v\n", err)
		}
	}
	return nil
}

// parseFromFile parses a formula from a file into a Formula struct.
func parseFromFile(file *os.File) (*types.SourceFormula, error) {
	scanner := bufio.NewScanner(file)
	formulaParser := &delegate.FormulaParser{Scanner: scanner}

	base := filepath.Base(file.Name())
	name := strings.TrimSuffix(base, ".rb")

	formula := &types.SourceFormula{Name: name}

	homepage, err := formulaParser.ParseField(homepagePattern, "homepage")
	if err != nil {
		return nil, err
	}
	formula.Homepage = homepage

	fields := []delegate.ParseStrategy{
		delegate.NewSLM("url", urlPattern, *formulaParser),
		delegate.NewSLM("mirror", mirrorPattern, *formulaParser),
		delegate.NewMLM("license", licensePattern, *formulaParser, isBeginLicenseSequence, hasUnopenedBrackets, cleanLicenseSequence),
		delegate.NewMLM("head", headURLPattern, *formulaParser, isBeginHeadSequence, isEndHeadSequence, cleanHeadSequence),
		delegate.NewMLM("dependency", dependencyPattern, *formulaParser, isBeginDependencySequence, isEndDependencySequence, cleanDependencySequence),
	}

	results, err := formulaParser.ParseFields(fields)
	if err != nil {
		log.Panicln("Error parsing fields:", err)
		return nil, err
	}

	formula.URL = results["url"].(string)
	if results["mirror"] != nil {
		formula.Mirror = results["mirror"].(string)
	}
	if results["license"] != nil {
		formula.License = results["license"].(string)
	}

	// Set the license to "pseudo" if it is empty.
	if formula.License == "" {
		formula.License = "pseudo"
	}

	if results["head"] != nil {
		if formulaHead, ok := results["head"].(*types.Head); ok {
			formula.Head = formulaHead
		} else {
			headURL := results["head"].(string)
			formula.Head = &types.Head{URL: headURL}
		}
		// log.Println("Head:", formula.Head)
	}

	dependencies := make([]*types.Dependency, 0)
	if results["dependency"] != nil {
		for _, dep := range results["dependency"].([][]string) {
			dependency := &types.Dependency{Name: dep[0]}
			if len(dep) > 1 {
				dependency.DepType = dep[1]
			}
			dependencies = append(dependencies, dependency)
		}
	}
	formula.Dependencies = dependencies

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return formula, nil
}

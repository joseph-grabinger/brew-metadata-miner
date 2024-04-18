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
	"main/parser/setup"
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

func (p *parser) Pipe() error {
	return p.writeFormulas()
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
			p.formulas[formula.Name] = formula

			log.Println("Successfully parsed formula:", formula)

			return nil
		}); err != nil {
			log.Printf("Error walking directory: %v\n", err)
		}
	}
	return nil
}

func (p *parser) writeFormulas() error {
	path := filepath.Join(p.config.OutputDir, "deps-brew.tsv")
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	// Close file on function exit and check its' returned error.
	defer func() error {
		if err := file.Close(); err != nil {
			return err
		}
		return nil
	}()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, formula := range p.formulas {
		// Write repo line type.
		line := formula.FormatRepoLine()
		_, err := writer.WriteString(line)
		if err != nil {
			return err
		}

		// Write dependency lines.
		for _, dep := range formula.Dependencies {
			f := p.formulas[dep.Name]
			if f == nil {
				log.Println("Dependency not found:", dep.Name)
				continue
			}

			line := formula.FormatDependencyLine(dep)
			_, err := writer.WriteString(line)
			if err != nil {
				return err
			}
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

	homepage, err := formulaParser.ParseField(setup.HomepagePattern, "homepage")
	if err != nil {
		return nil, err
	}
	formula.Homepage = homepage

	fields := setup.BuildStrategies(*formulaParser)

	results, err := formulaParser.ParseFields(fields)
	if err != nil {
		log.Panicln("Error parsing fields:", err)
		return nil, err
	}

	// Set the fields of the formula.
	formula.URL = results["url"].(string)
	if results["mirror"] != nil {
		formula.Mirror = results["mirror"].(string)
	}
	if results["license"] != nil {
		formula.License = results["license"].(string)
	}
	if results["head"] != nil {
		formula.Head = results["head"].(*types.Head)
	}
	if results["dependency"] != nil {
		formula.Dependencies = results["dependency"].([]*types.Dependency)
	} else {
		formula.Dependencies = make([]*types.Dependency, 0)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return formula, nil
}

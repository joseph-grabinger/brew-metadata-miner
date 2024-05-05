package parser

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
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
	// Match the fomula files.
	matches, err := filepath.Glob(p.config.CoreRepo.Dir + "/Formula/**/*.rb")
	if err != nil {
		return err
	}

	// Match alias formula files.
	aliasMatches, err := filepath.Glob(p.config.CoreRepo.Dir + "/Aliases/*")
	if err != nil {
		return err
	}

	matches = append(matches, aliasMatches...)

	for _, path := range matches {
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
		// Write package line.
		line := formula.FormatPackageLine()
		_, err := writer.WriteString(line)
		if err != nil {
			return err
		}

		// Write dependency lines.
		for _, dep := range formula.Dependencies {
			f := p.formulas[dep.Name]
			if f == nil {
				panic("Dependency " + dep.Name + " not found in formula " + formula.Name)
			}

			line := f.FormatDependencyLine(dep)
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
		log.Println("Error parsing fields:", err)
		return nil, err
	}

	// Set the fields of the formula.
	if results["url"] != nil {
		formula.Stable = results["url"].(*types.Stable)
	}
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
		formula.Dependencies = results["dependency"].(*types.Dependencies)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Check if the stable url contains a Ruby string interpolation and resolve it.
	// This is done here rather then in the cleanURLSequence function because the
	// variable used for the interpolation can have a global scope in the formula file.
	// The cleanURLSequence function could only resolve interpolations with a scope within the stable do block.
	found, resolved, err := checkForInterpolation(formula.Stable.URL, file)
	if err != nil {
		return nil, err
	}
	if found {
		formula.Stable.URL = resolved
	}

	return formula, nil
}

// checkForInterpolation checks if the given url contains a Ruby string interpolation.
// If it does, the interpolation is resolved using the given file.
// The function returns a boolean indicating if an interpolation was found and the resolved string.
func checkForInterpolation(url string, file *os.File) (bool, string, error) {
	regex := regexp.MustCompile(setup.InterpolationPattern)
	matches := regex.FindStringSubmatch(url)
	if len(matches) < 2 {
		return false, "", nil
	}

	varName := matches[1]

	// Seek the beginning of the file.
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return true, "", err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains the variable assignment.
		regex := regexp.MustCompile(setup.VarAssignmentPattern(varName))
		varMatches := regex.FindStringSubmatch(line)
		if len(varMatches) >= 2 {
			// Replace the interpolation with the variable value.
			resolved := strings.Replace(url, fmt.Sprintf("#{%s}", varName), varMatches[1], 1)
			return true, resolved, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return true, "", err
	}

	return true, "", fmt.Errorf("could not resolve interpolation in URL: %s", url)
}

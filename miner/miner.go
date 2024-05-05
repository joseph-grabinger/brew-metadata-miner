package miner

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
	"main/miner/parser"
	"main/miner/setup"
	"main/miner/types"
)

type miner struct {
	config *config.Config

	// A map of formulae, where the key is the name of the formula.
	formulae map[string]*types.Formula
}

// NewMiner creates a new parser.
func NewMiner(config *config.Config) *miner {
	return &miner{
		config:   config,
		formulae: make(map[string]*types.Formula),
	}
}

// ReadFormaulas reads all formulae from the core repository into the formulas map.
func (p *miner) ReadFormulae() error {
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
		sourceFormula, err := extractFromFile(file)
		if err != nil {
			log.Printf("Error parsing file %s: %v\n", path, err)
			return err
		}

		formula := types.FromSourceFormula(sourceFormula)
		p.formulae[formula.Name] = formula

		log.Println("Successfully parsed formula:", formula)
	}

	return nil
}

// WriteFormulae writes the formulae to the output file.
func (p *miner) WriteFormulae() error {
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

	for _, formula := range p.formulae {
		// Write package line.
		line := formula.FormatPackageLine()
		_, err := writer.WriteString(line)
		if err != nil {
			return err
		}

		// Write dependency lines.
		for _, dep := range formula.Dependencies {
			f := p.formulae[dep.Name]
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

// extractFromFile extracts a formula from a file and returns it as a Formula struct.
func extractFromFile(file *os.File) (*types.SourceFormula, error) {
	scanner := bufio.NewScanner(file)
	formulaParser := &parser.FormulaParser{Scanner: scanner}

	base := filepath.Base(file.Name())
	name := strings.TrimSuffix(base, ".rb")

	formula := &types.SourceFormula{Name: name}

	fields := setup.BuildStrategies(*formulaParser)

	results, err := formulaParser.ParseFields(fields)
	if err != nil {
		log.Println("Error parsing fields:", err)
		return nil, err
	}

	// Set the fields of the formula.
	if results["homepage"] != nil {
		formula.Homepage = results["homepage"].(string)
	}
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

package parser

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"main/config"
	"os"
	"path/filepath"
	"regexp"
)

type parser struct {
	config *config.Config

	// A map of formulas, where the key is the name of the formula.
	formulas map[string]*Formula
}

// NewParser creates a new parser.
func NewParser(config *config.Config) *parser {
	return &parser{
		config:   config,
		formulas: make(map[string]*Formula),
	}
}

// Parse parses the core repository and extracts the formulas.
func (p *parser) Parse() error {
	return p.readFormulas()

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
				log.Printf("Error opening file %s: %v\n", path, err)
				return err
			}
			defer file.Close()

			// Parse Formula from file.
			formula, err := parseFromFile(file)
			if err != nil {
				log.Printf("Error parsing file %s: %v\n", path, err)
				return err
			}

			// Add the formula to the formulas map.
			p.formulas[formula.Name] = formula

			log.Println("Successfully parsed formula: ", formula)

			return nil
		}); err != nil {
			log.Printf("Error walking directory: %v\n", err)
		}
	}
	return nil
}

// parseFromFile parses a formula from a file into a Formula struct.
func parseFromFile(file *os.File) (*Formula, error) {
	scanner := bufio.NewScanner(file)

	nameField := &field{name: "name", pattern: `class\s([a-zA-Z0-9]+)\s<\sFormula`}
	name, err := parseField(scanner, nameField)
	if err != nil {
		return nil, err
	}
	formula := &Formula{Name: name}

	homepageField := &field{name: "homepage", pattern: `homepage\s+"([^"]+)"`}
	homepage, err := parseField(scanner, homepageField)
	if err != nil {
		return nil, err
	}
	formula.Homepage = homepage

	fields := []*field{
		{
			name:    "url",
			pattern: `url\s+"([^"]+)"`,
		},
		{
			name:    "mirror",
			pattern: `mirror\s+"([^"]+)"`,
		},
		{
			name:    "license",
			pattern: `license\s+"([^"]+)"`,
		},
	}
	results, err := parseFields(scanner, fields)
	if err != nil {
		log.Panicln("Error parsing fields:", err)
		return nil, err
	}

	formula.URL = results[fields[0]]
	formula.Mirror = results[fields[1]]
	formula.License = results[fields[2]]
	if formula.License == "" {
		return formula, errors.New("no license found for formula")
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error scanning file:", err)
		return nil, err
	}

	return formula, nil
}

// parseField parses a specified field from a formula through a given scanner.
func parseField(scanner *bufio.Scanner, field *field) (string, error) {
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("Line: %s: %s\n", field.name, line)

		regex := regexp.MustCompile(field.pattern)
		matches := regex.FindStringSubmatch(line)

		if len(matches) >= 2 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("no %s found for formula", field.name)
}

// parseFields parses a specified list of fields from a formula through a given scanner.
func parseFields(scanner *bufio.Scanner, fields []*field) (map[*field]string, error) {
	results := make(map[*field]string)

	for scanner.Scan() {
		line := scanner.Text()
		log.Println("Line: <generic>: ", line)

		for _, f := range fields {
			pattern := f.pattern
			regex := regexp.MustCompile(pattern)
			matches := regex.FindStringSubmatch(line)

			if len(matches) >= 2 {
				results[f] = matches[1]
				log.Println("Matched: ", results[f])
				break
			}
		}
	}

	return results, nil
}

type field struct {
	name    string
	pattern string
}

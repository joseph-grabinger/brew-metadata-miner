package parser

import (
	"bytes"
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
	return &parser{config: config}
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

			// Read the content of the file.
			content, err := os.ReadFile(path)
			if err != nil {
				log.Printf("Error reading file %s: %v\n", path, err)
				return err
			}

			// Parse Formula from file.
			_, err = parseFromFile(content)
			if err != nil {
				log.Printf("Error parsing file %s: %v\n", path, err)
				return err
			}

			// Add the formula to the formulas map.
			//p.formulas[formula.Name] = formula

			return nil
		}); err != nil {
			log.Printf("Error walking directory: %v\n", err)
		}
	}
	return nil
}

// parseFromFile parses a formula from a file into a Formula struct.
func parseFromFile(content []byte) (*Formula, error) {
	// Find name in first line of content.
	namePattern := regexp.MustCompile(`class\s+([^\s\W]+)`)
	n := bytes.IndexByte(content, '\n')
	firstLine, content := content[:n], content[n+1:]
	nameMatches := namePattern.FindStringSubmatch(string(firstLine))
	if len(nameMatches) <= 1 {
		return nil, ErrInvalidFormula
	}

	formula := &Formula{Name: nameMatches[1]}
	log.Println(formula.Name)

	// Find homepage in remaining content.
	homepagePattern := regexp.MustCompile(`homepage\s+['"]([^'"]+)['"]`)
	homepageMatches := homepagePattern.FindStringSubmatch(string(content))
	if len(homepageMatches) <= 1 {
		return nil, ErrInvalidFormula
	}

	formula.Homepage = homepageMatches[1]

	return formula, nil
}

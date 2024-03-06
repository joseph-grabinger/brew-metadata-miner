package parser

import (
	"bufio"
	"errors"
	"log"
	"main/config"
	"os"
	"path/filepath"
	"regexp"
	"strings"
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

			// // Read the content of the file.
			// content, err := os.ReadFile(path)
			// if err != nil {
			// 	log.Printf("Error reading file %s: %v\n", path, err)
			// 	return err
			// }

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

	name, err := parseName(scanner)
	if err != nil {
		return nil, err
	}
	formula := &Formula{Name: name}

	homepage, err := parseHomepage(scanner)
	if err != nil {
		return nil, err
	}
	formula.Homepage = homepage

	url, err := parseURL(scanner)
	if err != nil {
		return nil, err
	}
	formula.URL = url

	mirror, err := parseMirror(scanner)
	if err == nil { // mirror is optional
		formula.Mirror = mirror
	}

	license, err := parseLicense(scanner)
	if err != nil {
		return nil, err
	}
	formula.License = license

	if err := scanner.Err(); err != nil {
		log.Println("Error scanning file:", err)
		return nil, err
	}

	return formula, nil
}

// parseName parses the name of a formula through a given scanner.
func parseName(scanner *bufio.Scanner) (string, error) {
	for scanner.Scan() {
		line := scanner.Text()
		log.Println("Line: name: ", line)

		if strings.HasPrefix(line, "class") && strings.Contains(line, "< Formula") {
			parts := strings.Fields(line)
			if len(parts) < 2 {
				log.Println("Invalid class line: ", line)
				return "", ErrInvalidFormula
			}

			return parts[1], nil
		}
	}

	return "", errors.New("no name found for formula")
}

// parsehomepage parses the homepage of a formula through a given scanner.
func parseHomepage(scanner *bufio.Scanner) (string, error) {
	for scanner.Scan() {
		line := scanner.Text()
		log.Println("Line: homepage: ", line)

		pattern := `homepage\s+"([^"]+)"`
		regex := regexp.MustCompile(pattern)
		matches := regex.FindStringSubmatch(line)

		if len(matches) >= 2 {
			return matches[1], nil
		}
	}

	return "", errors.New("no homepage found for formula")
}

// parseURL parses the URL of a formula through a given scanner.
func parseURL(scanner *bufio.Scanner) (string, error) {
	for scanner.Scan() {
		line := scanner.Text()
		log.Println("Line: url: ", line)

		pattern := `url\s+"([^"]+)"`
		regex := regexp.MustCompile(pattern)
		matches := regex.FindStringSubmatch(line)

		if len(matches) >= 2 {
			return matches[1], nil
		}
	}

	return "", errors.New("no url found for formula")
}

// parseMirror parses the mirror of a formula through a given scanner.
func parseMirror(scanner *bufio.Scanner) (string, error) {
	for scanner.Scan() {
		line := scanner.Text()
		log.Println("Line: mirror: ", line)

		pattern := `mirror\s+"([^"]+)"`
		regex := regexp.MustCompile(pattern)
		matches := regex.FindStringSubmatch(line)

		if len(matches) >= 2 {
			return matches[1], nil
		}
	}

	return "", errors.New("no mirror found for formula")
}

// parseLicense parses the license of a formula through a given scanner.
func parseLicense(scanner *bufio.Scanner) (string, error) {
	for scanner.Scan() {
		line := scanner.Text()
		log.Println("Line: license: ", line)

		pattern := `license\s+"([^"]+)"`
		regex := regexp.MustCompile(pattern)
		matches := regex.FindStringSubmatch(line)

		if len(matches) >= 2 {
			return matches[1], nil
		}
	}

	return "", errors.New("no license found for formula")
}

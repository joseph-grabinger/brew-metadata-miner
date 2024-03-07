package parser

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"main/config"
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

	nameField := &singleLineField{name: "name", pattern: `class\s([a-zA-Z0-9]+)\s<\sFormula`}
	name, err := parseField(scanner, nameField)
	if err != nil {
		return nil, err
	}
	formula := &Formula{Name: name}

	homepageField := &singleLineField{name: "homepage", pattern: `homepage\s+"([^"]+)"`}
	homepage, err := parseField(scanner, homepageField)
	if err != nil {
		return nil, err
	}
	formula.Homepage = homepage

	formulaParser := &FormulaParser{scanner: scanner}

	fields := []field{
		&singleLineField{name: "url", pattern: `url\s+"([^"]+)"`, strat: NewSLS(*formulaParser)},
		&singleLineField{name: "mirror", pattern: `mirror\s+"([^"]+)"`, strat: NewSLS(*formulaParser)},
		&multiLineField{field: &singleLineField{name: "license", pattern: `license\s+(:\w+|all_of\s*:\s*\[[^\]]+\]|any_of\s*:\s*\[[^\]]+\]|"[^"]+")`, strat: NewMLS(*formulaParser)}, isBeginSequence: isLicenseBeginSequence, isEndSequence: isLicenseEndSequence}, // `license\s+"([^"]+)"`
	}

	results, err := formulaParser.parseFields(fields)
	if err != nil {
		log.Panicln("Error parsing fields:", err)
		return nil, err
	}

	formula.URL = results[fields[0]]
	formula.Mirror = results[fields[1]]
	formula.License = results[fields[2]]
	if formula.License == "" {
		log.Println("NO license found for formula")
		formula.License = "No license found" // TODO: Change to a default license.
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error scanning file:", err)
		return nil, err
	}

	return formula, nil
}

// parseField parses a specified field from a formula through a given scanner.
func parseField(scanner *bufio.Scanner, field field) (string, error) {
	for scanner.Scan() {
		line := scanner.Text()
		log.Printf("Line: %s: %s\n", field.getName(), line)

		regex := regexp.MustCompile(field.getPattern())
		matches := regex.FindStringSubmatch(line)

		if len(matches) >= 2 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("no %s found for Formula", field.getName())
}

type field interface {
	getName() string
	getPattern() string
	getStrat() parseStrat
}

type singleLineField struct {
	name    string
	pattern string
	strat   parseStrat
}

func (slf singleLineField) getName() string {
	return slf.name
}

func (slf singleLineField) getPattern() string {
	return slf.pattern
}

func (slf singleLineField) getStrat() parseStrat {
	return slf.strat
}

type multiLineField struct {
	field           // embedding field struct
	isBeginSequence func(line string) bool
	isEndSequence   func(line string) bool
}

func (mlf multiLineField) getName() string {
	return mlf.field.getName()
}

func (mlf multiLineField) getPattern() string {
	return mlf.field.getPattern()
}

func (mlf multiLineField) getStrat() parseStrat {
	return mlf.field.getStrat()
}

type parseStrat interface {
	matchesField(field field, line string) bool
	extractField(field field, line string) (string, error)
}

type FormulaParser struct {
	scanner *bufio.Scanner
}

func (fp *FormulaParser) parseFields(fields []field) (map[field]string, error) {
	results := make(map[field]string)

	for fp.scanner.Scan() {
		line := fp.scanner.Text()

		for _, f := range fields {
			// Skip field if it has already been matched.
			if _, ok := results[f]; ok {
				continue
			}

			log.Printf("Line: <generic:%s>: %s\n", f.getName(), line)
			strat := f.getStrat()
			if strat.matchesField(f, line) {
				fieldValue, err := strat.extractField(f, line)
				if err != nil {
					return nil, err
				}
				results[f] = fieldValue
				log.Println("Matched: ", results[f])
				break
			}
		}
	}

	return results, nil
}

type singleLineStrategy struct {
	FormulaParser
	matches []string
}

func NewSLS(fp FormulaParser) *singleLineStrategy {
	return &singleLineStrategy{
		FormulaParser: fp,
		matches:       make([]string, 0),
	}
}

func (sls *singleLineStrategy) matchesField(field field, line string) bool {
	regex := regexp.MustCompile(field.getPattern())
	matches := regex.FindStringSubmatch(line)

	isMatch := len(matches) >= 2
	if isMatch {
		sls.matches = matches
	}
	return isMatch
}

func (sls *singleLineStrategy) extractField(field field, line string) (string, error) {
	return sls.matches[1], nil
}

type multiLineStrategy struct {
	FormulaParser
	matches []string
	opened  bool
}

func NewMLS(fp FormulaParser) *multiLineStrategy {
	return &multiLineStrategy{
		FormulaParser: fp,
		matches:       make([]string, 0),
		opened:        false,
	}
}

func (mls *multiLineStrategy) matchesField(field field, line string) bool {
	// Check for begin sequence.
	if field.(*multiLineField).isBeginSequence(line) {
		mls.opened = true
		mls.matches = append(mls.matches, line)
		return true
	}

	// Check for default field pattern.
	regex := regexp.MustCompile(field.getPattern())
	matches := regex.FindStringSubmatch(line)

	isMatch := len(matches) >= 2
	if isMatch {
		mls.matches = matches
	}
	return isMatch
}

func (mls *multiLineStrategy) extractField(field field, line string) (string, error) {
	// A not open sequence means the field's default pattern has been matched.
	// Thus, the field can be extracted from a single line.
	if !mls.opened {
		log.Println("Extracted at once: ", mls.matches[1])
		return mls.matches[1], nil
	}

	if mls.FormulaParser.scanner == nil {
		return "", errors.New("no scanner found for multi-line field")
	}

	for mls.FormulaParser.scanner.Scan() {
		line := mls.FormulaParser.scanner.Text()
		log.Println("Line: <genericMLS>: ", line)

		if field.(*multiLineField).isEndSequence(line) {
			mls.matches = append(mls.matches, line)
			return strings.Join(mls.matches, ""), nil
		}

		// Append line to matches since the sequence has been opened.
		mls.matches = append(mls.matches, line)
	}

	log.Println("Current matches: ", mls.matches)

	return "", fmt.Errorf("no %s found for formula", field.getName())
}

func isLicenseBeginSequence(line string) bool {
	openCount, closeCount := 0, 0
	for _, char := range line {
		switch char {
		case '[':
			openCount++
		case ']':
			closeCount++
		}
	}
	return openCount > closeCount
}

func isLicenseEndSequence(line string) bool {
	openCount, closeCount := 0, 0
	for _, char := range line {
		switch char {
		case '[':
			openCount++
		case ']':
			closeCount++
		}
	}
	return openCount < closeCount
}

package delegate

import (
	"fmt"
	"log"
	"regexp"
)

// parseStrategy is the strategy interface.
type ParseStrategy interface {
	matchesLine(line string) bool
	extractFromLine(line string) (interface{}, error)
	GetName() string
	GetPattern() string
}

// singleLineField acts a concrete component.
// It is a concrete implementation of the singleLineField interface.
type singleLineField struct {
	name    string
	pattern string
	FormulaParser
	matches []string
}

// NewSLF returns a pointer to a new instance of singleLineField.
func NewSLF(name string, pattern string, fp FormulaParser) *singleLineField {
	return &singleLineField{
		name:          name,
		pattern:       pattern,
		FormulaParser: fp,
		matches:       make([]string, 0),
	}
}

func (f *singleLineField) GetName() string {
	return f.name
}

func (f *singleLineField) GetPattern() string {
	return f.pattern
}

func (f *singleLineField) matchesLine(line string) bool {
	regex := regexp.MustCompile(f.pattern)
	matches := regex.FindStringSubmatch(line)

	isMatch := len(matches) >= 2
	if isMatch {
		f.matches = matches
	}
	return isMatch
}

func (f *singleLineField) extractFromLine(line string) (interface{}, error) {
	return f.matches[1], nil
}

// multiLineField acts a decorator.
// It embeds the Field interface and adds additional functionality.
type multiLineField struct {
	singleLineField
	isBeginSequence func(line string) bool
	isEndSequence   func(line string) bool
	cleanSequence   func(sequence []string) interface{}
	opened          bool
}

// NewMLF returns a pointer to a new instance of multiLineField.
func NewMLF(name string, pattern string, fp FormulaParser, isBeginSeq func(line string) bool, isEndSeq func(line string) bool, clean func([]string) interface{}) *multiLineField {
	return &multiLineField{
		singleLineField: *NewSLF(name, pattern, fp),
		isBeginSequence: isBeginSeq,
		isEndSequence:   isEndSeq,
		cleanSequence:   clean,
		opened:          false,
	}
}

func (f *multiLineField) matchesLine(line string) bool {
	// Check for begin sequence.
	if f.isBeginSequence(line) {
		f.opened = true
		f.matches = append(f.matches, line)
		return true
	}

	// Check for default field pattern.
	regex := regexp.MustCompile(f.pattern)
	matches := regex.FindStringSubmatch(line)

	isMatch := len(matches) >= 2
	if isMatch {
		f.matches = matches
	}
	return isMatch
}

func (f *multiLineField) extractFromLine(line string) (interface{}, error) {
	// A not open sequence means the field's default pattern has been matched.
	// Thus, the field can be extracted from a single line.
	if !f.opened {
		return f.matches[1], nil
	}

	for f.FormulaParser.Scanner.Scan() {
		line := f.FormulaParser.Scanner.Text()
		log.Println("Line: <genericMLS>: ", line)

		if f.isEndSequence(line) {
			f.matches = append(f.matches, line)
			cleaned := f.cleanSequence(f.matches)
			return cleaned, nil
		}

		// Append line to matches since the sequence has been opened.
		f.matches = append(f.matches, line)
	}

	return "", fmt.Errorf("no %s found for formula", f.name)
}

// sameLineMultiField acts as a decorator.
// It embeds the Field interface and adds additional functionality.
type sameLineMultiField struct {
	singleLineField
	additionalPatterns []string
	line               string
}

// NewSLMF returns a pointer to a new instance of sameLineMultiField.
func NNewSLMF(name string, pattern string, fp FormulaParser, additionalPatterns []string) *sameLineMultiField {
	return &sameLineMultiField{
		singleLineField:    *NewSLF(name, pattern, fp),
		additionalPatterns: additionalPatterns,
	}
}

func (f *sameLineMultiField) matchesLine(line string) bool {
	match := f.singleLineField.matchesLine(line)
	if match {
		f.line = line
	}
	return match
}

func (f *sameLineMultiField) extractFromLine(line string) (interface{}, error) {
	singleMatch, error := f.singleLineField.extractFromLine(f.line)
	if error != nil {
		return "", error
	}

	res := []string{singleMatch.(string)}

	// Check additional patterns.
	patterns := f.additionalPatterns
	if len(patterns) > 0 {
		for _, pattern := range patterns {
			regex := regexp.MustCompile(pattern)
			matches := regex.FindStringSubmatch(f.line)

			if len(matches) >= 2 {
				res = append(res, matches[1])
			}
		}
	}

	return res, nil
}

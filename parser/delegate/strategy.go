package delegate

import (
	"fmt"
	"log"
	"regexp"
)

// ParseStrategy represents a strategy for parsing lines and extracting relevant information.
// It declares all methods the context uses to execute a strategy.
type ParseStrategy interface {
	// matchesLine checks if the strategy matches the given line.
	matchesLine(line string) bool

	// extractFromLine extracts relevant information from the given line.
	// It returns the extracted information and an error if any.
	extractFromLine(line string) (interface{}, error)

	// getName returns the name of the field.
	getName() string

	// getPattern returns the pattern used by the strategy.
	getPattern() string
}

// singleLineMatcher acts as a concrete strategy.
type singleLineMatcher struct {
	// FormulaParser is the context for parsing fields.
	FormulaParser

	// name of the field.
	name string

	// pattern to match the field.
	pattern string

	// matches is the slice of matches found in the field.
	matches []string
}

// NewSLM creates a new instance of singleLineMatcher.
func NewSLM(name string, pattern string, fp FormulaParser) *singleLineMatcher {
	return &singleLineMatcher{
		name:          name,
		pattern:       pattern,
		FormulaParser: fp,
		matches:       make([]string, 0),
	}
}

func (f *singleLineMatcher) getName() string {
	return f.name
}

func (f *singleLineMatcher) getPattern() string {
	return f.pattern
}

// matchesLine checks if the given line matches the pattern defined in the singleLineMatcher.
// It returns true if there is a match, and false otherwise.
// If there is a match it also sets the matches property to the matches found in the line.
func (f *singleLineMatcher) matchesLine(line string) bool {
	regex := regexp.MustCompile(f.pattern)
	matches := regex.FindStringSubmatch(line)

	isMatch := len(matches) >= 2
	if isMatch {
		f.matches = matches
	}
	return isMatch
}

// extractFromLine returns the previously matched information at index 1 of the matches slice.
func (f *singleLineMatcher) extractFromLine(line string) (interface{}, error) {
	return f.matches[1], nil
}

// MultiLineMatcher acts as a concrete strategy that matches multiple lines of text.
// It embeds the singleLineMatcher and adds additional functionality.
type MultiLineMatcher struct {
	// Inherits properties and methods from the singleLineMatcher.
	singleLineMatcher

	// isBeginSequence checks if a line is the beginning of a sequence.
	isBeginSequence func(line string) bool

	// isEndSequence checks if a line is the end of a sequence.
	isEndSequence func(line string) bool

	// cleanSequence cleans and processes the matched sequence.
	// It is called with the matches after the end sequence has been matched.
	cleanSequence func(sequence []string) interface{}

	// Flag to indicate if a begin sequence has been matched.
	opened bool
}

// NewMLM creates a new instance of multiLineMatcher.
func NewMLM(name string, pattern string, fp FormulaParser, isBeginSeq func(line string) bool, isEndSeq func(line string) bool, clean func([]string) interface{}) *MultiLineMatcher {
	return &MultiLineMatcher{
		singleLineMatcher: *NewSLM(name, pattern, fp),
		isBeginSequence:   isBeginSeq,
		isEndSequence:     isEndSeq,
		cleanSequence:     clean,
		opened:            false,
	}
}

// MatchesLine checks if the given line contains a begin sequence.
// If a begin sequence is found, the sequence is opened and the line is appended to the matches slice.
// Else the line is checked against the default pattern of the field using the singleLineMatcher.
func (f *MultiLineMatcher) MatchesLine(line string) bool {
	// Check for begin sequence.
	if f.isBeginSequence(line) {
		log.Println("Begin sequence found: ", line)
		f.opened = true
		f.matches = append(f.matches, line)
		return true
	}

	// Check for default field pattern.
	return f.singleLineMatcher.matchesLine(line)
}

// ExtractFromLine returns the previously matched information at index 1
// of the matches slice if the opened flag is set.
// Else it matches all lines until the end sequence is found.
func (f *MultiLineMatcher) ExtractFromLine(line string) (interface{}, error) {
	// A not open sequence means the field's default pattern has been matched.
	// Thus, the field can be extracted from a single line.
	if !f.opened {
		return f.matches[1], nil
	}

	for f.FormulaParser.Scanner.Scan() {
		line := f.FormulaParser.Scanner.Text()
		log.Println("Line: <genericMLS>: ", line)

		if f.isEndSequence(line) {
			log.Println("End sequence found: ", line)
			f.matches = append(f.matches, line)
			cleaned := f.cleanSequence(f.matches)
			return cleaned, nil
		}

		// Append line to matches since the sequence has been opened.
		f.matches = append(f.matches, line)
	}

	return "", fmt.Errorf("no %s found for formula", f.name)
}

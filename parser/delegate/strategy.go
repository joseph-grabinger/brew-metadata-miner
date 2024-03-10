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

// multiLineMatcher acts as a concrete strategy that matches multiple lines of text.
// It embeds the singleLineMatcher and adds additional functionality.
type multiLineMatcher struct {
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
func NewMLM(name string, pattern string, fp FormulaParser, isBeginSeq func(line string) bool, isEndSeq func(line string) bool, clean func([]string) interface{}) *multiLineMatcher {
	return &multiLineMatcher{
		singleLineMatcher: *NewSLM(name, pattern, fp),
		isBeginSequence:   isBeginSeq,
		isEndSequence:     isEndSeq,
		cleanSequence:     clean,
		opened:            false,
	}
}

// matchesLine checks if the given line contains a begin sequence.
// If a begin sequence is found, the sequence is opened and the line is appended to the matches slice.
// Else the line is checked against the default pattern of the field using the singleLineMatcher.
func (f *multiLineMatcher) matchesLine(line string) bool {
	// Check for begin sequence.
	if f.isBeginSequence(line) {
		f.opened = true
		f.matches = append(f.matches, line)
		return true
	}

	// Check for default field pattern.
	return f.singleLineMatcher.matchesLine(line)
}

// extractFromLine returns the previously matched information at index 1
// of the matches slice if the opened flag is set.
// Else it matches all lines until the end sequence is found.
func (f *multiLineMatcher) extractFromLine(line string) (interface{}, error) {
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

// singleLineMultiMatcher acts as a concrete strategy that matches patterns within a single line of text.
// It embeds the singleLineMatcher and adds additional functionality.
type singleLineMultiMatcher struct {
	// Inherits properties and methods from the singleLineMatcher.
	singleLineMatcher

	// additionalPatterns are additional patterns to match within the line.
	additionalPatterns []string

	// line is the line of text that was matched.
	line string
}

// NewSLMM creates a new instance of singleLineMultiMatcher.
func NewSLMM(name string, pattern string, fp FormulaParser, additionalPatterns []string) *singleLineMultiMatcher {
	return &singleLineMultiMatcher{
		singleLineMatcher:  *NewSLM(name, pattern, fp),
		additionalPatterns: additionalPatterns,
	}
}

// matchesLine checks if the given line contains a match for the default pattern of the field
// and stores the matched line in its line field.
func (f *singleLineMultiMatcher) matchesLine(line string) bool {
	match := f.singleLineMatcher.matchesLine(line)
	if match {
		f.line = line
	}
	return match
}

// extractFromLine returns a list of all matches found in the line.
// This includes the default pattern and any additional patterns.
func (f *singleLineMultiMatcher) extractFromLine(line string) (interface{}, error) {
	singleMatch, error := f.singleLineMatcher.extractFromLine(f.line)
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

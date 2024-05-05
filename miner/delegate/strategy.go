package delegate

import (
	"fmt"
)

// ParseStrategy represents a strategy for parsing lines and extracting relevant information.
// It declares all methods the context uses to execute a strategy.
type ParseStrategy interface {
	// MatchesLine checks if the strategy matches the given line.
	MatchesLine(line string) bool

	// ExtractFromLine extracts relevant information from the given line.
	// It returns the extracted information and an error if any.
	ExtractFromLine(line string) (interface{}, error)

	// getName returns the name of the field.
	getName() string
}

// SingleLineMatcher acts as a concrete strategy.
type SingleLineMatcher[T any] struct {
	// FormulaParser is the context for parsing fields.
	FormulaParser

	// name of the field.
	name string

	// isDefaultPattern to match the field.
	isDefaultPattern func(string) (bool, []string)

	// matches is the slice of matches found in the field.
	matches []string
}

// NewSLM creates a new instance of singleLineMatcher.
func NewSLM[T any](name string, isDefaultPattern func(string) (bool, []string), fp FormulaParser) *SingleLineMatcher[T] {
	return &SingleLineMatcher[T]{
		name:             name,
		isDefaultPattern: isDefaultPattern,
		FormulaParser:    fp,
		matches:          make([]string, 0),
	}
}

func (f *SingleLineMatcher[T]) getName() string {
	return f.name
}

// MatchesLine checks if the given line matches the pattern defined in the singleLineMatcher.
// It returns true if there is a match, and false otherwise.
// If there is a match it also sets the matches property to the matches found in the line.
func (f *SingleLineMatcher[T]) MatchesLine(line string) bool {
	b, matches := f.isDefaultPattern(line)
	if b {
		f.matches = matches
	}
	return b
}

// ExtractFromLine returns the previously matched information at index 1 of the matches slice.
func (f *SingleLineMatcher[T]) ExtractFromLine(line string) (interface{}, error) {
	var i interface{} = f.matches[1]
	return i.(T), nil
}

// MultiLineMatcher acts as a concrete strategy that matches multiple lines of text.
// It embeds the singleLineMatcher and adds additional functionality.
type MultiLineMatcher[T any] struct {
	// Inherits properties and methods from the SingleLineMatcher.
	SingleLineMatcher[T]

	// isBeginSequence checks if a line is the beginning of a sequence.
	isBeginSequence func(line string) bool

	// isEndSequence checks if a line is the end of a sequence.
	isEndSequence func(line string) bool

	// cleanSequence cleans and processes the matched sequence.
	// It is called with the matches after the end sequence has been matched.
	cleanSequence func(sequence []string) T

	// Flag to indicate if a begin sequence has been matched.
	opened bool
}

// NewMLM creates a new instance of multiLineMatcher.
func NewMLM[T any](name string, isDefaultPattern func(string) (bool, []string), fp FormulaParser, isBeginSeq func(line string) bool, isEndSeq func(line string) bool, clean func([]string) T) *MultiLineMatcher[T] {
	return &MultiLineMatcher[T]{
		SingleLineMatcher: *NewSLM[T](name, isDefaultPattern, fp),
		isBeginSequence:   isBeginSeq,
		isEndSequence:     isEndSeq,
		cleanSequence:     clean,
		opened:            false,
	}
}

// MatchesLine checks if the given line contains a begin sequence.
// If a begin sequence is found, the sequence is opened and the line is appended to the matches slice.
// Else the line is checked against the default pattern of the field using the singleLineMatcher.
func (f *MultiLineMatcher[T]) MatchesLine(line string) bool {
	// Check for begin sequence.
	if f.isBeginSequence(line) {
		f.opened = true
		f.matches = append(f.matches, line)
		return true
	}

	// Check for default field pattern.
	return f.SingleLineMatcher.MatchesLine(line)
}

// ExtractFromLine returns the previously matched information at index 1
// of the matches slice if the opened flag is set.
// Else it matches all lines until the end sequence is found.
func (f *MultiLineMatcher[T]) ExtractFromLine(line string) (interface{}, error) {
	// A not open sequence means the field's default pattern has been matched.
	// Thus, the field can be extracted from a single line.
	if !f.opened {
		// Check if type is T.
		var i interface{} = f.matches[1]
		if v, ok := i.(T); ok {
			return v, nil
		}

		cleaned := f.cleanSequence([]string{f.matches[1]})
		return cleaned, nil
	}

	for f.FormulaParser.Scanner.Scan() {
		line := f.FormulaParser.Scanner.Text()

		if f.isEndSequence(line) {
			f.matches = append(f.matches, line)
			cleaned := f.cleanSequence(f.matches)
			return cleaned, nil
		}

		// Append line to matches since the sequence has been opened.
		f.matches = append(f.matches, line)
	}

	return nil, fmt.Errorf("no %s found for formula", f.name)
}

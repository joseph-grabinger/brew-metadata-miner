package delegate

import (
	"fmt"
	"log"
	"regexp"
)

// parseStrategy is the strategy interface.
type parseStrategy interface {
	matchesField(field Field, line string) bool
	extractField(field Field, line string) (interface{}, error)
}

// SingleLineStrategy is a concrete strategy.
type SingleLineStrategy struct {
	// FormulaParser is the context.
	FormulaParser
	matches []string
}

// NewSLS returns a pointer to a new instance of SingleLineStrategy.
func NewSLS(fp FormulaParser) *SingleLineStrategy {
	return &SingleLineStrategy{
		FormulaParser: fp,
		matches:       make([]string, 0),
	}
}

// matchesField checks if a given field's pattern matches a line.
func (sls *SingleLineStrategy) matchesField(field Field, line string) bool {
	regex := regexp.MustCompile(field.GetPattern())
	matches := regex.FindStringSubmatch(line)

	isMatch := len(matches) >= 2
	if isMatch {
		sls.matches = matches
	}
	return isMatch
}

// extractField returns the previously matched field.
// matchesField has to be called first.
func (sls *SingleLineStrategy) extractField(field Field, line string) (interface{}, error) {
	return sls.matches[1], nil
}

// MultiLineStrategy is a concrete strategy.
type MultiLineStrategy struct {
	// FormulaParser is the context.
	FormulaParser
	matches []string
	opened  bool
}

// NewMLS returns a pointer to a new instance of MultiLineStrategy.
func NewMLS(fp FormulaParser) *MultiLineStrategy {
	return &MultiLineStrategy{
		FormulaParser: fp,
		matches:       make([]string, 0),
		opened:        false,
	}
}

// matchesField checks if a given field's pattern matches a line.
// It first checks for the begin sequence.
// If the sequence is not found, it checks for the default field pattern.
func (mls *MultiLineStrategy) matchesField(field Field, line string) bool {
	// Check for begin sequence.
	if field.(*multiLineField).isBeginSequence(line) {
		mls.opened = true
		mls.matches = append(mls.matches, line)
		return true
	}

	// Check for default field pattern.
	regex := regexp.MustCompile(field.GetPattern())
	matches := regex.FindStringSubmatch(line)

	isMatch := len(matches) >= 2
	if isMatch {
		mls.matches = matches
	}
	return isMatch
}

// extractField returns the previously matched field if the sequence has been opened.
// If the sequence has been opened, it returns the entire multi-line sequence.
// matchesField has to be called first.
func (mls *MultiLineStrategy) extractField(field Field, line string) (interface{}, error) {
	// A not open sequence means the field's default pattern has been matched.
	// Thus, the field can be extracted from a single line.
	if !mls.opened {
		return mls.matches[1], nil
	}

	for mls.FormulaParser.Scanner.Scan() {
		line := mls.FormulaParser.Scanner.Text()
		log.Println("Line: <genericMLS>: ", line)

		if field.(*multiLineField).isEndSequence(line) {
			mls.matches = append(mls.matches, line)
			cleaned := field.(*multiLineField).cleanSequence(mls.matches)
			return cleaned, nil
		}

		// Append line to matches since the sequence has been opened.
		mls.matches = append(mls.matches, line)
	}

	return "", fmt.Errorf("no %s found for formula", field.GetName())
}

// SameLineMultiStrategy is a concrete strategy.
type SameLineMultiStrategy struct {
	SingleLineStrategy
	line string
}

// NewSLMS returns a pointer to a new instance of SameLineMultiStrategy.
func NewSLMS(fp FormulaParser) *SameLineMultiStrategy {
	return &SameLineMultiStrategy{
		SingleLineStrategy: *NewSLS(fp),
	}
}

// matchesField checks if a given field's pattern matches a line.
func (slms *SameLineMultiStrategy) matchesField(field Field, line string) bool {
	match := slms.SingleLineStrategy.matchesField(field, line)
	if match {
		slms.line = line
	}
	return match
}

// extractField returns the previously matched field.
// matchesField has to be called first.
func (slms *SameLineMultiStrategy) extractField(field Field, line string) (interface{}, error) {
	singleMatch, error := slms.SingleLineStrategy.extractField(field, slms.line)
	if error != nil {
		return "", error
	}

	res := []string{singleMatch.(string)}

	// Check additional patterns.
	patterns := field.(*sameLineMultiField).additionalPatterns
	if len(patterns) > 0 {
		for _, pattern := range patterns {
			regex := regexp.MustCompile(pattern)
			matches := regex.FindStringSubmatch(slms.line)

			if len(matches) >= 2 {
				res = append(res, matches[1])
			}
		}
	}

	return res, nil
}

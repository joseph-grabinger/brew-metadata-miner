package parser

import (
	"bufio"
)

// FormulaParser acts as context for parsing fields.
type FormulaParser struct {
	// Scanner is used to read the file line by line.
	Scanner *bufio.Scanner
}

// ParseFields parses the provided fields from a file.
// It returns a map of field names to their values.
func (fp *FormulaParser) ParseFields(fields []ParseStrategy) (map[string]interface{}, error) {
	results := make(map[string]interface{})

	for fp.Scanner.Scan() {
		line := fp.Scanner.Text()

		for _, f := range fields {
			// Skip field if it has already been matched.
			if _, ok := results[f.getName()]; ok {
				continue
			}

			if f.MatchesLine(line) {
				fieldValue, err := f.ExtractFromLine(line)
				if err != nil {
					return nil, err
				}
				results[f.getName()] = fieldValue
				break
			}
		}
	}

	return results, nil
}

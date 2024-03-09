package delegate

import (
	"bufio"
	"fmt"
	"log"
	"regexp"
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
			if _, ok := results[f.GetName()]; ok {
				continue
			}

			log.Printf("Line: <generic:%s>: %s\n", f.GetName(), line)
			if f.matchesLine(line) {
				fieldValue, err := f.extractFromLine(line)
				if err != nil {
					return nil, err
				}
				results[f.GetName()] = fieldValue
				log.Println("Matched: ", results[f.GetName()])
				break
			}
		}
	}

	return results, nil
}

// ParseField parses the field with the given name and pattern from a file.
// It returns the value of the field.
func (fp *FormulaParser) ParseField(pattern, name string) (string, error) {
	for fp.Scanner.Scan() {
		line := fp.Scanner.Text()
		log.Printf("Line: %s: %s\n", name, line)

		regex := regexp.MustCompile(pattern)
		matches := regex.FindStringSubmatch(line)

		if len(matches) >= 2 {
			return matches[1], nil
		}
	}

	return "", fmt.Errorf("no %s found for Formula", name)
}

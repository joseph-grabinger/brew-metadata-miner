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
// It returns a map of fields to their values.
func (fp *FormulaParser) ParseFields(fields []Field) (map[Field]string, error) {
	results := make(map[Field]string)

	for fp.Scanner.Scan() {
		line := fp.Scanner.Text()

		for _, f := range fields {
			// Skip field if it has already been matched.
			if _, ok := results[f]; ok {
				continue
			}

			log.Printf("Line: <generic:%s>: %s\n", f.GetName(), line)
			strat := f.GetStrat()
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

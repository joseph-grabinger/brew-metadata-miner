package writer

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"main/miner/types"
)

// WriteFormulae writes the given formulae to the specified outputDir.
func WriteFormulae(outputDir string, formulae map[string]*types.Formula) error {
	currentTime := time.Now()
	formattedDate := currentTime.Format("2006-01-02")

	fileName := fmt.Sprintf("deps-brew-%s.tsv", formattedDate)

	path := filepath.Join(outputDir, fileName)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	// Close file on function exit and check its' returned error.
	defer func() error {
		if err := file.Close(); err != nil {
			return err
		}
		return nil
	}()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

	for _, formula := range formulae {
		// Write package line.
		line := formula.FormatPackageLine()
		_, err := writer.WriteString(line)
		if err != nil {
			return err
		}

		// Write dependency lines.
		for _, dep := range formula.Dependencies {
			f := formulae[dep.Name]
			if f == nil {
				panic(fmt.Sprintf("Dependency %s not found in formula %s", dep.Name, formula.Name))
			}

			line := f.FormatDependencyLine(dep)
			_, err := writer.WriteString(line)
			if err != nil {
				return err
			}
		}
	}
	return nil

}

package reader

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"

	"main/config"
	"main/miner/parser"
	"main/miner/setup"
	"main/miner/types"
)

type reader struct {
	formulae map[string]*types.Formula
	mu       sync.Mutex
}

func (p *reader) addFormula(formula *types.Formula) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.formulae[formula.Name] = formula
}

// ReadFormulae reads all formulae from the core repository in parallel
// using the given number of workers. It returns a map of formulae where
// the key is the name of the formula and the first encountered error.
func ReadFormulae(coreRepoPath string, readerConfig config.ReaderConfig) (map[string]*types.Formula, error) {
	// Create a new reader.
	r := &reader{
		formulae: make(map[string]*types.Formula),
	}

	// Match the fomula files.
	matches, err := filepath.Glob(coreRepoPath + "/Formula/**/*.rb")
	if err != nil {
		return nil, err
	}

	// Match alias formula files.
	aliasMatches, err := filepath.Glob(coreRepoPath + "/Aliases/*")
	if err != nil {
		return nil, err
	}

	matches = append(matches, aliasMatches...)

	// Create a context with cancellation capability.
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel() // Ensure cancellation happens when function returns.

	// Create a channel to enqueue the files to be processed.
	taskCh := make(chan string)

	// Create channel to communicate errors.
	errCh := make(chan error)

	var wg sync.WaitGroup

	// Create a worker pool to process the files concurrently.
	for i := 0; i < readerConfig.MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for path := range taskCh {
				select {
				case <-ctx.Done():
					return
				default:
					if err := r.processFile(path, readerConfig.FallbackLicense, readerConfig.DeriveRepo); err != nil {
						cancel()
						errCh <- err
						return
					}
				}
			}
		}()
	}

	// Enqueue the files to be processed.
	go func() {
		for _, path := range matches {
			taskCh <- path
		}
		// Close the channel after all files have been enqueued.
		close(taskCh)
	}()

	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Return the first encountered error.
	for err := range errCh {
		return nil, err
	}

	return r.formulae, nil
}

func (p *reader) processFile(path string, fallbackLicense string, deriveRepo bool) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Parse Formula from file.
	sourceFormula, err := extractFromFile(file)
	if err != nil {
		log.Printf("Error parsing file %s: %v\n", path, err)
		return err
	}

	formula := types.FromSourceFormula(sourceFormula, fallbackLicense, deriveRepo)
	p.addFormula(formula)

	log.Println("Successfully parsed formula:", formula)
	return nil
}

// extractFromFile extracts a formula from a file and returns it as a Formula struct.
func extractFromFile(file *os.File) (*types.SourceFormula, error) {
	scanner := bufio.NewScanner(file)
	formulaParser := &parser.FormulaParser{Scanner: scanner}

	base := filepath.Base(file.Name())
	name := strings.TrimSuffix(base, ".rb")

	formula := &types.SourceFormula{Name: name}

	fields := setup.BuildStrategies(*formulaParser)

	results, err := formulaParser.ParseFields(fields)
	if err != nil {
		log.Println("Error parsing fields:", err)
		return nil, err
	}

	// Set the fields of the formula.
	if results["homepage"] != nil {
		formula.Homepage = results["homepage"].(string)
	}
	if results["url"] != nil {
		formula.Stable = results["url"].(*types.Stable)
	}
	if results["mirror"] != nil {
		formula.Mirror = results["mirror"].(string)
	}
	if results["license"] != nil {
		formula.License = results["license"].(string)
	}
	if results["head"] != nil {
		formula.Head = results["head"].(*types.Head)
	}
	if results["dependency"] != nil {
		formula.Dependencies = results["dependency"].(*types.Dependencies)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	// Check if the stable url contains a Ruby string interpolation and resolve it.
	// This is done here rather then in the cleanURLSequence function because the
	// variable used for the interpolation can have a global scope in the formula file.
	// The cleanURLSequence function could only resolve interpolations with a scope within the stable do block.
	found, resolved, err := checkForInterpolation(formula.Stable.URL, file)
	if err != nil {
		return nil, err
	}
	if found {
		formula.Stable.URL = resolved
	}

	return formula, nil
}

// checkForInterpolation checks if the given url contains a Ruby string interpolation.
// If it does, the interpolation is resolved using the given file.
// The function returns a boolean indicating if an interpolation was found and the resolved string.
func checkForInterpolation(url string, file *os.File) (bool, string, error) {
	regex := regexp.MustCompile(setup.InterpolationPattern)
	matches := regex.FindStringSubmatch(url)
	if len(matches) < 2 {
		return false, "", nil
	}

	varName := matches[1]

	// Seek the beginning of the file.
	if _, err := file.Seek(0, io.SeekStart); err != nil {
		return true, "", err
	}

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()

		// Check if the line contains the variable assignment.
		regex := regexp.MustCompile(setup.VarAssignmentPattern(varName))
		varMatches := regex.FindStringSubmatch(line)
		if len(varMatches) >= 2 {
			// Replace the interpolation with the variable value.
			resolved := strings.Replace(url, fmt.Sprintf("#{%s}", varName), varMatches[1], 1)
			return true, resolved, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return true, "", err
	}

	return true, "", fmt.Errorf("could not resolve interpolation in URL: %s", url)
}

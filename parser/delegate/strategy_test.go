package delegate_test

import (
	"bufio"
	"log"
	"os"
	"testing"

	"main/parser/delegate"
	"main/parser/setup"
	"main/parser/types"

	"github.com/stretchr/testify/assert"
)

var MlmDependencyTests = []struct {
	inputFilePath string
	expected      []*types.Dependency
}{
	{
		inputFilePath: "../../tmp/homebrew-core/Formula/p/pinentry.rb",
		expected: []*types.Dependency{
			{Name: "pkg-config", DepType: "build", SystemRequirement: ""},
			{Name: "libassuan", DepType: "", SystemRequirement: ""},
			{Name: "libgpg-error", DepType: "", SystemRequirement: ""},
			{Name: "libsecret", DepType: "", SystemRequirement: "linux"}, // on_linux
		},
	},
}

func TestMultiLineMatcherDependencies(t *testing.T) {
	for _, test := range MlmDependencyTests {
		file, err := os.Open(test.inputFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		formulaParser := &delegate.FormulaParser{Scanner: bufio.NewScanner(file)}

		mlm := setup.BuildDependencyMatcher(*formulaParser)

		for formulaParser.Scanner.Scan() {
			line := formulaParser.Scanner.Text()

			if mlm.MatchesLine(line) {
				fieldValue, err := mlm.ExtractFromLine(line)
				if err != nil {
					log.Fatal(err)
				}

				dependencies := make([]*types.Dependency, 0)
				if fieldValue != nil {
					dependencies = fieldValue.([]*types.Dependency)
				}

				assert.Equal(t, test.expected, dependencies, "expected: %v, got: %v", test.expected, dependencies)
				log.Println("Matched: ", dependencies)
				break
			}
		}
	}
}

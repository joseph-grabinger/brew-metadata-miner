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
	{
		inputFilePath: "../../tmp/homebrew-core/Formula/lib/libdill.rb",
		expected: []*types.Dependency{
			{Name: "autoconf", DepType: "build", SystemRequirement: ""},
			{Name: "automake", DepType: "build", SystemRequirement: ""},
			{Name: "libtool", DepType: "build", SystemRequirement: ""},
			{Name: "llvm", DepType: "test", SystemRequirement: "arm"}, // on_arm
		},
	},
	{
		inputFilePath: "../../tmp/homebrew-core/Formula/g/grafana.rb",
		expected: []*types.Dependency{
			{Name: "go", DepType: "build", SystemRequirement: ""},
			{Name: "node", DepType: "build", SystemRequirement: ""},
			{Name: "yarn", DepType: "build", SystemRequirement: ""},
			{Name: "python", DepType: "build", SystemRequirement: "linux, macos: < catalina"},           // uses_from_macos
			{Name: "zlib", DepType: "", SystemRequirement: "linux"},                                     // uses_from_macos
			{Name: "python-setuptools", DepType: "build", SystemRequirement: "linux, macos: <= mojave"}, // on_system
			{Name: "fontconfig", DepType: "", SystemRequirement: "linux"},                               // on_linux
			{Name: "freetype", DepType: "", SystemRequirement: "linux"},                                 // on_linux
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

				log.Println("Matched: ", dependencies)

				assert.Equal(t, test.expected, dependencies, "expected: %v, got: %v", test.expected, dependencies)
				break
			}
		}
	}
}

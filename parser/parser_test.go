package parser

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"main/config"
	"main/parser/types"

	"github.com/stretchr/testify/assert"
)

var ParseFromFileTests = []struct {
	inputFilePath string
	expected      *types.SourceFormula
}{
	{
		inputFilePath: "../tmp/homebrew-core/Formula/i/i686-elf-gcc.rb",
		expected: &types.SourceFormula{
			Name:     "i686-elf-gcc",
			Homepage: "https://gcc.gnu.org",
			URL:      "https://ftp.gnu.org/gnu/gcc/gcc-13.2.0/gcc-13.2.0.tar.xz",
			Mirror:   "https://ftpmirror.gnu.org/gcc/gcc-13.2.0/gcc-13.2.0.tar.xz",
			License:  `"GPL-3.0-or-later" => { with: "GCC-exception-3.1" }`,
			Head:     nil,
			Dependencies: []*types.Dependency{
				{Name: "gmp", DepType: ""},
				{Name: "i686-elf-binutils", DepType: ""},
				{Name: "libmpc", DepType: ""},
				{Name: "mpfr", DepType: ""},
			},
		},
	},
	{
		inputFilePath: "../tmp/homebrew-core/Formula/p/pike.rb",
		expected: &types.SourceFormula{
			Name:     "pike",
			Homepage: "https://pike.lysator.liu.se/",
			URL:      "https://pike.lysator.liu.se/pub/pike/latest-stable/Pike-v8.0.1738.tar.gz",
			Mirror:   "http://deb.debian.org/debian/pool/main/p/pike8.0/pike8.0_8.0.1738.orig.tar.gz",
			License:  `any_of: ["GPL-2.0-only", "LGPL-2.1-only", "MPL-1.1"]`,
			Head:     nil,
			Dependencies: []*types.Dependency{
				{Name: "gettext", DepType: ""},
				{Name: "gmp", DepType: ""},
				{Name: "jpeg-turbo", DepType: ""},
				{Name: "libtiff", DepType: ""},
				{Name: "nettle", DepType: ""},
				{Name: "pcre", DepType: ""},
				{Name: "webp", DepType: ""},
				{Name: "bzip2", DepType: "", SystemRequirement: "linux"},     // on_linux
				{Name: "krb5", DepType: "", SystemRequirement: "linux"},      // on_linux
				{Name: "libxcrypt", DepType: "", SystemRequirement: "linux"}, // on_linux
				// {Name: "libxslt", DepType: ""},      // on_linux
				{Name: "sqlite", DepType: "", SystemRequirement: "linux"},       // on_linux
				{Name: "zlib", DepType: "", SystemRequirement: "linux"},         // on_linux
				{Name: "gnu-sed", DepType: "build", SystemRequirement: "macos"}, // on_macos
				{Name: "libnsl", DepType: "", SystemRequirement: "linux"},       // on_linux
			},
		},
	},
	{
		inputFilePath: "../tmp/homebrew-core/Formula/s/srecord.rb",
		expected: &types.SourceFormula{
			Name:     "srecord",
			Homepage: "https://srecord.sourceforge.net/",
			URL:      "https://downloads.sourceforge.net/project/srecord/srecord/1.64/srecord-1.64.tar.gz",
			Mirror:   "",
			License:  `all_of: ["GPL-3.0-or-later", "LGPL-3.0-or-later"]`,
			Head:     nil,
			Dependencies: []*types.Dependency{
				{Name: "boost", DepType: "build"},
				{Name: "libtool", DepType: "build"},
				{Name: "libgcrypt", DepType: ""},
				{Name: "ghostscript", DepType: "build", SystemRequirement: "macos: >= sonoma, linux"}, // on_sonoma :or_newer && on_linux
				{Name: "groff", DepType: "build", SystemRequirement: "macos: >= ventura, linux"},      // on_ventura :or_newer && on_linux
			},
		},
	},
}

func TestParseFromFile(t *testing.T) {
	for _, test := range ParseFromFileTests {
		file, err := os.Open(test.inputFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		formula, err := parseFromFile(file)
		if err != nil {
			log.Fatal(err)
		}

		assert.Equal(t, test.expected, formula, "expected: %v, got: %v", test.expected, formula)
	}
}

func TestParse_Reliabity(t *testing.T) {
	config := &config.Config{}
	config.CoreRepo.Dir = "../tmp/homebrew-core"

	parser := NewParser(config)

	if err := parser.Parse(); err != nil {
		log.Fatal(err)
	}

	jsonLst := getJSONFromAPI()

	// Assert total number of formulas.
	assert.LessOrEqual(t, len(jsonLst), len(parser.formulas), "expected: %d formulas from API, got: %d from core repo", len(jsonLst), len(parser.formulas))

	for _, apiFormula := range jsonLst {
		name, ok := apiFormula["name"].(string)
		if !ok {
			t.Errorf("no name found in formula")
			continue
		}

		// Check if formula exists in parser.
		formula, ok := parser.formulas[name]
		if !ok {
			t.Errorf("formula %s not found", name)
			continue
		}

		// Assert licenses are equal.
		if apiFormula["license"] == nil {
			assert.EqualValues(t, "pseudo", formula.License, "expected: pseudo license of %s, got: %s", name, formula.License)
		} else {
			assert.Equal(t, apiFormula["license"], formula.License, "expected: %s as license of %s, got: %s", apiFormula["license"], name, formula.License)
		}

		if headUrl, ok := getNestedMapValue(apiFormula, "urls", "head", "url"); ok {
			assert.Equal(t, headUrl, formula.RepoURL, "expected: %s as head url of %s, got: %s", headUrl, name, formula.RepoURL)
		}

		// Assert dependencies are equal.
		if deps, ok := apiFormula["dependencies"].([]string); ok {
			normalFormulaDeps := getDependeciesByType(formula, "")
			assert.ElementsMatch(t, deps, normalFormulaDeps, "expected: %s as head url of %s, got: %s", deps, name, normalFormulaDeps)
		}

		// Assert build dependencies are equal.
		if buildDeps, ok := apiFormula["build_dependencies"].([]string); ok {
			buildFormulaDeps := getDependeciesByType(formula, "build")
			assert.ElementsMatch(t, buildDeps, buildFormulaDeps, "expected: %s as head url of %s, got: %s", buildDeps, name, buildFormulaDeps)
		}

		// Assert test dependencies are equal.
		if testDeps, ok := apiFormula["test_dependencies"].([]string); ok {
			testFormulaDeps := getDependeciesByType(formula, "test")
			assert.ElementsMatch(t, testDeps, testFormulaDeps, "expected: %s as head url of %s, got: %s", testDeps, name, testFormulaDeps)
		}

		// Assert system dependencies are equal.
		if systemDeps, ok := apiFormula["uses_from_macos"].([]string); ok {
			systemFormulaDeps := getDependeciesByType(formula, "system")
			assert.ElementsMatch(t, systemDeps, systemFormulaDeps, "expected: %s as head url of %s, got: %s", systemDeps, name, systemFormulaDeps)
		}
	}
}

// getJSONFromAPI returns a list of all formulas from the homebrew API.
func getJSONFromAPI() []map[string]interface{} {
	resp, err := http.Get("https://formulae.brew.sh/api/formula.json")
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var jsonLst []map[string]interface{}
	if err := json.Unmarshal(body, &jsonLst); err != nil {
		log.Fatal(err)
	}

	return jsonLst
}

// getDependeciesByType returns a list of dependencies of the given type.
func getDependeciesByType(formula *types.Formula, depType string) []string {
	deps := make([]string, 0)
	for _, dep := range formula.Dependencies {
		if dep.DepType == depType {
			deps = append(deps, dep.Name)
		}
	}
	return deps
}

func getNestedMapValue(m map[string]interface{}, keys ...string) (value string, ok bool) {
	var temp interface{} = m

	for _, key := range keys {
		tempMap, valid := temp.(map[string]interface{})
		if !valid {
			return "", false
		}

		temp, ok = tempMap[key]
		if !ok {
			return "", false
		}
	}

	value, ok = temp.(string)
	return
}

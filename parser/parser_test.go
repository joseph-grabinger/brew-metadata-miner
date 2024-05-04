package parser

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"slices"
	"strings"
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
			Stable: &types.Stable{
				URL: "https://ftp.gnu.org/gnu/gcc/gcc-13.2.0/gcc-13.2.0.tar.xz",
			},
			Mirror:  "https://ftpmirror.gnu.org/gcc/gcc-13.2.0/gcc-13.2.0.tar.xz",
			License: `"GPL-3.0-or-later" => { with: "GCC-exception-3.1" }`,
			Head:    nil,
			Dependencies: &types.Dependencies{
				Lst: []*types.Dependency{
					{Name: "gmp", DepType: []string{}},
					{Name: "i686-elf-binutils", DepType: []string{}},
					{Name: "libmpc", DepType: []string{}},
					{Name: "mpfr", DepType: []string{}},
				},
			},
		},
	},
	{
		inputFilePath: "../tmp/homebrew-core/Formula/p/pike.rb",
		expected: &types.SourceFormula{
			Name:     "pike",
			Homepage: "https://pike.lysator.liu.se/",
			Stable: &types.Stable{
				URL: "https://pike.lysator.liu.se/pub/pike/latest-stable/Pike-v8.0.1738.tar.gz",
			},
			Mirror:  "http://deb.debian.org/debian/pool/main/p/pike8.0/pike8.0_8.0.1738.orig.tar.gz",
			License: `any_of: ["GPL-2.0-only", "LGPL-2.1-only", "MPL-1.1"]`,
			Head:    nil,
			Dependencies: &types.Dependencies{
				Lst: []*types.Dependency{
					{Name: "gettext", DepType: []string{}},
					{Name: "gmp", DepType: []string{}},
					{Name: "jpeg-turbo", DepType: []string{}},
					{Name: "libtiff", DepType: []string{}},
					{Name: "nettle", DepType: []string{}},
					{Name: "pcre", DepType: []string{}},
					{Name: "webp", DepType: []string{}},
					{Name: "bzip2", DepType: []string{}, Restriction: "linux"},     // on_linux
					{Name: "krb5", DepType: []string{}, Restriction: "linux"},      // on_linux
					{Name: "libxcrypt", DepType: []string{}, Restriction: "linux"}, // on_linux
					// {Name: "libxslt", DepType: ""},      // on_linux
					{Name: "sqlite", DepType: []string{}, Restriction: "linux"},         // on_linux
					{Name: "zlib", DepType: []string{}, Restriction: "linux"},           // on_linux
					{Name: "gnu-sed", DepType: []string{"build"}, Restriction: "macos"}, // on_macos
					{Name: "libnsl", DepType: []string{}, Restriction: "linux"},         // on_linux
				},
			},
		},
	},
	{
		inputFilePath: "../tmp/homebrew-core/Formula/s/srecord.rb",
		expected: &types.SourceFormula{
			Name:     "srecord",
			Homepage: "https://srecord.sourceforge.net/",
			Stable: &types.Stable{
				URL: "https://downloads.sourceforge.net/project/srecord/srecord/1.64/srecord-1.64.tar.gz",
			},
			Mirror:  "",
			License: `all_of: ["GPL-3.0-or-later", "LGPL-3.0-or-later"]`,
			Head:    nil,
			Dependencies: &types.Dependencies{
				Lst: []*types.Dependency{
					{Name: "boost", DepType: []string{"build"}},
					{Name: "libtool", DepType: []string{"build"}},
					{Name: "libgcrypt", DepType: []string{}},
					{Name: "ghostscript", DepType: []string{"build"}, Restriction: "macos: >= sonoma or linux"}, // on_sonoma :or_newer && on_linux
					{Name: "groff", DepType: []string{"build"}, Restriction: "macos: >= ventura or linux"},      // on_ventura :or_newer && on_linux
				},
			},
		},
	},
	{
		inputFilePath: "../tmp/homebrew-core/Formula/g/geckodriver.rb",
		expected: &types.SourceFormula{
			Name:     "geckodriver",
			Homepage: "https://github.com/mozilla/geckodriver",
			Stable: &types.Stable{
				URL: "https://hg.mozilla.org/mozilla-central/archive/bc25087baba17c78246db06bcab71c299fd8f46f.zip/testing/geckodriver/",
				Dependencies: &types.Dependencies{
					Lst: []*types.Dependency{},
				},
			},
			Mirror:  "",
			License: `"MPL-2.0"`,
			Head: &types.Head{
				URL: "https://hg.mozilla.org/mozilla-central/",
			},
			Dependencies: &types.Dependencies{
				Lst: []*types.Dependency{
					{Name: "rust", DepType: []string{"build"}},
					{Name: "netcat", DepType: []string{"test"}, Restriction: "linux"},
					{Name: "unzip", DepType: []string{}, Restriction: "linux"},
				},
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

		assert.ElementsMatch(t, test.expected.Dependencies.Lst, formula.Dependencies.Lst, "expected: %v, got: %v", test.expected.Dependencies.Lst, formula.Dependencies.Lst)
		assert.Equal(t, test.expected.Dependencies.SystemRequirements, formula.Dependencies.SystemRequirements, "expected: %v, got: %v", test.expected.Dependencies.SystemRequirements, formula.Dependencies.SystemRequirements)
		test.expected.Dependencies, formula.Dependencies = nil, nil
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

		// Assert archive urls are equal.
		if archiveUrl, ok := getNestedMapValue(apiFormula, "urls", "stable", "url"); ok {
			if tag, ok := getNestedMapValue(apiFormula, "urls", "stable", "tag"); ok {
				archiveUrl = strings.TrimSuffix(archiveUrl, ".git")
				archiveUrl = strings.TrimRight(archiveUrl, "/")
				archiveUrl += "/tree/" + tag
			}
			assert.Equal(t, archiveUrl, formula.ArchiveURL, "expected: %s as archive url of %s, got: %s", archiveUrl, name, formula.ArchiveURL)
		}

		// Assert heads are equal.
		if headUrl, ok := getNestedMapValue(apiFormula, "urls", "head", "url"); ok {
			assert.Equal(t, headUrl, formula.RepoURL, "expected: %s as head url of %s, got: %s", headUrl, name, formula.RepoURL)
		}

		// Assert dependencies match.
		apiDeps := toStringSlice(apiFormula["dependencies"].([]interface{}))
		deps := getDependeciesByType(formula, "")
		assert.ElementsMatch(t, apiDeps, deps, "expected: %v as dependencies of %s, got: %v", apiDeps, name, deps)

		// Assert build dependencies match.
		apiBuildDpes := toStringSlice(apiFormula["build_dependencies"].([]interface{}))
		buildDeps := getDependeciesByType(formula, "build")
		assert.ElementsMatch(t, apiBuildDpes, buildDeps, "expected: %v as build dependencies of %s, got: %v", apiBuildDpes, name, buildDeps)

		// Assert test dependencies match.
		apiTestDeps := toStringSlice(apiFormula["test_dependencies"].([]interface{}))
		testDeps := getDependeciesByType(formula, "test")
		assert.ElementsMatch(t, apiTestDeps, testDeps, "expected: %v as test dependencies of %s, got: %v", apiTestDeps, name, testDeps)

		// Assert uses_from_macos dependencies match.
		apiUsesFromMacosDeps := toStringSlice(apiFormula["uses_from_macos"].([]interface{}))
		usesFromMacosDeps, complete := getCommonDependencies(apiUsesFromMacosDeps, formula.Dependencies)
		assert.True(t, complete, "expected: %v as uses_from_macos dependencies of %s, got: %v", apiUsesFromMacosDeps, name, usesFromMacosDeps)
		for _, dep := range usesFromMacosDeps {
			// Check if restriction is linux or empty.
			// Empty restircion is used if a formula decares a uses_from_macos dependency but also as a dependecy without restriction. E.g. python@3.8.rb
			assert.True(t, (strings.Contains(dep.Restriction, "linux") || dep.Restriction == ""), "expected: linux restriction for %s as uses_from_macos dependency of %s, got: %s", dep, name, dep.Restriction)
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

// getDependeciesByType returns a list of dependencies of the given type
// which either don't have a restriction or have the API's default restriction.
func getDependeciesByType(formula *types.Formula, depType string) []string {
	deps := make([]string, 0)
	for _, dep := range formula.Dependencies {
		// No dependecy type.
		if depType == "" && len(dep.DepType) == 0 && isDefaultRestriction(dep.Restriction) {
			deps = append(deps, dep.Name)
			continue
		}
		// Specified dependecy type.
		if slices.Contains(dep.DepType, depType) && isDefaultRestriction(dep.Restriction) {
			deps = append(deps, dep.Name)
		}
	}
	return deps
}

// isDefaultRestriction returns true if the given restriction is the default restriction.
// This is used to check if the restriction is the same as the one in the API,
// since the API interpretes certain restirctions as default.
func isDefaultRestriction(restr string) bool {
	// clang version >= 1500
	clangExp := regexp.MustCompile(`clang version >= [1-9][0-9]{3,}[^()]*`)

	return restr == "" ||
		restr == "macos" ||
		restr == "arm" ||
		(strings.Contains(restr, "macos: >=") && !strings.Contains(restr, "(")) ||
		(clangExp.MatchString(restr) && !strings.Contains(restr, "(")) ||
		strings.Contains(restr, "(macos and arm)") ||
		(strings.Contains(restr, "arm") && !strings.Contains(restr, "(") && strings.Count(restr, " and ") <= 1)
}

// getCommonDependencies returns a list of common dependencies between the given list of dependency names
// and the given dependecies.
// It also returns a boolean indicating if all the given dependency names are found in the dependencies.
func getCommonDependencies(s []string, deps []*types.Dependency) (common []*types.Dependency, complete bool) {
	common = make([]*types.Dependency, 0)
	for _, dep := range deps {
		if !slices.Contains(s, dep.Name) {
			continue
		}
		// Check if common slice already contains dep.
		// This is done to cope with the case where a dependency is listed multiple times but have different types. E.g. llnode.rb
		i := slices.IndexFunc(common, func(d *types.Dependency) bool { return d.Name == dep.Name })
		if i == -1 {
			common = append(common, dep)
			continue
		}
		// Use the dependency with linux or empty restrictions.
		if common[i].Restriction == "" || common[i].Restriction == "linux" {
			continue
		}
		if dep.Restriction == "" || dep.Restriction == "linux" {
			common[i] = dep
		}
	}
	return common, len(common) == len(s)
}

// getNestedMapValue returns the value of the given nested keys from the given map.
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

// toStringSlice converts a slice of interfaces to a slice of strings.
// If a map is found, the first key is used.
func toStringSlice(s []interface{}) []string {
	res := make([]string, 0)
	for _, v := range s {
		if str, ok := v.(string); ok {
			res = append(res, str)
			continue
		}
		if m, ok := v.(map[string]interface{}); ok {
			for key := range m {
				res = append(res, key)
				continue
			}
		}
	}
	return res
}

package parser

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"testing"

	"main/config"
	"main/parser/types"

	"github.com/stretchr/testify/assert"
)

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
			t.Log("name not found in formula")
			continue
		}

		// Check if formula exists in parser.
		formula, ok := parser.formulas[name]
		if !ok {
			t.Errorf("formula %s not found", apiFormula["name"].(string))
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

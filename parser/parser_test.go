package parser

import (
	"encoding/json"
	"io"
	"log"
	"main/config"
	"net/http"
	"testing"

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

	missingCount := 0

	for _, apiFormula := range jsonLst {
		name, ok := apiFormula["name"].(string)
		if !ok {
			t.Log("name not found in formula")
			continue
		}

		// Check if formula exists in parser.
		formula, ok := parser.formulas[name]
		if !ok {
			// t.Errorf("formula %s not found", apiFormula["name"].(string))
			missingCount++
			continue
		}

		// Assert licenses are equal.
		assert.Equal(t, apiFormula["license"], formula.license, "expected: %s as license of %s, got: %s", apiFormula["license"], name, formula.license)

		if headUrl, ok := getNestedMapValue(apiFormula, "urls", "head", "url"); ok {
			assert.Equal(t, headUrl, formula.repoURL, "expected: %s as head url of %s, got: %s", headUrl, name, formula.repoURL)
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

	t.Logf("missing formulas: %d", missingCount)
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
func getDependeciesByType(formula *formula, depType string) []string {
	deps := make([]string, 0)
	for _, dep := range formula.dependencies {
		if dep.depType == depType {
			deps = append(deps, dep.name)
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

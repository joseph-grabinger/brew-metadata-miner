package setup

import (
	"fmt"
	"regexp"
	"strings"

	"main/parser/types"
)

type requirementsHelper struct {
	requirements string
	depth        int
	// TODO check if depth is needed
}

// cleanDependencySequence returns a cleaned slice of dependencies from a given sequence.
// The slice is returned as an interface{} to be casted to []*types.Dependency.
func cleanDependencySequence(sequence []string) interface{} {
	// for i := range sequence {
	// 	log.Println(sequence[i])
	// }
	// log.Println("Cleaning sequence: ", sequence)

	reqHelper := &requirementsHelper{requirements: "", depth: 0}
	res := make([]*types.Dependency, 0)
	for i := range sequence {
		// Check for system dependency.
		regex := regexp.MustCompile(macOSSystemDependencyPattern)
		nameMatches := regex.FindStringSubmatch(sequence[i])
		if len(nameMatches) >= 2 {
			// TODO check if doable in one step
			depType := getDepType(sequence[i])

			req := "linux"
			if since := getOSRestriction(sequence[i]); since != "" {
				req += ", macos: < " + since
			}
			res = append(res, &types.Dependency{
				Name:              nameMatches[1],
				DepType:           depType,
				SystemRequirement: req,
			})
			continue
		}

		// Check for end.
		regex = regexp.MustCompile(endPattern)
		if regex.MatchString(sequence[i]) {
			reqHelper.depth--
			reqHelper.requirements = ""
			continue
		}

		// Check for the dependency name.
		regex = regexp.MustCompile(dependencyKeywordPattern)
		nameMatches = regex.FindStringSubmatch(sequence[i])
		if len(nameMatches) >= 2 {
			depType := getDepType(sequence[i])
			res = append(res, &types.Dependency{
				Name:              nameMatches[1],
				DepType:           depType,
				SystemRequirement: reqHelper.requirements,
			})
		}

		// Check for requirements.
		checkRequirements(sequence[i], reqHelper)
	}
	return res
}

// getDepType returns the dependency type from the given line.
// If no type is found, an empty string is returned.
func getDepType(line string) string {
	regex := regexp.MustCompile(dependencyTypePattern)
	typeMatches := regex.FindStringSubmatch(line)
	if len(typeMatches) >= 2 {
		return typeMatches[1]
	}
	return ""
}

// getOSRestriction returns the OS restriction from the given line.
// If no restriction is found, an empty string is returned.
func getOSRestriction(line string) string {
	regex := regexp.MustCompile(osRestrictionPattern)
	matches := regex.FindStringSubmatch(line)
	if len(matches) >= 2 {
		return matches[1]
	}
	return ""
}

// checkRequirements checks the given line for system requirements.
// If a requirement is found, it is added to the requirements string.
// System requirements include: on_system, on_linux, on_arm, and on_intel.
func checkRequirements(line string, reqHelper *requirementsHelper) {
	// Check for on_system.
	regex := regexp.MustCompile(onSystemPattern)
	if regex.MatchString(line) {
		reqHelper.depth++
		regex = regexp.MustCompile(onSystemExtractPattern)
		matches := regex.FindStringSubmatch(line)
		if len(matches) != 2 {
			panic("Invalid on_system pattern")
		}
		v, err := formatVersion(matches[1])
		if err != nil {
			panic(err)
		}
		reqHelper.requirements += "linux, macos: " + v
		return
	}

	// Check for on_linux.
	regex = regexp.MustCompile(onLinuxPattern)
	if regex.MatchString(line) {
		reqHelper.depth++
		reqHelper.requirements += "linux"
		return
	}

	// Check for on_macos.
	regex = regexp.MustCompile(onMacosPattern)
	if regex.MatchString(line) {
		reqHelper.depth++
		reqHelper.requirements += "macos"
		return
	}

	// Check for on_arm.
	regex = regexp.MustCompile(onArmPattern)
	if regex.MatchString(line) {
		reqHelper.depth++
		reqHelper.requirements += "arm"
		return
	}

	// Check for on_intel.
	regex = regexp.MustCompile(onIntelPattern)
	if regex.MatchString(line) {
		reqHelper.depth++
		reqHelper.requirements += "intel"
		return
	}

	// Check for on_macos versions e.g. "on_mojave :or_newer".
	regex = regexp.MustCompile(onMacOSVersionPattern)
	matches := regex.FindStringSubmatch(line)
	if len(matches) >= 3 {
		req := "macos: "
		v, res := matches[2], matches[3]
		if res != "" {
			fv, err := formatVersion(v + "_" + res)
			if err != nil {
				panic(err)
			}
			req += fv
		} else {
			req += v
		}
		reqHelper.depth++
		reqHelper.requirements += req
	}
}

// formatVersion returns a formatted string from the given version.
// If the string format is invalid, an error is returned.
// Example:
// "sierra_or_older" => "<= sierra" or
// "high_sierra_or_newer" => ">= high_sierra"
func formatVersion(version string) (string, error) {
	parts := strings.Split(version, "_")

	var v, res string
	switch len(parts) {
	case 3:
		v, res = parts[0], parts[2]
	case 4:
		v, res = parts[0]+"_"+parts[1], parts[3]
	default:
		return "", fmt.Errorf("invalid input string format")
	}

	if res != "older" && res != "newer" {
		return "", fmt.Errorf("invalid input string format")
	}

	r := strings.NewReplacer("older", "<=", "newer", ">=")
	return r.Replace(res) + " " + v, nil
}

// isBeginDependencySequence returns true if the given line
// is the beginning of a dependency sequence.
func isBeginDependencySequence(line string) bool {
	regex := regexp.MustCompile(beginDependencyPattern)
	return regex.MatchString(line)
}

// isEndDependencySequence returns true if the given line
// is the end of a dependency sequence.
func isEndDependencySequence(line string) bool {
	regex := regexp.MustCompile(endDependencyPattern)
	return !regex.MatchString(line)
}

package setup

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"main/parser/types"
	"main/stack"
)

// dependencySet represents a set of dependencies, stored in a map by their name.
type dependecySet map[string]*types.Dependency

// add adds a dependency to the set.
// If the dependency already exists, the system requirements are merged.
func (s dependecySet) add(dep *types.Dependency) {
	id := dep.Id()
	if d, ok := s[id]; ok {
		// Merge the dependency`s system requirements.
		log.Println("Merging system requirements for: ", dep.Name)
		if dep.SystemRequirement == "" || d.SystemRequirement == "" {
			return
		}

		if strings.Contains(dep.SystemRequirement, ", ") {
			dep.SystemRequirement = fmt.Sprintf("(%s)", dep.SystemRequirement)
		}
		if strings.Contains(d.SystemRequirement, ", ") {
			d.SystemRequirement = fmt.Sprintf("(%s)", d.SystemRequirement)
		}

		d.SystemRequirement = strings.Join([]string{d.SystemRequirement, dep.SystemRequirement}, ", ")
		return
	}
	s[id] = dep
}

// toSlice returns the set as a slice of dependencies.
func (s dependecySet) toSlice() []*types.Dependency {
	res := make([]*types.Dependency, 0)
	for _, v := range s {
		res = append(res, v)
	}
	return res
}

// skipSequence represents a sequence that should be skipped.
type skipSequence struct {
	begin string
	end   string
}

// skips is a slice of skip sequences.
type skips []skipSequence

func (s skips) shouldSkip(line string) (bool, *skipSequence) {
	for _, skip := range s {
		regex := regexp.MustCompile(skip.begin)
		if regex.MatchString(line) {
			return true, &skip
		}
	}
	return false, nil
}

// cleanDepSequence returns a cleaned slice of dependencies from a given sequence.
// The provided skips is used to skip certain lines.
// The numIgnoreEmpty is number of empty stack pops to ignore.
func cleanDepSequence(sequence []string, skips skips, numIgnoreEmpty int) []*types.Dependency {
	reqStack := stack.New[string]()
	set := make(dependecySet, 0)
	var skip *skipSequence
	for i := range sequence {
		// Check whether to skip the current line.
		if skip != nil {
			// Check for end sequence.
			regex := regexp.MustCompile(skip.end)
			if regex.MatchString(sequence[i]) {
				skip = nil
			}
			continue
		}

		// Check for skip sequence.
		if shouldSkip, s := skips.shouldSkip(sequence[i]); shouldSkip {
			skip = s
			continue
		}

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
			set.add(&types.Dependency{
				Name:              nameMatches[1],
				DepType:           depType,
				SystemRequirement: req,
			})

			continue
		}

		// Check for end.
		regex = regexp.MustCompile(endPatternGeneric)
		if regex.MatchString(sequence[i]) {
			_, err := reqStack.Pop()
			if err != nil {
				if numIgnoreEmpty != 0 {
					numIgnoreEmpty--
				} else {
					panic(err)
				}
			}
			continue
		}

		// Check for the dependency name.
		regex = regexp.MustCompile(dependencyKeywordPattern)
		nameMatches = regex.FindStringSubmatch(sequence[i])
		if len(nameMatches) >= 2 {
			depType := getDepType(sequence[i])
			set.add(&types.Dependency{
				Name:              nameMatches[1],
				DepType:           depType,
				SystemRequirement: strings.Join(reqStack.Values(), ", "),
			})
		}

		// Check for requirements.
		checkRequirements(sequence[i], reqStack)
	}
	return set.toSlice()
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
func checkRequirements(line string, reqStack *stack.Stack[string]) {
	// Check for on_system.
	regex := regexp.MustCompile(onSystemPattern)
	if regex.MatchString(line) {
		regex = regexp.MustCompile(onSystemExtractPattern)
		matches := regex.FindStringSubmatch(line)
		if len(matches) != 2 {
			panic("Invalid on_system pattern")
		}
		v, err := formatVersion(matches[1])
		if err != nil {
			panic(err)
		}
		reqStack.Push("linux, macos: " + v)
		return
	}

	// Check for on_linux.
	regex = regexp.MustCompile(onLinuxPattern)
	if regex.MatchString(line) {
		reqStack.Push("linux")
		return
	}

	// Check for on_macos.
	regex = regexp.MustCompile(onMacosPattern)
	if regex.MatchString(line) {
		reqStack.Push("macos")
		return
	}

	// Check for on_arm.
	regex = regexp.MustCompile(onArmPattern)
	if regex.MatchString(line) {
		reqStack.Push("arm")
		return
	}

	// Check for on_intel.
	regex = regexp.MustCompile(onIntelPattern)
	if regex.MatchString(line) {
		reqStack.Push("intel")
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
		reqStack.Push(req)
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

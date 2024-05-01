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
		// Merge the dependency`s system restrictions.
		log.Println("Merging system restirctions for: ", dep.Name)
		if dep.Restriction == "" || d.Restriction == "" {
			return
		}

		if strings.Contains(dep.Restriction, ", ") {
			dep.Restriction = fmt.Sprintf("(%s)", dep.Restriction)
		}
		if strings.Contains(d.Restriction, ", ") {
			d.Restriction = fmt.Sprintf("(%s)", d.Restriction)
		}

		d.Restriction = strings.Join([]string{d.Restriction, dep.Restriction}, ", ")
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
func cleanDepSequence(sequence []string, skips skips, numIgnoreEmpty int) *types.Dependencies {
	depResStack := stack.New[string]()     // Holds the dependecy restirctions.
	formulaReqStack := stack.New[string]() // Holds the formula requirements.
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

			res := "linux"
			if since := getOSRestriction(sequence[i]); since != "" {
				res += ", macos: < " + since
			}
			set.add(&types.Dependency{
				Name:        nameMatches[1],
				DepType:     depType,
				Restriction: res,
			})

			continue
		}

		// Check for end.
		regex = regexp.MustCompile(endPatternGeneric)
		if regex.MatchString(sequence[i]) {
			_, err := depResStack.Pop()
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
				Name:        nameMatches[1],
				DepType:     depType,
				Restriction: strings.Join(depResStack.Values(), ", "),
			})
		}

		// Check for restrictions.
		checkDependencyRestrictions(sequence[i], depResStack)

		// Check for formula requirements.
		checkFormulaRequirements(sequence[i], formulaReqStack)

		// Check for fails_with.
		regex = regexp.MustCompile(failsWithPattern)
		if regex.MatchString(sequence[i]) {
			// Add a new empty restriction which will be poped as soon as
			// the end statement of the fails_with block is reached.
			depResStack.Push("")
		}
	}
	return &types.Dependencies{
		Lst:                set.toSlice(),
		SystemRequirements: strings.Join(formulaReqStack.Values(), ", "),
	}
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

// checkFormulaRequirements checks the given line for formula requirements.
// If a requirement is found, it is added to the stack.
// Formula system requirements include: macos, maximum_macos, xcode, and arch.
func checkFormulaRequirements(line string, reqStack *stack.Stack[string]) {
	regex := regexp.MustCompile(formulaRequirementPattern)
	matches := regex.FindStringSubmatch(line)
	count := len(matches)
	if count < 2 {
		return
	}

	req := matches[1]

	// Leading colon indicates an OS requirement without version e.g. ":linux" or ":macos".
	if s, found := strings.CutPrefix(req, ":"); found {
		reqStack.Push(s)
		return
	}

	req = strings.TrimSuffix(req, ":")

	if count == 3 {
		switch req {
		case "macos":
			req += " >= " + strings.TrimPrefix(matches[2], ":") + " (or linux)"
		case "maximum_macos":
			req += " <= " + formatRequirements(matches[2]) + " (or linux)"
		case "xcode":
			if strings.Contains(matches[2], `"`) {
				// Indicates a min version.
				req += " >= " + formatRequirements(matches[2]) + " (on macos)"
			} else {
				req += " " + formatRequirements(matches[2]) + " (on macos)"
			}
		case "arch":
			req = formatRequirements(matches[2])
		default:
			log.Printf("Incomplete formula requirement: %s, %s\n", req, matches[2])
			return
		}
	}

	reqStack.Push(req)
}

// formatRequirements returns a formatted string from the given requirements.
// Example:
// "[:monterey, :build]" => "monterey build"
// ":catalina" => "catalina"
// "["15.0", :build]" => "15.0 build"
func formatRequirements(req string) string {
	if !strings.Contains(req, ",") {
		return strings.TrimPrefix(req, ":")
	}

	r := strings.NewReplacer("[", "", "]", "", ":", "", ",", "", `"`, "")
	return r.Replace(req)
}

// checkDependencyRestrictions checks the given line for dependecy restrictions.
// If a restriction is found, it is added to the stack.
// Dependecy restrictions include: on_system, on_linux, on_arm, and on_intel.
func checkDependencyRestrictions(line string, resStack *stack.Stack[string]) {
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
		resStack.Push("linux, macos: " + v)
		return
	}

	// Check for on_linux.
	regex = regexp.MustCompile(onLinuxPattern)
	if regex.MatchString(line) {
		resStack.Push("linux")
		return
	}

	// Check for on_macos.
	regex = regexp.MustCompile(onMacosPattern)
	if regex.MatchString(line) {
		resStack.Push("macos")
		return
	}

	// Check for on_arm.
	regex = regexp.MustCompile(onArmPattern)
	if regex.MatchString(line) {
		resStack.Push("arm")
		return
	}

	// Check for on_intel.
	regex = regexp.MustCompile(onIntelPattern)
	if regex.MatchString(line) {
		resStack.Push("intel")
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
		resStack.Push(req)
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

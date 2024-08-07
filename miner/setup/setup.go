package setup

import (
	"main/miner/parser"
	"main/miner/types"
)

// BuildStrategies returns a list of parse strategies.
// The list contains a strategy for each field, parsed from the formula file.
func BuildStrategies(fp parser.FormulaParser) []parser.ParseStrategy {
	return []parser.ParseStrategy{
		BuildHomepageMatcher(fp),
		BuildURLMatcher(fp),
		BuildStableURLMatcher(fp),
		BuildMirrorMatcher(fp),
		BuildLicenseMatcher(fp),
		BuildHeadMatcher(fp),
		BuildDependencyMatcher(fp),
	}
}

// BuildHomepageMatcher returns a SingleLineMatcher for the homepage field.
func BuildHomepageMatcher(fp parser.FormulaParser) *parser.SingleLineMatcher[string] {
	return parser.NewSLM[string]("homepage", isDefaultHomepagePattern, fp)
}

// BuildURLMatcher returns a MultiLineMatcher for the URL field.
func BuildURLMatcher(fp parser.FormulaParser) *parser.MultiLineMatcher[*types.Stable] {
	return parser.NewMLM[*types.Stable]("url", isDefaultURLPattern, fp, isBeginURLSequence, isEndURLSequence, cleanURLSequence)
}

// BuildStableURLMatcher returns a MultiLineMatcher for the stable URL field.
func BuildStableURLMatcher(fp parser.FormulaParser) *parser.MultiLineMatcher[*types.Stable] {
	return parser.NewMLM[*types.Stable]("url", isDefaultStableURLPattern, fp, isBeginStableURLSequence, isEndStableURLSequence, cleanURLSequence)
}

// BuildMirrorMatcher returns a SingleLineMatcher for the mirror field.
func BuildMirrorMatcher(fp parser.FormulaParser) *parser.SingleLineMatcher[string] {
	return parser.NewSLM[string]("mirror", isDefaultMirrorPattern, fp)
}

// BuildLicenseMatcher returns a MultiLineMatcher for the license field.
func BuildLicenseMatcher(fp parser.FormulaParser) *parser.MultiLineMatcher[string] {
	return parser.NewMLM[string]("license", isDefaultLicensePattern, fp, isBeginLicenseSequence, isEndLicenseSequence, cleanLicenseSequence)
}

// BuildHeadMatcher returns a MultiLineMatcher for the head field.
func BuildHeadMatcher(fp parser.FormulaParser) *parser.MultiLineMatcher[*types.Head] {
	return parser.NewMLM[*types.Head]("head", isDefaultHeadPattern, fp, isBeginHeadSequence, isEndHeadSequence, cleanHeadSequence)
}

// BuildDependencyMatcher returns a MultiLineMatcher for the dependency fields.
func BuildDependencyMatcher(fp parser.FormulaParser) *parser.MultiLineMatcher[*types.Dependencies] {
	return parser.NewMLM[*types.Dependencies]("dependency", isDefaultDependencyPattern, fp, isBeginDependencySequence, isEndDependencySequence, cleanDependencySequence)
}

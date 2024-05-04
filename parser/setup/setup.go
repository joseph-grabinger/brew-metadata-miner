package setup

import (
	"main/parser/delegate"
	"main/parser/types"
)

// BuildStrategies returns a list of parse strategies.
// The list contains a strategy for each field, parsed from the formula file.
func BuildStrategies(fp delegate.FormulaParser) []delegate.ParseStrategy {
	return []delegate.ParseStrategy{
		BuildURLMatcher(fp),
		BuildMirrorMatcher(fp),
		BuildLicenseMatcher(fp),
		BuildHeadMatcher(fp),
		BuildDependencyMatcher(fp),
	}
}

// BuildURLMatcher returns a SingleLineMatcher for the URL field.
func BuildURLMatcher(fp delegate.FormulaParser) *delegate.MultiLineMatcher[*types.Stable] {
	return delegate.NewMLM[*types.Stable]("url", isDefaultURLPattern, fp, isBeginURLSequence, isEndURLSequence, cleanURLSequence)
}

// BuildMirrorMatcher returns a SingleLineMatcher for the mirror field.
func BuildMirrorMatcher(fp delegate.FormulaParser) *delegate.SingleLineMatcher[string] {
	return delegate.NewSLM[string]("mirror", isDefaultMirrorPattern, fp)
}

// BuildLicenseMatcher returns a MultiLineMatcher for the license field.
func BuildLicenseMatcher(fp delegate.FormulaParser) *delegate.MultiLineMatcher[string] {
	return delegate.NewMLM[string]("license", isDefaultLicensePattern, fp, isBeginLicenseSequence, isEndLicenseSequence, cleanLicenseSequence)
}

// BuildHeadMatcher returns a MultiLineMatcher for the head field.
func BuildHeadMatcher(fp delegate.FormulaParser) *delegate.MultiLineMatcher[*types.Head] {
	return delegate.NewMLM[*types.Head]("head", isDefaultHeadPattern, fp, isBeginHeadSequence, isEndHeadSequence, cleanHeadSequence)
}

// BuildDependencyMatcher returns a MultiLineMatcher for the dependency fields.
func BuildDependencyMatcher(fp delegate.FormulaParser) *delegate.MultiLineMatcher[*types.Dependencies] {
	return delegate.NewMLM[*types.Dependencies]("dependency", isDefaultDependencyPattern, fp, isBeginDependencySequence, isEndDependencySequence, cleanDependencySequence)
}

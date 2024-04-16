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
func BuildURLMatcher(fp delegate.FormulaParser) *delegate.SingleLineMatcher[string] {
	return delegate.NewSLM[string]("url", urlPattern, fp)
}

// BuildMirrorMatcher returns a SingleLineMatcher for the mirror field.
func BuildMirrorMatcher(fp delegate.FormulaParser) *delegate.SingleLineMatcher[string] {
	return delegate.NewSLM[string]("mirror", mirrorPattern, fp)
}

// BuildLicenseMatcher returns a MultiLineMatcher for the license field.
func BuildLicenseMatcher(fp delegate.FormulaParser) *delegate.MultiLineMatcher[string] {
	return delegate.NewMLM[string]("license", licensePattern, fp, isBeginLicenseSequence, hasUnopenedBrackets, cleanLicenseSequence)
}

// BuildHeadMatcher returns a MultiLineMatcher for the head field.
func BuildHeadMatcher(fp delegate.FormulaParser) *delegate.MultiLineMatcher[*types.Head] {
	return delegate.NewMLM[*types.Head]("head", headURLPattern, fp, isBeginHeadSequence, isEndHeadSequence, cleanHeadSequence)
}

// BuildDependencyMatcher returns a MultiLineMatcher for the dependency fields.
func BuildDependencyMatcher(fp delegate.FormulaParser) *delegate.MultiLineMatcher[[]*types.Dependency] {
	return delegate.NewMLM[[]*types.Dependency]("dependency", dependencyPattern, fp, isBeginDependencySequence, isEndDependencySequence, cleanDependencySequence)
}

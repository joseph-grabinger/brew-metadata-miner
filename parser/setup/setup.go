package setup

import (
	"main/parser/delegate"
)

func BuildStrategies(fp delegate.FormulaParser) []delegate.ParseStrategy {
	return []delegate.ParseStrategy{
		delegate.NewSLM("url", urlPattern, fp),
		delegate.NewSLM("mirror", mirrorPattern, fp),
		delegate.NewMLM("license", licensePattern, fp, isBeginLicenseSequence, hasUnopenedBrackets, cleanLicenseSequence),
		delegate.NewMLM("head", headURLPattern, fp, isBeginHeadSequence, isEndHeadSequence, cleanHeadSequence),
		BuildDependencyMatcher(fp),
	}
}

func BuildDependencyMatcher(fp delegate.FormulaParser) *delegate.MultiLineMatcher {
	return delegate.NewMLM("dependency", dependencyPattern, fp, isBeginDependencySequence, isEndDependencySequence, cleanDependencySequence)
}

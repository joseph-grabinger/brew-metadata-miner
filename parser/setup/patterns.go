package setup

import "fmt"

// RegEx patterns for parsing Formula fields.
const (
	// HomepagePattern matches two consecutive spaces,
	// followed by the literal string "homepage",
	// followed by a URL enclosed in double quotes, which is captured.
	HomepagePattern = `^\s{2}homepage\s+"([^"]+)"`

	// urlPattern matches two consecutive spaces,
	// followed by the literal the string "url",
	// followed by a string enclosed in double quotes, which is captured.
	urlPattern = `^\s{2}url\s+"([^"]+)"`

	// urlBeginPattern matches two consecutive spaces,
	// followed by the literal string "stable do".
	urlBeginPattern = `^\s{2}stable\s+do`

	// stableUrlPattern matches four consecutive spaces,
	// followed by the literal string "url",
	// followed by a string enclosed in double quotes, which is captured.
	// Further, an optional trailing comma is matched.
	stableUrlPattern = `^\s{4}url\s+"([^"]+)"(,?)`

	// stableResourcePattern matches four consecutive spaces,
	// followed by the literal string "resource".
	stableResourcePattern = `^\s{4}resource\s+`

	// stablePatchPattern matches four or more consecutive spaces,
	// followed by the literal string "patch".
	stablePatchPattern = `^\s{4,}patch\s+`

	// tagPattern matches eight consecutive spaces,
	// followed by the literal string "tag:",
	// followed by one or more whitespace characters,
	// followed by a string enclosed in double quotes, which is captured.
	tagPattern = `^\s{8}tag:\s+"([^"]+)"`

	// mirrorPattern matches two consecutive spaces,
	// folowed by the literal string "mirror",
	// followed by a string enclosed in double quotes, which is captured.
	mirrorPattern = `^\s{2}mirror\s+"([^"]+)"`

	// licensePattern matches two consecutive spaces,
	// followed by the literal string "license",
	// followed by either a string enclosed in double quotes,
	// or the keyword "all_of" followed by a sequence of strings enclosed in square brackets,
	// or the keyword "any_of" followed by a sequence of strings enclosed in square brackets,
	// or the keyword "one_of" followed by a sequence of strings enclosed in square brackets,
	// optionally followed by the "=>" symbol and the "with:" literal,
	// followed by a string value enclosed in double quotes (license exeption).
	// Everything after the license keyword is captured.
	licensePattern = `^\s{2}license\s+(:\w+|all_of\s*:\s*\[[^\]]+\]|any_of\s*:\s*\[[^\]]+\]|one_of\s*:\s*\[[^\]]+\]|"[^"]+"(\s*=>\s*{\s*with:\s*"([^"]*)"\s*})?)`

	// licenseKeywordPattern matches the string "license" with
	// zero or more leading whitespace characters and
	// optional trailing whitespace characters.
	licenseKeywordPattern = `\s+license\s*`

	// trailingCommaPattern matches a sequence that does not start with a comment character "#",
	// followed by a comma and zero or more whitespace characters,
	// optionally followed by a comment that starts with the "#" character.
	trailingCommaPattern = `^[^#]*,\s*(?:#.*)?$`

	// headURLPattern matches two consecutive spaces,
	// followed by the literal string "head",
	// followed by a string enclosed in double quotes, which is captured.
	headURLPattern = `\s{2}head\s+"([^"]+)"`

	// headBlockURLPattern matches the string "url"
	// with four leading whitespace characters, followed by a string enclosed in double quotes.
	headBlockURLPattern = `^\s{4}url\s+"([^"]+)"`

	// beginHeadPattern matches two consecutive spaces,
	// followed by the string "head do".
	beginHeadPattern = `^\s{2}head do\s*`

	// dependencyPattern matches two consecutive spaces,
	// followed by the string "depends_on", and then a string enclosed in double quotes.
	dependencyPattern = `^\s{2}depends_on\s+"[^"]+"`

	// dependencyTypePattern matches the "=>" litertal,
	// followed by optional whitespace and an optional opening square bracket,
	// followed by a colon and a word which is captured.
	// Optionally, a comma and a colon followed by a word wich is also captured,
	// followed by a closing square bracket.
	dependencyTypePattern = `=>\s*\[?:(\w+)(?:,\s*:(\w+)\])?`

	// dependencyKeywordPattern matchesthe string "depends_on"
	// followed by one or more whitespace characters,
	// then a string enclosed in double quotes, which is captured.
	// Further, it is ensured that the line does not start with a comment character "#".
	dependencyKeywordPattern = `^[^#]*depends_on\s+"([^"]+)"`

	// macOSSystemDependencyPattern matches the string "uses_from_macos"
	// followed by one or more whitespace characters,
	// and then a string enclosed in double quotes, which is captured.
	// Further, it is ensured that the line does not start with a comment character "#".
	macOSSystemDependencyPattern = `^[^#]*uses_from_macos\s+"([^"]+)"`

	// osRestrictionPattern matches a sequence beginning with a comma,
	// followed by one or more whitespace characters, the string "since:",
	// one or more whitespace characters, and then one or more word characters
	// (equivalent to [a-zA-Z0-9_]), which are captured.
	osRestrictionPattern = `,\s+since:\s+:(\w+)`

	// beginDependencyPattern matches two consecutive spaces or a tab,
	// followed by either of the listed keywords: ("depends_on" or "uses_from_macos", etc.),
	// followed by one or more whitespace characters.
	beginDependencyPattern = `^(\s{2}|\t)(depends_on|uses_from_macos|on_macos|on_arm|on_intel|on_linux|on_system|on_mojave|on_catalina|on_big_sur|on_monterey|on_ventura|on_sonoma|on_el_capitan)\s+`

	// endDependencyPatternNegated matches lines that consist entirely of whitespace characters,
	// or a comment line (starts with zero or more spaces followed by '#'),
	// or a line that starts with two or more white spaces, followed by either of the listed keywords:
	// ("depends_on", "uses_from_macos", "on_arm", etc.).
	// Further, any line strting with four or more whitespace characters followed by "fails_with" or "resource" is also matched.
	// Further, any line starting with six or more whitespace characters followed by "url" or "sha256" is also matched.
	endDependencyPatternNegated = `^(\s{2,})(depends_on|uses_from_macos|on_macos|on_arm|on_intel|on_linux|on_system|on_mojave|on_catalina|on_big_sur|on_monterey|on_ventura|on_sonoma|on_el_capitan|end|if DevelopmentTools\.)|^[\s\t]*$|^\s*#.*$|^(\s{4,}(fails_with|resource))|^(\s{6,}(url|sha256))`

	// commentPattern matches matches a sequence that starts with the "#" character,
	// followed by any sequence of characters until the end of the line.
	commentPattern = `#.*$`

	// onSystemPattern matches a line beginning two or more whitespace characters,
	// followed by the literal string "on_system".
	onSystemPattern = `^(\s{2,})on_system`

	// onSystemExtractPattern matches a sequence beginning with the literal string ":linux,",
	// followed by one or more whitespace characters, the string literal "macos:",
	// one or more whitespace characters, and then one or more word characters
	// (equivalent to [a-zA-Z0-9_]), which are captured.
	onSystemExtractPattern = `:linux,\s+macos:\s+:(\w+)`

	// onLinuxPattern matches a line beginning with two or more whitespace characters,
	// followed by the literal string "on_linux".
	onLinuxPattern = `^(\s{2,})on_linux`

	// onMacosPattern matches a line beginning with two or more whitespace characters,
	onMacosPattern = `^(\s{2,})on_macos`

	// onMacOSVersionPattern matches a line beginning with two or more whitespace characters,
	// followed by the literal string "on_" and a macOS version, which is captured.
	// Optionally, the version may be followed by a colon and a word character
	// indicating a restriction, which is also captured.
	onMacOSVersionPattern = `^(\s{2,})on_(mojave|catalina|big_sur|monterey|ventura|sonoma|el_capitan)\s+(:\w+)?`

	// onArmPattern matches a line beginning with two or more whitespace characters,
	// followed by the literal string "on_arm".
	onArmPattern = `^(\s{2,})on_arm`

	// onIntelPattern matches a line beginning with two or more whitespace characters,
	// followed by the literal string "on_intel".
	onIntelPattern = `^(\s{2,})on_intel`

	// endPatternGeneric matches a line beginning with two or more whitespace characters,
	// followed by the literal string "end".
	endPatternGeneric = `^(\s{2,})end`

	// formulaRequirementPattern matches a sequence beginning with two or more whitespace characters,
	// followed by the literal string "depends_on",
	// followed by one or more word characters (including ':'), which are captured.
	// Further, an optional sequence of whitespace characters, word characters, numbers and
	// (',', '.', '"', ':', '[', ']') is matched and captured.
	formulaRequirementPattern = `^\s{2,}depends_on\s+([:\w]+)\s*([\s\w\d<>=:,."\[\]]+)?`

	// clangVersionPattern matches a sequence beginning with the literal string "if",
	// followed by one or more whitespace characters, a word, the word character "clang_build_version",
	// one or more whitespace characters, and then a comparison operator and a number, which are captured.
	// Further, it is ensured that the line does not start with a comment character "#".
	clangVersionPattern = `^[^#]*if\s+\w+\.clang_build_version\s+([<>]?=?=?\s+\d+)`

	// failsWithPattern matches a sequence beginning with two or more whitespace characters,
	// followed by the literal string "fails_with",
	// followed by an optional sequence of words, spaces, ', =, >, and : characters,
	// and then the literal string "do".
	failsWithPattern = `^\s{2,}fails_with[=>':\w\s]*\s+do`

	// resourcePattern matches a sequence beginning with two or more whitespace characters,
	// followed by the literal string "resource",
	// followed by a string enclosed in double quotes, one or more whitespcaes
	// and the literal string "do".
	resourcePattern = `^\s{2,}resource\s+"[^"]+"\s+do`

	// InterpolationPattern matches a sequence beginning with the literal "#{",
	// followed by one or more characters that are not the closing "}" character,
	// which is captured, and ending with the "}" character.
	// This extracts a variable used for Ruby string interpolation.
	InterpolationPattern = `#\{([^}]+)\}`
)

// endPattern returns a RegEx pattern matching a sequence beginning with
// the number of given leadingSpaces, followed by the literal string "end".
func endPattern(leadingSpaces int) string {
	return fmt.Sprintf(`^\s{%d}end`, leadingSpaces)
}

// VarAssignmentPattern returns a RegEx pattern matching a sequence beggining with
// two or more whitespace characters, followed by the given varName,
// followed by an assignment ("=" character) and a string enclosed in qutotes,
// which is captured.
func VarAssignmentPattern(varName string) string {
	return fmt.Sprintf(`^\s{2,}%s\s+=\s+"([^"]+)"`, varName)
}

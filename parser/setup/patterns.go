package setup

// RegEx patterns for parsing Formula fields.
const (
	// HomepagePattern matches the string "homepage"
	// followed by a URL enclosed in double quotes.
	HomepagePattern = `homepage\s+"([^"]+)"`

	// urlPattern matches the string "url"
	// followed by a string enclosed in double quotes.
	urlPattern = `url\s+"([^"]+)"`

	// mirrorPattern matches the string "mirror"
	// followed by a string enclosed in double quotes.
	mirrorPattern = `mirror\s+"([^"]+)"`

	// licensePattern matches the word "license"
	// followed by either a string enclosed in double quotes,
	// or the keyword "all_of" followed by a sequence of strings enclosed in square brackets,
	// or the keyword "any_of" followed by a sequence of strings enclosed in square brackets,
	// or the keyword "one_of" followed by a sequence of strings enclosed in square brackets,
	// optionally followed by the "=>" symbol and a hash with the "with" key and a string value enclosed in double quotes.
	licensePattern = `^\s+license\s+(:\w+|all_of\s*:\s*\[[^\]]+\]|any_of\s*:\s*\[[^\]]+\]|one_of\s*:\s*\[[^\]]+\]|"[^"]+"(\s*=>\s*{\s*with:\s*"([^"]*)"\s*})?)`

	// licenseKeywordPattern matches the string "license" with
	// zero or more leading whitespace characters and
	// optional trailing whitespace characters.
	licenseKeywordPattern = `\s+license\s*`

	// headURLPattern matches the string "head"
	// followed by a string enclosed in double quotes, with optional leading whitespace.
	headURLPattern = `\s*head\s+"([^"]+)"`

	// headVCSPattern match the string "using" followed by a colon,
	// optional whitespace, and then a sequence of alphanumeric characters.
	headVCSPattern = `using:\s*:(\w+)`

	// headBlockURLPattern matches the string "url"
	// with four leading whitespace characters, followed by a string enclosed in double quotes.
	headBlockURLPattern = `^\s{4}url\s+"([^"]+)"`

	// beginHeadPattern matches two consecutive spaces or a tab,
	// followed by the string "head do".
	beginHeadPattern = `^(\s{2}|\t)head do\s*$`

	// endHeadPattern matches two consecutive spaces or a tab,
	// followed by the string "end".
	endHeadPattern = `^(\s{2}|\t)end\s*$`

	// dependencyPattern matches two consecutive spaces or a tab,
	// followed by the string "depends_on", and then a string enclosed in double quotes.
	dependencyPattern = `^(\s{2}|\t)depends_on\s+"[^"]+"`

	// dependencyTypePattern matches the "=>" symbol
	// followed by optional whitespace and then a colon.
	dependencyTypePattern = `=>\s*:(\w+)`

	// dependencyKeywordPattern matchesthe string "depends_on"
	// followed by one or more whitespace characters,
	// then a string enclosed in double quotes.
	dependencyKeywordPattern = `depends_on\s+"([^"]+)"`

	// macOSSystemDependencyPattern matches the string "uses_from_macos"
	// followed by one or more whitespace characters,
	// and then a string enclosed in double quotes.
	macOSSystemDependencyPattern = `uses_from_macos\s+"([^"]+)"`

	// osRestrictionPattern matches a sequence beginning with a comma,
	// followed by one or more whitespace characters, the string "since:",
	// one or more whitespace characters, and then one or more word characters
	// (equivalent to [a-zA-Z0-9_]), which are captured.
	osRestrictionPattern = `,\s+since:\s+:(\w+)`

	// beginDependencyPattern matches two consecutive spaces or a tab,
	// followed by either of the listed keywords: ("depends_on" or "uses_from_macos", etc.),
	// followed by one or more whitespace characters.
	beginDependencyPattern = `^(\s{2}|\t)(depends_on|uses_from_macos|on_arm|on_intel|on_linux|on_system|on_mojave|on_catalina|on_big_sur|on_monterey|on_ventura|on_sonoma|on_el_capitan)\s+`

	// endDependencyPattern matches lines that consist entirely of whitespace characters,
	// or a comment line (starts with zero or more spaces followed by '#'),
	// or a line that starts with two white spaces, followed by either of the listed keywords:
	// ("depends_on", "uses_from_macos", "on_arm", etc.).
	endDependencyPattern = `^(\s{2,})(depends_on|uses_from_macos|on_arm|on_intel|on_linux|on_system|on_mojave|on_catalina|on_big_sur|on_monterey|on_ventura|on_sonoma|on_el_capitan|end)|^[\s\t]*$|^\s*#.*$`

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

	// endPattern matches a line beginning with two or more whitespace characters,
	// followed by the literal string "end".
	endPattern = `^(\s{2,})end`
)

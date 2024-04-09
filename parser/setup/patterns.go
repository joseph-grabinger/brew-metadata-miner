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

	// linuxDependencyPattern matches the string "on_linux"
	// followed by one or more whitespace characters,
	// and then a string enclosed in double quotes.
	linuxDependencyPattern = `on_linux|uses_from_macos\s+"([^"]+)"`

	// beginDependencyPattern matches two consecutive spaces or a tab,
	// followed by either "depends_on" or "uses_from_macos",
	//  and then a string enclosed in double quotes.
	beginDependencyPattern = `^(\s{2}|\t)(depends_on|uses_from_macos|on_linux)\s+"[^"]+"`

	// endDependencyPattern matches lines that consist entirely of whitespace characters
	// or two consecutive spaces or a tab, followed by either
	// "depends_on" or "uses_from_macos", followed by one or more whitespace characters.
	endDependencyPattern = `^(\s{2,})(depends_on|uses_from_macos|on_linux|end)|^[\s\t]*$`
	//`^(\s{2,})(depends_on|uses_from_macos|on_linux|end)\s+|^[\s\t]*$`

	// commentPattern matches matches a line that starts with the "#" character,
	// followed by any sequence of characters until the end of the line.
	commentPattern = `#.*$`
)

package parser

// RegEx patterns for parsing Formula fields.
const (
	// homepagePattern matches the string "homepage"
	// followed by a URL enclosed in double quotes.
	homepagePattern = `homepage\s+"([^"]+)"`

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
	// or the keyword "one_of" followed by a sequence of strings enclosed in square brackets.
	licensePattern = `license\s+(:\w+|all_of\s*:\s*\[[^\]]+\]|any_of\s*:\s*\[[^\]]+\]|one_of\s*:\s*\[[^\]]+\]|"[^"]+")`

	// licenseKeywordPattern matches the string "license" with
	// optional leading and trailing whitespace characters.
	licenseKeywordPattern = `\s*license\s*`

	// headURLPattern matches the string "head"
	// followed by a string enclosed in double quotes, with optional leading whitespace.
	headURLPattern = `\s*head\s+"([^"]+)"`

	// headVCSPattern match the string "using" followed by a colon,
	// optional whitespace, and then a sequence of alphanumeric characters.
	headVCSPattern = `using:\s*:(\w+)`

	// headBlockURLPattern matches the string "url"
	// with four leading whitespace characters, followed by a string enclosed in double quotes.
	headBlockURLPattern = `^\s{4}url\s+"([^"]+)"`

	// headDependencyPattern matches four consecutive spaces or two tabs,
	// followed by the string "depends_on", and then a string enclosed in double quotes.
	// headDependencyPattern = `^(\s{4}|\t{2})depends_on\s+"[^"]+"`

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

	// systemDependencyPattern matches the string "uses_from_macos"
	// followed by one or more whitespace characters,
	// and then a string enclosed in double quotes.
	systemDependencyPattern = `uses_from_macos\s+"([^"]+)"`

	// beginDependencyPattern matches two consecutive spaces or a tab,
	// followed by either "depends_on" or "uses_from_macos",
	//  and then a string enclosed in double quotes.
	beginDependencyPattern = `^(\s{2}|\t)(depends_on|uses_from_macos)\s+"[^"]+"`

	// endDependencyPattern matches lines that consist entirely of whitespace characters
	// or two consecutive spaces or a tab, followed by either
	// "depends_on" or "uses_from_macos", followed by one or more whitespace characters.
	endDependencyPattern = `^(\s{2}|\t)(depends_on|uses_from_macos)\s+|^[\s\t]*$`

	// commentPattern matches matches a line that starts with the "#" character,
	// followed by any sequence of characters until the end of the line.
	commentPattern = `#.*$`
)

// Known hosts for repository extraction.
const (
	// githubRepoPattern matches the URL of a Github repository.
	githubRepoPattern = `https://github.com/([a-zA-Z0-9_.-]+)\/([a-zA-Z0-9_.-]+)(/|\.git|\?.*)?$`

	// gitlabRepoPattern matches the URL of a Gitlab repository.
	gitlabRepoPattern = `https://gitlab.com/([a-zA-Z0-9_.-]+)\/([a-zA-Z0-9_.-]+)(/|\.git|\?.*)?$`

	// bitbucketRepoPattern matches the URL of a Bitbucket repository.
	bitbucketRepoPattern = `https://bitbucket.org/([a-zA-Z0-9_.-]+)\/([a-zA-Z0-9_.-]+)(/|\.git|\?.*)?$`

	// repoPattern represents a general repo pattern that matches the URL of any repository.
	repoPattern = `(https:\/\/[a-zA-Z0-9.-]+)\/([a-zA-Z0-9_.-]+)\/([a-zA-Z0-9_.-]+)`

	// githubArchivePattern matches the URL of a Github archive.
	githubArchivePattern = `(https://github.com/[a-zA-Z0-9_.-]+\/[a-zA-Z0-9_.-]+)\/(?:releases\/download|archive)\/.*` // archive\/refs\/tags

	// gitlabArchivePattern matches the URL of a Gitlab archive.
	gitlabArchivePattern = `(https://gitlab.com/[a-zA-Z0-9_.-]+\/[a-zA-Z0-9_.-]+)\/(-\/archive|uploads)\/.*`

	//bitbucketArchivePattern matches the URL of a Bitbucket archive.
	bitbucketArchivePattern = `(https://bitbucket.org/[a-zA-Z0-9_.-]+\/[a-zA-Z0-9_.-]+)\/(downloads|get)\/.*`
)

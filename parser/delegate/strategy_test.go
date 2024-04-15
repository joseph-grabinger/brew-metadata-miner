package delegate_test

import (
	"bufio"
	"log"
	"strings"
	"testing"

	"main/parser/delegate"
	"main/parser/setup"
	"main/parser/types"

	"github.com/stretchr/testify/assert"
)

var MlmDependencyTests = []struct {
	input    string
	expected []*types.Dependency
}{
	{
		input: `  depends_on "pkg-config" => :build
		depends_on "libassuan"
		depends_on "libgpg-error"
	  
		on_linux do
		  depends_on "libsecret"
		end
		
		def install`, // pinentry.rb
		expected: []*types.Dependency{
			{Name: "pkg-config", DepType: "build", SystemRequirement: ""},
			{Name: "libassuan", DepType: "", SystemRequirement: ""},
			{Name: "libgpg-error", DepType: "", SystemRequirement: ""},
			{Name: "libsecret", DepType: "", SystemRequirement: "linux"}, // on_linux
		},
	},
	{
		input: `  depends_on "autoconf" => :build
		depends_on "automake" => :build
		depends_on "libtool" => :build
	  
		on_arm do
		  # Using Apple clang to compile test results in executable that
		  # causes a segmentation fault, but LLVM clang or GCC seem to work.
		  # Issue ref: https://github.com/sustrik/libdill/issues/208
		  depends_on "llvm" => :test
		end
	  
		# Apply upstream commit to fix build with newer GCC.
		# Remove with next release.
		patch do`, // libdill.rb
		expected: []*types.Dependency{
			{Name: "autoconf", DepType: "build", SystemRequirement: ""},
			{Name: "automake", DepType: "build", SystemRequirement: ""},
			{Name: "libtool", DepType: "build", SystemRequirement: ""},
			{Name: "llvm", DepType: "test", SystemRequirement: "arm"}, // on_arm
		},
	},
	{
		input: `  depends_on "go" => :build
		depends_on "node" => :build
		depends_on "yarn" => :build
	  
		uses_from_macos "python" => :build, since: :catalina
		uses_from_macos "zlib"
	  
		on_system :linux, macos: :mojave_or_older do
		  # Workaround for old node-gyp that needs distutils.
		  # TODO: Remove when node-gyp is v10+
		  depends_on "python-setuptools" => :build
		end
	  
		on_linux do
		  depends_on "fontconfig"
		  depends_on "freetype"
		end
	  
		def install`, // grafana.rb
		expected: []*types.Dependency{
			{Name: "go", DepType: "build", SystemRequirement: ""},
			{Name: "node", DepType: "build", SystemRequirement: ""},
			{Name: "yarn", DepType: "build", SystemRequirement: ""},
			{Name: "python", DepType: "build", SystemRequirement: "linux, macos: < catalina"},           // uses_from_macos
			{Name: "zlib", DepType: "", SystemRequirement: "linux"},                                     // uses_from_macos
			{Name: "python-setuptools", DepType: "build", SystemRequirement: "linux, macos: <= mojave"}, // on_system
			{Name: "fontconfig", DepType: "", SystemRequirement: "linux"},                               // on_linux
			{Name: "freetype", DepType: "", SystemRequirement: "linux"},                                 // on_linux
		},
	},
	{
		input: `  depends_on "asciidoc" => :build
		depends_on "cmake" => :build
		depends_on "docbook-xsl" => :build
		depends_on "pkg-config" => :build
		depends_on "openssl@3"
		depends_on "pinentry"
	  
		uses_from_macos "curl"
		uses_from_macos "libxslt"
	  
		# Avoid crashes on Mojave's version of libcurl (https://github.com/lastpass/lastpass-cli/issues/427)
		on_mojave :or_newer do
		  depends_on "curl"
		end
	  
		def install`, // lastpass-cli.rb
		expected: []*types.Dependency{
			{Name: "asciidoc", DepType: "build", SystemRequirement: ""},
			{Name: "cmake", DepType: "build", SystemRequirement: ""},
			{Name: "docbook-xsl", DepType: "build", SystemRequirement: ""},
			{Name: "pkg-config", DepType: "build", SystemRequirement: ""},
			{Name: "openssl@3", DepType: "", SystemRequirement: ""},
			{Name: "pinentry", DepType: "", SystemRequirement: ""},
			{Name: "curl", DepType: "", SystemRequirement: "linux, macos: >= mojave"}, // uses_from_macos & on_mojave
			{Name: "libxslt", DepType: "", SystemRequirement: "linux"},                // uses_from_macos
		},
	},
	{
		input: `  on_macos do
		depends_on "coreutils" => :build
		depends_on "gcc" if DevelopmentTools.clang_build_version <= 1403
	
		on_arm do
		  depends_on "gcc"
		  depends_on macos: :ventura
		  fails_with :clang
		end
	  end
	
	  on_ventura do
		depends_on "gcc"
		fails_with :clang
	  end
	
	  # -ftree-loop-vectorize -flto=12 -s
	  fails_with :clang do`, // btop.rb
		expected: []*types.Dependency{
			{Name: "coreutils", DepType: "build", SystemRequirement: "macos"},                      // on_macos
			{Name: "gcc", DepType: "", SystemRequirement: "(macos, (macos, arm)), macos: ventura"}, // on_macos & on_macos, on_arm & on_ventura
		},
	},
	{
		input: `  depends_on "cmake" => :build
		depends_on "node"
		depends_on "python@3.12"
		depends_on "yuicompressor"
	  
		uses_from_macos "zlib"
	  
		# OpenJDK is needed as a dependency on Linux and ARM64 for google-closure-compiler,
		# an emscripten dependency, because the native GraalVM image will not work.
		on_macos do
		  on_arm do
			depends_on "openjdk"
		  end
		end
	  
		on_linux do
		  depends_on "openjdk"
		end
	  
		fails_with gcc: "5"`, // emscripten.rb
		expected: []*types.Dependency{
			{Name: "cmake", DepType: "build", SystemRequirement: ""},
			{Name: "node", DepType: "", SystemRequirement: ""},
			{Name: "python@3.12", DepType: "", SystemRequirement: ""},
			{Name: "yuicompressor", DepType: "", SystemRequirement: ""},
			{Name: "zlib", DepType: "", SystemRequirement: "linux"},                  // uses_from_macos
			{Name: "openjdk", DepType: "", SystemRequirement: "(macos, arm), linux"}, // uses_from_macos
		},
	},
	{
		input: `  depends_on "rust" => :build

		# Use llvm@15 to work around build failure with Clang 16 described in
		# https://github.com/rust-lang/rust-bindgen/issues/2312.
		# TODO: Switch back to 'uses_from_macos "llvm" => :build' when 'bindgen' is
		# updated to 0.62.0 or newer. There is a check in the 'install' method.
		on_macos do
		  depends_on "llvm@15" => :build if DevelopmentTools.clang_build_version >= 1500
		end
		on_linux do
		  depends_on "llvm@15" => :build # for libclang
		  depends_on "pkg-config" => :build
		  depends_on "openssl@3" # Uses Secure Transport on macOS
		end
	  
		# Backport fix for build error with Rust 1.71.0. Remove in the next release.
		patch do`, // grin-wallet.rb
		expected: []*types.Dependency{
			{Name: "rust", DepType: "build", SystemRequirement: ""},
			{Name: "llvm@15", DepType: "build", SystemRequirement: "macos, linux"}, // on_macos & on_linux
			{Name: "pkg-config", DepType: "build", SystemRequirement: "linux"},     // on_linux
			{Name: "openssl@3", DepType: "", SystemRequirement: "linux"},           // on_linux
		},
	},
	{
		input: `  on_macos do
		on_arm do
		  depends_on "gettext" => :build
		  on_mojave do
			depends_on "babl" => :test
		  end
		end
		on_intel do
		  depends_on "getmail" => :build
		end
	  end
	  def install`, // pseudo
		expected: []*types.Dependency{
			{Name: "gettext", DepType: "build", SystemRequirement: "macos, arm"},
			{Name: "babl", DepType: "test", SystemRequirement: "macos, arm, macos: mojave"},
			{Name: "getmail", DepType: "build", SystemRequirement: "macos, intel"},
		},
	},
}

func TestMultiLineMatcherDependencies(t *testing.T) {
	for _, test := range MlmDependencyTests {
		formulaParser := &delegate.FormulaParser{
			Scanner: bufio.NewScanner(strings.NewReader(test.input)),
		}

		mlm := setup.BuildDependencyMatcher(*formulaParser)

		for formulaParser.Scanner.Scan() {
			line := formulaParser.Scanner.Text()

			if mlm.MatchesLine(line) {
				fieldValue, err := mlm.ExtractFromLine(line)
				if err != nil {
					log.Fatal(err)
				}

				dependencies := make([]*types.Dependency, 0)
				if fieldValue != nil {
					dependencies = fieldValue.([]*types.Dependency)
				}

				log.Println("Matched: ", dependencies)

				assert.ElementsMatch(t, test.expected, dependencies, "expected: %v, got: %v", test.expected, dependencies)
				break
			}
		}
	}
}

var MlmLicenseTests = []struct {
	input    string
	expected string
}{
	{
		input:    `  license any_of: ["Apache-2.0", "MIT"]`, // wagyu.rb
		expected: "any_of: [\"Apache-2.0\", \"MIT\"]",
	},
	{
		input: `  license all_of: [
			"BSD-2-Clause", # file-dotlock.h
			"BSD-3-Clause",
			"BSD-4-Clause",
			"ISC",
			"HPND-sell-variant", # GSSAPI code
			"RSA-MD", # MD5 code
		  ]`, // s-nail.rb
		expected: "all_of: [\"BSD-2-Clause\",\"BSD-3-Clause\",\"BSD-4-Clause\",\"ISC\",\"HPND-sell-variant\",\"RSA-MD\",]",
	},
	{
		input: `  license all_of: [
			"BSD-2-Clause",
			"LGPL-2.0-only",
			"LGPL-2.0-or-later",
			any_of: ["LGPL-2.0-only", "LGPL-3.0-only"],
		  ]
		  head "https://invent.kde.org/frameworks/karchive.git", branch: "master"`, // karchive.rb
		expected: "all_of: [\"BSD-2-Clause\",\"LGPL-2.0-only\",\"LGPL-2.0-or-later\",any_of: [\"LGPL-2.0-only\", \"LGPL-3.0-only\"],]",
	},
	{
		input:    `  license :public_domain`, // latexml.rb
		expected: ":public_domain",
	},
	{
		input:    `  license all_of: ["MIT", :cannot_represent]`, // halibut.rb
		expected: "all_of: [\"MIT\", :cannot_represent]",
	},
	{
		input: `  license "GPL-2.0-only" => { with: "GCC-exception-2.0" }
		`, // libgit2@1.6.rb
		expected: "\"GPL-2.0-only\" => { with: \"GCC-exception-2.0\" }",
	},
	{
		input: `  license any_of: [
			"CDDL-1.1",
			{ "GPL-2.0-only" => { with: "Classpath-exception-2.0" } },
		  ]`, // payara.rb
		expected: "any_of: [\"CDDL-1.1\",{ \"GPL-2.0-only\" => { with: \"Classpath-exception-2.0\" } },]",
	},
}

func TestMultiLineMatcherLicense(t *testing.T) {
	for _, test := range MlmLicenseTests {
		formulaParser := &delegate.FormulaParser{
			Scanner: bufio.NewScanner(strings.NewReader(test.input)),
		}

		mlm := setup.BuildLicenseMatcher(*formulaParser)

		for formulaParser.Scanner.Scan() {
			line := formulaParser.Scanner.Text()

			if mlm.MatchesLine(line) {
				fieldValue, err := mlm.ExtractFromLine(line)
				if err != nil {
					log.Fatal(err)
				}

				var license string
				if fieldValue != nil {
					license = fieldValue.(string)
				}
				log.Println("Matched: ", license)

				assert.Equal(t, test.expected, license, "expected: %v, got: %v", test.expected, license)
				break
			}
		}
	}
}

var MlmHeadTests = []struct {
	input    string
	expected interface{}
}{
	{
		input:    `  head "https://github.com/EnzymeAD/Enzyme.git", branch: "main"`, // enzyme.rb
		expected: "https://github.com/EnzymeAD/Enzyme.git",
	},
	{
		input: `  head do
    url "https://github.com/bcgsc/abyss.git", branch: "master"
	
		depends_on "autoconf" => :build
		depends_on "automake" => :build
		depends_on "multimarkdown" => :build
  end
	`, // abyss.rb
		expected: &types.Head{
			URL: "https://github.com/bcgsc/abyss.git",
			Dependencies: []*types.Dependency{
				{Name: "autoconf", DepType: "build"},
				{Name: "automake", DepType: "build"},
				{Name: "multimarkdown", DepType: "build"},
			},
		},
	},
	{
		input:    `  head "http://hg.code.sf.net/p/optipng/mercurial", using: :hg`, // optipng.rb
		expected: "http://hg.code.sf.net/p/optipng/mercurial",
	},
}

func TestMultiLineMatcherHead(t *testing.T) {
	for _, test := range MlmHeadTests {
		formulaParser := &delegate.FormulaParser{
			Scanner: bufio.NewScanner(strings.NewReader(test.input)),
		}

		mlm := setup.BuildHeadMatcher(*formulaParser)

		for formulaParser.Scanner.Scan() {
			line := formulaParser.Scanner.Text()

			if mlm.MatchesLine(line) {
				head, err := mlm.ExtractFromLine(line)
				if err != nil {
					log.Fatal(err)
				}

				log.Println("Matched: ", head)

				assert.Equal(t, test.expected, head, "expected: %v, got: %v", test.expected, head)
				break
			}
		}
	}
}

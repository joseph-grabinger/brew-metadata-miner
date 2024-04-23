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

var MlmUrlTests = []struct {
	input    string
	expected *types.Stable
}{
	{
		input: `  # TODO: Remove '-fcommon' workaround and switch to 'sdl2' on next release
  stable do
    url "http://www.hampa.ch/pub/pce/pce-0.2.2.tar.gz"
    sha256 "a8c0560fcbf0cc154c8f5012186f3d3952afdbd144b419124c09a56f9baab999"
    depends_on "sdl12-compat"
  end`, // pce.rb
		expected: &types.Stable{
			URL: "http://www.hampa.ch/pub/pce/pce-0.2.2.tar.gz",
			Dependencies: []*types.Dependency{
				{Name: "sdl12-compat", DepType: ""},
			},
		},
	},
	{
		input: `
  stable do
    url "https://github.com/mpalmer/action-validator/archive/refs/tags/v0.6.0.tar.gz"
    sha256 "bdec75f6383a887986192685538a736c88be365505e950aab262977c8845aa88"

    # always pull the HEAD commit hash
    resource "schemastore" do
      url "https://github.com/SchemaStore/schemastore.git",
          revision: "7bf746bd90d7e88cd11f0a9dc4bc34c91fbbf7b4"
    end
  end
	  `, // action-validator.rb
		expected: &types.Stable{
			URL:          "https://github.com/mpalmer/action-validator/archive/refs/tags/v0.6.0.tar.gz",
			Dependencies: []*types.Dependency{},
		},
	},
	{
		input: `
  stable do
    url "https://github.com/ArtifexSoftware/ghostpdl-downloads/releases/download/gs10030/ghostpdl-10.03.0.tar.xz"
    sha256 "854fd1958711b9b5108c052a6d552b906f1e3ebf3262763febf347d77618639d"

    on_macos do
      # 1. Prevent dependent rebuilds on minor version bumps.
      # Reported upstream at:
      #   https://bugs.ghostscript.com/show_bug.cgi?id=705907
      patch :DATA
    end
  end
		`, // ghostscript.rb
		expected: &types.Stable{
			URL:          "https://github.com/ArtifexSoftware/ghostpdl-downloads/releases/download/gs10030/ghostpdl-10.03.0.tar.xz",
			Dependencies: []*types.Dependency{},
		},
	},
	{
		input: `
  stable do
    url "https://github.com/HaxeFoundation/neko/archive/refs/tags/v2-3-0/neko-2.3.0.tar.gz"
    sha256 "850e7e317bdaf24ed652efeff89c1cb21380ca19f20e68a296c84f6bad4ee995"

    depends_on "pcre"

    on_linux do
      depends_on "gtk+" # On mac, neko uses carbon. On Linux it uses gtk2
    end

    # Don't redefine MSG_NOSIGNAL -- https://github.com/HaxeFoundation/neko/pull/217
    patch do
      url "https://github.com/HaxeFoundation/neko/commit/24a5e8658a104ae0f3afe66ef1906bb7ef474bfa.patch?full_index=1"
      sha256 "1a707e44b7c1596c4514e896211356d1b35d4e4b578b14b61169a7be47e91ccc"
    end

    # Fix -Wimplicit-function-declaration issue in libs/ui/ui.c
    # https://github.com/HaxeFoundation/neko/pull/218
    patch do
      url "https://github.com/HaxeFoundation/neko/commit/908149f06db782f6f1aa35723d6a403472a2d830.patch?full_index=1"
      sha256 "3e9605cccf56a2bdc49ff6812eb56f3baeb58e5359601a8215d1b704212d2abb"
    end

    # Fix -Wimplicit-function-declaration issue in libs/std/process.c
    # https://github.com/HaxeFoundation/neko/pull/219
    patch do
      url "https://github.com/HaxeFoundation/neko/commit/1a4bfc62122aef27ce4bf27122ed6064399efdc4.patch?full_index=1"
      sha256 "7fbe2f67e076efa2d7aa200456d4e5cc1e06d21f78ac5f2eed183f3fcce5db96"
    end

    # Fix mariadb-connector-c CMake error: "Flow control statements are not properly nested."
    # https://github.com/HaxeFoundation/neko/pull/225
    patch do
      url "https://github.com/HaxeFoundation/neko/commit/660fba028af1b77be8cb227b8a44cc0ef16aba79.patch?full_index=1"
      sha256 "7b0a60494eaef7c67cd15e5d80d867fee396ac70e99000603fba0dc3cd5e1158"
    end

    # Fix m1 specifics
    # https://github.com/HaxeFoundation/neko/pull/224
    patch do
      url "https://github.com/HaxeFoundation/neko/commit/ff5da9b0e96cc0eabc44ad2c10b7a92623ba49ee.patch?full_index=1"
      sha256 "ac843dfc7585535f3b08fee2b22e667fa6c38e62dcf8374cdfd1d8fcbdbcdcfd"
    end
  end`, // neko.rb
		expected: &types.Stable{
			URL: "https://github.com/HaxeFoundation/neko/archive/refs/tags/v2-3-0/neko-2.3.0.tar.gz",
			Dependencies: []*types.Dependency{
				{Name: "pcre", DepType: ""},
				{Name: "gtk+", DepType: "", SystemRequirement: "linux"},
			},
		},
	},
	{
		input: `  stable do
    url "https://github.com/unisonweb/unison.git",
        tag:      "release/M5j",
        revision: "7778bdc1a1e97e82a6ae3910a7ed10074297ff27"
    version "M5j"

    resource "local-ui" do
      url "https://github.com/unisonweb/unison-local-ui/archive/refs/tags/release/M5j.tar.gz"
      version "M5j"
      sha256 "99f8dd4c86b1cae263f16b2e04ace88764a8a1b138cead4756ceaadb7899c338"
    end
  end`, // unisonlang.rb
		expected: &types.Stable{
			URL:          "https://github.com/unisonweb/unison/tree/release/M5j",
			Dependencies: []*types.Dependency{},
		},
	},
	{
		input: `class Icecast < Formula
  desc "Streaming MP3 audio server"
  homepage "https://icecast.org/"
  url "https://downloads.xiph.org/releases/icecast/icecast-2.4.4.tar.gz", using: :homebrew_curl
  mirror "https://ftp.osuosl.org/pub/xiph/releases/icecast/icecast-2.4.4.tar.gz"
  sha256 "49b5979f9f614140b6a38046154203ee28218d8fc549888596a683ad604e4d44"
  revision 2`, // icecast.rb
		expected: &types.Stable{
			URL: "https://downloads.xiph.org/releases/icecast/icecast-2.4.4.tar.gz",
		},
	},
}

func TestMultiLineMatcherURL(t *testing.T) {
	for _, test := range MlmUrlTests {
		formulaParser := &delegate.FormulaParser{
			Scanner: bufio.NewScanner(strings.NewReader(test.input)),
		}

		mlm := setup.BuildURLMatcher(*formulaParser)

		for formulaParser.Scanner.Scan() {
			line := formulaParser.Scanner.Text()

			if mlm.MatchesLine(line) {
				stable, err := mlm.ExtractFromLine(line)
				if err != nil {
					log.Fatal(err)
				}

				log.Println("Matched: ", stable)

				assert.Equal(t, test.expected, stable, "expected: %v, got: %v", test.expected, stable)
				break
			}
		}
	}
}

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
	{
		input: `class Djvu2pdf < Formula
		desc "Small tool to convert Djvu files to PDF files"
		homepage "https://0x2a.at/site/projects/djvu2pdf/"
		url "https://0x2a.at/site/projects/djvu2pdf/djvu2pdf-0.9.2.tar.gz"
		sha256 "afe86237bf4412934d828dfb5d20fe9b584d584ef65b012a893ec853c1e84a6c"
	  
		livecheck do
		  url :homepage
		  regex(/href=.*?djvu2pdf[._-]v?(\d+(?:\.\d+)+)\.t/i)
		end
	  
		bottle do
		  sha256 cellar: :any_skip_relocation, all: "712580b5fb3dc550722146cdce9ce17c9928565a29433bd0697cb231691e566f"
		end
	  
		depends_on "djvulibre"
		depends_on "ghostscript"
	  
		def install
		  bin.install "djvu2pdf"
		  man1.install "djvu2pdf.1.gz"
		end
	  
		test do
		  system "#{bin}/djvu2pdf", "-h"
		end
	  end
	  `, // djvu2pdf.rb
		expected: "",
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
	expected *types.Head
}{
	{
		input:    `  head "https://github.com/EnzymeAD/Enzyme.git", branch: "main"`, // enzyme.rb
		expected: &types.Head{URL: "https://github.com/EnzymeAD/Enzyme.git"},
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
		expected: &types.Head{URL: "http://hg.code.sf.net/p/optipng/mercurial"},
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

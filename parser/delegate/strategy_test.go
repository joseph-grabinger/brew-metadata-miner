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
			Dependencies: &types.Dependencies{
				Lst: []*types.Dependency{
					{Name: "sdl12-compat", DepType: []string{}},
				},
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
			URL: "https://github.com/mpalmer/action-validator/archive/refs/tags/v0.6.0.tar.gz",
			Dependencies: &types.Dependencies{
				Lst: []*types.Dependency{},
			},
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
			URL: "https://github.com/ArtifexSoftware/ghostpdl-downloads/releases/download/gs10030/ghostpdl-10.03.0.tar.xz",
			Dependencies: &types.Dependencies{
				Lst: []*types.Dependency{},
			},
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
			Dependencies: &types.Dependencies{
				Lst: []*types.Dependency{
					{Name: "pcre", DepType: []string{}},
					{Name: "gtk+", DepType: []string{}, Restriction: "linux"},
				},
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
			URL: "https://github.com/unisonweb/unison/tree/release/M5j",
			Dependencies: &types.Dependencies{
				Lst: []*types.Dependency{},
			},
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
	{
		input: `  stable do
    url "https://github.com/chakra-core/ChakraCore/archive/refs/tags/v1.11.24.tar.gz"
    sha256 "b99e85f2d0fa24f2b6ccf9a6d2723f3eecfe986a9d2c4d34fa1fd0d015d0595e"

    depends_on arch: :x86_64 # https://github.com/chakra-core/ChakraCore/issues/6860

    # Fix build with modern compilers.
    # Remove with 1.12.
    patch do
      url "https://raw.githubusercontent.com/Homebrew/formula-patches/204ce95fb69a2cd523ccb0f392b7cce4f791273a/chakra/clang10.patch"
      sha256 "5337b8d5de2e9b58f6908645d9e1deb8364d426628c415e0e37aa3288fae3de7"
    end

    # Support Python 3.
    # Remove with 1.12.
    patch do
      url "https://raw.githubusercontent.com/Homebrew/formula-patches/308bb29254605f0c207ea4ed67f049fdfe5ec92c/chakra/python3.patch"
      sha256 "61c61c5376bc28ac52ec47e6d4c053eb27c04860aa4ba787a78266840ce57830"
    end
  end
	`,
		expected: &types.Stable{
			URL: "https://github.com/chakra-core/ChakraCore/archive/refs/tags/v1.11.24.tar.gz",
			Dependencies: &types.Dependencies{
				Lst:                []*types.Dependency{},
				SystemRequirements: "x86_64",
			},
		},
	},
	{
		input: `  url "https://github.com/zyantific/zydis.git",
      tag:      "v4.1.0",
      revision: "569320ad3c4856da13b9dbf1f0d9e20bda63870e"
  license "MIT"`, // zydis.rb
		expected: &types.Stable{
			URL: "https://github.com/zyantific/zydis/tree/v4.1.0",
			Dependencies: &types.Dependencies{
				Lst:                []*types.Dependency{},
				SystemRequirements: "",
			},
		},
	},
	{
		input: `  stable do
		url "https://downloads.xiph.org/releases/theora/libtheora-1.1.1.tar.bz2", using: :homebrew_curl
		mirror "https://ftp.osuosl.org/pub/xiph/releases/theora/libtheora-1.1.1.tar.bz2"
		sha256 "b6ae1ee2fa3d42ac489287d3ec34c5885730b1296f0801ae577a35193d3affbc"
	
		# Fix -flat_namespace being used on Big Sur and later.
		patch do
		  url "https://raw.githubusercontent.com/Homebrew/formula-patches/03cf8088210822aa2c1ab544ed58ea04c897d9c4/libtool/configure-pre-0.4.2.418-big_sur.diff"
		  sha256 "83af02f2aa2b746bb7225872cab29a253264be49db0ecebb12f841562d9a2923"
		end
	  end`, // theora.rb
		expected: &types.Stable{
			URL: "https://downloads.xiph.org/releases/theora/libtheora-1.1.1.tar.bz2",
			Dependencies: &types.Dependencies{
				Lst:                []*types.Dependency{},
				SystemRequirements: "",
			},
		},
	},
	{
		input: `  url "https://gitlab.gnome.org/GNOME/phodav.git", tag: "v3.0", revision: "d733fd853f0664ad8035b1b85604c62de0e97098"
		license "LGPL-2.1-only"`, // phodav.rb
		expected: &types.Stable{
			URL: "https://gitlab.gnome.org/GNOME/phodav/tree/v3.0",
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
	expected *types.Dependencies
}{
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
		expected: &types.Dependencies{
			Lst: []*types.Dependency{
				{Name: "go", DepType: []string{"build"}, Restriction: ""},
				{Name: "node", DepType: []string{"build"}, Restriction: ""},
				{Name: "yarn", DepType: []string{"build"}, Restriction: ""},
				{Name: "python", DepType: []string{"build"}, Restriction: "linux or macos: < catalina"},           // uses_from_macos
				{Name: "zlib", DepType: []string{}, Restriction: "linux"},                                         // uses_from_macos
				{Name: "python-setuptools", DepType: []string{"build"}, Restriction: "linux or macos: <= mojave"}, // on_system
				{Name: "fontconfig", DepType: []string{}, Restriction: "linux"},                                   // on_linux
				{Name: "freetype", DepType: []string{}, Restriction: "linux"},                                     // on_linux
			},
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
		expected: &types.Dependencies{
			Lst: []*types.Dependency{
				{Name: "asciidoc", DepType: []string{"build"}, Restriction: ""},
				{Name: "cmake", DepType: []string{"build"}, Restriction: ""},
				{Name: "docbook-xsl", DepType: []string{"build"}, Restriction: ""},
				{Name: "pkg-config", DepType: []string{"build"}, Restriction: ""},
				{Name: "openssl@3", DepType: []string{}, Restriction: ""},
				{Name: "pinentry", DepType: []string{}, Restriction: ""},
				{Name: "curl", DepType: []string{}, Restriction: "linux or macos: >= mojave"}, // uses_from_macos & on_mojave
				{Name: "libxslt", DepType: []string{}, Restriction: "linux"},                  // uses_from_macos
			},
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
		expected: &types.Dependencies{
			Lst: []*types.Dependency{
				{Name: "coreutils", DepType: []string{"build"}, Restriction: "macos"},                                                     // on_macos
				{Name: "gcc", DepType: []string{}, Restriction: "(macos and clang version <= 1403) or (macos and arm) or macos: ventura"}, // on_macos & on_macos, on_arm & on_ventura
			},
			SystemRequirements: "macos >= ventura (or linux)",
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
		expected: &types.Dependencies{
			Lst: []*types.Dependency{
				{Name: "cmake", DepType: []string{"build"}, Restriction: ""},
				{Name: "node", DepType: []string{}, Restriction: ""},
				{Name: "python@3.12", DepType: []string{}, Restriction: ""},
				{Name: "yuicompressor", DepType: []string{}, Restriction: ""},
				{Name: "zlib", DepType: []string{}, Restriction: "linux"},                       // uses_from_macos
				{Name: "openjdk", DepType: []string{}, Restriction: "(macos and arm) or linux"}, // uses_from_macos
			},
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
		expected: &types.Dependencies{
			Lst: []*types.Dependency{
				{Name: "rust", DepType: []string{"build"}, Restriction: ""},
				{Name: "llvm@15", DepType: []string{"build"}, Restriction: "(macos and clang version >= 1500) or linux"}, // on_macos & on_linux
				{Name: "pkg-config", DepType: []string{"build"}, Restriction: "linux"},                                   // on_linux
				{Name: "openssl@3", DepType: []string{}, Restriction: "linux"},                                           // on_linux
			},
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
		expected: &types.Dependencies{
			Lst: []*types.Dependency{
				{Name: "gettext", DepType: []string{"build"}, Restriction: "macos and arm"},
				{Name: "babl", DepType: []string{"test"}, Restriction: "macos and arm and macos: mojave"},
				{Name: "getmail", DepType: []string{"build"}, Restriction: "macos and intel"},
			},
		},
	},
	{
		input: `  depends_on xcode: ["15.0", :build]
		depends_on arch: :arm64
		depends_on :macos
		depends_on macos: :ventura
		uses_from_macos "swift"

		def install`, // whisperkit-cli.rb
		expected: &types.Dependencies{
			Lst: []*types.Dependency{
				{Name: "swift", DepType: []string{}, Restriction: "linux"}, // uses_from_macos
			},
			SystemRequirements: `xcode >= 15.0 build (on macos), arm64, macos, macos >= ventura (or linux)`,
		},
	},
	{
		input: `  depends_on "autoconf" => :build
		depends_on "automake" => :build
		depends_on "cmake" => :build
		depends_on "libtool" => :build
		depends_on "pkg-config" => :build
		depends_on "openssl@3"
		depends_on "python@3.12"

		on_macos do
		  depends_on xcode: :build
		  depends_on macos: :catalina
		end

		def install`, // retdec.rb
		expected: &types.Dependencies{
			Lst: []*types.Dependency{
				{Name: "autoconf", DepType: []string{"build"}, Restriction: ""},
				{Name: "automake", DepType: []string{"build"}, Restriction: ""},
				{Name: "cmake", DepType: []string{"build"}, Restriction: ""},
				{Name: "libtool", DepType: []string{"build"}, Restriction: ""},
				{Name: "pkg-config", DepType: []string{"build"}, Restriction: ""},
				{Name: "openssl@3", DepType: []string{}, Restriction: ""},
				{Name: "python@3.12", DepType: []string{}, Restriction: ""},
			},
			SystemRequirements: `xcode build (on macos), macos >= catalina (or linux)`,
		},
	},
	{
		input: `  depends_on "cmake" => [:build, :test]
		uses_from_macos "expat"
		uses_from_macos "libxml2"
		uses_from_macos "tcl-tk"
		uses_from_macos "zlib"

		on_macos do
		  on_arm do
			if DevelopmentTools.clang_build_version == 1316
			  depends_on "llvm" => :build

			  # clang: error: unable to execute command: Segmentation fault: 11
			  # clang: error: clang frontend command failed due to signal (use -v to see invocation)
			  # Apple clang version 13.1.6 (clang-1316.0.21.2)
			  fails_with :clang
			end
		  end
		end

		on_linux do
		  depends_on "libaec"
		  depends_on "mesa-glu"
		end

		fails_with gcc: "5"`, // vkt.rb
		expected: &types.Dependencies{
			Lst: []*types.Dependency{
				{Name: "cmake", DepType: []string{"build", "test"}, Restriction: ""},
				{Name: "expat", DepType: []string{}, Restriction: "linux"},                                         // uses_from_macos
				{Name: "libxml2", DepType: []string{}, Restriction: "linux"},                                       // uses_from_macos
				{Name: "tcl-tk", DepType: []string{}, Restriction: "linux"},                                        // uses_from_macos
				{Name: "zlib", DepType: []string{}, Restriction: "linux"},                                          // uses_from_macos
				{Name: "llvm", DepType: []string{"build"}, Restriction: "macos and arm and clang version == 1316"}, // on_macos
				{Name: "libaec", DepType: []string{}, Restriction: "linux"},                                        // on_linux
				{Name: "mesa-glu", DepType: []string{}, Restriction: "linux"},                                      // on_linux
			},
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

				dependencies := &types.Dependencies{}
				if fieldValue != nil {
					dependencies = fieldValue.(*types.Dependencies)
				}

				log.Println("Matched: ", dependencies)

				assert.ElementsMatch(t, test.expected.Lst, dependencies.Lst, "expected: %v, got: %v", test.expected.Lst, dependencies.Lst)
				assert.Equal(t, test.expected.SystemRequirements, dependencies.SystemRequirements, "expected: %v, got: %v", test.expected.SystemRequirements, dependencies.SystemRequirements)
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
				{Name: "autoconf", DepType: []string{"build"}},
				{Name: "automake", DepType: []string{"build"}},
				{Name: "multimarkdown", DepType: []string{"build"}},
			},
		},
	},
	{
		input:    `  head "http://hg.code.sf.net/p/optipng/mercurial", using: :hg`, // optipng.rb
		expected: &types.Head{URL: "http://hg.code.sf.net/p/optipng/mercurial"},
	},
	{
		input: `  license "BSD-3-Clause"
		head "https://svn.code.sf.net/p/spimsimulator/code/"
	`, // spim.rb
		expected: &types.Head{URL: "https://svn.code.sf.net/p/spimsimulator/code/"},
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

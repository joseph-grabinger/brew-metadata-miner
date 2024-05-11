package reader

import (
	"log"
	"os"
	"testing"

	"main/miner/types"

	"github.com/stretchr/testify/assert"
)

var extractFromFileTests = []struct {
	inputFilePath string
	expected      *types.SourceFormula
}{
	{
		inputFilePath: "../../test-data/i686-elf-gcc.rb",
		expected: &types.SourceFormula{
			Name:     "i686-elf-gcc",
			Homepage: "https://gcc.gnu.org",
			Stable: &types.Stable{
				URL: "https://ftp.gnu.org/gnu/gcc/gcc-13.2.0/gcc-13.2.0.tar.xz",
			},
			Mirror:  "https://ftpmirror.gnu.org/gcc/gcc-13.2.0/gcc-13.2.0.tar.xz",
			License: `"GPL-3.0-or-later" => { with: "GCC-exception-3.1" }`,
			Head:    nil,
			Dependencies: &types.Dependencies{
				Lst: []*types.Dependency{
					{Name: "gmp", DepType: []string{}},
					{Name: "i686-elf-binutils", DepType: []string{}},
					{Name: "libmpc", DepType: []string{}},
					{Name: "mpfr", DepType: []string{}},
				},
			},
		},
	},
	{
		inputFilePath: "../../test-data/pike.rb",
		expected: &types.SourceFormula{
			Name:     "pike",
			Homepage: "https://pike.lysator.liu.se/",
			Stable: &types.Stable{
				URL: "https://pike.lysator.liu.se/pub/pike/latest-stable/Pike-v8.0.1738.tar.gz",
			},
			Mirror:  "http://deb.debian.org/debian/pool/main/p/pike8.0/pike8.0_8.0.1738.orig.tar.gz",
			License: `any_of: ["GPL-2.0-only", "LGPL-2.1-only", "MPL-1.1"]`,
			Head:    nil,
			Dependencies: &types.Dependencies{
				Lst: []*types.Dependency{
					{Name: "gettext", DepType: []string{}},
					{Name: "gmp", DepType: []string{}},
					{Name: "jpeg-turbo", DepType: []string{}},
					{Name: "libtiff", DepType: []string{}},
					{Name: "nettle", DepType: []string{}},
					{Name: "pcre", DepType: []string{}},
					{Name: "webp", DepType: []string{}},
					{Name: "bzip2", DepType: []string{}, Restriction: "linux"},          // on_linux
					{Name: "krb5", DepType: []string{}, Restriction: "linux"},           // on_linux
					{Name: "libxcrypt", DepType: []string{}, Restriction: "linux"},      // on_linux
					{Name: "sqlite", DepType: []string{}, Restriction: "linux"},         // on_linux
					{Name: "zlib", DepType: []string{}, Restriction: "linux"},           // on_linux
					{Name: "gnu-sed", DepType: []string{"build"}, Restriction: "macos"}, // on_macos
					{Name: "libnsl", DepType: []string{}, Restriction: "linux"},         // on_linux
				},
			},
		},
	},
	{
		inputFilePath: "../../test-data/srecord.rb",
		expected: &types.SourceFormula{
			Name:     "srecord",
			Homepage: "https://srecord.sourceforge.net/",
			Stable: &types.Stable{
				URL: "https://downloads.sourceforge.net/project/srecord/srecord/1.64/srecord-1.64.tar.gz",
			},
			Mirror:  "",
			License: `all_of: ["GPL-3.0-or-later", "LGPL-3.0-or-later"]`,
			Head:    nil,
			Dependencies: &types.Dependencies{
				Lst: []*types.Dependency{
					{Name: "boost", DepType: []string{"build"}},
					{Name: "libtool", DepType: []string{"build"}},
					{Name: "libgcrypt", DepType: []string{}},
					{Name: "ghostscript", DepType: []string{"build"}, Restriction: "macos: >= sonoma or linux"}, // on_sonoma :or_newer && on_linux
					{Name: "groff", DepType: []string{"build"}, Restriction: "macos: >= ventura or linux"},      // on_ventura :or_newer && on_linux
				},
			},
		},
	},
	{
		inputFilePath: "../../test-data/geckodriver.rb",
		expected: &types.SourceFormula{
			Name:     "geckodriver",
			Homepage: "https://github.com/mozilla/geckodriver",
			Stable: &types.Stable{
				URL: "https://hg.mozilla.org/mozilla-central/archive/bc25087baba17c78246db06bcab71c299fd8f46f.zip/testing/geckodriver/",
				Dependencies: &types.Dependencies{
					Lst: []*types.Dependency{},
				},
			},
			Mirror:  "",
			License: `"MPL-2.0"`,
			Head: &types.Head{
				URL: "https://hg.mozilla.org/mozilla-central/",
			},
			Dependencies: &types.Dependencies{
				Lst: []*types.Dependency{
					{Name: "rust", DepType: []string{"build"}},
					{Name: "netcat", DepType: []string{"test"}, Restriction: "linux"},
					{Name: "unzip", DepType: []string{}, Restriction: "linux"},
				},
			},
		},
	},
}

func TestExtractFromFile(t *testing.T) {
	for _, test := range extractFromFileTests {
		file, err := os.Open(test.inputFilePath)
		if err != nil {
			log.Fatal(err)
		}
		defer file.Close()

		formula, err := extractFromFile(file)
		if err != nil {
			log.Fatal(err)
		}

		assert.ElementsMatch(t, test.expected.Dependencies.Lst, formula.Dependencies.Lst, "expected: %v, got: %v", test.expected.Dependencies.Lst, formula.Dependencies.Lst)
		assert.Equal(t, test.expected.Dependencies.SystemRequirements, formula.Dependencies.SystemRequirements, "expected: %v, got: %v", test.expected.Dependencies.SystemRequirements, formula.Dependencies.SystemRequirements)
		test.expected.Dependencies, formula.Dependencies = nil, nil
		assert.Equal(t, test.expected, formula, "expected: %v, got: %v", test.expected, formula)
	}
}

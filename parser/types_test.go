package parser

import (
	"testing"
)

var formatLicenseTests = []struct {
	input    string
	expected string
}{
	{
		input:    `all_of: ["BSD-2-Clause","LGPL-2.0-only","LGPL-2.0-or-later",any_of: ["LGPL-2.0-only", "LGPL-3.0-only"],]`,
		expected: "BSD-2-Clause and LGPL-2.0-only and LGPL-2.0-or-later and (LGPL-2.0-only or LGPL-3.0-only)",
	},
	{
		input:    `all_of: ["GPL-2.0-or-later", "LGPL-2.1-or-later"]`,
		expected: "GPL-2.0-or-later and LGPL-2.1-or-later",
	},
	{
		input:    `any_of: ["GPL-2.0-or-later", "LGPL-2.1-or-later", "MIT"]`,
		expected: "GPL-2.0-or-later or LGPL-2.1-or-later or MIT",
	},
	{
		input:    `all_of: ["0BSD", "LGPL-2.1-or-later", "GPL-2.0-or-later", "GPL-3.0-or-later", ]`,
		expected: "0BSD and LGPL-2.1-or-later and GPL-2.0-or-later and GPL-3.0-or-later",
	},
	{
		input:    `one_of: [:public_domain, :cannot_represent]`,
		expected: "Public Domain or Cannot Represent",
	},
	{
		input:    `all_of: ["GPL-2.0-or-later",{ any_of: ["GPL-2.0-only", "Artistic-2.0"] },]`,
		expected: "GPL-2.0-or-later and (GPL-2.0-only or Artistic-2.0)",
	},
	{
		input:    `MIT`,
		expected: "MIT",
	},
	{
		input:    ``,
		expected: "pseudo",
	},
	{
		input:    `"LGPL-2.1-only" => { with: "OCaml-LGPL-linking-exception" }`,
		expected: "LGPL-2.1-only with OCaml-LGPL-linking-exception",
	},
	{
		input:    `"GPL-2.0-or-later" => {with: "Classpath-exception-2.0",}`,
		expected: "GPL-2.0-or-later with Classpath-exception-2.0",
	},
	{
		input:    `all_of: ["BSD-3-Clause", "GFDL-1.3-no-invariants-only", "GPL-2.0-only", "GPL-3.0-only" => { with: "Qt-GPL-exception-1.0" }, "LGPL-3.0-only"]`,
		expected: "BSD-3-Clause and GFDL-1.3-no-invariants-only and GPL-2.0-only and (GPL-3.0-only with Qt-GPL-exception-1.0) and LGPL-3.0-only",
	},
	{
		input:    `all_of: ["Unlicense", "Zlib", "MIT", "BSL-1.0", "BSD-3-Clause", "Apache-2.0", "BSD-2-Clause", "Apache-2.0" => { with: "LLVM-exception" }]`,
		expected: "Unlicense and Zlib and MIT and BSL-1.0 and BSD-3-Clause and Apache-2.0 and BSD-2-Clause and (Apache-2.0 with LLVM-exception)",
	},
}

func TestFormatLicense(t *testing.T) {
	for _, test := range formatLicenseTests {
		sf := &sourceFormula{
			license: test.input,
		}
		license := sf.formatLicense()
		if license != test.expected {
			t.Errorf("expected: %s, got: %s", test.expected, license)
		}
	}
}

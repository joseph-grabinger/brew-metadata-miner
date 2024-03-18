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
		input:    `one_of: [:public_domain, :cannot_represent]`,
		expected: "Public Domain or Cannot Represent",
	},
	{
		input:    `MIT`,
		expected: "MIT",
	},
	{
		input:    ``,
		expected: "pseudo",
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

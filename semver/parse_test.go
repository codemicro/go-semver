package semver

import (
	"reflect"
	"strings"
	"testing"
)

type parseTests []struct {
	name    string
	args    string
	want    *Version
	wantErr bool
}

func (tests parseTests) Run(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Parse() = %v, want %v", got, tt.want)
			}
		})
	}
}

func msp(i string) []string {
	return strings.Split(i, ".")
}

// most test cases taken from https://github.com/zmarko/semver/blob/master/test/semver200_parser_tests.cpp

func TestParseVersionCore(t *testing.T) {
	parseTests{
		// must have a major, minor and patch build version
		{name: "Valid version core", args: "1.0.0", want: &Version{Major: 1, Minor: 0, Patch: 0, Stable: true}},
		{name: "Valid version core", args: "1.2.3", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: true}},
		{name: "Valid version core", args: "65535.65534.65533", want: &Version{Major: 65535, Minor: 65534, Patch: 65533, Stable: true}},
		{name: "Incomplete version core", args: "1", wantErr: true},
		{name: "Incomplete version core", args: "1.1", wantErr: true},

		// must not have leading zeros
		{name: "Major version leading zero", args: "01.0.0", wantErr: true},
		{name: "Minor version leading zero", args: "1.01.0", wantErr: true},
		{name: "Patch version leading zero", args: "1.0.01", wantErr: true},

		// must be non-negative integers
		{name: "Negative major version", args: "-1.0.0", wantErr: true},
		{name: "Negative minor version", args: "1.-1.0", wantErr: true},
		{name: "Negative patch version", args: "1.1.-1", wantErr: true},
		{name: "Invalid major version", args: "a.0.0", wantErr: true},
		{name: "Invalid minor version", args: "1.a.0", wantErr: true},
		{name: "Invalid minor version", args: "1.0.a", wantErr: true},

		// Incomplete
		{name: "Incomplete version core", args: "1.0-banana", wantErr: true},
		{name: "Incomplete version core", args: "1.0+banana", wantErr: true},
	}.Run(t)
}

func TestParsePrerelease(t *testing.T) {
	parseTests{
		// contains one or more dot-separated ids with distinct numeric and mixed ids
		{name: "Valid text", args: "1.0.0-alpha", want: &Version{Major: 1, Minor: 0, Patch: 0, Stable: false, Prerelease: msp("alpha")}},
		{name: "Valid text", args: "1.2.3-test", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: false, Prerelease: msp("test")}},
		{name: "Valid digits", args: "1.2.3-321", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: false, Prerelease: msp("321")}},
		{name: "Valid digit", args: "1.0.0-0", want: &Version{Major: 1, Minor: 0, Patch: 0, Stable: false, Prerelease: msp("0")}},
		{name: "Valid multipart", args: "1.2.3-test.1", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: false, Prerelease: msp("test.1")}},
		{name: "Valid multipart", args: "1.2.3-1.test", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: false, Prerelease: msp("1.test")}},
		{name: "Valid multipart", args: "1.2.3-test.123456", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: false, Prerelease: msp("test.123456")}},
		{name: "Valid multipart", args: "1.2.3-123456.test", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: false, Prerelease: msp("123456.test")}},
		{name: "Long prerelease", args: "1.2.3-1.a.22.bb.333.ccc.4444.dddd.55555.fffff", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: false, Prerelease: msp("1.a.22.bb.333.ccc.4444.dddd.55555.fffff")}},

		// contain only alphanumerics and hyphen
		{name: "Valid alphanumerics and hyphen", args: "1.2.3-test-1-2-3-CAP", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: false, Prerelease: msp("test-1-2-3-CAP")}},
		{name: "Invalid with hash symbol", args: "1.2.3-test#1", wantErr: true},
		{name: "Invalid with copyright symbol", args: "1.2.3-test.©2015", wantErr: true},
		{name: "Invalid with cyrillic", args: "1.2.3-ћирилица-1", wantErr: true},

		// ids must not be empty
		{name: "Empty", args: "1.2.3-", wantErr: true},
		{name: "Empty after dot", args: "1.2.3-test.", wantErr: true},
		{name: "Empty after dot", args: "1.2.3-test.", wantErr: true},
		{name: "Empty after two dots", args: "1.2.3-test..", wantErr: true},
		{name: "Empty in between two dots with nonempty after", args: "1.2.3-test..1", wantErr: true},

		// numeric ids must not have leading 0
		{name: "Numeric with leading zero", args: "1.2.3-01", wantErr: true},
		{name: "Numeric with leading zero", args: "1.2.3-test.0023", wantErr: true},
		{name: "Alphanumeric with leading zero", args: "1.2.3-test.01a", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: false, Prerelease: msp("test.01a")}},
		{name: "Alphanumeric with leading zero", args: "1.2.3-test.01-s", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: false, Prerelease: msp("test.01-s")}},
	}.Run(t)
}

func TestParseBuild(t *testing.T) {
	parseTests{
		// contains one or more dot separated ids
		{name: "Valid text", args: "1.0.0+test", want: &Version{Major: 1, Minor: 0, Patch: 0, Stable: true, Build: msp("test")}},
		{name: "Valid text", args: "1.2.3+test", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: true, Build: msp("test")}},
		{name: "Valid digits", args: "1.2.3+321", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: true, Build: msp("321")}},
		{name: "Valid digit", args: "1.0.0+0", want: &Version{Major: 1, Minor: 0, Patch: 0, Stable: true, Build: msp("0")}},
		{name: "Valid multipart", args: "1.2.3+test.1", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: true, Build: msp("test.1")}},
		{name: "Valid multipart", args: "1.2.3+1.test", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: true, Build: msp("1.test")}},
		{name: "Valid multipart", args: "1.2.3+test.123456", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: true, Build: msp("test.123456")}},
		{name: "Valid multipart", args: "1.2.3+123456.test", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: true, Build: msp("123456.test")}},
		{name: "Long build", args: "1.2.3+1.a.22.bb.333.ccc.4444.dddd.55555.fffff", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: true, Build: msp("1.a.22.bb.333.ccc.4444.dddd.55555.fffff")}},

		// contain only alphanumerics and hyphen
		{name: "Valid alphanumerics and hyphen", args: "1.2.3+test-1-2-3-CAP", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: true, Build: msp("test-1-2-3-CAP")}},
		{name: "Invalid with hash symbol", args: "1.2.3+test#1", wantErr: true},
		{name: "Invalid with copyright symbol", args: "1.2.3+test.©2015", wantErr: true},
		{name: "Invalid with cyrillic", args: "1.2.3+ћирилица-1", wantErr: true},

		// ids must not be empty
		{name: "Empty", args: "1.2.3+", wantErr: true},
		{name: "Empty after dot", args: "1.2.3+test.", wantErr: true},
		{name: "Empty after dot", args: "1.2.3+test.", wantErr: true},
		{name: "Empty after two dots", args: "1.2.3+test..", wantErr: true},
		{name: "Empty in between two dots with nonempty after", args: "1.2.3+test..1", wantErr: true},

		// numeric ids must not have leading 0
		{name: "Numeric with leading zero", args: "1.2.3+01", wantErr: true},
		{name: "Numeric with leading zero", args: "1.2.3+test.0023", wantErr: true},
		{name: "Alphanumeric with leading zero", args: "1.2.3+test.01a", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: true, Build: msp("test.01a")}},
		{name: "Alphanumeric with leading zero", args: "1.2.3+test.01-s", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: true, Build: msp("test.01-s")}},
	}.Run(t)
}

func TestParseSequence(t *testing.T) {
	parseTests{
		{name: "Prerelease and build", args: "1.2.3-r4+b5", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: false, Build: msp("b5"), Prerelease: msp("r4")}},
		{name: "Build only", args: "1.2.3+b4-r5", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: true, Build: msp("b4-r5")}},
	}.Run(t)
}

func TestParseEdgeCases(t *testing.T) {
	parseTests{
		{name: "Empty prerelease ID then build", args: "1.2.3-r4.+b5", wantErr: true},
		{name: "Prerelease ID then build with empty ID", args: "1.2.3-r4+b5.", wantErr: true},

		{name: "Valid with prerelease and build", args: "1.2.3-alpha+build.314", want: &Version{Major: 1, Minor: 2, Patch: 3, Stable: false, Build: msp("build.314"), Prerelease: msp("alpha")}},

		{name: "Check Stable flag", args: "0.2.3", want: &Version{Major: 0, Minor: 2, Patch: 3, Stable: false}},
	}.Run(t)
}

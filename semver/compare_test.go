package semver

import (
	"fmt"
	"testing"
)

func mkv(x string) *Version {
	y, err := Parse(x)
	if err != nil {
		panic(fmt.Errorf("%s: %v", x, err))
	}
	return y
}

type compareTests []struct {
	name   string
	fields *Version
	args   *Version
	want   int
}

func (tests compareTests) Run(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fields.CompareTo(tt.args); got != tt.want {
				t.Errorf("Version.CompareTo() = %v, want %v", got, tt.want)
			}
		})
	}
}

// most test cases taken from https://github.com/zmarko/semver/blob/master/test/semver200_comparator_tests.cpp

func TestVersion_CompareToEqual(t *testing.T) {
	compareTests{
		// normal precedence
		{fields: mkv("1.2.3"), args: mkv("1.2.3"), want: 0},

		// same prereleases
		{fields: mkv("1.0.0-alpha"), args: mkv("1.0.0-alpha"), want: 0},
		{fields: mkv("1.0.0-alpha.1"), args: mkv("1.0.0-alpha.1"), want: 0},
		{fields: mkv("1.0.0-1"), args: mkv("1.0.0-1"), want: 0},

		// build is ignored
		{fields: mkv("1.0.0"), args: mkv("1.0.0+build.1.2.3"), want: 0},
		{fields: mkv("1.0.0+ZZZ"), args: mkv("1.0.0+build.1.2.3"), want: 0},
		{fields: mkv("1.0.0+100"), args: mkv("1.0.0+200"), want: 0},
	}.Run(t)
}

func TestVersion_CompareToGreaterThan(t *testing.T) {
	compareTests{
		// normal precedence
		{fields: mkv("0.0.2"), args: mkv("0.0.1"), want: 1},
		{fields: mkv("0.2.0"), args: mkv("0.0.3"), want: 1},
		{fields: mkv("0.2.0"), args: mkv("0.1.3"), want: 1},
		{fields: mkv("2.0.0"), args: mkv("0.0.1"), want: 1},
		{fields: mkv("2.0.0"), args: mkv("0.3.1"), want: 1},
		{fields: mkv("2.0.0"), args: mkv("1.3.1"), want: 1},

		// normal and prerelease precedence
		{fields: mkv("1.0.0"), args: mkv("1.0.0-alpha"), want: 1},
		{fields: mkv("1.0.0"), args: mkv("1.0.0-99"), want: 1},
		{fields: mkv("1.0.0"), args: mkv("1.0.0-ZZ"), want: 1},

		// prerelease precedence with numeric id
		{fields: mkv("1.0.0-1"), args: mkv("1.0.0-0"), want: 1},
		{fields: mkv("1.0.0-10"), args: mkv("1.0.0-1"), want: 1},
		{fields: mkv("1.0.0-alpha.3"), args: mkv("1.0.0-alpha.1"), want: 1},

		// prerelease precedence with alphanumeric identifier
		{fields: mkv("1.0.0-1"), args: mkv("1.0.0-0"), want: 1},
		{fields: mkv("1.0.0-Z"), args: mkv("1.0.0-A"), want: 1},
		{fields: mkv("1.0.0-Z"), args: mkv("1.0.0-1"), want: 1},
		{fields: mkv("1.0.0-alpha-3"), args: mkv("1.0.0-alpha-1"), want: 1},
		{fields: mkv("1.0.0-alpha-3"), args: mkv("1.0.0-alpha-100"), want: 1},
	}.Run(t)
}

func TestVersion_CompareToLessThan(t *testing.T) {
	compareTests{
		// prerelease precedence miscellaneous
		{fields: mkv("1.0.0-alpha"), args: mkv("1.0.0-alpha.1"), want: -1},
		{fields: mkv("1.0.0-alpha.1"), args: mkv("1.0.0-alpha.beta"), want: -1},
		{fields: mkv("1.0.0-alpha.beta"), args: mkv("1.0.0-beta"), want: -1},
		{fields: mkv("1.0.0-beta"), args: mkv("1.0.0-beta.2"), want: -1},
		{fields: mkv("1.0.0-beta.2"), args: mkv("1.0.0-beta.11"), want: -1},
		{fields: mkv("1.0.0-beta.11"), args: mkv("1.0.0-rc.1"), want: -1},
		{fields: mkv("1.0.0-rc.1"), args: mkv("1.0.0"), want: -1},
	}.Run(t)
}

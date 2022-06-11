package semver

import (
	"fmt"
	"strconv"
	"strings"
)

func isStringNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

type Version struct {
	Major, Minor, Patch int
	Prerelease, Build []string
	Stable bool
}

func (v *Version) String() string {
	var (
		prerelease string
		build string
	)

	if len(v.Prerelease) != 0 {
		prerelease = "-" + strings.Join(v.Prerelease, ".")
	}

	if len(v.Build) != 0 {
		build = "+" + strings.Join(v.Build, ".")
	}

	return fmt.Sprintf("%d.%d.%d%s%s", v.Major, v.Minor, v.Patch, prerelease, build)
}

type Slice []*Version

func (s Slice) Len() int {
	return len(s)
}

func (s Slice) Less(i, j int) bool {
	return s[i].CompareTo(s[j]) == -1
}

func (s Slice) Swap(i, j int) {
	s[j], s[i] = s[i], s[j]
}
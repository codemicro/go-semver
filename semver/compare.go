package semver

import (
	"strconv"
)

// compareVersionCore compares the major, minor and patch versions of a semantic version. 1 is v > vx, -1 is v < vx,
// 0 is v == vx
func compareVersionCore(v, vx *Version) int {
	switch {
	case v.Major > vx.Major:
		return 1
	case v.Major < vx.Major:
		return -1
	case v.Minor > vx.Minor:
		return 1
	case v.Minor < vx.Minor:
		return -1
	case v.Patch > vx.Patch:
		return 1
	case v.Patch < vx.Patch:
		return -1
	default:
		return 0
	}
}

func compareNumericIdentifiers(v, vx string) int {
	vi, _ := strconv.Atoi(v)
	vxi, _ := strconv.Atoi(vx)

	switch {
	case vi > vxi:
		return 1
	case vi < vxi:
		return -1
	default:
		return 0
	}
}

func compareAlphanumericIdentifiers(v, vx string) int {
	switch {
	case v > vx:
		return 1
	case v < vx:
		return -1
	default:
		return 0
	}
}

// compareIdentifiers compares two prerelease identifiers, v and vx. A return value of 0 indicates equality, -1 indicates v < vx and 1 indicates v > vx.
func compareIdentifiers(v, vx string) int {
	vn := isStringNumeric(v)
	vxn := isStringNumeric(vx)

	switch {
	case vn && vxn:
		return compareNumericIdentifiers(v, vx)
	case !vn && !vxn:
		return compareAlphanumericIdentifiers(v, vx)
	case !vn && vxn:
		return 1
	case vn && !vxn:
		return -1
	default:
		return 0 // I'm not sure this can ever be reached but oh well :)
	}
}

// CompareTo compares two instances of Version. If the return value is 0, the versions are equal. 1 is v > vx, -1 is
// v < vx, 0 is v == vx
func (v *Version) CompareTo(vx *Version) int {

	versionCore := compareVersionCore(v, vx)

	if versionCore != 0 {
		return versionCore
	}

	// If we're here, this means the version core is equal and we need to compare the prerelease versions

	var (
		lv  = len(v.Prerelease)
		lvx = len(vx.Prerelease)
	)

	// "When major, minor, and patch are equal, a pre-release version has lower precedence than a normal version"
	if lv == 0 || lvx == 0 {
		switch {
		case lv < lvx:
			return 1
		case lv > lvx:
			return -1
		default:
			return 0
		}
	}

	fromSlice := func(x []string, i int) string {
		if i >= len(x) {
			return ""
		}
		return x[i]
	}

	vi := func(i int) string { return fromSlice(v.Prerelease, i) }
	vxi := func(i int) string { return fromSlice(vx.Prerelease, i) }

	var c int
	for {
		vc := vi(c)
		vxc := vxi(c)

		if vc == "" && vxc == "" {
			return 0
		}

		if vc == "" {
			return -1
		} else if vxc == "" {
			return 1
		}

		comp := compareIdentifiers(vc, vxc)
		if comp != 0 {
			return comp
		}

		c += 1
	}
}

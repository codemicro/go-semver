package semver

import (
	"errors"
	"strings"
)

type filterFunction func(v *Version) bool

func Filter(filter string, options Slice) (Slice, error) {

	filterFunc, err := parseFilter(filter)
	if err != nil {
		return nil, err
	}

	var n int
	for _, x := range options {
		// true denotes an item to keep
		if filterFunc(x) {
			options[n] = x
			n += 1
		}
	}
	options = options[:n]

	return options, nil
}

var (
	ErrorNoFilter             = errors.New("semver: Filter: no filter provided")
	ErrorEmptyFilter          = errors.New("semver: Filter: empty filter")
	ErrorPrereleaseDisallowed = errors.New("semver: Filter: ^ filter must not have prerelease identifiers")

	allowableFilterPrefixes = []string{"^", "~", ">=", "<=", ">", "<", "="}
)

func parseFilter(filter string) (filterFunction, error) {

	if filter == "" {
		return nil, ErrorNoFilter
	}

	if strings.Index(filter, "||") != -1 {
		// handle uses of ||

		splitOrSegments := strings.Split(filter, "||")
		var filterFunctions []filterFunction
		for _, block := range splitOrSegments {
			ff, err := parseFilter(strings.TrimSpace(block))
			if err != nil {
				return nil, err
			}
			filterFunctions = append(filterFunctions, ff)
		}

		return func(v *Version) bool {
			for _, ff := range filterFunctions {
				if ff(v) {
					return true
				}
			}
			return false
		}, nil

	}

	splitFilter := strings.Split(filter, " ")

	var filterFunctions []filterFunction

	for _, rawFilter := range splitFilter {
		if rawFilter == "" {
			return nil, ErrorEmptyFilter
		}

		var prefix string
		for _, possiblePrefix := range allowableFilterPrefixes {
			if strings.HasPrefix(rawFilter, possiblePrefix) {
				prefix = possiblePrefix
				break
			}
		}

		// If the prefix is rubbish, it'll be caught here as an error and returned
		parsedFilterVersion, err := Parse(rawFilter[len(prefix):])
		if err != nil {
			return nil, err
		}

		// this if-block must be placed below the call to Parse since it will affect what is and isn't cut off the
		// start of the `rawFilter` string
		if prefix == "" {
			// no prefix means match the version exactly
			prefix = "="
		}

		var ffunc filterFunction
		switch prefix {
		case "^":
			// match the same version and any newer versions that don't increment the first non-zero segment of the version
			// `^2.2.1` can be expanded out as `>=2.2.1 <3.0.0`

			if len(parsedFilterVersion.Prerelease) != 0 {
				return nil, ErrorPrereleaseDisallowed
			}

			upperbound := new(Version)

			if parsedFilterVersion.Major != 0 {
				upperbound.Major = parsedFilterVersion.Major + 1
			} else if parsedFilterVersion.Minor != 0 {
				upperbound.Major = parsedFilterVersion.Major
				upperbound.Minor = parsedFilterVersion.Minor + 1
			} else if parsedFilterVersion.Patch != 0 {
				upperbound.Major = parsedFilterVersion.Major
				upperbound.Minor = parsedFilterVersion.Minor
				upperbound.Patch = parsedFilterVersion.Patch + 1
			}

			ffunc = func(v *Version) bool {
				return parsedFilterVersion.CompareTo(v) <= 0 && upperbound.CompareTo(v) == 1 && len(v.Prerelease) == 0
			}

		case "~":
			// include everything greater than or equal to the stated version in the current minor range
			ffunc = func(v *Version) bool {
				if parsedFilterVersion.Minor != v.Minor || parsedFilterVersion.Major != v.Major {
					return false
				}
				return parsedFilterVersion.CompareTo(v) <= 0
			}
		case ">":
			ffunc = func(v *Version) bool {
				// v > parsedFilterVersion
				return parsedFilterVersion.CompareTo(v) == -1 && len(v.Prerelease) == 0
			}
		case "<":
			ffunc = func(v *Version) bool {
				// v < parsedFilterVersion
				return parsedFilterVersion.CompareTo(v) == 1 && len(v.Prerelease) == 0
			}
		case ">=":
			ffunc = func(v *Version) bool {
				// v >= parsedFilterVersion
				return parsedFilterVersion.CompareTo(v) <= 0 && len(v.Prerelease) == 0
			}
		case "<=":
			ffunc = func(v *Version) bool {
				// v <= parsedFilterVersion
				return parsedFilterVersion.CompareTo(v) >= 0 && len(v.Prerelease) == 0
			}
		case "=":
			ffunc = func(v *Version) bool {
				// v == parsedFilterVersion
				return parsedFilterVersion.CompareTo(v) == 0
			}
		default:
			panic("this should never happen")
		}

		filterFunctions = append(filterFunctions, ffunc)
	}

	return func(v *Version) bool {
		for _, ff := range filterFunctions {
			if !ff(v) {
				return false
			}
		}
		return true
	}, nil
}

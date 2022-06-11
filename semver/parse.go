package semver

import (
	"errors"
	"fmt"
	"strconv"
)

func isLetter(char rune) bool {
	return (65 <= char && char <= 90) || (97 <= char && char <= 122)
}

func isDigit(char rune) bool {
	return '0' <= char && char <= '9'
}

func isAlphanumericIdentifier(char rune) bool {
	return isDigit(char) || (isLetter(char) || char == '-')
}

var (
	ErrorIncompleteVersionCore     = errors.New("semver: Parse: incomplete version core")
	ErrorLeadingZero               = errors.New("semver: Parse: leading zeros are disallowed")
	ErrorEmptyPrereleaseIdentifier = errors.New("semver: Parse: empty prerelease identifier")
	ErrorEmptyBuildIdentifier      = errors.New("semver: Parse: empty build identifier")
)

func Parse(in string) (*Version, error) {

	type parseState uint8
	const (
		versionCore parseState = iota
		prerelease
		build
	)

	var state parseState
	var index int

	peek := func(offset int) rune {
		n := index + offset
		if n >= len(in) || n < 0 {
			return 0
		}
		return rune(in[n])
	}

	consume := func() rune {
		if index >= len(in) {
			return 0
		}
		x := in[index]
		index += 1
		return rune(x)
	}

	version := new(Version)
	var buf []rune

	for index < len(in) {
		switch state {
		case versionCore:

			writeBufTo := func(x *int) { *x, _ = strconv.Atoi(string(buf)) }
			errorUnknownCharacter := func() error { return fmt.Errorf("parse: unrecognised character '%s' in version core", string(peek(0))) }

			var component int
			// TODO: these nested for loops could probably be removed somehow. At present, however, this is not-trivial
			//  because of the requirement in some cases to run until `peek(0)` returns 0.
			for {

				if isDigit(peek(0)) {

					if peek(0) == '0' && len(buf) == 0 && isDigit(peek(1)) {
						return nil, ErrorLeadingZero
					}

					buf = append(buf, consume())
				} else if peek(0) == 0 {
					// end of input, nothing more to parse

					if component != 2 {
						return nil, ErrorIncompleteVersionCore
					}

					writeBufTo(&version.Patch)
					buf = nil

					break
				} else if peek(0) == '.' {

					consume()

					if component == 0 { // major number
						writeBufTo(&version.Major)
					} else if component == 1 { // minor number
						writeBufTo(&version.Minor)
					}

					component += 1
					buf = nil

				} else if len(buf) != 0 && peek(-1) != '.' {
					if peek(0) == '-' {
						// moving on to prerelease section

						if component != 2 {
							return nil, ErrorIncompleteVersionCore
						} else if peek(1) == 0 {
							return nil, ErrorEmptyPrereleaseIdentifier
						}

						writeBufTo(&version.Patch)
						buf = nil
						consume()
						state = prerelease
						break
					} else if peek(0) == '+' {
						// moving on to build section

						if component != 2 {
							return nil, ErrorIncompleteVersionCore
						} else if peek(1) == 0 {
							return nil, ErrorEmptyBuildIdentifier
						}

						writeBufTo(&version.Patch)
						buf = nil
						consume()
						state = build
						break
					} else {
						return nil, errorUnknownCharacter()
					}
				} else {
					return nil, errorUnknownCharacter()
				}
			}

		case prerelease:
			// dot separated prerelease identifiers, runs until end or '+'

			writeBuf := func() {
				version.Prerelease = append(version.Prerelease, string(buf))
			}

			for {
				if isAlphanumericIdentifier(peek(0)) || isDigit(peek(0)) {
					buf = append(buf, consume())
				} else if peek(0) == '.' {
					if peek(1) == '.' || peek(1) == '+' {
						return nil, ErrorEmptyPrereleaseIdentifier
					}
					writeBuf()
					consume()
					buf = nil
				} else if peek(0) == 0 {

					if len(buf) == 0 || buf[len(buf)-1] == '.' {
						return nil, ErrorEmptyPrereleaseIdentifier
					}

					// end
					writeBuf()
					buf = nil
					break
				} else if peek(0) == '+' {
					consume()
					writeBuf()
					buf = nil
					state = build
					break
				} else {
					return nil, fmt.Errorf("parse: unrecognised character '%s' in pre-release", string(peek(0)))
				}
			}

			for _, x := range version.Prerelease {
				if isStringNumeric(x) && x[0] == byte('0') && len(x) > 1 { // leading zeros on numeric ids disallowed
					return nil, ErrorLeadingZero
				}
			}

		case build:
			// dot separated build identifiers, runs until end

			writeBuf := func() {
				version.Build = append(version.Build, string(buf))
			}

			for {

				if isAlphanumericIdentifier(peek(0)) || isDigit(peek(0)) {
					buf = append(buf, consume())
				} else if peek(0) == '.' {
					if peek(1) == '.' {
						return nil, ErrorEmptyBuildIdentifier
					}
					writeBuf()
					consume()
					buf = nil
				} else if peek(0) == 0 {

					if len(buf) == 0 || buf[len(buf)-1] == '.' {
						return nil, ErrorEmptyBuildIdentifier
					}

					// end
					writeBuf()
					buf = nil
					break
				} else {
					return nil, fmt.Errorf("parse: unrecognised character '%s' in build", string(peek(0)))
				}
			}

			for _, x := range version.Build {
				if isStringNumeric(x) && x[0] == byte('0') && len(x) > 1 { // leading zeros on numeric ids disallowed
					return nil, ErrorLeadingZero
				}
			}

		}
	}

	version.Stable = version.Major != 0 && len(version.Prerelease) == 0

	return version, nil
}

func ParseMultiple(rawVersions []string) (Slice, error) {
	var x Slice
	for _, rawVersion := range rawVersions {
		parsedVersion, err := Parse(rawVersion)
		if err != nil {
			return nil, err
		}
		x = append(x, parsedVersion)
	}
	return x, nil
}

func MustParse(in string) *Version {
	v, err := Parse(in)
	if err != nil {
		panic(err)
	}
	return v
}
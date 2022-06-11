package semver

func Format(x string) (string, error) {
	v, err := Parse(x)
	if err != nil {
		return "", err
	}
	return v.String(), nil
}

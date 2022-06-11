package semver

import (
	"database/sql/driver"
	"fmt"
)

// Scan implements sql.Scanner
func (v *Version) Scan(data interface{}) error {

	var versionString string

	switch data := data.(type) {
	case nil:
		return nil
	case string:
		if data == "" {
			return nil
		}

		versionString = data

	case []byte:
		if len(data) == 0 {
			return nil
		}

		versionString = string(data)
	default:
		return fmt.Errorf("Scan: unable to scan type %T into semantic version", data)
	}

	vx, err := Parse(versionString)
	if err != nil {
		return nil
	}

	if v == nil {
		v = vx
	}

	*v = *vx

	return nil
}

// Value implements driver.Valuer
func (v *Version) Value() (driver.Value, error) {
	if v == nil {
		return nil, fmt.Errorf("Value: cannot take value of nil semantic version")
	}
	return v.String(), nil
}

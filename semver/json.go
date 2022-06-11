package semver

import (
	"bytes"
	"encoding/json"
	"strconv"
)

// MarshalJSON implements json.Marshaler
func (v *Version) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.String())
}

// UnmarshalJSON implements json.Unmarshaler
func (v *Version) UnmarshalJSON(x []byte) error {

	// "By convention, to approximate the behavior of Unmarshal itself, Unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op."
	// - https://golang.org/pkg/encoding/json/#Unmarshaler
	if bytes.Equal(x, []byte("null")) {
		return nil
	}

	inx, err := strconv.Unquote(string(x))
	if err != nil {
		return err
	}
	vx, err := Parse(inx)
	if err != nil {
		return err
	}
	*v = *vx
	return nil
}

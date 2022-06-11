package semver

import (
	"encoding/json"
	"testing"
)

type structWithVersion struct {
	A *Version
}

func TestVersion_MarshalJSON(t *testing.T) {
	version, _ := Parse("5.3.6-hello+world.4")
	x := structWithVersion{A: version}

	o, err := json.Marshal(x)

	got := string(o)
	want := `{"A":"5.3.6-hello+world.4"}`

	if err != nil {
		t.Fatalf("(*Version).MarshalJSON() = %v, want <nil>", err)
	} else if got != want {
		t.Fatalf("(*Version).MarshalJSON() = %v, want %v", got, want)
	}
}

func TestVersion_UnmarshalJSON(t *testing.T) {

	t.Run(`unmarshal with []byte("null") as no-op`, func (t *testing.T) {
		// "By convention, to approximate the behavior of Unmarshal itself, Unmarshalers implement UnmarshalJSON([]byte("null")) as a no-op."
		// - https://golang.org/pkg/encoding/json/#Unmarshaler
		var x *Version
		if err := x.UnmarshalJSON([]byte("null")); err != nil {
			t.Fatalf("unmarshal with []byte(\"null\"): (*Version).UnmarshalJSON() = %v, want <nil>", err)
		}
	})

	t.Run("", func(t *testing.T) {
		version, _ := Parse("5.3.6-hello+world.4")
		x := structWithVersion{}

		err := json.Unmarshal([]byte(`{"A":"5.3.6-hello+world.4"}`), &x)

		if err != nil {
			t.Fatalf("(*Version).UnmarshalJSON() = %v, want <nil>", err)
		} else if x.A.CompareTo(version) != 0 {
			t.Fatalf("(*Version).UnmarshalJSON(): output not equal to input")
		}
	})
}
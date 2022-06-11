package semver

import "testing"

func TestVersion_Scan(t *testing.T) {
	
	testVersionString := "5.3.6-hello.w45.13+world.775.r5"
	testVersion, err := Parse(testVersionString)
	if err != nil {
		t.Fatalf("unable to parse %s - %s", testVersionString, err)
	}
	
	{
		t.Run("string", func(t *testing.T) {
			v := new(Version)
			if err := v.Scan(testVersionString); err != nil {
				t.Fatalf("(*Version).Scan() = %s, want = <nil>", err)
			}
			if testVersion.CompareTo(v) != 0 {
				t.Fatal("comparison between test version and parsed version is not equal")
			}
		})
	}

	{
		t.Run("empty string", func(t *testing.T) {
			v := new(Version)
			if err := v.Scan(""); err != nil {
				t.Fatalf("(*Version).Scan() = %s, want = <nil>", err)
			}
		})
	}

	{
		t.Run("bytes", func(t *testing.T) {
			v := new(Version)
			if err := v.Scan([]byte(testVersionString)); err != nil {
				t.Fatalf("(*Version).Scan() = %s, want = <nil>", err)
			}
			if testVersion.CompareTo(v) != 0 {
				t.Fatal("comparison between test version and parsed version is not equal")
			}
		})
	}

	{
		t.Run("empty bytes", func(t *testing.T) {
			v := new(Version)
			if err := v.Scan([]byte{}); err != nil {
				t.Fatalf("(*Version).Scan() = %s, want = <nil>", err)
			}
		})
	}

	{
		t.Run("unknown type", func(t *testing.T) {
			v := new(Version)
			if err := v.Scan(15); err == nil {
				t.Fatal("(*Version).Scan() = <nil>, want an error")
			}
		})
	}

}
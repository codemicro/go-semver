package semver

import (
	"reflect"
	"sort"
	"testing"
)

type versionStringTest struct {
	name   string
	fields *Version
	want   string
}

func x(y string) versionStringTest {
	z := mkv(y)
	return versionStringTest{fields: z, want: y}
}

func TestVersion_String(t *testing.T) {
	tests := []versionStringTest{
		x("1.0.0"),
		x("0.0.0-alpha+kljasdfkh.sdfkfkjl"),
		x("1.2.3-123.hello.aaa"),
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Version{
				Major:      tt.fields.Major,
				Minor:      tt.fields.Minor,
				Patch:      tt.fields.Patch,
				Prerelease: tt.fields.Prerelease,
				Build:      tt.fields.Build,
				Stable:     tt.fields.Stable,
			}
			if got := v.String(); got != tt.want {
				t.Errorf("Version.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSliceSort(t *testing.T) {
	// The Slice type implements sort.Interface - this tests that implementation

	parsed, _ := ParseMultiple([]string{"1.6.3", "2.6.2", "0.3.1", "1.6.3-alpha+shldsfkjh"})
	sort.Sort(parsed)
	var sortedVersions []string
	for _, x := range parsed {
		sortedVersions = append(sortedVersions, x.String())
	}

	expectedVersions := []string{"0.3.1", "1.6.3-alpha+shldsfkjh", "1.6.3", "2.6.2"}
	if !reflect.DeepEqual(expectedVersions, sortedVersions) {
		t.Errorf("Sorted slice is %v, want %v", sortedVersions, expectedVersions)
	}
}

func TestSlice_Less(t *testing.T) {
	// Less must describe a transitive ordering:

	var (
		i = 0
		j = 1
		k = 2
	)

	//  - if both Less(i, j) and Less(j, k) are true, then Less(i, k) must be true as well.
	{
		parsed, _ := ParseMultiple([]string{"1.0.0", "2.0.0", "3.0.0"})
		if !parsed.Less(i, j) {
			t.Fatalf("Slice.Less(i, j) = false, want true\n(i = %s, j = %s)", parsed[i].String(), parsed[j].String())
		}
		if !parsed.Less(j, k) {
			t.Fatalf("Slice.Less(j, k) = false, want true\n(j = %s, k = %s)", parsed[j].String(), parsed[k].String())
		}
		if !parsed.Less(i, k) {
			t.Fatalf("Slice.Less(i, k) = false, want true\n(i = %s, k = %s)", parsed[i].String(), parsed[k].String())
		}
	}

	//  - if both Less(i, j) and Less(j, k) are false, then Less(i, k) must be false as well.
	{
		parsed, _ := ParseMultiple([]string{"1.0.0", "1.0.0", "1.0.0"})
		if parsed.Less(i, j) {
			t.Fatalf("Slice.Less(i, j) = true, want false\n(i = %s, j = %s)", parsed[i].String(), parsed[j].String())
		}
		if parsed.Less(j, k) {
			t.Fatalf("Slice.Less(j, k) = true, want false\n(j = %s, k = %s)", parsed[j].String(), parsed[k].String())
		}
		if parsed.Less(i, k) {
			t.Fatalf("Slice.Less(i, k) = true, want false\n(i = %s, k = %s)", parsed[i].String(), parsed[k].String())
		}
	}
}

func TestSlice_Len(t *testing.T) {
	parsed, _ := ParseMultiple([]string{"1.0.0", "1.0.0", "1.0.0"})
	l := len(parsed)
	f := parsed.Len()
	if l != f {
		t.Fatalf("Slice.Len() = %d, len(parsed) = %d; are different, should match", f, l)
	}
}

func TestSlice_Swap(t *testing.T) {
	parsed, _ := ParseMultiple([]string{"1.0.0", "1.0.0"})
	var (
		i = 0
		j = 1

		ix = parsed[i]
		jx = parsed[j]
	)

	parsed.Swap(i, j)

	if parsed[i] != jx || parsed[j] != ix {
		t.Fatalf("Slice.Swap() is not correctly swapping values")
	}
}
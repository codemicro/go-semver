package semver

import (
	"reflect"
	"sort"
	"testing"
)

type filterTests []struct {
	name    string
	args    filterTestArgs
	want    Slice
	wantErr bool
}

type filterTestArgs struct {
	filter  string
	options Slice
}

func (tests filterTests) Run(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Filter(tt.args.filter, tt.args.options)
			if (err != nil) != tt.wantErr {
				t.Errorf("Filter() error = %v, wantErr %v (got %v)", err, tt.wantErr, got)
				return
			}
			if !reflect.DeepEqual(got, tt.want) && len(got) != 0 && len(tt.want) != 0 {
				t.Errorf("Filter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func ft(filter string) filterTestArgs {
	return filterTestArgs{
		filter:  filter,
		options: mustParseMultiple("0.0.1", "0.0.2", "0.0.3", "0.1.0", "0.2.0", "0.2.1", "0.2.2", "0.3.0", "0.3.1", "0.3.2", "0.4.0", "0.4.1", "0.4.2", "0.5.0-rc.1", "0.5.0", "0.5.1", "0.5.2", "0.6.0", "0.6.1", "0.7.0", "0.8.0", "0.8.1", "0.8.2", "0.9.0", "0.9.1", "0.9.2", "0.10.0", "1.0.0-rc.1", "1.0.0-rc.2", "1.0.0-rc.3", "1.0.0", "1.0.1", "1.1.0", "1.1.1", "1.2.0", "1.2.1", "1.3.0", "1.3.1", "2.0.0", "2.1.0", "2.2.0", "2.2.1", "2.3.0", "2.4.0", "2.4.1", "3.0.0", "3.0.1", "3.1.0", "3.2.0", "3.3.0", "3.3.1", "3.4.0", "3.5.0", "3.6.0", "1.0.2", "3.7.0", "2.4.2", "3.8.0", "3.9.0", "3.9.1", "3.9.2", "3.9.3", "3.10.0", "3.10.1", "4.0.0", "4.0.1", "4.1.0", "4.2.0", "4.2.1", "4.3.0", "4.4.0", "4.5.0", "4.5.1", "4.6.0", "4.6.1", "4.7.0", "4.8.0", "4.8.1", "4.8.2", "4.9.0", "4.10.0", "4.11.0", "4.11.1", "4.11.2", "4.12.0", "4.13.0", "4.13.1", "4.14.0", "4.14.1", "4.14.2", "4.15.0", "4.16.0", "4.16.1", "4.16.2", "4.16.3", "4.16.4", "4.16.5", "4.16.6", "4.17.0", "4.17.1", "4.17.2", "4.17.3", "4.17.4", "4.17.5", "4.17.9", "4.17.10", "4.17.11", "4.17.12", "4.17.13", "4.17.14", "4.17.15", "4.17.16", "4.17.17", "4.17.18", "4.17.19", "4.17.20", "4.17.21"),
	}
}

func mustParseMultiple(vers ...string) Slice {
	x, err := ParseMultiple(vers)
	if err != nil {
		panic(err)
	}
	sort.Sort(x)
	return x
}

func TestFilterByNormal(t *testing.T) {
	filterTests{
		{args: ft("~2.2.0"), want: mustParseMultiple("2.2.0", "2.2.1")},
		{args: ft("~10.0.0"), want: mustParseMultiple()},
		{args: ft("~0.5.0"), want: mustParseMultiple("0.5.0", "0.5.1", "0.5.2")},
		{args: ft("~0.5.0-rc.1"), want: mustParseMultiple("0.5.0-rc.1", "0.5.0", "0.5.1", "0.5.2")},

		{args: ft("^2.2.1"), want: mustParseMultiple("2.2.1", "2.3.0", "2.4.0", "2.4.1", "2.4.2")},
		{args: ft("^2.0.0"), want: mustParseMultiple("2.0.0", "2.1.0", "2.2.0", "2.2.1", "2.3.0", "2.4.0", "2.4.1", "2.4.2")},
		{args: ft("^0.1.0"), want: mustParseMultiple("0.1.0")},
		{args: ft("^1.0.0"), want: mustParseMultiple("1.0.0", "1.0.1", "1.1.0", "1.1.1", "1.2.0", "1.2.1", "1.3.0", "1.3.1", "1.0.2")},
		{args: ft("^0.0.1"), want: mustParseMultiple("0.0.1")},

		{args: ft("1.0.0-rc.1"), want: mustParseMultiple("1.0.0-rc.1")},

		{args: ft(">4.11.1"), want: mustParseMultiple("4.11.2", "4.12.0", "4.13.0", "4.13.1", "4.14.0", "4.14.1", "4.14.2", "4.15.0", "4.16.0", "4.16.1", "4.16.2", "4.16.3", "4.16.4", "4.16.5", "4.16.6", "4.17.0", "4.17.1", "4.17.2", "4.17.3", "4.17.4", "4.17.5", "4.17.9", "4.17.10", "4.17.11", "4.17.12", "4.17.13", "4.17.14", "4.17.15", "4.17.16", "4.17.17", "4.17.18", "4.17.19", "4.17.20", "4.17.21")},
		{args: ft(">=4.11.1"), want: mustParseMultiple("4.11.1", "4.11.2", "4.12.0", "4.13.0", "4.13.1", "4.14.0", "4.14.1", "4.14.2", "4.15.0", "4.16.0", "4.16.1", "4.16.2", "4.16.3", "4.16.4", "4.16.5", "4.16.6", "4.17.0", "4.17.1", "4.17.2", "4.17.3", "4.17.4", "4.17.5", "4.17.9", "4.17.10", "4.17.11", "4.17.12", "4.17.13", "4.17.14", "4.17.15", "4.17.16", "4.17.17", "4.17.18", "4.17.19", "4.17.20", "4.17.21")},

		{args: ft("<2.4.1"), want: mustParseMultiple("0.0.1", "0.0.2", "0.0.3", "0.1.0", "0.2.0", "0.2.1", "0.2.2", "0.3.0", "0.3.1", "0.3.2", "0.4.0", "0.4.1", "0.4.2", "0.5.0", "0.5.1", "0.5.2", "0.6.0", "0.6.1", "0.7.0", "0.8.0", "0.8.1", "0.8.2", "0.9.0", "0.9.1", "0.9.2", "0.10.0", "1.0.0", "1.0.1", "1.1.0", "1.1.1", "1.2.0", "1.2.1", "1.3.0", "1.3.1", "2.0.0", "2.1.0", "2.2.0", "2.2.1", "2.3.0", "2.4.0", "1.0.2")},
		{args: ft("<=2.4.1"), want: mustParseMultiple("0.0.1", "0.0.2", "0.0.3", "2.4.1", "0.1.0", "0.2.0", "0.2.1", "0.2.2", "0.3.0", "0.3.1", "0.3.2", "0.4.0", "0.4.1", "0.4.2", "0.5.0", "0.5.1", "0.5.2", "0.6.0", "0.6.1", "0.7.0", "0.8.0", "0.8.1", "0.8.2", "0.9.0", "0.9.1", "0.9.2", "0.10.0", "1.0.0", "1.0.1", "1.1.0", "1.1.1", "1.2.0", "1.2.1", "1.3.0", "1.3.1", "2.0.0", "2.1.0", "2.2.0", "2.2.1", "2.3.0", "2.4.0", "1.0.2")},

		{args: ft("=0.0.1"), want: mustParseMultiple("0.0.1")},
		{args: ft("0.0.1"), want: mustParseMultiple("0.0.1")},

		// test AND and OR
		{args: ft("0.0.1 0.0.3"), want: Slice{}},
		{args: ft("0.0.1 || 0.0.3"), want: mustParseMultiple("0.0.1", "0.0.3")},
	}.Run(t)
}

func TestFilterByAbnormal(t *testing.T) {
	filterTests{
		{name: "Invalid prefix", args: ft("z1.0.0"), wantErr: true},
		{name: "Empty filter", args: ft(""), wantErr: true},
		{name: "Double space", args: ft("0.0.0  0.0.0"), wantErr: true},
		{name: "Prerelease and caret", args: ft("^1.0.0-abcdefg"), wantErr: true},
	}.Run(t)
}

//func Test_allowUnstable(t *testing.T) {
//	tests := []struct {
//		name string
//		args *Version
//		want bool
//	}{
//		// If the major version of v is 0, this function returns true
//		{name: "Major version zero", args: mkv("0.1.2"), want: true},
//		// If the major version of v is not 0 and v has 1 or more prerelease identifiers, this function returns true
//		{name: "Has prerelease identifiers", args: mkv("0.1.2-potato-cake"), want: true},
//		// Else, the function returns false
//		{name: "Normal", args: mkv("6.2.0")},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			if got := allowUnstable(tt.args); got != tt.want {
//				t.Errorf("allowUnstable() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

# go-semver

*Golang library to parse and compare [v2.0.0 semantic version numbers](https://semver.org/spec/v2.0.0.html).*

---

I built this package for my A-Level Computer Science NEA and never ended up using it. It's a thoroughly-tested library for parsing and comparing semantic versions.

This library is complicant with version 2.0.0 of the Semantic Versioning specification, has no external dependencies and is JSON and SQL ready.

## Installing

```
go get -u github.com/codemicro/go-semver/semver
```

## Usage

### Parse

```go
v, err := semver.Parse("1.6.3-alpha")
if err != nil {
	// handle err
}
// v == &semver.Version{Major:1, Minor:6, Patch:3, Prerelease:[]string{"alpha"}, Build:[]string(nil), Stable:true}

vs, _ := semver.ParseMultiple([]string{"1.0.0", "1.1.0"})
```

### Validate

```go
var valid bool
if _, err := semver.Parse("1.0.0"); err == nil {
    valid = true
}
```

### Compare

```go
a, _ := semver.Parse("1.0.0")
b, _ := semver.Parse("1.0.1")

n := a.CompareTo(b) // n == 1 means a greater than b, n == 0 means a equal to b, n == -1 means a less than b
```

### Filtering

```go
vers, _ := semver.ParseMultiple([]string{"1.0.0", "2.0.0", "2.1.0", "3.0.0"})

v, err := semver.Filter("^2.0.0", vers)   
if err != nil {
// handle err
}
// v == [2.0.0, 2.1.0]
```

#### Specifying version ranges

* `^` - include everything greater than or equal to the stated version that doesn't increment the first non-zero item of the version core
  * eg `^2.1.0` includes version `2.1.0` and any newer `2.x.x` versions
  * eg `^0.3.0` will match only versions `0.3.0` and any newer `0.3.x` versions
  * For example, `^2.2.1` can be expanded out as `>=2.2.1 <3.0.0`
* `~` - include everything greater than or equal to the stated version in the current minor range
  * eg `~2.2.0` will match version `2.2.0` and any newer `2.2.x` but not `2.3.x`
* `>` `<` `=` `>=` `<=` for version comparisons - specify a range of versions
  * eg `>2.1.0` matches anything greater than `2.1.0`
* `||` - include multiple sets of version specifiers
  * eg `^2.0.0 <2.2.0 || > 2.3.0` matches versions that satisfy both `^2.0.0 <2.2.0` and `>2.3.0`

Version numbers must be in their complete form, for example `2` will not work - it must be `2.0.0`.

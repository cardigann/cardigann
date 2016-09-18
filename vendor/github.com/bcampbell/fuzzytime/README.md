# fuzzytime

A date/time parsing package for Go.

Documentation:
[![GoDoc](https://godoc.org/github.com/bcampbell/fuzzytime?status.png)](https://godoc.org/github.com/bcampbell/fuzzytime)

Example:
```
func ExampleExtract() {

inputs := []string{
"Wed Apr 16 17:32:51 NZST 2014",
"2010-02-01T13:14:43Z", // an iso 8601 form
"no date or time info here",
"Published on March 10th, 1999 by Brian Credability",
"2:51pm",
}

for _, inp := range inputs {
dt := Extract(inp)
fmt.Println(dt.ISOFormat())
}

// Output:
// 2014-04-16T17:32:51+12:00
// 2010-02-01T13:14:43Z
//
// 1999-03-10
// T14:51
}
```



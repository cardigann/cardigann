/*

Package fuzzytime helps with the parsing and representation of dates and times.

Fuzzytime defines types (Date, Time, DateTime) which have optional fields.
So you can represent a date with a year and month, but no day.


A quick parsing example:

    package main

    import (
        "fmt"
        "github.com/bcampbell/fuzzytime"
    )

    func main() {

        inputs := []string{
            "Wed Apr 16 17:32:51 NZST 2014",
            "2010-02-01T13:14:43Z", // an iso 8601 form
            "no date or time info here",
            "Published on March 10th, 1999 by Brian Credability",
            "2:51pm",
            "April 2004",
        }

        for _, s := range inputs {
            dt, _, _ := fuzzytime.Extract(s)
            fmt.Println(dt.ISOFormat())
        }


    }

This should output:

    2014-04-16T17:32:51+12:00
    2010-02-01T13:14:43Z

    1999-03-10
    T14:51
    2004-04

Timezones, once resolved, are stored as an offset from UTC (in seconds).

Sometimes dates and times are ambiguous and can't be parsed without
extra information (eg "dd/mm/yy" vs "mm/dd/yy"). The default behaviour when
such a data is encountered is for Extract() function to just return an error.
This can be overriden by using a Context struct, which provides
functions to perform the decisions required in otherwise-ambiguous cases.


*/
package fuzzytime

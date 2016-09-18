package fuzzytime

import (
	"fmt"
)

// A Date represents a year/month/day set where any of the three may be
// unset.
// default initialisation (ie Date{}) is a valid but empty Date.
type Date struct {
	year, month, day int // internally, we'll say 0=undefined
}

// Year returns the year (result undefined if field unset)
func (d *Date) Year() int { return d.year }

// Month returns the month (result undefined if field unset)
func (d *Date) Month() int { return d.month }

// Day returns the day (result undefined if field unset)
func (d *Date) Day() int { return d.day }

// SetYear sets the year field
func (d *Date) SetYear(year int) { d.year = year }

// SetMonth sets the month field
func (d *Date) SetMonth(month int) { d.month = month }

// SetDay sets the day field
func (d *Date) SetDay(day int) { d.day = day }

// HasYear returns true if the year is set
func (d *Date) HasYear() bool { return d.year != 0 }

// HasMonth returns true if the month is set
func (d *Date) HasMonth() bool { return d.month != 0 }

// HasDay returns trus if the day is set
func (d *Date) HasDay() bool { return d.day != 0 }

// Equals returns true if dates match. Fields present in one date but
// not the other are considered mismatches.
func (d *Date) Equals(other *Date) bool {
	// TODO: should check if fields are set before comparing
	if d.year == other.year && d.month == other.month && d.day == other.day {
		return true
	}
	return false
}

// Conflicts returns true if date d conflicts with the other date.
// Missing fields are not considered so, for example "2013-01-05"
// doesn't conflict with "2013-01"
func (d *Date) Conflicts(other *Date) bool {
	if d.HasYear() && other.HasYear() && d.Year() != other.Year() {
		return true
	}
	if d.HasMonth() && other.HasMonth() && d.Month() != other.Month() {
		return true
	}
	if d.HasDay() && other.HasDay() && d.Day() != other.Day() {
		return true
	}
	return false
}

// Merge copies all fields set in other into d.
// any fields unset in other are left unchanged in d.
func (d *Date) Merge(other *Date) {
	if other.HasYear() {
		d.SetYear(other.Year())
	}
	if other.HasMonth() {
		d.SetMonth(other.Month())
	}
	if other.HasDay() {
		d.SetDay(other.Day())
	}
}

// Empty tests if date is blank (ie all fields unset)
func (d *Date) Empty() bool {
	if d.HasYear() || d.HasMonth() || d.HasDay() {
		return false
	}
	return true
}

// sane returns true if date isn't obviously bogus
func (d *Date) sane() bool {
	if d.HasMonth() {
		if d.Month() < 1 || d.Month() > 12 {
			return false
		}
	}
	if d.HasDay() {
		// TODO: adjust for month! (and leapyears!)
		if d.Day() < 1 || d.Day() > 31 {
			return false
		}
	}
	return true
}

// String returns "YYYY-MM-DD" with question marks in place of
// any missing values
func (d *Date) String() string {
	var year, month, day = "????", "??", "??"
	if d.HasYear() {
		year = fmt.Sprintf("%04d", d.Year())
	}
	if d.HasMonth() {
		month = fmt.Sprintf("%02d", d.Month())
	}
	if d.HasDay() {
		day = fmt.Sprintf("%02d", d.Day())
	}

	return year + "-" + month + "-" + day
}

// ISOFormat returns "YYYY-MM-DD", "YYYY-MM" or "YYYY" depending on which
// fields are available (or "" if year is missing).
func (d *Date) ISOFormat() string {
	if d.HasYear() {
		if d.HasMonth() {
			if d.HasDay() {
				return fmt.Sprintf("%04d-%02d-%02d", d.Year(), d.Month(), d.Day())
			}
			return fmt.Sprintf("%04d-%02d", d.Year(), d.Month())
		}
		return fmt.Sprintf("%04d", d.Year())
	}
	return ""
}

// NewDate creates a Date with all fields set
func NewDate(y, m, d int) *Date {
	return &Date{y, m, d}
}

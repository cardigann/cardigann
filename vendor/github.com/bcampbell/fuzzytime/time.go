package fuzzytime

import (
	"fmt"
)

const (
	hourFlag   int = 0x01
	minuteFlag int = 0x02
	secondFlag int = 0x04
	tzFlag     int = 0x08
)

// Time represents a set of time fields, any of which may be unset.
// The default initialisation (ie Time{}) produces a Time with all fields unset.
type Time struct {
	set      int // flags to show which fields are set
	hour     int
	minute   int
	second   int
	tzOffset int // offset from UTC, in seconds
}

// Hour returns the hour (result undefined if field unset)
func (t *Time) Hour() int { return t.hour }

// Minute returns the minute (result undefined if field unset)
func (t *Time) Minute() int { return t.minute }

// Second returns the second (result undefined if field unset)
func (t *Time) Second() int { return t.second }

// TZOffset returns the offset from UTC, in seconds. Result
// undefined if field is unset.
func (t *Time) TZOffset() int { return t.tzOffset }

// SetHour sets the hour (0-23)
func (t *Time) SetHour(hour int) { t.hour = hour; t.set |= hourFlag }

// SetMinute sets the Minute field (0-59)
func (t *Time) SetMinute(minute int) { t.minute = minute; t.set |= minuteFlag }

// SetSecond sets the Second field (0-59)
func (t *Time) SetSecond(second int) { t.second = second; t.set |= secondFlag }

// SetTZOffset sets the timezone offset from UTC, in seconds
func (t *Time) SetTZOffset(tzOffset int) { t.tzOffset = tzOffset; t.set |= tzFlag }

// HasHour returns true if the hour is set
func (t *Time) HasHour() bool { return (t.set & hourFlag) != 0 }

// HasMinute returns true if the minute is set
func (t *Time) HasMinute() bool { return (t.set & minuteFlag) != 0 }

// HasSecond returns true if the second is set
func (t *Time) HasSecond() bool { return (t.set & secondFlag) != 0 }

// HasTZOffset returns true if the timezone offset field is set
func (t *Time) HasTZOffset() bool { return (t.set & tzFlag) != 0 }

// Equals returns true if the two times have the same fields set
// and match exactly. Fields present in one time but not the other
// are considered mismatches.
func (t *Time) Equals(other *Time) bool {
	if t.set != other.set {
		return false
	}
	if t.HasHour() && t.hour != other.hour {
		return false
	}
	if t.HasMinute() && t.minute != other.minute {
		return false
	}
	if t.HasSecond() && t.second != other.second {
		return false
	}
	if t.HasTZOffset() && t.tzOffset != other.tzOffset {
		return false
	}
	return true
}

// Conflicts returns true if time t conflicts with the other time.
// Missing fields are not considered so, for example "10:59:01"
// doesn't conflict with "10:59"
func (t *Time) Conflicts(other *Time) bool {
	if t.HasHour() && other.HasHour() && t.Hour() != other.Hour() {
		return true
	}
	if t.HasMinute() && other.HasMinute() && t.Minute() != other.Minute() {
		return true
	}
	if t.HasSecond() && other.HasSecond() && t.Second() != other.Second() {
		return true
	}
	if t.HasTZOffset() && other.HasTZOffset() && t.TZOffset() != other.TZOffset() {
		return true
	}

	return false // no conflict
}

// String returns "hh:mm:ss+tz", with question marks in place of
// any missing values (except for timezone, which will be blank if missing)
func (t *Time) String() string {
	var hour, minute, second, tz = "??", "??", "??", ""
	if t.HasHour() {
		hour = fmt.Sprintf("%02d", t.Hour())
	}
	if t.HasMinute() {
		minute = fmt.Sprintf("%02d", t.Minute())
	}
	if t.HasSecond() {
		second = fmt.Sprintf("%02d", t.Second())
	}
	if t.HasTZOffset() {
		tz = OffsetToTZ(t.TZOffset())
	}
	return hour + ":" + minute + ":" + second + tz
}

// Empty tests if time is blank (ie all fields unset)
func (t *Time) Empty() bool {
	return t.set == 0
}

// ISOFormat returns the most precise possible ISO-formatted time
func (t *Time) ISOFormat() string {
	var out string
	if t.HasHour() {
		if t.HasMinute() {
			if t.HasSecond() {
				out = fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second())
			} else {
				out = fmt.Sprintf("%02d:%02d", t.Hour(), t.Minute())
			}
		} else {
			out = fmt.Sprintf("%02d", t.Hour())
		}
		if t.HasTZOffset() {
			out += OffsetToTZ(t.TZOffset())
		}
	}
	return out
}

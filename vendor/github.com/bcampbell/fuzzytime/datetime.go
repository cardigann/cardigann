package fuzzytime

// DateTime represents a set of fields for date and time, any of which may
// be unset. The default initialisation is a valid empty datetime with no
// fields set.
type DateTime struct {
	Date
	Time
}

// Equals returns true if dates and times match
func (dt *DateTime) Equals(other *DateTime) bool {
	return dt.Date.Equals(&other.Date) && dt.Time.Equals(&other.Time)
}

// String returns "YYYY-MM-DD hh:mm:ss tz" with question marks in place of
// any missing values (except for timezone, which will be blank if missing)
func (dt *DateTime) String() string {
	return dt.Date.String() + " " + dt.Time.String()
}

// Empty tests if datetime is blank (ie all fields unset)
func (dt *DateTime) Empty() bool {
	return dt.Time.Empty() && dt.Date.Empty()
}

// ISOFormat returns the most precise-possible datetime given the available
// data.
// Aims for "YYYY-MM-DDTHH:MM:SS+TZ" but will drop off
// higher-precision components as required eg "YYYY-MM"
func (dt *DateTime) ISOFormat() string {
	if dt.Time.Empty() {
		// just the date.
		return dt.Date.ISOFormat()
	}
	return dt.Date.ISOFormat() + "T" + dt.Time.ISOFormat()
}

// HasFullDate returns true if Year, Month and Day are all set
func (dt *DateTime) HasFullDate() bool {
	return dt.HasYear() && dt.HasMonth() && dt.HasDay()
}

// Conflicts returns true if the two datetimes conflict.
// Note that this is not the same as the two being equal - one
// datetime can be more precise than the other. They are only in
// conflict if they have different values set for the same field.
// eg "2012-01-01T03:34:10" doesn't conflict with "03:34"
func (dt *DateTime) Conflicts(other *DateTime) bool {
	return dt.Time.Conflicts(&other.Time) || dt.Date.Conflicts(&other.Date)
}

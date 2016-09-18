package fuzzytime

// lookup.go contains lookup data (eg month names)

// TODO: look into CLDR - http://cldr.unicode.org/index
//       provides locale-specific names and format patterns

// useful reference for month abbreviations:
// http://library.princeton.edu/departments/tsd/katmandu/reference/months.html

/*
dayLookup = {
    'mon': 'mon', 'monday': 'mon',
    'tue': 'tue', 'tuesday': 'tue',
    'wed': 'wed', 'wednesday': 'wed',
    'thu': 'thu', 'thursday': 'thu',
    'fri': 'fri', 'friday': 'fri',
    'sat': 'sat', 'saturday': 'sat',
    'sun': 'sun', 'sunday': 'sun',
    # es
    'lunes': 'mon',
    'martes': 'tue',
    'miércoles': 'wed',
    'jueves': 'thu',
    'viernes': 'fri',
    'sábado': 'sat',
    'domingo': 'sun',
}
*/

var monthLookup = map[string]int{
	"jan": 1,
	"feb": 2,
	"mar": 3,
	"apr": 4,
	"may": 5,
	"jun": 6,
	"jul": 7,
	"aug": 8,
	"sep": 9,
	"oct": 10,
	"nov": 11,
	"dec": 12,

	"january":  1,
	"february": 2,
	"march":    3,
	"april":    4,
	// "may": 5,
	"june":      6,
	"july":      7,
	"august":    8,
	"september": 9,
	"october":   10,
	"november":  11,
	"december":  12,

	"01": 1,
	"02": 2,
	"03": 3,
	"04": 4,
	"05": 5,
	"06": 6,
	"07": 7,
	"08": 8,
	"09": 9,
	"10": 10,
	"11": 11,
	"12": 12,
	"1":  1,
	"2":  2,
	"3":  3,
	"4":  4,
	"5":  5,
	"6":  6,
	"7":  7,
	"8":  8,
	"9":  9,
	// "10":10,
	// "11":11,
	// "12":12,

	// es - full
	"enero":      1,
	"febrero":    2,
	"marzo":      3,
	"abril":      4,
	"mayo":       5,
	"junio":      6,
	"julio":      7,
	"agosto":     8,
	"septiembre": 9,
	"octubre":    10,
	"noviembre":  11,
	"diciembre":  12,

	// es - abbreviations
	//"enero": 1,
	//	"feb":    2,
	//"marzo": 3,
	"abr": 4,
	//"mayo":  5,
	//	"jun":    6,
	//	"jul":    7,
	//"agosto": 8,
	//	"sept":   9,
	"set": 9,
	//	"oct":    10,
	//	"nov":    11,
	"dic": 12,
}

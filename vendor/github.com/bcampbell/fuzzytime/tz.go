package fuzzytime

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// TZInfo holds info about a timezone offset
type TZInfo struct {
	// Name of the timezone eg "BST", "UTC", "NZDT"
	Name string
	// Offset from UTC, in ISO8601 form [+-]<HH>[:<MM>]
	Offset string
	// Locale contains comma-separated country identifiers
	// to help resolve ambiguities
	Locale string
}

// lookup table of common abbreviations of timezones
// source: http://en.wikipedia.org/wiki/List_of_time_zone_abbreviations
//
// Note that the names are not unambiguous... (eg BST: Britain or Bangladesh?)
var tzTable = map[string][]TZInfo{
	"ACDT": {{"ACDT", "+10:30", ""}}, //Australian Central Daylight Time
	"ACST": {{"ACST", "+09:30", ""}}, //Australian Central Standard Time
	"ACT":  {{"ACT", "+08", ""}},     //ASEAN Common Time
	"ADT":  {{"ADT", "-03", ""}},     //Atlantic Daylight Time
	"AEDT": {{"AEDT", "+11", ""}},    //Australian Eastern Daylight Time
	"AEST": {{"AEST", "+10", ""}},    //Australian Eastern Standard Time
	"AFT":  {{"AFT", "+04:30", ""}},  //Afghanistan Time
	"AKDT": {{"AKDT", "-08", ""}},    //Alaska Daylight Time
	"AKST": {{"AKST", "-09", ""}},    //Alaska Standard Time
	"AMST": {{"AMST", "+05", ""}},    //Armenia Summer Time
	"AMT":  {{"AMT", "+04", ""}},     //Armenia Time
	"ART":  {{"ART", "-03", ""}},     //Argentina Time
	"AST": {{"AST", "+03", "BH,IQ,JO,KW,SA,YE,QA"}, //"Arab Standard Time (Kuwait, Riyadh)"
		{"AST", "-04", "AW,BB,BM,VG,CA,CW,DO,GD,SX,TT,PR,VI"}}, //Atlantic Standard Time (https://en.wikipedia.org/wiki/Atlantic_Standard_Time_Zone)
	"AWDT":  {{"AWDT", "+09", ""}},  //Australian Western Daylight Time
	"AWST":  {{"AWST", "+08", ""}},  //Australian Western Standard Time
	"AZOST": {{"AZOST", "-01", ""}}, //Azores Standard Time
	"AZT":   {{"AZT", "+04", ""}},   //Azerbaijan Time
	"BDT":   {{"BDT", "+08", ""}},   //Brunei Time
	"BIOT":  {{"BIOT", "+06", ""}},  //British Indian Ocean Time
	"BIT":   {{"BIT", "-12", ""}},   //Baker Island Time
	"BOT":   {{"BOT", "-04", ""}},   //Bolivia Time
	"BRT":   {{"BRT", "-03", ""}},   //Brasilia Time
	"BST": {{"BST", "+06", "BD"}, //Bangladesh Standard Time
		{"BST", "+01", "GB"}}, //British Summer Time (British Standard Time from Feb 1968 to Oct 1971)
	"BTT":   {{"BTT", "+06", ""}},      //Bhutan Time
	"CAT":   {{"CAT", "+02", ""}},      //Central Africa Time
	"CCT":   {{"CCT", "+06:30", ""}},   //Cocos Islands Time
	"CDT":   {{"CDT", "-05", ""}},      //Central Daylight Time (North America)
	"CEDT":  {{"CEDT", "+02", ""}},     //Central European Daylight Time
	"CEST":  {{"CEST", "+02", ""}},     //Central European Summer Time (Cf. HAEC)
	"CET":   {{"CET", "+01", ""}},      //Central European Time
	"CHADT": {{"CHADT", "+13:45", ""}}, //Chatham Daylight Time
	"CHAST": {{"CHAST", "+12:45", ""}}, //Chatham Standard Time
	"CIST":  {{"CIST", "-08", ""}},     //Clipperton Island Standard Time
	"CKT":   {{"CKT", "-10", ""}},      //Cook Island Time
	"CLST":  {{"CLST", "-03", ""}},     //Chile Summer Time
	"CLT":   {{"CLT", "-04", ""}},      //Chile Standard Time
	"COST":  {{"COST", "-04", ""}},     //Colombia Summer Time
	"COT":   {{"COT", "-05", ""}},      //Colombia Time
	"CST": {{"CST", "-06", "US,CA,JM,BZ,MX"}, //Central Standard Time (North America)
		{"CST", "+08", "CN,HK,MO,TW"}, //China Standard Time
		{"CST", "+09:30", "AU"}},      //Central Standard Time (Australia)
	"CT":   {{"CT", "+08", ""}},   //China Time
	"CVT":  {{"CVT", "-01", ""}},  //Cape Verde Time
	"CXT":  {{"CXT", "+07", ""}},  //Christmas Island Time
	"CHST": {{"CHST", "+10", ""}}, //Chamorro Standard Time
	"DFT":  {{"DFT", "+01", ""}},  //AIX specific equivalent of Central European Time
	"EAST": {{"EAST", "-06", ""}}, //Easter Island Standard Time
	"EAT":  {{"EAT", "+03", ""}},  //East Africa Time
	"ECT": {{"ECT", "-04", "AI,AG,BB,DM,GD,MS,KN,LC,VC,TT,VG,JM"}, //Eastern Caribbean Time (does not recognise DST)
		{"ECT", "-05", "EC"}}, //Ecuador Time
	"EDT":  {{"EDT", "-04", ""}},  //Eastern Daylight Time (North America)
	"EEDT": {{"EEDT", "+03", ""}}, //Eastern European Daylight Time
	"EEST": {{"EEST", "+03", ""}}, //Eastern European Summer Time
	"EET":  {{"EET", "+02", ""}},  //Eastern European Time
	"EST":  {{"EST", "-05", ""}},  //Eastern Standard Time (North America)
	"FET":  {{"FET", "+03", ""}},  //Further-eastern_European_Time
	"FJT":  {{"FJT", "+12", ""}},  //Fiji Time
	"FKST": {{"FKST", "-03", ""}}, //Falkland Islands Summer Time
	"FKT":  {{"FKT", "-04", ""}},  //Falkland Islands Time
	"GALT": {{"GALT", "-06", ""}}, //Galapagos Time
	"GET":  {{"GET", "+04", ""}},  //Georgia Standard Time
	"GFT":  {{"GFT", "-03", ""}},  //French Guiana Time
	"GILT": {{"GILT", "+12", ""}}, //Gilbert Island Time
	"GIT":  {{"GIT", "-09", ""}},  //Gambier Island Time
	"GMT":  {{"GMT", "Z", ""}},    //Greenwich Mean Time
	"GST": {{"GST", "-02", "GS"}, //South Georgia and the South Sandwich Islands
		{"GST", "+04", "AE,OM"}}, //Gulf Standard Time
	"GYT":  {{"GYT", "-04", ""}},     //Guyana Time
	"HADT": {{"HADT", "-09", ""}},    //Hawaii-Aleutian Daylight Time
	"HAEC": {{"HAEC", "+02", ""}},    //Heure Avancée d'Europe Centrale francised name for CEST
	"HAST": {{"HAST", "-10", ""}},    //Hawaii-Aleutian Standard Time
	"HKT":  {{"HKT", "+08", ""}},     //Hong Kong Time
	"HMT":  {{"HMT", "+05", ""}},     //Heard and McDonald Islands Time
	"HST":  {{"HST", "-10", ""}},     //Hawaii Standard Time
	"ICT":  {{"ICT", "+07", ""}},     //Indochina Time
	"IDT":  {{"IDT", "+03", ""}},     //Israeli Daylight Time
	"IRKT": {{"IRKT", "+08", ""}},    //Irkutsk Time
	"IRST": {{"IRST", "+03:30", ""}}, //Iran Standard Time
	"IST": {{"IST", "+05:30", "IN,LK"}, //Indian Standard Time
		{"IST", "+01", "IE"},  //Irish Summer Time
		{"IST", "+02", "IL"}}, //Israel Standard Time
	"JST":  {{"JST", "+09", ""}},     //Japan Standard Time
	"KRAT": {{"KRAT", "+07", ""}},    //Krasnoyarsk Time
	"KST":  {{"KST", "+09", ""}},     //Korea Standard Time
	"LHST": {{"LHST", "+10:30", ""}}, //Lord Howe Standard Time
	"LINT": {{"LINT", "+14", ""}},    //Line Islands Time
	"MAGT": {{"MAGT", "+11", ""}},    //Magadan Time
	"MDT":  {{"MDT", "-06", ""}},     //Mountain Daylight Time (North America)
	"MET":  {{"MET", "+01", ""}},     //Middle European Time Same zone as CET
	"MEST": {{"MEST", "+02", ""}},    //Middle European Saving Time Same zone as CEST
	"MIT":  {{"MIT", "-09:30", ""}},  //Marquesas Islands Time
	"MSK":  {{"MSK", "+04", ""}},     //Moscow Time
	"MST": {{"MST", "+08", "MY"}, //Malaysian Standard Time
		{"MST", "-07", "CA,MX,US"}, //Mountain Standard Time (North America)
		{"MST", "+06:30", "MM"}},   //Myanmar Standard Time
	"MUT":  {{"MUT", "+04", ""}},    //Mauritius Time
	"MYT":  {{"MYT", "+08", ""}},    //Malaysia Time
	"NDT":  {{"NDT", "-02:30", ""}}, //Newfoundland Daylight Time
	"NFT":  {{"NFT", "+11:30", ""}}, //Norfolk Time[1]
	"NPT":  {{"NPT", "+05:45", ""}}, //Nepal Time
	"NST":  {{"NST", "-03:30", ""}}, //Newfoundland Standard Time
	"NT":   {{"NT", "-03:30", ""}},  //Newfoundland Time
	"NZDT": {{"NZDT", "+13", ""}},   //New Zealand Daylight Time
	"NZST": {{"NZST", "+12", ""}},   //New Zealand Standard Time
	"OMST": {{"OMST", "+06", ""}},   //Omsk Time
	"PDT":  {{"PDT", "-07", ""}},    //Pacific Daylight Time (North America)
	"PETT": {{"PETT", "+12", ""}},   //Kamchatka Time
	"PHOT": {{"PHOT", "+13", ""}},   //Phoenix Island Time
	"PKT":  {{"PKT", "+05", ""}},    //Pakistan Standard Time
	"PST": {{"PST", "-08", "CA,MX,US"}, //Pacific Standard Time (North America)
		{"PST", "+08", "PH"}}, //Philippine Standard Time
	"RET":  {{"RET", "+04", ""}},    //Réunion Time
	"SAMT": {{"SAMT", "+04", ""}},   //Samara Time
	"SAST": {{"SAST", "+02", ""}},   //South African Standard Time
	"SBT":  {{"SBT", "+11", ""}},    //Solomon Islands Time
	"SCT":  {{"SCT", "+04", ""}},    //Seychelles Time
	"SGT":  {{"SGT", "+08", ""}},    //Singapore Time
	"SLT":  {{"SLT", "+05:30", ""}}, //Sri Lanka Time
	"SST": {{"SST", "-11", "WS,AS"}, //Samoa Standard Time
		{"SST", "+08", "SG"}}, //Singapore Standard Time
	"TAHT": {{"TAHT", "-10", ""}},   //Tahiti Time
	"THA":  {{"THA", "+07", ""}},    //Thailand Standard Time
	"UTC":  {{"UTC", "Z", ""}},      //Coordinated Universal Time
	"UYST": {{"UYST", "-02", ""}},   //Uruguay Summer Time
	"UYT":  {{"UYT", "-03", ""}},    //Uruguay Standard Time
	"VET":  {{"VET", "-04:30", ""}}, //Venezuelan Standard Time
	"VLAT": {{"VLAT", "+10", ""}},   //Vladivostok Time
	"WAT":  {{"WAT", "+01", ""}},    //West Africa Time
	"WEDT": {{"WEDT", "+01", ""}},   //Western European Daylight Time
	"WEST": {{"WEST", "+01", ""}},   //Western European Summer Time
	"WET":  {{"WET", "Z", ""}},      //Western European Time
	"WST":  {{"WST", "+08", ""}},    //Western Standard Time
	"YAKT": {{"YAKT", "+09", ""}},   //Yakutsk Time
	"YEKT": {{"YEKT", "+05", ""}},   //Yekaterinburg Time
}

// OffsetToTZ converts an offset in seconds from UTC into an ISO8601-style
// offset (like "+HH:MM")
func OffsetToTZ(secs int) string {
	if secs == 0 {
		return "Z"
	}
	sign := '+'
	if secs < 0 {
		sign = '-'
		secs = -secs
	}
	mins := (secs / 60) % 60
	hours := secs / (60 * 60)
	return fmt.Sprintf("%c%02d:%02d", sign, hours, mins)
}

var isoTZRE = regexp.MustCompile(`([-+])(\d{2})(?:[:]?(\d{2}))?`)

// TZToOffset parses an ISO8601 timezone offset ("Z", "[+-]HH" "[+-]HH[:]?MM" etc...)
// and returns the offset from UTC in seconds
func TZToOffset(s string) (int, error) {
	if s == "Z" {
		return 0, nil
	}
	m := isoTZRE.FindStringSubmatch(s)
	if m == nil {
		return 0, errors.New("bad timezone")
	}

	var hours, mins int
	hours, err := strconv.Atoi(m[2])
	if err != nil {
		return 0, err
	}
	// has mins?
	if m[3] != "" {
		mins, err = strconv.Atoi(m[3])
		if err != nil {
			return 0, err
		}
	}

	switch m[1] {
	case "+":
		return 60*60*hours + 60*mins, nil
	case "-":
		return -((60 * 60 * hours) + (60 * mins)), nil
	}

	return 0, errors.New("bad timezone")
}

// FindTimeZone returns timezones with the matching name (eg "BST")
// Some timezone names are ambiguous (eg "BST"), so all the matching
// ones will be returned. It's up to the caller to disambiguate them.
// To aid in this, ambiguous timezones include a list of country
// locale codes ("US", "AU" etc) in where they are used.
func FindTimeZone(name string) []TZInfo {
	name = strings.ToUpper(name)
	matches, got := tzTable[name]
	if got {
		return matches
	}
	return []TZInfo{}
}

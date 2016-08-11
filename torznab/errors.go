package torznab

import (
	"encoding/xml"
	"net/http"
)

type err struct {
	Code        int
	Description string
}

func (e err) Error() string {
	return e.Description
}

var (
	ErrIncorrectUserCreds     = err{100, "Incorrect user credentials"}
	ErrAccountSuspended       = err{101, "Account suspended"}
	ErrInsufficientPrivs      = err{102, "Insufficient privileges/not authorized"}
	ErrRegistrationDenied     = err{103, "Registration denied"}
	ErrRegistrationsAreClosed = err{104, "Registrations are closed"}
	ErrEmailAddressTaken      = err{105, "Invalid registration (Email Address Taken)"}
	ErrEmailAddressBadFormat  = err{106, "Invalid registration (Email Address Bad Format)"}
	ErrRegistrationFailed     = err{107, "Registration Failed (Data error)"}
	ErrMissingParameter       = err{200, "Missing parameter"}
	ErrIncorrectParameter     = err{201, "Incorrect parameter"}
	ErrNoSuchFunction         = err{202, "No such function. (Function not defined in this specification)."}
	ErrFunctionNotAvailable   = err{203, "Function not available. (Optional function is not implemented)."}
	ErrNoSuchItem             = err{300, "No such item."}
	ErrItemAlreadyExists      = err{300, "Item already exists."}
	ErrUnknownError           = err{900, "Unknown error"}
	ErrAPIDisabled            = err{910, "API Disabled"}
)

func Error(w http.ResponseWriter, description string, err err) {
	var resp = struct {
		XMLName     struct{} `xml:"error"`
		Code        int      `xml:"code"`
		Description string   `xml:"description"`
	}{
		Code:        err.Code,
		Description: description,
	}
	x, mErr := xml.MarshalIndent(resp, "", "  ")
	if mErr != nil {
		http.Error(w, mErr.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/xml")
	w.WriteHeader(http.StatusBadGateway)
	w.Write(x)
}

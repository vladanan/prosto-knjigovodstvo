package vet

import (
	"net/http"
	"regexp"
	"strings"

	"github.com/vladanan/prosto/src/controllers/clr"
)

var l = clr.GetErrorLogger()

func ValidateEmailAddress(email string) error {
	// validacija za UPIS NOVOG KORISNIKA a-zA-Z09 .,+-*:!?() min char 8 max 32 ISTO URADITI I NA FE UZ ARGUMENTS I JS

	// https://gobyexample.com/regular-expressions
	// https://pkg.go.dev/regexp
	// https://regex101.com/

	// matched, err := regexp.MatchString(`[^a-zA-Z\d]`, recordData.Email)
	re, err := regexp.Compile(`[^a-zA-Z\d@\.\-_]`)
	if err != nil {
		return l(nil, 8, err)
	}

	// var m bool
	// if m = strings.ContainsAny(email, "@."); !m {
	// 	return l(r, 8, clr.NewAPIError(http.StatusBadRequest, "malformed request syntax for mail address m1"))
	// }
	// // napraviti funkciju za validaciju i sanitaciju za mejl itd.
	// if m = strings.ContainsAny(email, ",:;!?&#$%=\"'*+()[]<>{}/\\"); m {
	// 	return l(r, 8, clr.NewAPIError(http.StatusBadRequest, "malformed request syntax for mail address m2"))
	// }

	if len(email) < 8 ||
		len(email) > 32 ||
		!strings.ContainsAny(email, "@") ||
		!strings.ContainsAny(email, ".") ||
		re.MatchString(email) {

		return clr.NewAPIError(http.StatusNotAcceptable, "Malformed_email_address")

	}

	return nil
}

func ValidateUserName(userName string) error {

	re, err := regexp.Compile(`[^a-zA-Z\d@\.\-_]`)
	if err != nil {
		return l(nil, 7, err)
	}

	if len(userName) < 8 ||
		len(userName) > 32 ||
		re.MatchString(userName) {
		return clr.NewAPIError(http.StatusNotAcceptable, "Malformed_user_name")
	}

	return nil
}

func ValidatePassword(pass string) error {

	re, err := regexp.Compile(`[^a-zA-Z\d@\.\-_]`)
	if err != nil {
		return l(nil, 7, err)
	}

	if len(pass) < 8 ||
		len(pass) > 32 ||
		re.MatchString(pass) {
		return clr.NewAPIError(http.StatusNotAcceptable, "Malformed_password")
	}

	return nil
}

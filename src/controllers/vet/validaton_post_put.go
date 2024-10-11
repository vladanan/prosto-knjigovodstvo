package vet

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/schema"

	"github.com/vladanan/prosto/src/controllers/clr"
	"github.com/vladanan/prosto/src/controllers/utils"
	"github.com/vladanan/prosto/src/models"
)

// Set a Decoder instance as a package global, because it caches
// meta-data about structs, and an instance can be shared safely.
var decoder = schema.NewDecoder()

func ValidateFormPostAndPutData(tableDb string, r *http.Request) (any, error) {
	switch tableDb {
	case "g_pitanja_c_testovi":
		return nil, nil
	case "mi_users":
		reqData, err := ValidateSignUpData(r)
		if err != nil {
			return nil, err
		}
		return reqData, nil
	case "g_user_blog":
		return nil, nil
	case "v_settings":
		return nil, nil
	default:
		return nil, l(nil, 8, fmt.Errorf("malformed table name"))
	}
}

func ValidateSignUpData(r *http.Request) (models.User, error) {
	u := models.User{}

	if err := r.ParseForm(); err != nil {
		log.Println(err)
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}

	// za html elemente form-a gleda se name a ne id tag i pisu se velikim pocetnim slovima
	var f models.SignUpFormData

	// MOZE SE DESITI DA PUKNE JER JE NEKO POLJE U FORM NEISPRAVNO I DOVEDE DO GRESKE U CHECKERR KADA POKUSAVA DA PRIKAZE err.Error()
	if err := decoder.Decode(&f, r.PostForm); err != nil {
		log.Println(err)
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}

	// kada se stampa ceo f onda redirekcija iz templ stderr u neki drugi terminal hoce da brljavi nesto i ispisuje kod i delimican struct
	// log.Println("parsovan i dekodiran post form:", f, f.User_email1)

	// validacija forma za ispravnost ssr bycrypt key
	if err := utils.Check64UrlKey(f.Bk_, "", 1); err != nil {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Session_expired")
	}

	// validacija forma za dva ista i svi prisutni
	if f.User_email1 != f.User_email2 {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Email_not_same")
	}

	if f.Password1 != f.Password2 {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Password_not_same")
	}

	if f.User_email1 == "" || f.User_email2 == "" || f.Name == "" || f.Password1 == "" || f.Password2 == "" || f.Bk_ == "" {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}

	//bot pokusava da salje masu formulara i daje mu se greska kao za nepotpun form da se ne otkrije da je popunjavanje tog input polja greska
	if f.Fk_ != "" {
		l(r, 7, fmt.Errorf("bot attempt with sign up input field: "+f.Fk_))
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}
	if f.Confirm == "" {
		l(r, 7, fmt.Errorf("bot attempt with empty sign up checkbox"))
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}

	// provera da li ima isti takav email u db
	if _, err := models.CheckEmail(f.User_email1, r); err != nil {
		return models.User{}, err
	}

	if err := ValidateEmailAddress(f.User_email1); err != nil {
		return models.User{}, err
	}

	// provera da li ima isti takav user name u db
	if err := models.CheckUserName(f.Name, r); err != nil {
		return u, err
	}

	if err := ValidateUserName(f.Name); err != nil {
		return models.User{}, err
	}

	if err := ValidatePassword(f.Password1); err != nil {
		log.Println("provera lozinke")
		return models.User{}, err
	}

	u.Email = f.User_email1
	u.User_name = f.Name
	u.Hash_lozinka = f.Password1

	return u, nil
}

func ValidateSignInData(r *http.Request, f models.SignInFormData) error {

	// log.Println("parsovan i dekodiran post form:", f, f.User_email1)

	// validacija forma za ispravnost ssr bycrypt key
	if err := utils.Check64UrlKey(f.Bk_, "", 1); err != nil {
		return clr.NewAPIError(http.StatusNotAcceptable, "Session_expired")
	}

	// validacija forma za svi prisutni
	if f.User_email1 == "" || f.Password1 == "" || f.Bk_ == "" {
		return clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}

	//bot pokusava da salje masu formulara i daje mu se greska kao za nepotpun form da se ne otkrije da je popunjavanje tog input polja greska
	if f.Fk_ != "" {
		l(r, 7, fmt.Errorf("bot attempt with sign up input field: "+f.Fk_))
		return clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}
	if f.Confirm == "" {
		l(r, 7, fmt.Errorf("bot attempt with empty sign up checkbox"))
		return clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}

	if err := ValidateEmailAddress(f.User_email1); err != nil {
		return err
	}

	if err := ValidatePassword(f.Password1); err != nil {
		return err
	}

	return nil
}

func ValidateForgottenPasswordData(r *http.Request, f models.ForgottenPasswordData) (models.User, error) {

	// log.Println("parsovan i dekodiran post form:", f, f.User_email1)

	// validacija forma za ispravnost ssr bycrypt key i fake key koji treba da zavara botove
	if err := utils.Check64UrlKey(f.Bk_, "", 1); err != nil {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Session_expired")
	}

	// validacija forma za dva ista
	if f.User_email1 != f.User_email2 {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Email_not_same")
	}

	//bot pokusava da salje masu formulara i daje mu se greska kao za nepotpun form da se ne otkrije da je popunjavanje tog input polja greska
	if f.User_email1 == "" || f.User_email2 == "" || f.Bk_ == "" {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}
	if f.Fk_ != "" {
		l(r, 7, fmt.Errorf("bot attempt with sign up input field: "+f.Fk_))
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}
	if f.Confirm == "" {
		l(r, 7, fmt.Errorf("bot attempt with empty sign up checkbox"))
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}

	if err := ValidateEmailAddress(f.User_email1); err != nil {
		return models.User{}, err
	}

	// CheckEmail vraca gresku (i User struct) ako vec ima tog mejla a nil ako nema tog mejla a ovde nam to bas treba tako da ako nema greske to znaci da nema tog mejla u db a to onda znaci da se otkazuje slanje verifikacionog linka na taj mejl jer neko to moze da zloupotrebi da bombarduje mejlove ljudi koji nisu clanovi sajta porukama sa sajta
	// provera da li ima isti takav email u db
	if user, err := models.CheckEmail(f.User_email1, r); err == nil {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Email_or_password_wrong")
	} else {
		return user, nil
	}

}

func ValidateChangeEmailData(r *http.Request, f models.ChangeEmailData) (models.User, error) {

	// log.Println("parsovan i dekodiran post form:", f, f.User_email1)

	// validacija forma za ispravnost ssr bycrypt key
	if err := utils.Check64UrlKey(f.Bk_, "", 1); err != nil {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Session_expired")
	}

	// validacija forma za dva ista
	if f.NewUser_email1 != f.NewUser_email2 {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Email_not_same")
	}

	//bot pokusava da salje masu formulara i daje mu se greska kao za nepotpun form da se ne otkrije da je popunjavanje tog input polja greska
	if f.User_email1 == "" || f.NewUser_email1 == "" || f.NewUser_email2 == "" || f.Password1 == "" || f.Bk_ == "" {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}
	if f.Fk_ != "" {
		l(r, 7, fmt.Errorf("bot attempt with sign up input field: "+f.Fk_))
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}
	if f.Confirm == "" {
		l(r, 7, fmt.Errorf("bot attempt with empty sign up checkbox"))
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}

	if err := ValidateEmailAddress(f.User_email1); err != nil {
		return models.User{}, err
	}
	if err := ValidateEmailAddress(f.NewUser_email1); err != nil {
		return models.User{}, err
	}
	if err := ValidateEmailAddress(f.NewUser_email2); err != nil {
		return models.User{}, err
	}

	_, err := models.CheckEmail(f.User_email1, r)
	if err == nil {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Email_or_password_wrong")
	} else {
		if user, err := models.CheckEmail(f.NewUser_email1, r); err != nil {
			return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Email_already_exists")
		} else {
			return user, nil
		}
	}

}

func ValidateChangeNameData(r *http.Request, f models.ChangeNameData) (models.User, error) {

	// log.Println("parsovan i dekodiran post form:", f, f.User_email1)

	// validacija forma za ispravnost ssr bycrypt key
	if err := utils.Check64UrlKey(f.Bk_, "", 1); err != nil {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Session_expired")
	}

	//bot pokusava da salje masu formulara i daje mu se greska kao za nepotpun form da se ne otkrije da je popunjavanje tog input polja greska
	if f.User_email1 == "" || f.Name == "" || f.Password1 == "" || f.Bk_ == "" {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}
	if f.Fk_ != "" {
		l(r, 7, fmt.Errorf("bot attempt with sign up input field: "+f.Fk_))
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}
	if f.Confirm == "" {
		l(r, 7, fmt.Errorf("bot attempt with empty sign up checkbox"))
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}

	if err := ValidateEmailAddress(f.User_email1); err != nil {
		return models.User{}, err
	}
	if err := ValidateUserName(f.Name); err != nil {
		return models.User{}, err
	}
	if err := ValidatePassword(f.Password1); err != nil {
		return models.User{}, err
	}

	// CheckEmail vraca gresku (i User struct) ako vec ima tog mejla a nil ako nema tog mejla
	user, err := models.CheckEmail(f.User_email1, r)
	if err == nil {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Email_or_password_wrong")
	} else {
		return user, nil
	}

}

func ValidateChangePasswordData(r *http.Request, f models.ChangePasswordData) (models.User, error) {

	// log.Println("parsovan i dekodiran post form:", f, f.User_email1)

	// validacija forma za ispravnost ssr bycrypt key
	if err := utils.Check64UrlKey(f.Bk_, "", 1); err != nil {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Session_expired")
	}

	if f.Password1 != f.Password2 {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Password_not_same")
	}

	//bot pokusava da salje masu formulara i daje mu se greska kao za nepotpun form da se ne otkrije da je popunjavanje tog input polja greska
	if f.User_email1 == "" || f.Password0 == "" || f.Password1 == "" || f.Password2 == "" || f.Bk_ == "" {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}
	if f.Fk_ != "" {
		l(r, 7, fmt.Errorf("bot attempt with sign up input field: "+f.Fk_))
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}
	if f.Confirm == "" {
		l(r, 7, fmt.Errorf("bot attempt with empty sign up checkbox"))
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Form_error")
	}

	if err := ValidateEmailAddress(f.User_email1); err != nil {
		return models.User{}, err
	}
	if err := ValidatePassword(f.Password0); err != nil {
		return models.User{}, err
	}
	if err := ValidatePassword(f.Password1); err != nil {
		return models.User{}, err
	}

	// CheckEmail vraca gresku (i User struct) ako vec ima tog mejla a nil ako nema tog mejla
	user, err := models.CheckEmail(f.User_email1, r)
	if err == nil {
		return models.User{}, clr.NewAPIError(http.StatusNotAcceptable, "Email_or_password_wrong")
	} else {
		return user, nil
	}

}

func DecodeBodyAndValidatePostAndPutData1(tableDb string, dec *json.Decoder) (any, error) {
	switch tableDb {
	case "g_pitanja_c_testovi":
		var recordData models.Test
		if err := dec.Decode(&recordData); err != nil {
			return nil, l(nil, 4, err)
		}
		if err := validateTestData(recordData); err != nil {
			return nil, err
		}
		return recordData, nil
	case "mi_users":
		var recordData models.User
		if err := dec.Decode(&recordData); err != nil {
			return nil, l(nil, 4, err)
		}
		if err := validateUserData1(recordData); err != nil {
			return nil, err
		}
		return recordData, nil
	case "g_user_blog":
		var recordData models.Note
		if err := dec.Decode(&recordData); err != nil {
			return nil, l(nil, 4, err)
		}
		if err := validateNoteData(recordData); err != nil {
			return nil, err
		}
		return recordData, nil
	case "v_settings":
		var recordData models.Settings
		if err := dec.Decode(&recordData); err != nil {
			return nil, l(nil, 4, err)
		}
		if err := validateSettingsData(recordData); err != nil {
			return nil, err
		}
		return recordData, nil
	default:
		return nil, l(nil, 8, fmt.Errorf("malformed table name"))
	}
}

func validateUserData1(recordData models.User) error {

	if err := ValidateEmailAddress(recordData.Email); err != nil {
		return err
	}

	if err := ValidateUserName(recordData.User_name); err != nil {
		return err
	}

	if err := ValidatePassword(recordData.Hash_lozinka); err != nil {
		return err
	}

	return nil
}

func validateTestData(recordData models.Test) error {
	fmt.Print(recordData)
	return nil
}
func validateNoteData(recordData models.Note) error {
	fmt.Print(recordData)
	return nil
}
func validateSettingsData(recordData models.Settings) error {
	fmt.Print(recordData)
	return nil
}

/*

kim tec
011/44-44-560 kim tek
011/331-35-68

ektro
035/8245-834
slavke đurđević b1/3

džin pc & klime

aleksandar stev. prvovenćanog
065/245-5700

cold clima m. tepića 17
069/40-54-455

*/

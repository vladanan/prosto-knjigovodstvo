// Package routes služi da obrađuje zahvete iz main
package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/a-h/templ"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/vladanan/prosto/src/controllers/clr"
	"github.com/vladanan/prosto/src/controllers/i18n"
	"github.com/vladanan/prosto/src/controllers/utils"
	"github.com/vladanan/prosto/src/controllers/vet"
	"github.com/vladanan/prosto/src/models"
	"github.com/vladanan/prosto/src/views"
	"github.com/vladanan/prosto/src/views/cp"
	"github.com/vladanan/prosto/src/views/dashboard"
	"github.com/vladanan/prosto/src/views/firma"
	"github.com/vladanan/prosto/src/views/site"
)

func RouterSite(r *mux.Router) {
	r.HandleFunc("/", Index)
	r.HandleFunc("/pausal", Pausal)
	r.HandleFunc("/user_portal", UserPortal)
	r.HandleFunc("/privacy", Privacy)
	r.HandleFunc("/terms", Terms)
	r.HandleFunc("/arrangements", Arrangements)
}

func RouterUsers(r *mux.Router) {
	r.HandleFunc("/sign_up", Sign_up)
	r.HandleFunc("/sign_in", Sign_in)
	r.HandleFunc("/sign_out", Sign_out)
	r.HandleFunc("/forgotten_password", ForgottenPassword)
	r.HandleFunc("/dashboard", Dashboard)
}

func RouterFirma(r *mux.Router) {
	r.HandleFunc("/fakture", Fakture)
	r.HandleFunc("/kpo", Kpo)
	r.HandleFunc("/zurnal", Zurnal)
	r.HandleFunc("/klijenti", Klijenti)
	r.HandleFunc("/artikli", Artikli)
}

func RouterI18n(r *mux.Router) {
	r.HandleFunc("/sh", SetSh)
	r.HandleFunc("/en", SetEn)
	r.HandleFunc("/es", SetEs)
}

// static sa funkcijom koja pravi niz mux PathPrefix handlera jer bi se inače main zagušio sa vazdan njih
// za svaki folder gde se koristi path sa promeljivima kao što je r.HandleFunc("/auth/vmk/{key}"
// https://stackoverflow.com/questions/15834278/serving-static-content-with-a-root-url-with-the-gorilla-toolkits
func ServeStatic(router *mux.Router, staticDirectory string) {
	staticPaths := map[string]string{
		"/":    "" + staticDirectory,
		"vmk":  "/auth/vmk" + staticDirectory,
		"fp":   "/auth/fp" + staticDirectory,
		"du":   "/auth/du" + staticDirectory,
		"cm":   "/auth/cm" + staticDirectory,
		"cn":   "/auth/cn" + staticDirectory,
		"cp":   "/auth/cp" + staticDirectory,
		"auth": "/auth" + staticDirectory,
		// "qapi": "/questions" + staticDirectory,
	}
	for _, pathValue := range staticPaths {
		// pathPrefix := "/" + pathName + "/"
		router.PathPrefix(pathValue).Handler(http.StripPrefix(pathValue, http.FileServer(http.Dir("assets"))))
	}
}

func apiCallGet[
	T models.Test | models.User | models.Note | models.Settings | models.UserData](
	table, field, record string, r *http.Request) ([]T, error) {

	// log.Println("api call get:", table, field, record)
	switch {
	case table != "" && field != "" && record != "":
		table = "/" + table
		field = "/" + field
		record = "/" + record
	case table != "" && field == "" && record == "":
		table = "/" + table
	default:
		return nil, l(r, 7, clr.NewAPIError(http.StatusNotAcceptable, "Path_elements_missing"))
	}

	var url string
	if os.Getenv("PRODUCTION") == "FALSE" {
		url = "http://127.0.0.1:7331/api/v" + table + field + record
	} else {
		url = "https://vezbamo.onrender.com/api/v" + table + field + record
	}

	resp, err := http.Get(url)
	if err != nil {
		// error if the request fails, such as if the requested URL is not found, or if the server is not reachable
		return nil, l(r, 8, err)
	} else {
		data, err := io.ReadAll(resp.Body)
		dec := json.NewDecoder(bytes.NewReader(data))
		if err != nil {
			return nil, l(r, 0, err)
		} else if resp.StatusCode == http.StatusOK {
			// successful request should return a 200 OK status, if not we should log and then exit with error
			var model []T
			if err := dec.Decode(&model); err != nil {
				return nil, l(r, 4, err)
			}
			defer resp.Body.Close()
			return model, nil
		} else {
			var apiError clr.APIError
			if err := dec.Decode(&apiError); err != nil {
				log.Println(err)
				errMsg := strings.ReplaceAll(string(data), "\n", "")
				return nil, l(r, 7, fmt.Errorf(errMsg))
			}
			defer resp.Body.Close()
			return nil, apiError
		}
	}

}

////**** SITE

// var key = []byte(os.Getenv("SESSION_KEY"))
// var store = sessions.NewCookieStore(key)
var store *sessions.CookieStore

func Index(w http.ResponseWriter, r *http.Request) {
	views.Index(r).Render(r.Context(), w)
}

func GoTo404(w http.ResponseWriter, r *http.Request) {
	site.Page404().Render(r.Context(), w)
}

func Pausal(w http.ResponseWriter, r *http.Request) {
	site.Pausal(r).Render(r.Context(), w)
}

func UserPortal(w http.ResponseWriter, r *http.Request) {
	// if notes, err := apiCallGet[models.Note]("note", "", "", r); err != nil {
	// 	smtu(w, r, l(r, 7, err))
	// } else {
	site.UserPortal(r, []models.Note{}).Render(r.Context(), w)
	// }
}

func Privacy(w http.ResponseWriter, r *http.Request) {
	site.Privacy(r).Render(r.Context(), w)
}

func Terms(w http.ResponseWriter, r *http.Request) {
	site.Terms(r).Render(r.Context(), w)
}

func Arrangements(w http.ResponseWriter, r *http.Request) {
	cp.Arrangements().Render(r.Context(), w)
}

////**** USERS 123

func Sign_up(w http.ResponseWriter, r *http.Request) {
	if already_authenticated, _, _, err := getSesionData(r); err == nil && !already_authenticated {
		dashboard.Sign_up(r).Render(r.Context(), w)
	} else {
		views.Index(r).Render(r.Context(), w)
	}
}

// Send message to user on fe
func smtu(w http.ResponseWriter, r *http.Request, message any) {
	var errorMsg string
	switch msg := message.(type) {
	case error:
		errorMsg = clr.CheckErr(msg).Msg
	case string:
		errorMsg = msg
	default:
		log.Println("u smtu je poslato nešto što nije ni string ni error")
		errorMsg = "Internal_server_error_1"
	}
	dashboard.MessageForUser(r, errorMsg).Render(r.Context(), w)
}

func getSesionData(r *http.Request) (bool, string, string, error) {

	key := i18n.SessionKey("session")
	session, err := i18n.GetSessionFromContext(r.Context(), key)
	if err != nil {
		return false, "", "", l(r, 8, err)
	}

	var alreadyAuthenticated bool
	authMap := session.Values["authenticated"]
	if authMap != nil {
		alreadyAuthenticated = authMap.(bool)
	}

	emailMap := session.Values["user_email"]
	email := ""
	if emailMap != nil {
		email = emailMap.(string)
	}

	nameMap := session.Values["name"]
	name := ""
	if nameMap == nil {
		// fmt.Println("nema mail:", session.Values["user_email"])
	} else {
		// fmt.Println("ima mail:", session.Values["user_email"])
		name = nameMap.(string)
	}

	return alreadyAuthenticated, email, name, nil
}

// Genericki form Pack za response forms
type formPack[T models.SignUpFormData | models.SignInFormData | models.ForgottenPasswordData | models.ChangeEmailData | models.ChangeNameData | models.ChangePasswordData] struct {
	w      http.ResponseWriter
	r      *http.Request
	f      T
	fTempl func(*http.Request, T, string) templ.Component
}

// Genericka form Response func koja koristi form Pack i error
func formResp[T models.SignUpFormData | models.SignInFormData | models.ForgottenPasswordData | models.ChangeEmailData | models.ChangeNameData | models.ChangePasswordData](fp formPack[T], err error) {
	fp.fTempl(fp.r, fp.f, clr.CheckErr(err).Msg).Render(fp.r.Context(), fp.w)
}

func Sign_up_post(w http.ResponseWriter, r *http.Request) {
	fp := formPack[models.SignUpFormData]{
		w: w,
		r: r,
		f: models.SignUpFormData{
			User_email1: r.FormValue("user_email1"),
			User_email2: r.FormValue("user_email2"),
			Name:        r.FormValue("name"),
			Password1:   r.FormValue("password1"),
			Password2:   r.FormValue("password2"),
			Confirm:     r.FormValue("confirm"),
			Submit:      r.FormValue("submit"),
			Wait:        r.FormValue("wait"),
			Bk_:         r.FormValue("bk_"),
			Fk_:         r.FormValue("fk_"),
		},
		fTempl: cp.Sign_up,
	}

	if _, err := vet.ValidateSignUpData(r); err != nil {
		formResp(fp, err)
	} else if err := models.AddUser(fp.f, r); err != nil {
		formResp(fp, err)
	} else {
		clr.GetStringLogger()(r, 4, "web: created user for email: "+fp.f.User_email1)
		cp.MessageForUser(r, "User_registered").Render(r.Context(), w)
	}

	// TREBACE ZA KASNIJE LEPO SREDJEN POST FORM API POZIV
	// // https://www.sohamkamani.com/golang/http-client/
	// var apiUrl string
	// if os.Getenv("PRODUCTION") == "FALSE" {
	// 	apiUrl = "http://127.0.0.1:7331/api/v/krsnc_usrs"
	// } else {
	// 	apiUrl = "https://vezbamo.onrender.com/api/v/krsnc_usrs"
	// }
	// // api post poziv i kroz error wraper dobija se apiError ili system err i na to se salje odgovarajuci sign up form
	// resp, err := http.PostForm(apiUrl, r.Form)
	// if err != nil {
	// 	formRender(fp, l(r, 4, err))
	// }
	// defer resp.Body.Close()
	// if data, err := io.ReadAll(resp.Body); err != nil {
	// 	formRender(fp, l(r, 4, err))
	// } else if resp.StatusCode != http.StatusOK {
	// 	var apiErr clr.APIError
	// 	if err := json.NewDecoder(bytes.NewReader(data)).Decode(&apiErr); err != nil {
	// 		// greska koja nije apiError recimo kada flood rate limiter odbije api/auth request
	// 		log.Println(err)
	// 		// prilikom pretvaranja []byte u string dodaju se automatski navodnici i novi red
	// 		// navodnici ne smetaju za pravljenje json jer se automatski escapuju ali novi red ga ubija
	// 		formRender(fp, l(r, 4, fmt.Errorf("web: "+resp.Status+" "+strings.ReplaceAll(string(data), "\n", ""))))
	// 	} else {
	// 		l(r, 4, fmt.Errorf("%v, url: %v", apiErr.Error(), apiUrl))
	// 		formRender(fp, apiErr)
	// 	}
	// } else {

	// 	clr.GetStringLogger()(r, 4, "web: "+resp.Status+" "+strings.ReplaceAll(string(data), "\n", ""))
	// 	cp.MessageForUser(r, "User_registered").Render(r.Context(), w)
	// }

}

func CheckLinkFromEmailRegister(w http.ResponseWriter, r *http.Request) {
	// https://stackoverflow.com/questions/45378566/gorilla-mux-optional-query-values
	// deo iz query URL.Query i FormValue ne rade na isti način pogotovo ako u r ima body i multipart form
	// fmt.Print("CheckLinkFromEmail: url vars and queries:", vars, r.URL.Query()["email"][0], r.FormValue("email"), "\n")
	// fmt.Println("ceo url query", r.URL.Query())
	vars := mux.Vars(r)
	key := vars["key"]
	email := r.URL.Query()["user_email"][0]
	if err := vet.ValidateEmailAddress(email); err != nil {
		smtu(w, r, err)
	} else if err := models.VerifyEmail(key, email, r); err != nil {
		smtu(w, r, err)
	} else {
		smtu(w, r, "Mail_verified")
	}
}

func SendNewSecretInputField(w http.ResponseWriter, r *http.Request) {
	// proverava da li je prvobitni kod u formu za desetice minuta ispravan pre nego sto posalje novi kod za jedinice minuta
	if err := utils.Check64UrlKey(r.FormValue("bk_"), "", 10); err != nil {
		cp.ErrorForUser(r, "Session_expired").Render(r.Context(), w)
	} else {
		// salje se novo polje sa kodom za jedinice minuta
		cp.NewSecretInputField().Render(r.Context(), w)
	}
}

func CheckUserName(w http.ResponseWriter, r *http.Request) {
	userName := r.URL.Query()["name"][0]

	if err := models.CheckUserName(userName, r); err != nil {
		if _, ok := err.(clr.APIError); ok {
			w.Write([]byte("<span class='text-sm'>&#x274c;</span>"))
		} else {
			l(r, 4, err)
			w.WriteHeader(500)
			cp.ErrorForUser(r, clr.CheckErr(err).Msg).Render(r.Context(), w)
			// w.Write([]byte("<span class='text-sm'>&#128681;</span>")) // crvena zastavica
			// w.Write([]byte("<span class='text-sm'>&#9888;</span>")) // zuti trougao
		}
	} else {
		w.Write([]byte("&check;"))
		// w.Write([]byte("&#10004;"))
		// w.Write([]byte("&#9989;"))
	}
}

func Sign_in(w http.ResponseWriter, r *http.Request) {
	if already_authenticated, _, _, err := getSesionData(r); err == nil && !already_authenticated {
		dashboard.Sign_in(r).Render(r.Context(), w)

		// resp, err := http.Get("http://127.0.0.1:7331/sign_in_remote")
		// if err != nil {
		// 	// error if the request fails, such as if the requested URL is not found, or if the server is not reachable
		// 	http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError)
		// 	return
		// } else {
		// 	data, err := io.ReadAll(resp.Body)
		// 	// dec := json.NewDecoder(bytes.NewReader(data))
		// 	if err != nil {
		// 		http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError)
		// 		return
		// 	} else if resp.StatusCode == http.StatusOK {
		// 		// successful request should reeturn a 200 OK status, if not we should log and then exit with error
		// 		// log.Println("data sa be", string(data))
		// 		// w.Write(data)
		// 		dashboard.Sign_in(r, string(data)).Render(r.Context(), w)
		// 		defer resp.Body.Close()
		// 		return
		// 	} else {
		// 		log.Println("error sa be", string(data))
		// 		defer resp.Body.Close()
		// 		return
		// 	}
		// }

	} else {
		views.Index(r).Render(r.Context(), w)
	}
}

func Sign_in_post(w http.ResponseWriter, r *http.Request) {
	// fmt.Println("SING IN POST:", r.Body, r.MultipartForm, r.URL, r.PostForm, r.Form)
	// fmt.Println("SIGN IN POST form:", r.FormValue("mail"), r.FormValue("password"))

	key := i18n.SessionKey("session")
	session, err := i18n.GetSessionFromContext(r.Context(), key)
	if err != nil {
		http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError)
		return
	}

	// v := url.Values{}
	// // sign_up
	// v.Set("user_email1", r.FormValue("user_email1"))
	// v.Set("password1", r.FormValue("password1"))
	// v.Set("confirm", r.FormValue("confirm"))
	// v.Set("submit", r.FormValue("submit"))
	// v.Set("wait", r.FormValue("wait"))
	// // sign_in

	// v.Set("bk_", r.FormValue("bk_"))
	// v.Set("fk_", r.FormValue("fk_"))
	// resp, err := http.PostForm("http://127.0.0.1:7331/auth/sign_in_post", v)
	// if err != nil {
	// 	log.Print(err)
	// 	return
	// }
	// if resp.StatusCode != http.StatusFound {
	// 	session.Values["authenticated"] = false
	// 	session.Values["user_email"] = ""
	// } else {
	// 	session.Values["authenticated"] = true
	// 	session.Values["user_email"] = r.FormValue("user_email1")
	// 	session.Values["name"] = "neki tip"
	// 	session.Values["mode"] = "user"
	// 	data, err := io.ReadAll(resp.Body)
	// 	// dec := json.NewDecoder(bytes.NewReader(data))
	// 	if err != nil {
	// 		// http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError)
	// 		log.Println("dashboard greska sa be", string(data))
	// 		w.Write(data)
	// 		defer resp.Body.Close()
	// 		return
	// 	} else {
	// 		// successful request should reeturn a 200 OK status, if not we should log and then exit with error
	// 		log.Println("dashboard sa be", string(data))
	// 		w.WriteHeader(302) //da bi htmx u templ formu drugacije reagovao na dobijeni form i dobijeni dashboard
	// 		w.Write(data)
	// 		// dashboard.Sign_in(r, string(data)).Render(r.Context(), w)
	// 		defer resp.Body.Close()
	// 		return
	// 	}
	// }

	si := models.SignInFormData{
		User_email1: r.FormValue("user_email1"),
		Password1:   r.FormValue("password1"),
		Confirm:     r.FormValue("confirm"),
		Submit:      r.FormValue("submit"),
		Wait:        r.FormValue("wait"),
		Bk_:         r.FormValue("bk_"),
		Fk_:         r.FormValue("fk_"),
	}
	// log.Print(si)
	fp := formPack[models.SignInFormData]{
		w:      w,
		r:      r,
		f:      si,
		fTempl: cp.Sign_in,
	}

	if err := vet.ValidateSignInData(r, si); err != nil {
		formResp(fp, err)
	} else {

		user, err := models.AuthenticateUser(si.User_email1, si.Password1, r)
		if err != nil {
			formResp(fp, l(r, 4, err))
			session.Values["authenticated"] = false
			session.Values["user_email"] = ""
			if err = session.Save(r, w); err != nil {
				// http.Error(w, l(r, 4, err).Error(), http.StatusInternalServerError)
				formResp(fp, l(r, 4, err))
			}
		} else {
			session.Values["authenticated"] = true
			session.Values["user_email"] = user.Email
			session.Values["name"] = user.User_name
			session.Values["mode"] = user.Mode
			// Save it before we write to the response/return from the handler.
			if err = session.Save(r, w); err != nil {
				formResp(fp, l(r, 4, err))
			}
			w.WriteHeader(302) //da bi htmx u templ formu drugacije reagovao na dobijeni form i dobijeni dashboard
			// dashboard.Dashboard(r, user).Render(r.Context(), w)
			Dashboard(w, r)
		}
	}

}

func AutoLoginDemo(w http.ResponseWriter, r *http.Request) {
	key := i18n.SessionKey("session")
	session, err := i18n.GetSessionFromContext(r.Context(), key)
	if err != nil {
		http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError)
		return
	}

	// key := $2a$12$nmIsb4TJ/NfWaxwCdRrPWuRkfx790PVX0pURqm2D9QujBpxwh6WZO
	email := "y.emailbox-proba@yahoo.com"
	password := "y.emailbox-proba@yahoo.com"
	user, err := models.AuthenticateUser(email, password, r)
	if err != nil {
		smtu(w, r, l(r, 7, err))
		session.Values["authenticated"] = false
		session.Values["user_email"] = ""
		if err = session.Save(r, w); err != nil {
			http.Error(w, l(r, 4, err).Error(), http.StatusInternalServerError)
		}
	} else {
		session.Values["authenticated"] = true
		session.Values["user_email"] = user.Email
		session.Values["name"] = user.User_name
		session.Values["mode"] = user.Mode
		// Save it before we write to the response/return from the handler.
		if err = session.Save(r, w); err != nil {
			http.Error(w, l(r, 4, err).Error(), http.StatusInternalServerError)
		}
		log.Println("Sign_in_post: autentikacija JE PROŠLA")
		// dashboard.Dashboard(r, user).Render(r.Context(), w)
		Dashboard(w, r)
	}
}

func AutoLoginUser(w http.ResponseWriter, r *http.Request) {
	key := i18n.SessionKey("session")
	session, err := i18n.GetSessionFromContext(r.Context(), key)
	if err != nil {
		http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError)
		return
	}

	// key := $2a$12$K6NV/GHc7qkMoabnsVMONOOhOoWQkIYnr6s5kty70VIPPpw.CbU0G
	email := "vladan_zasve@yahoo.com"
	password := "321654987"
	user, err := models.AuthenticateUser(email, password, r)
	if err != nil {
		smtu(w, r, l(r, 7, err))
		session.Values["authenticated"] = false
		session.Values["user_email"] = ""
		if err = session.Save(r, w); err != nil {
			http.Error(w, l(r, 4, err).Error(), http.StatusInternalServerError)
		}
	} else {
		session.Values["authenticated"] = true
		session.Values["user_email"] = user.Email
		session.Values["name"] = user.User_name
		session.Values["mode"] = user.Mode
		// Save it before we write to the response/return from the handler.
		if err = session.Save(r, w); err != nil {
			http.Error(w, l(r, 4, err).Error(), http.StatusInternalServerError)
		}
		log.Println("Sign_in_post: autentikacija JE PROŠLA")
		// dashboard.Dashboard(r, user).Render(r.Context(), w)
		Dashboard(w, r)
	}
}

func AutoLoginAdmin(w http.ResponseWriter, r *http.Request) {
	key := i18n.SessionKey("session")
	session, err := i18n.GetSessionFromContext(r.Context(), key)
	if err != nil {
		http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError)
		return
	}

	// key := $2a$12$c4wyJyZrY4E4O4yn1wJrOeRoygscTdvIV79ia/L7Juw5BhVhFw2Mi
	email := "vladan.andjelkovic@gmail.com"
	password := "vezbamo.2015"
	user, err := models.AuthenticateUser(email, password, r)
	if err != nil {
		smtu(w, r, l(r, 7, err))
		session.Values["authenticated"] = false
		session.Values["user_email"] = ""
		if err = session.Save(r, w); err != nil {
			http.Error(w, l(r, 4, err).Error(), http.StatusInternalServerError)
		}
	} else {
		session.Values["authenticated"] = true
		session.Values["user_email"] = user.Email
		session.Values["name"] = user.User_name
		session.Values["mode"] = user.Mode
		// Save it before we write to the response/return from the handler.
		if err = session.Save(r, w); err != nil {
			http.Error(w, l(r, 4, err).Error(), http.StatusInternalServerError)
		}
		log.Println("Sign_in_post: autentikacija JE PROŠLA")
		// dashboard.Dashboard(r, user).Render(r.Context(), w)
		Dashboard(w, r)
	}
}

func ForgottenPassword(w http.ResponseWriter, r *http.Request) {
	if already_authenticated, _, _, err := getSesionData(r); err == nil && !already_authenticated {
		dashboard.Forgotten_password(r).Render(r.Context(), w)
	} else {
		views.Index(r).Render(r.Context(), w)
	}
}

func ForgottenPasswordSendMail(w http.ResponseWriter, r *http.Request) {
	fp := formPack[models.ForgottenPasswordData]{
		w: w,
		r: r,
		f: models.ForgottenPasswordData{
			User_email1: r.FormValue("user_email1"),
			User_email2: r.FormValue("user_email2"),
			Confirm:     r.FormValue("confirm"),
			Submit:      r.FormValue("submit"),
			Wait:        r.FormValue("wait"),
			Bk_:         r.FormValue("bk_"),
			Fk_:         r.FormValue("fk_"),
		},
		fTempl: cp.Forgotten_password,
	}

	if user, err := vet.ValidateForgottenPasswordData(r, fp.f); err != nil {
		formResp(fp, err)
	} else {
		// key za verifikaciju se dobija od mejla i stringa za datum odsecenog na desetice minuta
		urlKey := utils.Get64UrlKey(fp.f.User_email1, 10)
		// kreira se enkriptovan teskt za mejl (da se ne bi mejl sa vremenski valdinim url key zloupottrebio da se resetuje neciji tudji nalog)
		// ali tako da nije enkodovan samo mejl da se oteza falsifikovanje
		apiAdd := "forgotten-passwd" // 16byte
		key := utils.DateEncryptionKey(apiAdd)
		msg := fp.f.User_email1 + "||" + user.User_name
		if msgKey, err := utils.MsgEncrypt(msg, apiAdd, key); err != nil {
			formResp(fp, l(r, 4, err))
		} else {
			// kreira se url za mejl
			var url string
			if os.Getenv("PRODUCTION") == "FALSE" {
				url = "http://127.0.0.1:7331/auth/fp/" + urlKey + "?fpm=" + msgKey
			} else {
				url = "https://vezbamo.onrender.com/auth/fp/" + urlKey + "?fpm=" + msgKey
			}
			// salje se mejl 123
			email := utils.MailForgottenPasswordRequest{
				R:        r,
				Email:    fp.f.User_email1,
				UserName: user.User_name,
				Url:      url,
			}
			if err := utils.SendEmail(email); err != nil {
				formResp(fp, l(r, 4, err))
			} else {
				cp.MessageForUser(r, "Fp_verify_sent").Render(r.Context(), w)
			}
		}
	}

}

func CheckLinkFromEmailFP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlKey := vars["key"]
	msgKey := r.URL.Query()["fpm"][0]
	apiAdd := "forgotten-passwd" // 16byte
	key := utils.DateEncryptionKey(apiAdd)

	if plaintexts, err := utils.MsgDecrypt(msgKey, apiAdd, key); err != nil {
		smtu(w, r, (l(r, 4, err)))
	} else if plaintexts[0] != apiAdd { // verifikovanje dekriptovanih poruka
		smtu(w, r, clr.NewAPIError(http.StatusNotAcceptable, "Email_link_not_verified"))
	} else {
		email := plaintexts[1]
		userName := plaintexts[2]
		// provera validnosti vremenskog url kljuca iz mejla
		if err := utils.Check64UrlKey(urlKey, email, 10); err != nil {
			smtu(w, r, clr.NewAPIError(http.StatusNotAcceptable, "Session_expired"))
		} else if newPass, newPassKey, err := utils.GetNewPassword(""); err != nil { // sve ok i pravi se novi random pass i key za db
			smtu(w, r, err)
		} else if err := models.ReplaceForgottenPassword(email, userName, newPass, newPassKey, r); err != nil {
			smtu(w, r, err)
		} else {
			clr.GetStringLogger()(r, 4, "new random password and key: "+newPass+" "+newPassKey)
			smtu(w, r, "Fp_ok")
		}
	}

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	if already_authenticated, _, _, err := getSesionData(r); err == nil && !already_authenticated {
		views.Index(r).Render(r.Context(), w)
	} else {
		dashboard.Delete_user(r).Render(r.Context(), w)
	}
}

func DeleteUserSendMail(w http.ResponseWriter, r *http.Request) {
	fp := formPack[models.SignInFormData]{
		w: w,
		r: r,
		f: models.SignInFormData{
			User_email1: r.FormValue("user_email1"),
			Password1:   r.FormValue("password1"),
			Confirm:     r.FormValue("confirm"),
			Submit:      r.FormValue("submit"),
			Wait:        r.FormValue("wait"),
			Bk_:         r.FormValue("bk_"),
			Fk_:         r.FormValue("fk_"),
		},
		fTempl: cp.Delete_User,
	}

	if err := vet.ValidateSignInData(r, fp.f); err != nil {
		formResp(fp, err)
	} else if user, err := models.AuthenticateUser(fp.f.User_email1, fp.f.Password1, r); err != nil {
		formResp(fp, err)
	} else {
		urlKey := utils.Get64UrlKey(fp.f.User_email1, 10)
		apiAdd := "delete-user-perm" // 16byte
		key := utils.DateEncryptionKey(apiAdd)
		msg := fp.f.User_email1 + "||" + user.User_name
		if msgKey, err := utils.MsgEncrypt(msg, apiAdd, key); err != nil {
			formResp(fp, l(r, 4, err))
		} else {
			var url string
			if os.Getenv("PRODUCTION") == "FALSE" {
				url = "http://127.0.0.1:7331/auth/du/" + urlKey + "?dum=" + msgKey
			} else {
				url = "https://vezbamo.onrender.com/auth/du/" + urlKey + "?dum=" + msgKey
			}
			email := utils.DeleteUserRequest{
				R:        r,
				Email:    fp.f.User_email1,
				UserName: user.User_name,
				Url:      url,
			}
			if err := utils.SendEmail(email); err != nil {
				formResp(fp, l(r, 4, err))
			} else {
				cp.MessageForUser(r, "DU_verify_sent").Render(r.Context(), w)
			}
		}
	}

}

func CheckLinkFromEmailDU(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlKey := vars["key"]
	msgKey := r.URL.Query()["dum"][0]
	apiAdd := "delete-user-perm"
	key := utils.DateEncryptionKey(apiAdd)

	if plaintexts, err := utils.MsgDecrypt(msgKey, apiAdd, key); err != nil {
		smtu(w, r, (l(r, 4, err)))
	} else if plaintexts[0] != apiAdd {
		smtu(w, r, clr.NewAPIError(http.StatusNotAcceptable, "Email_link_not_verified"))
	} else {
		email := plaintexts[1]
		userName := plaintexts[2]
		if err := utils.Check64UrlKey(urlKey, email, 10); err != nil {
			smtu(w, r, clr.NewAPIError(http.StatusNotAcceptable, "Session_expired"))
		} else if err := models.DeleteUser(email, userName, r); err != nil {
			smtu(w, r, err)
		} else {
			key := i18n.SessionKey("session")
			session, err := i18n.GetSessionFromContext(r.Context(), key)
			if err != nil {
				l(r, 8, err)
			}
			session.Values["authenticated"] = false
			session.Values["user_email"] = ""
			session.Values["name"] = ""
			if err := session.Save(r, w); err != nil {
				l(r, 8, err)
			}
			clr.GetStringLogger()(r, 4, "web: deleted user for email: "+email)
			smtu(w, r, "DU_ok")
		}
	}

}

func ChangeEmail(w http.ResponseWriter, r *http.Request) {
	if already_authenticated, _, _, err := getSesionData(r); err == nil && !already_authenticated {
		views.Index(r).Render(r.Context(), w)
	} else {
		dashboard.Change_email(r).Render(r.Context(), w)
	}
}

func ChangeEmailSendMail(w http.ResponseWriter, r *http.Request) {
	fp := formPack[models.ChangeEmailData]{
		w: w,
		r: r,
		f: models.ChangeEmailData{
			User_email1:    r.FormValue("user_email1"),
			NewUser_email1: r.FormValue("new_user_email1"),
			NewUser_email2: r.FormValue("new_user_email2"),
			Password1:      r.FormValue("password1"),
			Confirm:        r.FormValue("confirm"),
			Submit:         r.FormValue("submit"),
			Wait:           r.FormValue("wait"),
			Bk_:            r.FormValue("bk_"),
			Fk_:            r.FormValue("fk_"),
		},
		fTempl: cp.Change_email,
	}

	if _, err := vet.ValidateChangeEmailData(r, fp.f); err != nil {
		formResp(fp, err)
	} else {
		user, err := models.AuthenticateUser(fp.f.User_email1, fp.f.Password1, r)
		if err != nil {
			formResp(fp, err)
		} else {
			urlKey := utils.Get64UrlKey(fp.f.User_email1, 10)
			apiAdd := "change-email-usr" // 16byte
			key := utils.DateEncryptionKey(apiAdd)
			msg := fp.f.User_email1 + "||" + fp.f.NewUser_email1 + "||" + user.User_name
			if msgKey, err := utils.MsgEncrypt(msg, apiAdd, key); err != nil {
				formResp(fp, l(r, 4, err))
			} else {
				var url string
				if os.Getenv("PRODUCTION") == "FALSE" {
					url = "http://127.0.0.1:7331/auth/cm/" + urlKey + "?cmm=" + msgKey
				} else {
					url = "https://vezbamo.onrender.com/auth/cm/" + urlKey + "?cmm=" + msgKey
				}
				email := utils.ChangeEmailRequest{
					R:        r,
					Email:    fp.f.User_email1,
					NewEmail: fp.f.NewUser_email1,
					UserName: user.User_name,
					Url:      url,
				}
				if err := utils.SendEmail(email); err != nil {
					formResp(fp, l(r, 4, err))
				} else {
					cp.MessageForUser(r, "CM_verify_sent").Render(r.Context(), w)
				}
			}
		}
	}

}

func CheckLinkFromEmailCM(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlKey := vars["key"]
	msgKey := r.URL.Query()["cmm"][0]
	apiAdd := "change-email-usr"
	key := utils.DateEncryptionKey(apiAdd)

	if plaintexts, err := utils.MsgDecrypt(msgKey, apiAdd, key); err != nil {
		smtu(w, r, (l(r, 4, err)))
	} else if plaintexts[0] != apiAdd {
		smtu(w, r, clr.NewAPIError(http.StatusNotAcceptable, "Email_link_not_verified"))
	} else {
		email := plaintexts[1]
		newEmail := plaintexts[2]
		userName := plaintexts[3]
		if err := utils.Check64UrlKey(urlKey, email, 10); err != nil {
			smtu(w, r, clr.NewAPIError(http.StatusNotAcceptable, "Session_expired"))
		} else if err := models.ChangeEmail(email, newEmail, userName, r); err != nil {
			smtu(w, r, err)
		} else {
			key := i18n.SessionKey("session")
			session, err := i18n.GetSessionFromContext(r.Context(), key)
			if err != nil {
				l(r, 8, err)
			}
			session.Values["authenticated"] = false
			session.Values["user_email"] = ""
			session.Values["name"] = ""
			if err := session.Save(r, w); err != nil {
				l(r, 8, err)
			}
			clr.GetStringLogger()(r, 4, "web: changed: "+email+" to: "+newEmail)
			smtu(w, r, "CM_ok")
		}
	}

}

func ChangeName(w http.ResponseWriter, r *http.Request) {
	if already_authenticated, _, _, err := getSesionData(r); err == nil && !already_authenticated {
		views.Index(r).Render(r.Context(), w)
	} else {
		dashboard.Change_name(r).Render(r.Context(), w)
	}
}

func ChangeNameSendMail(w http.ResponseWriter, r *http.Request) {
	fp := formPack[models.ChangeNameData]{
		w: w,
		r: r,
		f: models.ChangeNameData{
			User_email1: r.FormValue("user_email1"),
			Name:        r.FormValue("name"),
			Password1:   r.FormValue("password1"),
			Confirm:     r.FormValue("confirm"),
			Submit:      r.FormValue("submit"),
			Wait:        r.FormValue("wait"),
			Bk_:         r.FormValue("bk_"),
			Fk_:         r.FormValue("fk_"),
		},
		fTempl: cp.Change_name,
	}

	if _, err := vet.ValidateChangeNameData(r, fp.f); err != nil {
		formResp(fp, err)
	} else {
		user, err := models.AuthenticateUser(fp.f.User_email1, fp.f.Password1, r)
		if err != nil {
			formResp(fp, err)
		} else {
			urlKey := utils.Get64UrlKey(fp.f.User_email1, 10)
			apiAdd := "change-name-user" // 16byte
			key := utils.DateEncryptionKey(apiAdd)
			msg := fp.f.User_email1 + "||" + user.User_name + "||" + fp.f.Name
			if msgKey, err := utils.MsgEncrypt(msg, apiAdd, key); err != nil {
				formResp(fp, l(r, 4, err))
			} else {
				var url string
				if os.Getenv("PRODUCTION") == "FALSE" {
					url = "http://127.0.0.1:7331/auth/cn/" + urlKey + "?cnm=" + msgKey
				} else {
					url = "https://vezbamo.onrender.com/auth/cn/" + urlKey + "?cnm=" + msgKey
				}
				email := utils.ChangeNameRequest{
					R:           r,
					Email:       fp.f.User_email1,
					NewUserName: fp.f.Name,
					UserName:    user.User_name,
					Url:         url,
				}
				if err := utils.SendEmail(email); err != nil {
					formResp(fp, l(r, 4, err))
				} else {
					cp.MessageForUser(r, "CN_verify_sent").Render(r.Context(), w)
				}
			}
		}
	}

}

func CheckLinkFromEmailCN(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlKey := vars["key"]
	msgKey := r.URL.Query()["cnm"][0]
	apiAdd := "change-name-user"
	key := utils.DateEncryptionKey(apiAdd)

	if plaintexts, err := utils.MsgDecrypt(msgKey, apiAdd, key); err != nil {
		smtu(w, r, (l(r, 4, err)))
	} else if plaintexts[0] != apiAdd {
		smtu(w, r, clr.NewAPIError(http.StatusNotAcceptable, "Email_link_not_verified"))
	} else {
		email := plaintexts[1]
		userName := plaintexts[2]
		newUserName := plaintexts[3]
		if err := utils.Check64UrlKey(urlKey, email, 10); err != nil {
			smtu(w, r, clr.NewAPIError(http.StatusNotAcceptable, "Session_expired"))
		} else if err := models.ChangeName(email, userName, newUserName, r); err != nil {
			smtu(w, r, err)
		} else {
			key := i18n.SessionKey("session")
			session, err := i18n.GetSessionFromContext(r.Context(), key)
			if err != nil {
				l(r, 8, err)
			}
			// session.Values["authenticated"] = false
			// session.Values["user_email"] = ""
			session.Values["name"] = newUserName
			if err := session.Save(r, w); err != nil {
				l(r, 8, err)
			}
			clr.GetStringLogger()(r, 4, "web: changed: "+userName+" to: "+newUserName)
			smtu(w, r, "CN_ok")
		}
	}

}

func ChangePassword(w http.ResponseWriter, r *http.Request) {
	if already_authenticated, _, _, err := getSesionData(r); err == nil && !already_authenticated {
		views.Index(r).Render(r.Context(), w)
	} else {
		dashboard.Change_password(r).Render(r.Context(), w)
	}
}

func ChangePasswordSendMail(w http.ResponseWriter, r *http.Request) {
	fp := formPack[models.ChangePasswordData]{
		w: w,
		r: r,
		f: models.ChangePasswordData{
			User_email1: r.FormValue("user_email1"),
			Password0:   r.FormValue("password0"),
			Password1:   r.FormValue("password1"),
			Password2:   r.FormValue("password2"),
			Confirm:     r.FormValue("confirm"),
			Submit:      r.FormValue("submit"),
			Wait:        r.FormValue("wait"),
			Bk_:         r.FormValue("bk_"),
			Fk_:         r.FormValue("fk_"),
		},
		fTempl: cp.Change_password,
	}

	if _, err := vet.ValidateChangePasswordData(r, fp.f); err != nil {
		formResp(fp, err)
	} else {
		user, err := models.AuthenticateUser(fp.f.User_email1, fp.f.Password0, r)
		if err != nil {
			formResp(fp, err)
		} else {
			urlKey := utils.Get64UrlKey(fp.f.User_email1, 10)
			apiAdd := "change-passwordu" // 16byte
			key := utils.DateEncryptionKey(apiAdd)
			msg := fp.f.User_email1 + "||" + fp.f.Password0 + "||" + fp.f.Password1
			if msgKey, err := utils.MsgEncrypt(msg, apiAdd, key); err != nil {
				formResp(fp, l(r, 4, err))
			} else {
				var url string
				if os.Getenv("PRODUCTION") == "FALSE" {
					url = "http://127.0.0.1:7331/auth/cp/" + urlKey + "?cpm=" + msgKey
				} else {
					url = "https://vezbamo.onrender.com/auth/cp/" + urlKey + "?cpm=" + msgKey
				}
				email := utils.ChangePasswordRequest{
					R:        r,
					Email:    fp.f.User_email1,
					UserName: user.User_name,
					Url:      url,
				}
				if err := utils.SendEmail(email); err != nil {
					formResp(fp, l(r, 4, err))
				} else {
					cp.MessageForUser(r, "CP_verify_sent").Render(r.Context(), w)
				}
			}
		}
	}

}

func CheckLinkFromEmailCP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	urlKey := vars["key"]
	msgKey := r.URL.Query()["cpm"][0]
	apiAdd := "change-passwordu"
	key := utils.DateEncryptionKey(apiAdd)

	if plaintexts, err := utils.MsgDecrypt(msgKey, apiAdd, key); err != nil {
		smtu(w, r, (l(r, 4, err)))
	} else if plaintexts[0] != apiAdd {
		smtu(w, r, clr.NewAPIError(http.StatusNotAcceptable, "Email_link_not_verified"))
	} else {
		email := plaintexts[1]
		password0 := plaintexts[2]
		password1 := plaintexts[3]
		if err := utils.Check64UrlKey(urlKey, email, 10); err != nil {
			smtu(w, r, clr.NewAPIError(http.StatusNotAcceptable, "Session_expired"))
		} else if err := models.ChangePassword(email, password0, password1, r); err != nil {
			smtu(w, r, err)
		} else {
			key := i18n.SessionKey("session")
			session, err := i18n.GetSessionFromContext(r.Context(), key)
			if err != nil {
				l(r, 8, err)
			}
			session.Values["authenticated"] = false
			session.Values["user_email"] = ""
			session.Values["name"] = ""
			if err := session.Save(r, w); err != nil {
				l(r, 8, err)
			}
			clr.GetStringLogger()(r, 4, "web: changed password for: "+email)
			smtu(w, r, "CP_ok")
		}
	}

}

func Dashboard(w http.ResponseWriter, r *http.Request) {
	if already_authenticated, user_email, _, err := getSesionData(r); err != nil {
		http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError)
	} else if already_authenticated {
		users, err := apiCallGet[models.User]("krsnc_usrs", "mail", user_email, r)
		if err != nil {
			log.Println("nakon user api")
			smtu(w, r, l(r, 7, err))
		} else if data, err := apiCallGet[models.UserData]("data", "mail", user_email, r); err != nil {
			log.Println("nakon user data api")
			smtu(w, r, l(r, 7, err))
		} else {
			dashboard.Dashboard(r, users[0], data[0]).Render(r.Context(), w)
		}
	} else {
		smtu(w, r, "UnWelcome")
	}
}

func Sign_out(w http.ResponseWriter, r *http.Request) {
	key := i18n.SessionKey("session")
	session, err := i18n.GetSessionFromContext(r.Context(), key)
	if err != nil {
		l(r, 8, err)
	}
	session.Values["authenticated"] = false
	session.Values["user_email"] = ""
	session.Values["name"] = ""
	session.Values["mode"] = ""
	if err := session.Save(r, w); err != nil {
		l(r, 8, err)
	}
	views.Index(r).Render(r.Context(), w)
}

////**** FIRMA

func Fakture(w http.ResponseWriter, r *http.Request) {
	if already_authenticated, email, _, err := getSesionData(r); err == nil && !already_authenticated {
		views.Index(r).Render(r.Context(), w)
	} else {
		if data, err := apiCallGet[models.UserData]("data", "mail", email, r); err != nil {
			smtu(w, r, l(r, 7, err))
		} else {
			firma.Fakture(r, data[0]).Render(r.Context(), w)
		}
	}
}

func Kpo(w http.ResponseWriter, r *http.Request) {
	if already_authenticated, email, _, err := getSesionData(r); err == nil && !already_authenticated {
		views.Index(r).Render(r.Context(), w)
	} else {
		if data, err := apiCallGet[models.UserData]("data", "mail", email, r); err != nil {
			smtu(w, r, l(r, 7, err))
		} else {
			firma.Kpo(r, data[0]).Render(r.Context(), w)
		}
	}
}

func Zurnal(w http.ResponseWriter, r *http.Request) {
	if already_authenticated, email, _, err := getSesionData(r); err == nil && !already_authenticated {
		views.Index(r).Render(r.Context(), w)
	} else {
		if data, err := apiCallGet[models.UserData]("data", "mail", email, r); err != nil {
			smtu(w, r, l(r, 7, err))
		} else {
			firma.Zurnal(r, data[0]).Render(r.Context(), w)
		}
	}
}

func Klijenti(w http.ResponseWriter, r *http.Request) {
	if already_authenticated, email, _, err := getSesionData(r); err == nil && !already_authenticated {
		views.Index(r).Render(r.Context(), w)
	} else {
		if data, err := apiCallGet[models.UserData]("data", "mail", email, r); err != nil {
			smtu(w, r, l(r, 7, err))
		} else {
			firma.Klijenti(r, data[0]).Render(r.Context(), w)
		}
	}
}

func Artikli(w http.ResponseWriter, r *http.Request) {
	if already_authenticated, email, _, err := getSesionData(r); err == nil && !already_authenticated {
		views.Index(r).Render(r.Context(), w)
	} else {
		if data, err := apiCallGet[models.UserData]("data", "mail", email, r); err != nil {
			smtu(w, r, l(r, 7, err))
		} else {
			firma.Artikli(r, data[0]).Render(r.Context(), w)
		}
	}
}

////**** i18n

func SetSh(w http.ResponseWriter, r *http.Request) {
	key := i18n.SessionKey("session")
	session, err := i18n.GetSessionFromContext(r.Context(), key)
	if err != nil {
		http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError)
		return
	}
	session.Values["language"] = "sh"
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError)
		return
	}
}

func SetEn(w http.ResponseWriter, r *http.Request) {
	key := i18n.SessionKey("session")
	session, err := i18n.GetSessionFromContext(r.Context(), key)
	if err != nil {
		http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError)
		return
	}
	session.Values["language"] = "en"
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError)
		return
	}
}

func SetEs(w http.ResponseWriter, r *http.Request) {
	key := i18n.SessionKey("session")
	session, err := i18n.GetSessionFromContext(r.Context(), key)
	if err != nil {
		http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError)
		return
	}
	session.Values["language"] = "es"
	err = session.Save(r, w)
	if err != nil {
		http.Error(w, l(r, 8, err).Error(), http.StatusInternalServerError)
		return
	}
}

// func SetBrowserLang(w http.ResponseWriter, r *http.Request) {
// 	key := i18n.SessionKey("session")
//	session, err := i18n.GetSessionFromContext(r.Context(), key)
// 	if err != nil {
// 		// fmt.Println("browser greška get sessio")
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	session.Values["language"] = ""
// 	err2 := session.Save(r, w)
// 	if err2 != nil {
// 		// fmt.Println("brower greška save sessio")
// 		http.Error(w, err2.Error(), http.StatusInternalServerError)
// 		return
// 	}
// }

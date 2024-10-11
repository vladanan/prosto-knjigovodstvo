// Package routes služi da obrađuje zahvete iz main
package routes

import (
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rs/cors"

	"github.com/vladanan/prosto/src/controllers/api"
	"github.com/vladanan/prosto/src/controllers/clr"

	"github.com/vladanan/prosto/src/models"
)

// func check(e error) {
// 	if e != nil {
// 		panic(e)
// 	}
// }

func RouterCRUD(s *mux.Router) {
	// CRUD
	vh := api.NewVezbamoHandler(models.DB{})
	s.HandleFunc("/v/{table}", clr.CheckFunc(vh.HandlePostOne)).Methods("POST")

	s.HandleFunc("/v/{table}/{field}/{record}", clr.CheckFunc(vh.HandleGet)).Methods("GET")
	s.HandleFunc("/v/{table}", clr.CheckFunc(vh.HandleGet)).Methods("GET")

	s.HandleFunc("/v/{table}/{field}/{record}", clr.CheckFunc(vh.HandlePutOne)).Methods("PUT")
	s.HandleFunc("/v/{table}/{field}/{record}", clr.CheckFunc(vh.HandleDeleteOne)).Methods("DELETE")

	ch := api.NewEoneHandler(models.DB{})
	s.HandleFunc("/c/eone/billing", clr.CheckFunc(ch.HandleGetBilling))

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:4200"},
	})
	s.Handle("/locations", c.Handler(http.HandlerFunc(GetLocationsForAngularFE)))
	s.Handle("/locations/", c.Handler(http.HandlerFunc(GetLocationsForAngularFE)))
}

func RouterAuth(s *mux.Router) {
	s.HandleFunc("/sign_up_post", Sign_up_post)
	s.HandleFunc("/sign_in_post", Sign_in_post)
	s.HandleFunc("/auto_login_demo", AutoLoginDemo)
	s.HandleFunc("/auto_login_user", AutoLoginUser)
	s.HandleFunc("/auto_login_admin", AutoLoginAdmin)
	s.HandleFunc("/forgotten_password_send_email", ForgottenPasswordSendMail)
	s.HandleFunc("/delete_user", DeleteUser)
	s.HandleFunc("/delete_user_send_email", DeleteUserSendMail)
	s.HandleFunc("/change_email", ChangeEmail)
	s.HandleFunc("/change_email_send_email", ChangeEmailSendMail)
	s.HandleFunc("/change_name", ChangeName)
	s.HandleFunc("/change_name_send_email", ChangeNameSendMail)
	s.HandleFunc("/change_password", ChangePassword)
	s.HandleFunc("/change_password_send_email", ChangePasswordSendMail)

	// samo query koji ima u sebi tačno određene promenljive može da prođe
	s.HandleFunc("/vmk/{key}", CheckLinkFromEmailRegister).Queries("user_email", "") // , "user", "vladan")
	// isto kao i ono gore:
	// vmk := r.PathPrefix("/vmk").Subrouter()
	// vmk.HandleFunc("/{key}", CheckLinkFromEmail).Queries("email", "")
	s.HandleFunc("/check_name/", CheckUserName).Queries("name", "")
	s.HandleFunc("/confirm/", SendNewSecretInputField)
	s.HandleFunc("/fp/{key}", CheckLinkFromEmailFP).Queries("fpm", "")
	s.HandleFunc("/du/{key}", CheckLinkFromEmailDU).Queries("dum", "")
	s.HandleFunc("/cm/{key}", CheckLinkFromEmailCM).Queries("cmm", "")
	s.HandleFunc("/cn/{key}", CheckLinkFromEmailCN).Queries("cnm", "")
	s.HandleFunc("/cp/{key}", CheckLinkFromEmailCP).Queries("cpm", "")

}

func GetLocationsForAngularFE(w http.ResponseWriter, r *http.Request) {
	fmt.Println("\nget locations", r.URL)
	dat, err := os.ReadFile("src/models/locations.json")
	if err != nil {
		l(r, 7, err)
	}
	// fmt.Println("dat: ", string(dat))

	w.Header().Set("Content-Type", "application/json")

	if _, err = io.WriteString(w, string(dat)); err != nil {
		w.Write([]byte(l(r, 7, err).Error()))
	}

}

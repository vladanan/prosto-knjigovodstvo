package models

import (
	"encoding/base64"
	"net"
	"net/http"
	"time"

	"fmt"

	"github.com/vladanan/prosto/src/controllers/clr"
	"github.com/vladanan/prosto/src/controllers/utils"

	"os"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type SignUpFormData struct {
	User_email1 string
	User_email2 string
	Name        string
	Password1   string
	Password2   string
	Confirm     string
	Submit      string
	Wait        string
	Bk_         string
	Fk_         string
}

type SignInFormData struct {
	User_email1 string
	Password1   string
	Confirm     string
	Submit      string
	Wait        string
	Bk_         string
	Fk_         string
}

type ForgottenPasswordData struct {
	User_email1 string
	User_email2 string
	Confirm     string
	Submit      string
	Wait        string
	Bk_         string
	Fk_         string
}

type ChangeEmailData struct {
	User_email1    string
	NewUser_email1 string
	NewUser_email2 string
	Password1      string
	Confirm        string
	Submit         string
	Wait           string
	Bk_            string
	Fk_            string
}

type ChangeNameData struct {
	User_email1 string
	Name        string
	Password1   string
	Confirm     string
	Submit      string
	Wait        string
	Bk_         string
	Fk_         string
}

type ChangePasswordData struct {
	User_email1 string
	Password0   string
	Password1   string
	Password2   string
	Confirm     string
	Submit      string
	Wait        string
	Bk_         string
	Fk_         string
}

func GetVerifyKey2(text string, r *http.Request) (string, error) {
	// text se pretvarau byte [], pravi crypto key
	ciphertext, err := bcrypt.GenerateFromPassword([]byte(text), 12)
	if err != nil {
		return "", l(r, 8, err)
	}
	// kada se ciphertext koristi bez zamena / i . onda ne može da se koristi kao url jer / dovodi do toga da je url pogrešan
	// zato se ciphertext prebaci u base64 prilagodjen za url-ove
	return base64.RawURLEncoding.EncodeToString(ciphertext), nil
}

// Koristi se da se proveri da li vec postiji korisnik sa tim mejlom prilikom registracije na sajt
// ali i da se dobije user_name ako se ima email
// i ako mejl ne postoji onda nema greske i vraca naravno praznan User
// ali ako postoji onda vraca gresku da postoji i tog Usera
func CheckEmail(email string, r *http.Request) (User, error) {

	// BAZA, UZIMANJE USERA
	conn, err := getDBConn(r)
	if err != nil {
		return User{}, l(r, 8, err)
	}
	defer freeConn(conn, r)
	rows, err := conn.Query(r.Context(), "SELECT * FROM mi_users where email=$1;", email)
	if err != nil {
		return User{}, l(r, 8, err)
	}
	pgxUsers, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		return User{}, l(r, 8, err)
	}

	if len(pgxUsers) == 0 { // nema korisnika sa takvim user_name
		return User{}, nil
	} else {
		return pgxUsers[0], clr.NewAPIError(http.StatusNotAcceptable, "Email_already_exists")
	}

}

func CheckUserName(userName string, r *http.Request) error {

	// BAZA, UZIMANJE USERA
	conn, err := getDBConn(r)
	if err != nil {
		return l(r, 8, err)
	}
	defer freeConn(conn, r)
	rows, err := conn.Query(r.Context(), "SELECT * FROM mi_users where user_name=$1;", userName)
	if err != nil {
		return l(r, 8, err)
	}
	pgxUsers, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		return l(r, 8, err)
	}

	if len(pgxUsers) == 0 { // nema korisnika sa takvim user_name
		return nil
	} else {
		return clr.NewAPIError(http.StatusNotAcceptable, "User_name_already_exists")
	}

}

/*
ADD USER
*/
func AddUser(form SignUpFormData, r *http.Request) error {

	password := []byte(form.Password1)

	// KREIRANJE ŠIFRI ZA SIGN IN I ZA MAIL VERIFIKACIJU
	//https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-index
	//https://gowebexamples.com/password-hashing/
	// GenerateFromPassword does not accept passwords longer than 72 bytes, which is the longest password bcrypt will operate on.
	// nije dobro da se key za proveru mejla pravi na osnovu lozinke
	// da se ne bi desilo da neko proba da rekonstruiše password iz poslatog linka
	// najbolje samo iz mejla jer je mejl svakako već poznat onome ko ima link a novi link svakako ne može sam da generiše
	ciphertextSignIn, err := bcrypt.GenerateFromPassword(password, 12)
	if err != nil {
		return l(r, 8, err)
	}
	keyForVerifyLink := utils.Get64UrlKey(form.User_email1, 0)

	// DB
	conn, err := getDBConn(r)
	if err != nil {
		return l(r, 8, err)
	}
	defer freeConn(conn, r)
	// err := godotenv.Load(".env")
	// if err != nil {
	// 	return User{}, l(r, 8, err)
	// }

	// PROVERA DA LI ima istog takvog hedera u ostalim zapisima deset minuta unazad tj. koji imaju created_at deset minuta stariji od sadašnjeg trenutka
	// fmt.Print("AddUser: r header: ", to_map(bytearray_headers)["X-Forwarded-For"][0], "\n")
	// https://stackoverflow.com/questions/23320945/postgresql-select-if-string-contains
	// https://stackoverflow.com/questions/45849494/how-do-i-search-for-a-specific-string-in-a-json-postgres-data-type-column
	// razni sql pokušaji i varijante sa i bez promenljivih i jsonb polja:
	// SELECT id FROM TAG_TABLE WHERE position(tag_name in 'aaaaaaaaaaa')>0;
	// rows2, err := conn.Query(r.Context(), "SELECT * FROM mi_users where position($1 IN created_at_headers) > 0;", to_map(bytearray_headers)["X-Forwarded-For"][0])
	// rows2, err := conn.Query(r.Context(), `with  vars as (select '127.0.0.1' as var1) select * from  mi_users,  vars where jsonb_path_exists(created_at_headers,'$.X-Forwarded-For ? (@ == var1)');`, to_map(bytearray_headers)["X-Forwarded-For"][0])

	// Ne mogu da se koriste prepared statements niti sql promenljive u upitima za sadržaj X-Forwarded-For heder u jsonb polju (može sa string concatenation ali to je opasno) tako da mora prvo da se pokupe upisi u poslednjih 10 min i zatim da se kod svih uporedi X-Forwarded-For sa aktuelnim

	// UZIMANJE PROMENLJIVIH IZ ENV I DB ZA ATTEMPT TIME LIMIT
	s, err := GetEnvDbSettings(r)
	if err != nil {
		return l(r, 8, err)
	}

	// uzimanje svih usera koji imaju created_at_time u poslednjih n min
	rows, err := conn.Query(
		r.Context(),
		`select * from  mi_users where (now() :: timestamp - created_at_time) < interval `+s.Same_ip_sign_up_time_limit_string)
	// rows2, err := conn.Query(r.Context(), `select * from  mi_users where (now() :: timestamp - created_at_time) < interval '3m'`)
	if err != nil {
		return l(r, 8, err)
	}
	pgxUser, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		return l(r, 8, err)
	}

	if len(pgxUser) > 0 { // ako ima više od 0 usera koji imaju created_at_time u poslednjih n min onda se proverava ip

		for _, item := range pgxUser {

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				return l(r, 8, err)
			}
			// if HeadersToMap([]byte(item.Created_at_headers))["X-Forwarded-For"][0] == r.Header.Get("X-Forwarded-For") {
			if item.Req_ip == ip {

				l(r, 8, fmt.Errorf(
					"ima upis u posledjih %vmin i JESTE isti ip:%v",
					s.Same_ip_sign_up_time_limit_string,
					// HeadersToMap([]byte(item.Created_at_headers))["X-Forwarded-For"][0]))
					item.Req_ip))

				return clr.NewAPIError(http.StatusNotAcceptable, "Too_many_attempts")

			}
			// else {
			// 	clr.GetStringLogger()(r, 0, fmt.Sprintf(
			// 		"ima upis u posledjih %vmin ali NIJE isti ip:%v",
			// 		same_ip_sign_up_time_limit,
			// 		item.Req_ip))
			// 		HeadersToMap([]byte(item.Created_at_headers))["X-Forwarded-For"][0]))
			// }

		}

	}
	// else {
	// 	log.Println("nema upisa od pre:", same_ip_sign_up_time_limit)
	// }
	// }

	// budući da nema upisa sa istog ip u posledjih n minuta user se upisuje u db
	_, err = conn.Exec(r.Context(), `INSERT INTO mi_users
		(
			hash_lozinka,
			email,
			user_name,
			user_mode,
			user_level,
			basic,
			js,
			c,
			verified_email,
			created_at_headers
		)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);`,
		ciphertextSignIn,
		form.User_email1,
		form.Name,
		"user",
		0,
		true, true, false,
		keyForVerifyLink,
		r.Header,
	)
	if err != nil {

		return l(r, 8, err)

	} else {

		// log.Println("insert rezultat:", commandTag)
		// AKO JE SVE OKEJ ŠALJE SE EMAIL PREKO EMAIL KLIJENTA

		var urlDomainForEmail string
		if os.Getenv("PRODUCTION") == "FALSE" {
			urlDomainForEmail = "http://127.0.0.1:7331/auth/vmk/" + keyForVerifyLink + "?user_email=" + form.User_email1 //"vladan_zasve@yahoo.com"
		} else {
			urlDomainForEmail = "https://prosto-knjigovodstvo.onrender.com/auth/vmk/" + keyForVerifyLink + "?user_email=" + form.User_email1
		}

		email := utils.MailRegister{
			R:        r,
			Email:    form.User_email1,
			UserName: form.Name,
			Url:      urlDomainForEmail,
		}
		if err := utils.SendEmail(email); err != nil {
			return l(r, 4, err)
		}

		return nil

	}

}

/*
AUTHENTICATE EMAIL
*/
func VerifyEmail(key, email string, r *http.Request) error {

	conn, err := getDBConn(r)
	if err != nil {
		return l(r, 8, err)
	}
	defer freeConn(conn, r)

	// proverava se i da li mejl i 64 url key odgovaraju jedno drugom
	if err = utils.Check64UrlKey(key, email, 0); err != nil {
		return clr.NewAPIError(http.StatusNotAcceptable, "Email_link_not_verified")
	}

	// proverava se i da li takav par postoji u bazi
	rows, err := conn.Query(r.Context(), "SELECT * FROM mi_users where verified_email=$1 and email=$2;", key, email)
	if err != nil {
		return l(r, 8, err)
	}

	pgxKey, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		return l(r, 8, err)
	}

	if len(pgxKey) > 0 {
		_, err := conn.Exec(r.Context(), `UPDATE mi_users SET verified_email=$1 where verified_email=$2;`,
			"verified",
			pgxKey[0].Verified_email,
		)
		if err != nil {
			return l(r, 8, err)
		} else {
			// log.Printf("insert result: %v\n", commandTag)
			return nil
		}

	} else {
		return clr.NewAPIError(http.StatusNotAcceptable, "Email_link_not_verified")
	}

}

/*
AUTHENTICATE USER
*/
func AuthenticateUser(email string, passwordStr string, r *http.Request) (User, error) {

	password := []byte(passwordStr)

	// BAZA, UZIMANJE USERA
	conn, err := getDBConn(r)
	if err != nil {
		return User{}, l(r, 8, err)
	}
	defer freeConn(conn, r)

	rows, err := conn.Query(r.Context(), "SELECT * FROM mi_users where email=$1;", email)
	if err != nil {
		return User{}, l(r, 8, err)
	}
	pgxUsers, err := pgx.CollectRows(rows, pgx.RowToStructByName[User])
	if err != nil {
		return User{}, l(r, 8, err)
	}

	// UZIMANJE PROMENLJIVIH IZ ENV I DB ZA BAD ATTEMPT LIMITE
	s, err := GetEnvDbSettings(r)
	if err != nil {
		return User{}, l(r, 8, err)
	}

	var user User
	if len(pgxUsers) == 0 { // nema korisnika sa takvim mejlom ali se to ne odaje nego se piše i lozinka
		return User{}, clr.NewAPIError(http.StatusNotAcceptable, "Email_or_password_wrong") //, email, password_str
	} else {
		user = pgxUsers[0]
	}

	if user.Verified_email != "verified" { // ako ima mejla proverava se verifikacija

		return User{}, clr.NewAPIError(http.StatusNotAcceptable, "Email_not_verified")

		//
	} else if int64(user.Bad_sign_in_attempts) < s.Bad_sign_in_attempts_limit {
		// mejl je verifikovan i ide se na proveru broja neuspelih pokušaja:
		// ako je broj neuspelih pokušaja manji od limita upisuje se pokušaj i ide se na proveru lozinke

		_, err := conn.Exec(r.Context(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`,
			user.Bad_sign_in_attempts+1,
			email,
		)
		if err != nil {
			return User{}, l(r, 8, err)
		}
		// log.Println("added 1 to current bad sign in attempts:", user.Bad_sign_in_attempts)

		_, err = conn.Exec(r.Context(), `UPDATE mi_users SET bad_sign_in_time=$1 where email=$2;`,
			time.Now(),
			email)
		if err != nil {
			return User{}, l(r, 8, err)
		}
		// log.Println("last bad sign time has been set")

		// https://pkg.go.dev/golang.org/x/crypto/bcrypt#pkg-index
		// https://gowebexamples.com/password-hashing/
		// provera lozinke:
		err = bcrypt.CompareHashAndPassword([]byte(user.Hash_lozinka), password)

		if err != nil { // LOŠA LOZINKA

			return User{}, clr.NewAPIError(http.StatusNotAcceptable, "Wrong_password")

		} else { // SVE JE OKEJ (ver, mejl, pokušaji, pass) UPISUJE SE U BAZU SVE ŠTO TREBA I PODACI ŠALJU RUTERU

			_, err = conn.Exec(r.Context(), `UPDATE mi_users SET last_sign_in_time=$1 where email=$2;`, time.Now(), email)
			if err != nil {
				return User{}, l(r, 8, err)
			}

			_, err = conn.Exec(r.Context(), `UPDATE mi_users SET last_sign_in_headers=$1 where email=$2;`, r.Header, email)
			if err != nil {
				return User{}, l(r, 8, err)
			}

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				return User{}, l(r, 8, err)
			}
			_, err = conn.Exec(r.Context(), `UPDATE mi_users SET req_ip=$1 where email=$2;`, ip, email)
			if err != nil {
				return User{}, l(r, 8, err)
			}

			_, err = conn.Exec(r.Context(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`, 0, email)
			if err != nil {
				return User{}, l(r, 8, err)
			}
			// log.Println("zeroed bad sign in attempts:", user.Bad_sign_in_attempts)
			// log.Println("authenticate user: prošlo je")
			user.Hash_lozinka = ""
			return user, nil
		}

	} else { // ako je broj neuselih pokušaja veći od limita gleda se da li je prošlo više vremena od limita

		if time.Since(user.Bad_sign_in_time).Minutes() < float64(s.Bad_sign_in_time_limit) { //

			// fmt.Print("previše pokušaja za sign in, pokušati za minuta: ", float64(bad_sign_in_time_limit)-time.Since(user.Bad_sign_in_time).Minutes())
			return User{}, clr.NewAPIError(http.StatusNotAcceptable, "Too_many_attempts")

			//
		} else { // kada je prošlo dovoljno vremena resetuje se broj neuspelih pokušaja

			_, err := conn.Exec(r.Context(), `UPDATE mi_users SET bad_sign_in_attempts=$1 where email=$2;`, 0, email)
			if err != nil {
				return User{}, l(r, 8, err)
			}
			// log.Println("time is up, zeroing bad sign in attempts:", user.Bad_sign_in_attempts)

			return User{}, clr.NewAPIError(http.StatusNotAcceptable, "Sign_in_open")

		}
	}

}

/*
FORGOT PASSWORD
*/
func ReplaceForgottenPassword(email, userName, newPassword, newPasswordKey string, r *http.Request) error {

	// log.Println("db new pass i key", newPassword, newPasswordKey)

	// DB
	conn, err := getDBConn(r)
	if err != nil {
		return l(r, 8, err)
	}
	defer freeConn(conn, r)

	_, err = conn.Exec(r.Context(), `UPDATE mi_users SET hash_lozinka=$1, updated_at=$2 WHERE email=$3;`, newPasswordKey, time.Now(), email)
	if err != nil {
		return l(r, 8, err)
	} else {

		// log.Println("insert rezultat:", commandTag)
		// AKO JE SVE OKEJ ŠALJE SE EMAIL PREKO EMAIL KLIJENTA
		var urlDomainForSignIn string
		if os.Getenv("PRODUCTION") == "FALSE" {
			urlDomainForSignIn = "http://127.0.0.1:7331/sign_in"
		} else {
			urlDomainForSignIn = "https://prosto-knjigovodstvo.onrender.com/sign_in"
		}

		email := utils.MailFPSendNewPassword{
			R:        r,
			Email:    email,
			UserName: userName,
			Password: newPassword,
			Url:      urlDomainForSignIn,
		}
		if err := utils.SendEmail(email); err != nil {
			return l(r, 4, err)
		}

		return nil

	}

}

/*
DELETE USER
*/
func DeleteUser(email, userName string, r *http.Request) error {

	// DB
	conn, err := getDBConn(r)
	if err != nil {
		return l(r, 8, err)
	}
	defer freeConn(conn, r)

	_, err = conn.Exec(r.Context(), `DELETE FROM mi_users WHERE email=$1;`, email)
	if err != nil {
		return l(r, 8, err)
	} else {

		// log.Println("insert rezultat:", commandTag)
		// AKO JE SVE OKEJ ŠALJE SE EMAIL da je user izbrisan

		email := utils.DeletedUser{
			R:        r,
			Email:    email,
			UserName: userName,
		}
		if err := utils.SendEmail(email); err != nil {
			return l(r, 4, err)
		}

		return nil

	}

}

/*
CHANGE EMAIL
*/
func ChangeEmail(email, newEmail, userName string, r *http.Request) error {

	// DB
	conn, err := getDBConn(r)
	if err != nil {
		return l(r, 8, err)
	}
	defer freeConn(conn, r)

	keyForVerifyLink := utils.Get64UrlKey(newEmail, 0)

	_, err = conn.Exec(r.Context(), `UPDATE mi_users SET email=$1, updated_at=$2, verified_email=$3 WHERE email=$4;`, newEmail, time.Now(), keyForVerifyLink, email)
	if err != nil {
		return l(r, 8, err)
	} else {

		// log.Println("insert rezultat:", commandTag)
		// AKO JE SVE OKEJ ŠALJE SE EMAIL da je mejl promenjen

		email := utils.ChangedEmail{
			R:        r,
			Email:    email,
			UserName: userName,
		}
		if err := utils.SendEmail(email); err != nil {
			return l(r, 4, err)
		}

		// i salje se mejl na novi emejl da se verifikuje
		var urlDomainForEmail string
		if os.Getenv("PRODUCTION") == "FALSE" {
			urlDomainForEmail = "http://127.0.0.1:7331/auth/vmk/" + keyForVerifyLink + "?user_email=" + newEmail
		} else {
			urlDomainForEmail = "https://prosto-knjigovodstvo.onrender.com/auth/vmk/" + keyForVerifyLink + "?user_email=" + newEmail
		}

		verifyEmail := utils.MailRegister{
			R:        r,
			Email:    newEmail,
			UserName: userName,
			Url:      urlDomainForEmail,
		}
		if err := utils.SendEmail(verifyEmail); err != nil {
			return l(r, 4, err)
		}

		return nil

	}

}

/*
CHANGE NAME
*/
func ChangeName(email, userName, newUserName string, r *http.Request) error {

	// DB
	conn, err := getDBConn(r)
	if err != nil {
		return l(r, 8, err)
	}
	defer freeConn(conn, r)

	// keyForVerifyLink := utils.Get64UrlKey(email, 10)

	_, err = conn.Exec(r.Context(), `UPDATE mi_users SET user_name=$1, updated_at=$2 WHERE email=$3 and user_name=$4;`, newUserName, time.Now(), email, userName)
	if err != nil {
		return l(r, 8, err)
	} else {
		return nil
	}

}

/*
CHANGE PASSWORD
*/
func ChangePassword(email, password0, password1 string, r *http.Request) error {

	// DB
	conn, err := getDBConn(r)
	if err != nil {
		return l(r, 8, err)
	}
	defer freeConn(conn, r)

	if rows, err := conn.Query(r.Context(), "SELECT * FROM mi_users where email=$1;", email); err != nil {
		return l(r, 8, err)
	} else if pgxUsers, err := pgx.CollectRows(rows, pgx.RowToStructByName[User]); err != nil {
		return l(r, 8, err)
	} else if len(pgxUsers) == 0 { // nema korisnika sa takvim mejlom ali se to ne odaje nego se piše i lozinka
		return clr.NewAPIError(http.StatusNotAcceptable, "Email_or_password_wrong") //, email, password_str
	} else {
		// provera lozinke:
		err = bcrypt.CompareHashAndPassword([]byte(pgxUsers[0].Hash_lozinka), []byte(password0))
		if err != nil {
			return clr.NewAPIError(http.StatusNotAcceptable, "Wrong_password")
		} else {
			ciphertext, err := bcrypt.GenerateFromPassword([]byte(password1), 12)
			if err != nil {
				return l(r, 8, err)
			}
			_, err = conn.Exec(r.Context(), `UPDATE mi_users SET hash_lozinka=$1, updated_at=$2 WHERE email=$3;`, ciphertext, time.Now(), email)
			if err != nil {
				return l(r, 8, err)
			} else {
				return nil
			}
		}
	}

}

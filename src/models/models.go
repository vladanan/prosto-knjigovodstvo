package models

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vladanan/prosto/src/controllers/clr"
)

type DB struct{}

var l = clr.GetErrorLogger()

// Šalje konekciju uz pomoć odabrane funkcije, bilo local pool, bilo direktno lokal ili remote
func getDBConn(r *http.Request) (*pgx.Conn, error) {

	// un-komentuje se jedna od dve funkcije koja se hoće i menja se tip konekcije u return definiciji

	return getDBConnFromDirectLocalOrRemotePool(r) //pgx.Conn

	// if r == nil {
	// 	fmt.Println("")
	// }

	// return nil, nil

	// return getLocalPoolDBConnFromRequestContext(r) //pgxpool.Conn

}

// Pribavlja db konekciju bilo iz lokalne baze bilo udaljene supa (koja ima svoj pool)
// ako je u dev onda sa lokalne baze a ako je na prod onda sa supa
// STAVLJENO JE DA I U DEV I U PROD BAZA IDE IZ SUPA DA NE BI POSTAVLJAO PG DB
// NA MALOM ACERU I DA BI PRATIO REALNU BRZINU RADA PROD DB
// A KASNIJE CE MOZDA BITI POTREBNO DA SE RADI MIGRACIJA I SYNC NA DEV
// KADA PROJEKAT ODMAKNE I NEMA ZEZANJA SA DB
func getDBConnFromDirectLocalOrRemotePool(r *http.Request) (*pgx.Conn, error) {

	// zahtevi za supa sa dev treba oko 1100-1200ms ali se izvrsavaju i sa ogranicenjem na 1000 sto znaci da ne ide sve na sql nego nesto i na sever i browser, padaju kada je ogranicenje na 900ms, stavio sam za svaki slucaj 1200
	// stavljeno 2200 jer su sa malog acera pucali zahtevi zbog sporosti kao i zbog sporosti odgovora sa supabase
	ctx, cancel := context.WithTimeout(r.Context(), 2500*time.Millisecond)
	defer cancel()

	// log.Println("db get conn from local db or from supa, production:", os.Getenv("PRODUCTION"))
	if os.Getenv("PRODUCTION") == "FALSE" {
		conn, err := pgx.Connect(ctx, os.Getenv("FEDORA_CONNECTION_STRING"))
		// conn, err := pgx.Connect(ctx, os.Getenv("SUPABASE_CONNECTION_STRING"))
		if err != nil {
			return nil, l(r, 8, err)
		}
		return conn, nil
	} else {
		conn, err := pgx.Connect(ctx, os.Getenv("SUPABASE_CONNECTION_STRING"))
		if err != nil {
			return nil, l(r, 8, err)
		}
		return conn, nil
	}
}

// posredni tip za key, value pair u context Value
type ConnKey string

// Pribavlja pointer db konekcije iz lokalnog poola kroz request context
// ako je pointer za konekciju ubačen kroz middleware za api sub-ruter u main
// u main mora da se un-komentuju ove dve linije:
// * pool = createLocalPool()
// * apirouter.Use(insertLocalPoolConnMiddleware)
func getLocalPoolDBConnFromRequestContext(r *http.Request) (*pgxpool.Conn, error) {
	// log.Println("db get conn from context pool")
	f := func(ctx context.Context, k ConnKey) (*pgxpool.Conn, error) {
		if v := ctx.Value(k); v != nil {
			// log.Println("found value:", v)
			if p, ok := v.(*pgxpool.Conn); ok {
				// log.Println("jeste pointer za pool conn")
				return p, nil
			} else {
				return nil, fmt.Errorf("poslat je pointer koji nije tip pgxpool.Conn")
			}
		} else {
			return nil, fmt.Errorf("request context key nije validan: %v", k)
		}
	}

	key := ConnKey("pool_conn")
	conn, err := f(r.Context(), key)
	if err != nil {
		return nil, l(r, 8, err)
	}
	return conn, nil
}

// gasi db konekciju u zavisnosti kog je tipa
func freeConn(conn any, r *http.Request) {
	switch c := conn.(type) {
	case *pgxpool.Conn:
		c.Release()
	case *pgx.Conn:
		c.Close(r.Context())
	default:
		l(r, 8, fmt.Errorf("no db connections released/closed!"))
	}
}

type User struct {
	U_id                 int       `db:"u_id"`
	Created_at_time      time.Time `db:"created_at_time"`
	Hash_lozinka         string    `db:"hash_lozinka"`
	Email                string    `db:"email"`
	User_name            string    `db:"user_name"`
	Mode                 string    `db:"user_mode"`
	Level                string    `db:"user_level"`
	Basic                bool      `db:"basic"`
	Js                   bool      `db:"js"`
	C                    bool      `db:"c"`
	Payment_date         time.Time `db:"payment_date"`
	Payment_expire       time.Time `db:"payment_expire"`
	Payment_amount       int       `db:"payment_amount"`
	Payment_currency     string    `db:"payment_currency"`
	Verified_email       string    `db:"verified_email"`
	Last_sign_in_time    time.Time `db:"last_sign_in_time"`
	Last_sign_in_headers string    `db:"last_sign_in_headers"`
	Req_ip               string    `db:"req_ip"`
	Created_at_headers   string    `db:"created_at_headers"`
	Bad_sign_in_attempts int       `db:"bad_sign_in_attempts"`
	Bad_sign_in_time     time.Time `db:"bad_sign_in_time"`
	Updated_at           time.Time `db:"updated_at"`
}

type Settings struct {
	S_id                              int       `db:"s_id"`
	Updated_at                        time.Time `db:"updated_at"`
	Bad_sign_in_attempts_limit        int64     `db:"bad_sign_in_attempts_limit"`
	Bad_sign_in_time_limit            int64     `db:"bad_sign_in_time_limit"`
	Same_ip_sign_up_time_limit        int64     `db:"same_ip_sign_up_time_limit"`
	Same_ip_sign_up_time_limit_string string    `db:"-"`
}

func GetEnvDbSettings(r *http.Request) (Settings, error) {
	var s Settings
	conn, err := getDBConn(r)
	if err != nil {
		return Settings{}, l(r, 8, err)
	}
	defer freeConn(conn, r)

	rows, err := conn.Query(r.Context(), "SELECT * FROM v_settings where s_id=1;")
	if err != nil {
		return Settings{}, l(r, 8, err)
	}
	pgxSettings, err := pgx.CollectRows(rows, pgx.RowToStructByName[Settings])
	if err != nil {
		return Settings{}, l(r, 8, err)
	}

	var same_ip_sign_up_time_limit = "2m"
	SAME_IP_SIGN_UP_TIME_LIMIT := os.Getenv("SAME_IP_SIGN_UP_TIME_LIMIT")
	if SAME_IP_SIGN_UP_TIME_LIMIT == "" || SAME_IP_SIGN_UP_TIME_LIMIT == "0" {
		SAME_IP_SIGN_UP_TIME_LIMIT = "0m"
	}
	db_same_ip_sign_up_time_limit := strconv.Itoa(int(pgxSettings[0].Same_ip_sign_up_time_limit))

	if db_same_ip_sign_up_time_limit == "" || db_same_ip_sign_up_time_limit == "0" {
		db_same_ip_sign_up_time_limit = "0m"
	}
	if SAME_IP_SIGN_UP_TIME_LIMIT != "0m" {
		same_ip_sign_up_time_limit = "'" + SAME_IP_SIGN_UP_TIME_LIMIT + "m'"
	} else if db_same_ip_sign_up_time_limit != "0m" {
		same_ip_sign_up_time_limit = "'" + db_same_ip_sign_up_time_limit + "m'"
	}
	s.Same_ip_sign_up_time_limit_string = same_ip_sign_up_time_limit

	var bad_sign_in_attempts_limit int64 = 2
	var bad_sign_in_time_limit int64 = 8
	BAD_SIGN_IN_ATTEMPTS_LIMIT, err := strconv.ParseInt(os.Getenv("BAD_SIGN_IN_ATTEMPTS_LIMIT"), 0, 8)
	if err != nil {
		BAD_SIGN_IN_ATTEMPTS_LIMIT = 0
	}
	BAD_SIGN_IN_TIME_LIMIT, err := strconv.ParseInt(os.Getenv("BAD_SIGN_IN_TIME_LIMIT"), 0, 8)
	if err != nil {
		BAD_SIGN_IN_TIME_LIMIT = 0
	}

	db_bad_sign_in_attempts_limit := pgxSettings[0].Bad_sign_in_attempts_limit
	db_bad_sign_in_time_limit := pgxSettings[0].Bad_sign_in_time_limit
	if BAD_SIGN_IN_ATTEMPTS_LIMIT != 0 {
		bad_sign_in_attempts_limit = BAD_SIGN_IN_ATTEMPTS_LIMIT
	} else if db_bad_sign_in_attempts_limit != 0 {
		bad_sign_in_attempts_limit = db_bad_sign_in_attempts_limit
	}
	if BAD_SIGN_IN_TIME_LIMIT != 0 {
		bad_sign_in_time_limit = BAD_SIGN_IN_TIME_LIMIT
	} else if db_bad_sign_in_time_limit != 0 {
		bad_sign_in_time_limit = db_bad_sign_in_time_limit
	}

	s.Bad_sign_in_attempts_limit = bad_sign_in_attempts_limit
	s.Bad_sign_in_time_limit = bad_sign_in_time_limit

	return s, nil
}

type UserData struct {
	Ud_id                  int          `db:"ud_id"`
	U_id                   int          `db:"u_id"`
	U_email                string       `db:"u_email"`
	Obveznik               string       `db:"obveznik"`
	Sediste                string       `db:"sediste"`
	SifraPoreskogObveznika string       `db:"sifra_poreskog_pbveznika"`
	SifraDelatnosti        string       `db:"sifra_delatnosti"`
	Firma                  Firma        `db:"firma"`
	Settings               UserSettings `db:"settings"`
	Created_at             time.Time    `db:"created_at"`
	Updated_at             time.Time    `db:"updated_at"`
}
type UserSettings struct {
	LogoPath string `json:"logo_path"`
}

type Firma struct {
	Fi_id            int    `db:"fi_id"`
	U_id             int    `db:"u_id"`
	NazivFirmaRadnje string `db:"naziv_firma_radnje"`
	PIB              int    `db:"pib"`
	MB               int    `db:"mb"`
	Tr               string `db:"tr"`
	Adresa           string `db:"adresa"`
	Fiksni           string `db:"fiksni"`
	Mobilni          string `db:"mobilni"`
	Email            string `db:"email"`
	Link             string `db:"link"`
}

type Faktura struct {
	Fa_id      int           `db:"fa_id"`
	U_id       int           `db:"u_id"`
	Podaci     PodaciFakture `db:"podaci_fakture"`
	Created_at time.Time     `db:"created_at"`
	Zakljucena bool          `db:"zakljucena"`
}

type PodaciFakture struct {
	User          Firma     `json:"user"`
	Klijent       Firma     `json:"klijent"`
	DatumRacuna   time.Time `json:"datum_racuna"`
	MestoRacuna   string    `json:"mesto_racuna"`
	DatumPrometa  time.Time `json:"datum_prometa"`
	MestoPrometa  string    `json:"mesto_prometa"`
	Oznaka        string    `json:"oznaka"`
	Valuta        time.Time `json:"valuta"`
	Stavke        []Stavka  `json:"stavke"`
	Specifikacija string    `json:"specifikacija"`
	Porez         string    `json:"porez"`
	Instrukcije   string    `json:"instrukcije"`
}

type Stavka struct {
	Rb     int    `json:"rb"`
	Artikl Artikl `json:"artikl"`
	Kom    int    `json:"kom"`
}

type Artikl struct {
	A_id  int    `db:"a_id"`
	U_id  int    `db:"u_id"`
	Sifra int    `db:"sifra"`
	Naziv string `db:"naziv"`
	Tip   string `db:"tip"`
	Cena  int    `db:"cena"`
}

type Kpo struct {
	Kp_id     int       `db:"kp_id"`
	U_id      int       `db:"u_id"`
	Rb        int       `db:"rb"`
	Datum     time.Time `db:"datum"`
	Opis      string    `db:"opis"`
	Proizvodi int       `db:"proizvodi"`
	Usluge    int       `db:"usluge"`
}

type Zurnal struct {
	Z_id         int       `db:"z_id"`
	U_id         int       `db:"u_id"`
	Datum        time.Time `db:"datum"`
	KlijentNaziv string    `db:"klijent_naziv"`
	Opis         string    `db:"opis"`
}

// ///////////////////////////////////////////////////////////
type UserData1 struct {
	Ud_id      int       `db:"ud_id"`
	U_email    string    `db:"u_email"`
	Klijenti   string    `db:"klijenti"`
	Artikli    string    `db:"artikli"`
	Fakture    string    `db:"fakture"`
	Kpo        string    `db:"kpo"`
	Zurnal     string    `db:"zurnal"`
	Firma      Firma1    `db:"firma"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
}
type Firma1 struct {
	PIB                    int    `json:"pib"`
	Obveznik               string `json:"obveznik"`
	FirmaRadnje            string `json:"firmaRadnje"`
	Sediste                string `json:"sediste"`
	SifraPoreskogObveznika string `json:"sifra_poreskog_pbveznika"`
	SifraDelatnosti        string `json:"sifra_delatnosti"`
}
type Faktura1 struct {
	Id      string    `json:"id"`
	Datum   time.Time `json:"datum"`
	Klijent string    `json:"klijent"`
	Stavke  []Stavka  `json:"stavke"`
}
type Klijent1 struct {
	Naziv   string `json:"naziv"`
	Adresa  string `json:"adresa"`
	Telefon string `json:"telefon"`
	Email   string `json:"email"`
}

type Billing struct {
	Id                 int32     `db:"id"`
	Client_id          int       `db:"client_id"`
	Client_name        string    `db:"client_name"`
	Messages_sent      int       `db:"messages_sent"`
	Charge_per_message float32   `db:"charge_per_message"`
	Sms_cost           float32   `db:"sms_cost"`
	Create_time        time.Time `db:"create_time"`
}

type Note struct {
	B_id        int       `db:"b_id,omit"`
	Ime_tag     string    `db:"ime_tag"`
	Mejl        string    `db:"mejl"`
	Tema        string    `db:"tema"`
	Poruka      string    `db:"poruka"`
	User_id     string    `db:"user_id"`   // id usera ako je ulogovan
	User_mail   string    `db:"user_mail"` // mejl usera ako je ulogovan
	From_url    string    `db:"from_url"`  // ovde može da se upisuju headers ali će to možda da remeti next/react sajt
	Datum_upisa time.Time `db:"datum_upisa"`
}

type Test struct {
	G_id             int       `db:"g_id"`
	Tip              string    `db:"tip"`
	Obrazovni_profil any       `db:"obrazovni_profil"`
	Razred           any       `db:"razred"`
	Predmet          string    `db:"predmet"`
	Oblast           string    `db:"oblast"`
	Link1            any       `db:"link1"`
	Link2            any       `db:"link2"`
	Link3            any       `db:"link3"`
	User_id          any       `db:"user_id"`
	From_url         any       `db:"from_url"`
	Datum_upisa      time.Time `db:"datum_upisa"`
	Pitanje_1        any       `db:"pitanje_1"`
	Odg_1_1          any       `db:"odg_1_1"`
	Odg_1_2          any       `db:"odg_1_2"`
	Odg_1_3          any       `db:"odg_1_3"`
	Odg_1_4          any       `db:"odg_1_4"`
	R_1              any       `db:"r_1"`
	Pitanje_2        any       `db:"pitanje_2"`
	Odg_2_1          any       `db:"odg_2_1"`
	Odg_2_2          any       `db:"odg_2_2"`
	Odg_2_3          any       `db:"odg_2_3"`
	Odg_2_4          any       `db:"odg_2_4"`
	R_2              any       `db:"r_2"`
	Pitanje_3        any       `db:"pitanje_3"`
	Odg_3_1          any       `db:"odg_3_1"`
	Odg_3_2          any       `db:"odg_3_2"`
	Odg_3_3          any       `db:"odg_3_3"`
	Odg_3_4          any       `db:"odg_3_4"`
	R_3              any       `db:"r_3"`
	Pitanje_4        any       `db:"pitanje_4"`
	Odg_4_1          any       `db:"odg_4_1"`
	Odg_4_2          any       `db:"odg_4_2"`
	Odg_4_3          any       `db:"odg_4_3"`
	Odg_4_4          any       `db:"odg_4_4"`
	R_4              any       `db:"r_4"`
	Pitanje_5        any       `db:"pitanje_5"`
	Odg_5_1          any       `db:"odg_5_1"`
	Odg_5_2          any       `db:"odg_5_2"`
	Odg_5_3          any       `db:"odg_5_3"`
	Odg_5_4          any       `db:"odg_5_4"`
	R_5              any       `db:"r_5"`
	Pitanje_6        any       `db:"pitanje_6"`
	Odg_6_1          any       `db:"odg_6_1"`
	Odg_6_2          any       `db:"odg_6_2"`
	Odg_6_3          any       `db:"odg_6_3"`
	Odg_6_4          any       `db:"odg_6_4"`
	R_6              any       `db:"r_6"`
	Pitanje_7        any       `db:"pitanje_7"`
	Odg_7_1          any       `db:"odg_7_1"`
	Odg_7_2          any       `db:"odg_7_2"`
	Odg_7_3          any       `db:"odg_7_3"`
	Odg_7_4          any       `db:"odg_7_4"`
	R_7              any       `db:"r_7"`
	Pitanje_8        any       `db:"pitanje_8"`
	Odg_8_1          any       `db:"odg_8_1"`
	Odg_8_2          any       `db:"odg_8_2"`
	Odg_8_3          any       `db:"odg_8_3"`
	Odg_8_4          any       `db:"odg_8_4"`
	R_8              any       `db:"r_8"`
	Pitanje_9        any       `db:"pitanje_9"`
	Odg_9_1          any       `db:"odg_9_1"`
	Odg_9_2          any       `db:"odg_9_2"`
	Odg_9_3          any       `db:"odg_9_3"`
	Odg_9_4          any       `db:"odg_9_4"`
	R_9              any       `db:"r_9"`
	Pitanje_10       any       `db:"pitanje_10"`
	Odg_10_1         any       `db:"odg_10_1"`
	Odg_10_2         any       `db:"odg_10_2"`
	Odg_10_3         any       `db:"odg_10_3"`
	Odg_10_4         any       `db:"odg_10_4"`
	R_10             any       `db:"r_10"`
	Pitanje_11       any       `db:"pitanje_11"`
	Odg_11_1         any       `db:"odg_11_1"`
	Odg_11_2         any       `db:"odg_11_2"`
	Odg_11_3         any       `db:"odg_11_3"`
	Odg_11_4         any       `db:"odg_11_4"`
	R_11             any       `db:"r_11"`
	Pitanje_12       any       `db:"pitanje_12"`
	Odg_12_1         any       `db:"odg_12_1"`
	Odg_12_2         any       `db:"odg_12_2"`
	Odg_12_3         any       `db:"odg_12_3"`
	Odg_12_4         any       `db:"odg_12_4"`
	R_12             any       `db:"r_12"`
	Pitanje_13       any       `db:"pitanje_13"`
	Odg_13_1         any       `db:"odg_13_1"`
	Odg_13_2         any       `db:"odg_13_2"`
	Odg_13_3         any       `db:"odg_13_3"`
	Odg_13_4         any       `db:"odg_13_4"`
	R_13             any       `db:"r_13"`
	Pitanje_14       any       `db:"pitanje_14"`
	Odg_14_1         any       `db:"odg_14_1"`
	Odg_14_2         any       `db:"odg_14_2"`
	Odg_14_3         any       `db:"odg_14_3"`
	Odg_14_4         any       `db:"odg_14_4"`
	R_14             any       `db:"r_14"`
	Pitanje_15       any       `db:"pitanje_15"`
	Odg_15_1         any       `db:"odg_15_1"`
	Odg_15_2         any       `db:"odg_15_2"`
	Odg_15_3         any       `db:"odg_15_3"`
	Odg_15_4         any       `db:"odg_15_4"`
	R_15             any       `db:"r_15"`
	Pitanje_16       any       `db:"pitanje_16"`
	Odg_16_1         any       `db:"odg_16_1"`
	Odg_16_2         any       `db:"odg_16_2"`
	Odg_16_3         any       `db:"odg_16_3"`
	Odg_16_4         any       `db:"odg_16_4"`
	R_16             any       `db:"r_16"`
	Pitanje_17       any       `db:"pitanje_17"`
	Odg_17_1         any       `db:"odg_17_1"`
	Odg_17_2         any       `db:"odg_17_2"`
	Odg_17_3         any       `db:"odg_17_3"`
	Odg_17_4         any       `db:"odg_17_4"`
	R_17             any       `db:"r_17"`
	Pitanje_18       any       `db:"pitanje_18"`
	Odg_18_1         any       `db:"odg_18_1"`
	Odg_18_2         any       `db:"odg_18_2"`
	Odg_18_3         any       `db:"odg_18_3"`
	Odg_18_4         any       `db:"odg_18_4"`
	R_18             any       `db:"r_18"`
	Pitanje_19       any       `db:"pitanje_19"`
	Odg_19_1         any       `db:"odg_19_1"`
	Odg_19_2         any       `db:"odg_19_2"`
	Odg_19_3         any       `db:"odg_19_3"`
	Odg_19_4         any       `db:"odg_19_4"`
	R_19             any       `db:"r_19"`
	Pitanje_20       any       `db:"pitanje_20"`
	Odg_20_1         any       `db:"odg_20_1"`
	Odg_20_2         any       `db:"odg_20_2"`
	Odg_20_3         any       `db:"odg_20_3"`
	Odg_20_4         any       `db:"odg_20_4"`
	R_20             any       `db:"r_20"`
	Pitanje_21       any       `db:"pitanje_21"`
	Odg_21_1         any       `db:"odg_21_1"`
	Odg_21_2         any       `db:"odg_21_2"`
	Odg_21_3         any       `db:"odg_21_3"`
	Odg_21_4         any       `db:"odg_21_4"`
	R_21             any       `db:"r_21"`
	Pitanje_22       any       `db:"pitanje_22"`
	Odg_22_1         any       `db:"odg_22_1"`
	Odg_22_2         any       `db:"odg_22_2"`
	Odg_22_3         any       `db:"odg_22_3"`
	Odg_22_4         any       `db:"odg_22_4"`
	R_22             any       `db:"r_22"`
	Pitanje_23       any       `db:"pitanje_23"`
	Odg_23_1         any       `db:"odg_23_1"`
	Odg_23_2         any       `db:"odg_23_2"`
	Odg_23_3         any       `db:"odg_23_3"`
	Odg_23_4         any       `db:"odg_23_4"`
	R_23             any       `db:"r_23"`
	Pitanje_24       any       `db:"pitanje_24"`
	Odg_24_1         any       `db:"odg_24_1"`
	Odg_24_2         any       `db:"odg_24_2"`
	Odg_24_3         any       `db:"odg_24_3"`
	Odg_24_4         any       `db:"odg_24_4"`
	R_24             any       `db:"r_24"`
	Pitanje_25       any       `db:"pitanje_25"`
	Odg_25_1         any       `db:"odg_25_1"`
	Odg_25_2         any       `db:"odg_25_2"`
	Odg_25_3         any       `db:"odg_25_3"`
	Odg_25_4         any       `db:"odg_25_4"`
	R_25             any       `db:"r_25"`
	Pitanje_26       any       `db:"pitanje_26"`
	Odg_26_1         any       `db:"odg_26_1"`
	Odg_26_2         any       `db:"odg_26_2"`
	Odg_26_3         any       `db:"odg_26_3"`
	Odg_26_4         any       `db:"odg_26_4"`
	R_26             any       `db:"r_26"`
	Pitanje_27       any       `db:"pitanje_27"`
	Odg_27_1         any       `db:"odg_27_1"`
	Odg_27_2         any       `db:"odg_27_2"`
	Odg_27_3         any       `db:"odg_27_3"`
	Odg_27_4         any       `db:"odg_27_4"`
	R_27             any       `db:"r_27"`
	Pitanje_28       any       `db:"pitanje_28"`
	Odg_28_1         any       `db:"odg_28_1"`
	Odg_28_2         any       `db:"odg_28_2"`
	Odg_28_3         any       `db:"odg_28_3"`
	Odg_28_4         any       `db:"odg_28_4"`
	R_28             any       `db:"r_28"`
	Pitanje_29       any       `db:"pitanje_29"`
	Odg_29_1         any       `db:"odg_29_1"`
	Odg_29_2         any       `db:"odg_29_2"`
	Odg_29_3         any       `db:"odg_29_3"`
	Odg_29_4         any       `db:"odg_29_4"`
	R_29             any       `db:"r_29"`
	Pitanje_30       any       `db:"pitanje_30"`
	Odg_30_1         any       `db:"odg_30_1"`
	Odg_30_2         any       `db:"odg_30_2"`
	Odg_30_3         any       `db:"odg_30_3"`
	Odg_30_4         any       `db:"odg_30_4"`
	R_30             any       `db:"r_30"`
}
type Test2 struct {
	G_id             int8      `db:"g_id"`
	Tip              string    `db:"tip"`
	Obrazovni_profil string    `db:"obrazovni_profil"`
	Razred           string    `db:"razred"`
	Predmet          string    `db:"predmet"`
	Oblast           string    `db:"oblast"`
	Link1            string    `db:"link1"`
	Link2            string    `db:"link2"`
	Link3            string    `db:"link3"`
	User_id          string    `db:"user_id"`
	From_url         string    `db:"from_url"`
	Datum_upisa      time.Time `db:"datum_upisa"`
	Pitanje_1        string    `db:"pitanje_1"`
	Odg_1_1          string    `db:"odg_1_1"`
	Odg_1_2          string    `db:"odg_1_2"`
	Odg_1_3          string    `db:"odg_1_3"`
	Odg_1_4          string    `db:"odg_1_4"`
	R_1              string    `db:"r_1"`
	Pitanje_2        string    `db:"pitanje_2"`
	Odg_2_1          string    `db:"odg_2_1"`
	Odg_2_2          string    `db:"odg_2_2"`
	Odg_2_3          string    `db:"odg_2_3"`
	Odg_2_4          string    `db:"odg_2_4"`
	R_2              string    `db:"r_2"`
	Pitanje_3        string    `db:"pitanje_3"`
	Odg_3_1          string    `db:"odg_3_1"`
	Odg_3_2          string    `db:"odg_3_2"`
	Odg_3_3          string    `db:"odg_3_3"`
	Odg_3_4          string    `db:"odg_3_4"`
	R_3              string    `db:"r_3"`
	Pitanje_4        string    `db:"pitanje_4"`
	Odg_4_1          string    `db:"odg_4_1"`
	Odg_4_2          string    `db:"odg_4_2"`
	Odg_4_3          string    `db:"odg_4_3"`
	Odg_4_4          string    `db:"odg_4_4"`
	R_4              string    `db:"r_4"`
	Pitanje_5        string    `db:"pitanje_5"`
	Odg_5_1          string    `db:"odg_5_1"`
	Odg_5_2          string    `db:"odg_5_2"`
	Odg_5_3          string    `db:"odg_5_3"`
	Odg_5_4          string    `db:"odg_5_4"`
	R_5              string    `db:"r_5"`
	Pitanje_6        string    `db:"pitanje_6"`
	Odg_6_1          string    `db:"odg_6_1"`
	Odg_6_2          string    `db:"odg_6_2"`
	Odg_6_3          string    `db:"odg_6_3"`
	Odg_6_4          string    `db:"odg_6_4"`
	R_6              string    `db:"r_6"`
	Pitanje_7        string    `db:"pitanje_7"`
	Odg_7_1          string    `db:"odg_7_1"`
	Odg_7_2          string    `db:"odg_7_2"`
	Odg_7_3          string    `db:"odg_7_3"`
	Odg_7_4          string    `db:"odg_7_4"`
	R_7              string    `db:"r_7"`
	Pitanje_8        string    `db:"pitanje_8"`
	Odg_8_1          string    `db:"odg_8_1"`
	Odg_8_2          string    `db:"odg_8_2"`
	Odg_8_3          string    `db:"odg_8_3"`
	Odg_8_4          string    `db:"odg_8_4"`
	R_8              string    `db:"r_8"`
	Pitanje_9        string    `db:"pitanje_9"`
	Odg_9_1          string    `db:"odg_9_1"`
	Odg_9_2          string    `db:"odg_9_2"`
	Odg_9_3          string    `db:"odg_9_3"`
	Odg_9_4          string    `db:"odg_9_4"`
	R_9              string    `db:"r_9"`
	Pitanje_10       string    `db:"pitanje_10"`
	Odg_10_1         string    `db:"odg_10_1"`
	Odg_10_2         string    `db:"odg_10_2"`
	Odg_10_3         string    `db:"odg_10_3"`
	Odg_10_4         string    `db:"odg_10_4"`
	R_10             string    `db:"r_10"`
	Pitanje_11       string    `db:"pitanje_11"`
	Odg_11_1         string    `db:"odg_11_1"`
	Odg_11_2         string    `db:"odg_11_2"`
	Odg_11_3         string    `db:"odg_11_3"`
	Odg_11_4         string    `db:"odg_11_4"`
	R_11             string    `db:"r_11"`
	Pitanje_12       string    `db:"pitanje_12"`
	Odg_12_1         string    `db:"odg_12_1"`
	Odg_12_2         string    `db:"odg_12_2"`
	Odg_12_3         string    `db:"odg_12_3"`
	Odg_12_4         string    `db:"odg_12_4"`
	R_12             string    `db:"r_12"`
	Pitanje_13       string    `db:"pitanje_13"`
	Odg_13_1         string    `db:"odg_13_1"`
	Odg_13_2         string    `db:"odg_13_2"`
	Odg_13_3         string    `db:"odg_13_3"`
	Odg_13_4         string    `db:"odg_13_4"`
	R_13             string    `db:"r_13"`
	Pitanje_14       string    `db:"pitanje_14"`
	Odg_14_1         string    `db:"odg_14_1"`
	Odg_14_2         string    `db:"odg_14_2"`
	Odg_14_3         string    `db:"odg_14_3"`
	Odg_14_4         string    `db:"odg_14_4"`
	R_14             string    `db:"r_14"`
	Pitanje_15       string    `db:"pitanje_15"`
	Odg_15_1         string    `db:"odg_15_1"`
	Odg_15_2         string    `db:"odg_15_2"`
	Odg_15_3         string    `db:"odg_15_3"`
	Odg_15_4         string    `db:"odg_15_4"`
	R_15             string    `db:"r_15"`
	Pitanje_16       string    `db:"pitanje_16"`
	Odg_16_1         string    `db:"odg_16_1"`
	Odg_16_2         string    `db:"odg_16_2"`
	Odg_16_3         string    `db:"odg_16_3"`
	Odg_16_4         string    `db:"odg_16_4"`
	R_16             string    `db:"r_16"`
	Pitanje_17       string    `db:"pitanje_17"`
	Odg_17_1         string    `db:"odg_17_1"`
	Odg_17_2         string    `db:"odg_17_2"`
	Odg_17_3         string    `db:"odg_17_3"`
	Odg_17_4         string    `db:"odg_17_4"`
	R_17             string    `db:"r_17"`
	Pitanje_18       string    `db:"pitanje_18"`
	Odg_18_1         string    `db:"odg_18_1"`
	Odg_18_2         string    `db:"odg_18_2"`
	Odg_18_3         string    `db:"odg_18_3"`
	Odg_18_4         string    `db:"odg_18_4"`
	R_18             string    `db:"r_18"`
	Pitanje_19       string    `db:"pitanje_19"`
	Odg_19_1         string    `db:"odg_19_1"`
	Odg_19_2         string    `db:"odg_19_2"`
	Odg_19_3         string    `db:"odg_19_3"`
	Odg_19_4         string    `db:"odg_19_4"`
	R_19             string    `db:"r_19"`
	Pitanje_20       string    `db:"pitanje_20"`
	Odg_20_1         string    `db:"odg_20_1"`
	Odg_20_2         string    `db:"odg_20_2"`
	Odg_20_3         string    `db:"odg_20_3"`
	Odg_20_4         string    `db:"odg_20_4"`
	R_20             string    `db:"r_20"`
	Pitanje_21       string    `db:"pitanje_21"`
	Odg_21_1         string    `db:"odg_21_1"`
	Odg_21_2         string    `db:"odg_21_2"`
	Odg_21_3         string    `db:"odg_21_3"`
	Odg_21_4         string    `db:"odg_21_4"`
	R_21             string    `db:"r_21"`
	Pitanje_22       string    `db:"pitanje_22"`
	Odg_22_1         string    `db:"odg_22_1"`
	Odg_22_2         string    `db:"odg_22_2"`
	Odg_22_3         string    `db:"odg_22_3"`
	Odg_22_4         string    `db:"odg_22_4"`
	R_22             string    `db:"r_22"`
	Pitanje_23       string    `db:"pitanje_23"`
	Odg_23_1         string    `db:"odg_23_1"`
	Odg_23_2         string    `db:"odg_23_2"`
	Odg_23_3         string    `db:"odg_23_3"`
	Odg_23_4         string    `db:"odg_23_4"`
	R_23             string    `db:"r_23"`
	Pitanje_24       string    `db:"pitanje_24"`
	Odg_24_1         string    `db:"odg_24_1"`
	Odg_24_2         string    `db:"odg_24_2"`
	Odg_24_3         string    `db:"odg_24_3"`
	Odg_24_4         string    `db:"odg_24_4"`
	R_24             string    `db:"r_24"`
	Pitanje_25       string    `db:"pitanje_25"`
	Odg_25_1         string    `db:"odg_25_1"`
	Odg_25_2         string    `db:"odg_25_2"`
	Odg_25_3         string    `db:"odg_25_3"`
	Odg_25_4         string    `db:"odg_25_4"`
	R_25             string    `db:"r_25"`
	Pitanje_26       string    `db:"pitanje_26"`
	Odg_26_1         string    `db:"odg_26_1"`
	Odg_26_2         string    `db:"odg_26_2"`
	Odg_26_3         string    `db:"odg_26_3"`
	Odg_26_4         string    `db:"odg_26_4"`
	R_26             string    `db:"r_26"`
	Pitanje_27       string    `db:"pitanje_27"`
	Odg_27_1         string    `db:"odg_27_1"`
	Odg_27_2         string    `db:"odg_27_2"`
	Odg_27_3         string    `db:"odg_27_3"`
	Odg_27_4         string    `db:"odg_27_4"`
	R_27             string    `db:"r_27"`
	Pitanje_28       string    `db:"pitanje_28"`
	Odg_28_1         string    `db:"odg_28_1"`
	Odg_28_2         string    `db:"odg_28_2"`
	Odg_28_3         string    `db:"odg_28_3"`
	Odg_28_4         string    `db:"odg_28_4"`
	R_28             string    `db:"r_28"`
	Pitanje_29       string    `db:"pitanje_29"`
	Odg_29_1         string    `db:"odg_29_1"`
	Odg_29_2         string    `db:"odg_29_2"`
	Odg_29_3         string    `db:"odg_29_3"`
	Odg_29_4         string    `db:"odg_29_4"`
	R_29             string    `db:"r_29"`
	Pitanje_30       string    `db:"pitanje_30"`
	Odg_30_1         string    `db:"odg_30_1"`
	Odg_30_2         string    `db:"odg_30_2"`
	Odg_30_3         string    `db:"odg_30_3"`
	Odg_30_4         string    `db:"odg_30_4"`
	R_30             string    `db:"r_30"`
}

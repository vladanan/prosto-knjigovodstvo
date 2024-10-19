package models

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/vladanan/prosto/src/controllers/clr"
)

type Tim struct {
	Table string
	Id    string
	Mail  string
}

var ApiToDb = map[string]Tim{
	"test": {
		Table: "g_pitanja_c_testovi",
		Id:    "g_id",
		Mail:  "user_id",
	},
	"krsnc_usrs": {
		Table: "mi_users",
		Id:    "u_id",
		Mail:  "email",
	},
	"data": {
		Table: "user_data",
		Id:    "ud_id",
		Mail:  "u_email",
	},
	"note": {
		Table: "g_user_blog",
		Id:    "b_id",
		Mail:  "user_mail",
	},
	"setting": {
		Table: "v_settings",
		Id:    "s_id",
		Mail:  "",
	},
}

// Any struct to map/string funkcija, trenutno ne vraća map nego string
func getReturnData(recordData any) string {
	uMap := map[string]interface{}{}
	// otkriva type na osnovu strukture samih podataka
	tip := reflect.TypeOf(recordData)
	// orkriva sva polja u tipu
	fields := reflect.VisibleFields(tip)
	// početak povratne poruke
	msg := strings.ToLower(tip.Name()) + " command successful for: "
	for i, e := range fields {
		// uzima se kompletan opis polja, dole je primer za Email
		field := fmt.Sprint(tip.FieldByIndex(e.Index))
		// {Email  string db:"email" 48 [3] false}
		// iz potpunog opisa vadi se samo ime polja: Email
		field = strings.Split(field, "  ")[0]
		field = strings.ReplaceAll(field, "{", "")
		// samo polja koja nisu prazna se stavljaju u map i string iako nam map trenutno ne treba
		if !reflect.ValueOf(recordData).Field(i).IsZero() {
			if field == "Hash_lozinka" {
				continue // ne želimo da se lozinka korisnika pojavljuje u našim logovima
			}
			uMap[field] = reflect.ValueOf(recordData).Field(i)
			msg = msg + " " + fmt.Sprint(reflect.ValueOf(recordData).Field(i))
		}
		// fmt.Print(uMap["Email"])
	}
	return msg
}

func (db DB) PostOne(recordData any, r *http.Request) (string, error) {

	// lock := sync.RWMutex{}

	conn, err := getDBConn(r)
	if err != nil {
		return "", l(r, 8, err) //
	}
	defer freeConn(conn, r)

	switch data := recordData.(type) {

	case Test:
		commandTag, err := conn.Exec(r.Context(), `INSERT INTO g_pitanja_c_testovi 
			(
				tip,
				obrazovni_profil,
				razred,
				predmet,
				oblast
			)
				VALUES ($1, $2, $3, $4, $5);`,
			data.Tip,
			data.Obrazovni_profil,
			data.Razred,
			data.Predmet,
			data.Oblast,
		)
		if err != nil {
			return "", l(r, 8, err)
		}
		if commandTag.String() != "INSERT 0 1" {
			return "", l(r, 0, clr.NewAPIError(http.StatusBadRequest, "Insert-update-delete_operation_failed"))
		} else {
			// return clr.GetStringLogger()(r, 0, getReturnData(data)), nil
			return getReturnData(data), nil
		}

	case User:
		return "", l(r, 4, fmt.Errorf("kod za post user je u zasbnom api"))
		// lock.Lock()
		// if err := AddUser(data.Email, data.User_name, data.Hash_lozinka, r); err != nil {
		// 	lock.Unlock()
		// 	return "", err
		// } else {
		// 	lock.Unlock()
		// 	return getReturnData(data), nil
		// }

	case Note:
		commandTag, err := conn.Exec(r.Context(), `INSERT INTO g_user_blog 
			(
				ime_tag,
				mejl,
				tema,
				poruka,
				user_id,
				user_mail,
				from_url
			)
				VALUES ($1, $2, $3, $4, $5, $6, $7);`,
			data.Ime_tag,
			data.Mejl,
			data.Tema,
			data.Poruka,
			data.User_id,
			data.User_mail,
			data.From_url,
		)
		if err != nil {
			return "", l(r, 8, err)
		}
		if commandTag.String() != "INSERT 0 1" {
			return "", l(r, 0, clr.NewAPIError(http.StatusBadRequest, "Insert-update-delete_operation_failed"))
		} else {
			return getReturnData(data), nil
		}

	default:
		return "", l(r, 8, fmt.Errorf("post record ne pripada nijednom tipu"))
	}

}

func (db DB) Get(table string, field string, record any, r *http.Request) (any, error) {

	conn, err := getDBConn(r)
	if err != nil {
		return nil, l(r, 8, err)
	}
	defer freeConn(conn, r)

	var pgxData any
	var rows pgx.Rows

	//https://stackoverflow.com/questions/61704842/how-to-scan-a-queryrow-into-a-struct-with-pgx
	if field == "" {
		rows, err = conn.Query(r.Context(), "SELECT * FROM "+table+";")
		if err != nil {
			return nil, l(r, 8, err)
		}
	} else {
		rows, err = conn.Query(r.Context(), "SELECT * FROM "+table+" WHERE "+field+"=$1;", record)
		if err != nil {
			return nil, l(r, 8, err)
		}
	}

	// var datum time.Time
	// row := conn.QueryRow(r.Context(), "select datum_upisa from g_user_blog where b_id=$1", 7)
	// // log.Println("red:", row)
	// row.Scan(&datum)
	// if err != nil {
	// 	log.Print(err)
	// }
	// log.Println("datum:", datum)

	// for rows.Next() {
	// 	if val, err := rows.Values(); err != nil {
	// 		fmt.Println("rows greška:", err)
	// 		// return nil, l(r, 8, err)
	// 	} else {
	// 		fmt.Println("row:", fmt.Sprint(val))
	// 	}
	// }

	switch table {

	case "g_pitanja_c_testovi":
		pgxData, err = pgx.CollectRows(rows, pgx.RowToStructByName[Test])
		// CollectRows automatski zatvara rows nakon učitavanja svih rows tako da nema potrebe da se zatvara rows kao kada se rucno radi sa njima
		if err != nil {
			return nil, l(r, 8, err)
		}
		if fmt.Sprint(pgxData) == "[]" {
			pgxData = []Test{{}}
		}

	case "mi_users":
		pgxData, err = pgx.CollectRows(rows, pgx.RowToStructByName[User])
		if err != nil {
			return nil, l(r, 8, err)
		}
		if fmt.Sprint(pgxData) == "[]" {
			pgxData = []User{{}}
		}

	case "user_data":
		pgxData, err = pgx.CollectRows(rows, pgx.RowToStructByName[UserData])
		if err != nil {
			return nil, l(r, 8, err)
		}
		if fmt.Sprint(pgxData) == "[]" {
			pgxData = []UserData{{}}
		}

	case "g_user_blog":
		// https://github.com/jackc/pgx/issues/186 metode da se neka polja isključe ne rade kod mene
		pgxData, err = pgx.CollectRows(rows, pgx.RowToStructByName[Note])
		if err != nil {
			return nil, l(r, 8, err)
		}
		if fmt.Sprint(pgxData) == "[]" {
			pgxData = []Note{{}}
		}

	case "v_settings":
		pgxData, err = pgx.CollectRows(rows, pgx.RowToStructByName[Settings])
		if err != nil {
			return nil, l(r, 8, err)
		}
		if fmt.Sprint(pgxData) == "[]" {
			pgxData = []Settings{{}}
		}
	default:
		return nil, l(r, 8, fmt.Errorf("malformed table name"))
	}

	return pgxData, nil

}

func (db DB) PutOne(table string, field string, record any, recordData any, r *http.Request) (string, error) {

	conn, err := getDBConn(r)
	if err != nil {
		return "", l(r, 8, err) //
	}
	defer freeConn(conn, r)

	switch data := recordData.(type) {

	case Test:
		// za update napraviti kod koji na osnovu poslatih polja za izmenu i već postojećih napravi skroz novi upis za isti id tako da se izbegnu kompleksni (string contactenation) query i kompleksan kod JER SE NE ZNA UNAPRED KOJA POLJA ĆE KORISNIK DA MENJA A KOJA NE
		commandTag, err := conn.Exec(r.Context(), `UPDATE `+table+` SET
			tip=$1,
			obrazovni_profil=$2,
			razred=$3,
			predmet=$4,
			oblast=$5
			WHERE `+field+`=$6;`,
			data.Tip,
			data.Obrazovni_profil,
			data.Razred,
			data.Predmet,
			data.Oblast,
			record,
		)
		if err != nil {
			return "", l(r, 8, err)
		}
		if commandTag.String() != "UPDATE 1" {
			return "", l(r, 0, clr.NewAPIError(http.StatusBadRequest, "Insert-update-delete_operation_failed"))
		} else {
			return getReturnData(data), nil
		}

	case Note:
		// za update napraviti kod koji na osnovu poslatih polja za izmenu i već postojećih napravi skroz novi upis za isti id tako da se izbegnu kompleksni (string contactenation) query i kompleksan kod JER SE NE ZNA UNAPRED KOJA POLJA ĆE KORISNIK DA MENJA A KOJA NE
		commandTag, err := conn.Exec(r.Context(), `UPDATE `+table+` SET
			ime_tag=$1,
			mejl=$2,
			tema=$3,
			poruka=$4,
			user_id=$5
			WHERE `+field+`=$6;`,
			data.Ime_tag,
			data.Mejl,
			data.Tema,
			data.Poruka,
			data.User_id,
			record,
		)
		if err != nil {
			return "", l(r, 8, err)
		}
		// fmt.Println("put commandTag:", commandTag)
		if commandTag.String() != "UPDATE 1" {
			return "", l(r, 0, clr.NewAPIError(http.StatusBadRequest, "Insert-update-delete_operation_failed"))
		} else {
			return getReturnData(data), nil
		}

	default:
		return "", l(r, 8, fmt.Errorf("put record ne pripada nijednom tipu"))
	}

}

func (db DB) DeleteOne(table string, field string, record any, r *http.Request) (string, error) {

	conn, err := getDBConn(r)
	if err != nil {
		return "", l(r, 8, err)
	}
	defer freeConn(conn, r)

	commandTag, err := conn.Exec(r.Context(), "DELETE FROM "+table+" WHERE "+field+"=$1;", record)
	// log.Println("delete test:", commandTag, commandTag.String())
	if err != nil {
		return "", l(r, 8, err)
	}
	if commandTag.String() == "DELETE 0" {
		return "", l(r, 0, clr.NewAPIError(http.StatusBadRequest, "Insert-update-delete_operation_failed"))
	} else {
		return "record deleted", nil
	}

}

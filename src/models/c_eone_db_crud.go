package models

import (
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func (db DB) GetBilling(r *http.Request) (any, error) {

	conn, err := getDBConn(r)
	if err != nil {
		return nil, l(r, 8, err)
	}
	defer freeConn(conn, r)

	rows, err := conn.Query(r.Context(), "SELECT * FROM billing ORDER BY id ASC;")
	if err != nil {
		return nil, l(r, 8, err)
	}

	pgxData, err := pgx.CollectRows(rows, pgx.RowToStructByName[Billing])
	if err != nil {
		return nil, l(r, 8, err)
	}

	if fmt.Sprint(pgxData) == "[]" {
		pgxData = []Billing{{}}
	}

	return pgxData, nil

}

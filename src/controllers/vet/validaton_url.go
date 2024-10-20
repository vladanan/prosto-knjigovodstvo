package vet

import (
	"net/http"
	"strconv"

	"github.com/vladanan/prosto/src/controllers/clr"
	"github.com/vladanan/prosto/src/models"
)

func ValidateURLKeys(
	tableApi, fieldApi, recordApi string) (
	tableDb, fieldDb string, recordDb any, err error) {

	// log.Println("ValidateURLKeys tableapi, fieldapi, record api:", tableApi, fieldApi, recordApi)

	if fieldApi == "id" && recordApi != "" {
		var err error
		recordDb, err = strconv.Atoi(recordApi)
		if err != nil {
			return "", "", nil, clr.NewAPIError(http.StatusNotAcceptable, "Malformed_id")
		}
	}

	if fieldApi == "mail" && recordApi != "" {
		if err := ValidateEmailAddress(recordApi); err != nil {
			return "", "", nil, err
		}
		recordDb = recordApi
	}
	// log.Println("ValidateURLKeys recordD:", recordDb)

	for a := range models.ApiToDb {
		// log.Println("deo od map:", a)
		if a == tableApi {
			tableDb = models.ApiToDb[a].Table
			switch fieldApi {
			case "id":
				fieldDb = models.ApiToDb[a].Id
			case "mail":
				fieldDb = models.ApiToDb[a].Mail
			case "":
				// log.Println("field prazan")
				fieldDb = ""
			default:
				// log.Println("field neispravan")
				return "", "", nil, clr.NewAPIError(http.StatusNotAcceptable, "Malformed_field_name")
			}

		}
	}

	// log.Println(tableApi, fieldApi, recordApi, tableDb, fieldDb)

	if tableDb == "" {
		return "", "", nil, clr.NewAPIError(http.StatusNotAcceptable, "Malformed_data_type")
	}

	return
}

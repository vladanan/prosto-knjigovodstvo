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

	for a := range models.ApiToDb {
		// fmt.Println("deo od map:", a)
		if a == tableApi {
			tableDb = models.ApiToDb[a].Table
			switch fieldApi {
			case "id":
				fieldDb = models.ApiToDb[a].Id
			case "mail":
				fieldDb = models.ApiToDb[a].Mail
			case "":
				fieldDb = ""
			default:
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

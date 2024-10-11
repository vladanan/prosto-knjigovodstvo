package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/vladanan/prosto/src/controllers/clr"
	"github.com/vladanan/prosto/src/controllers/vet"
	"github.com/vladanan/prosto/src/models"
)

type VezbamoHandler struct {
	db models.DB
}

func NewVezbamoHandler(db models.DB) *VezbamoHandler {
	return &VezbamoHandler{db: db}
}

func (h *VezbamoHandler) HandlePostOne(w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	tableApi := vars["table"]
	tableDb, _, _, err := vet.ValidateURLKeys(tableApi, "", "")
	if err != nil {
		return err
	}

	// body, err := io.ReadAll(r.Body)
	// if err != nil {
	// 	return l(r, 7, err)
	// }

	// dec := json.NewDecoder(bytes.NewReader(body))
	// recordData, err := vet.DecodeBodyAndValidatePostAndPutData1(tableDb, dec)
	recordData, err := vet.ValidateFormPostAndPutData(tableDb, r)
	if err != nil {
		return err
	}

	if returnData, err := h.db.PostOne(recordData, r); err != nil {
		return err
	} else {
		return clr.WriteJSON(w, 200, returnData)
	}

}

func (h *VezbamoHandler) HandleGet(w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	tableApi := vars["table"]
	fieldApi := vars["field"]
	recordApi := vars["record"]

	tableDb, fieldDb, recordDb, err := vet.ValidateURLKeys(tableApi, fieldApi, recordApi)
	if err != nil {
		return err
	}

	data, err := h.db.Get(tableDb, fieldDb, recordDb, r)
	if err != nil {
		return err
	} else {
		return clr.WriteJSON(w, 200, data)
		// io.WriteString(w, string(h.db.Get()))
	}

}

func (h *VezbamoHandler) HandlePutOne(w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	tableApi := vars["table"]
	fieldApi := vars["field"]
	recordApi := vars["record"]
	tableDb, fieldDb, recordDb, err := vet.ValidateURLKeys(tableApi, fieldApi, recordApi)
	if err != nil {
		return err
	}
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return l(r, 7, err)
	}

	dec := json.NewDecoder(bytes.NewReader(body))
	recordData, err := vet.DecodeBodyAndValidatePostAndPutData1(tableDb, dec)
	if err != nil {
		return err
	}

	if returnData, err := h.db.PutOne(tableDb, fieldDb, recordDb, recordData, r); err != nil {
		return l(r, 8, err)
	} else {
		return clr.WriteJSON(w, 200, returnData)
	}

}

func (h *VezbamoHandler) HandleDeleteOne(w http.ResponseWriter, r *http.Request) error {

	vars := mux.Vars(r)
	tableApi := vars["table"]
	fieldApi := vars["field"]
	recordApi := vars["record"]

	tableDb, fieldDb, recordDb, err := vet.ValidateURLKeys(tableApi, fieldApi, recordApi)
	if err != nil {
		return err
	}

	if returnData, err := h.db.DeleteOne(tableDb, fieldDb, recordDb, r); err != nil {
		return l(r, 8, err)
	} else {
		return clr.WriteJSON(w, 200, returnData)
	}

}

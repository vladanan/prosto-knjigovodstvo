package api

import (
	"net/http"

	"github.com/vladanan/prosto/src/controllers/clr"
	"github.com/vladanan/prosto/src/models"
)

type EoneHandler struct {
	db models.DB
}

func NewEoneHandler(db models.DB) *EoneHandler {
	return &EoneHandler{db: db}
}

func (h *EoneHandler) HandleGetBilling(w http.ResponseWriter, r *http.Request) error {

	data, err := h.db.GetBilling(r)
	if err != nil {
		return err
	} else {
		return clr.WriteJSON(w, 200, data)
	}

}

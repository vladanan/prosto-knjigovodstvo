package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/vladanan/prosto/src/controllers/clr"
	"github.com/vladanan/prosto/src/models"
)

func TestGetTests(t *testing.T) {
	// https://go.dev/doc/code
	// https://www.cloudbees.com/blog/testing-http-handlers-go

	tp := []struct {
		path string
		pass bool
	}{
		{"", true},
		{"/", false},
		{"/8m", false},
		{"/0", true},
		{"/7", true},
		{"/22", false},
		{"/37", false},
		{"/38", false},
	}

	tm := []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodPut}

	for _, p := range tp {
		for _, m := range tm {

			path := fmt.Sprintf("/test%s", p.path)
			r, err := http.NewRequest(m, path, nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			vh := NewVezbamoHandler(models.DB{})

			m := mux.NewRouter()
			m.HandleFunc("/test", clr.CheckFunc(vh.HandleGet)).Methods(http.MethodGet)
			m.HandleFunc("/test/{id}", clr.CheckFunc(vh.HandleGet))
			m.ServeHTTP(rr, r)

			if status := rr.Code; status != http.StatusOK {
				var red string
				if r.Method == http.MethodGet {
					red = clr.BgRed
					// r.Method = clr.BgRed + http.MethodGet
				}
				t.Errorf("%surL: %s\t%s\t%v%s", red, r.URL.Path, r.Method, status, clr.Reset)
			}

			// list_string := rr.Body.String() // r.Body is a *bytes.Buffer
			// dec := json.NewDecoder(strings.NewReader(list_string))
			// var all_tests []models.Test
			// if err := dec.Decode(&all_tests); err != nil {
			// 	t.Error("json error")
			// }

			// if all_tests == nil {
			// 	t.Errorf("data array of %v, want %v", len(all_tests), "> 0")
			// }

			// fmt.Println("test path i method:", r.URL.Path, r.Method)

		}
		t.Errorf("")
	}

}

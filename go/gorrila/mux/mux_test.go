package mux

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/sqlcommenter/go/core"
	"github.com/gorilla/mux"
)

func TestSQLCommenterMiddleware(t *testing.T) {
	framework := "gorrila/mux"
	route := "GET--/test/{id}"
	action := "github.com/google/sqlcommenter/go/gorrila/mux.TestSQLCommenterMiddleware.func1"

	mockHandler := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		_framework := ctx.Value(core.Framework)
		_route := ctx.Value(core.Route)
		_action := ctx.Value(core.Action)

		if _framework != framework {
			t.Errorf("mismatched framework - got: %s, want: %s", _framework, framework)
		}

		if _route != route {
			t.Errorf("mismatched route - got: %s, want: %s", _route, route)
		}

		if _action != action {
			t.Errorf("mismatched action - got: %s, want: %s", _action, action)
		}
	}

	router := mux.NewRouter()
	router.Use(SQLCommenterMiddleware)
	router.HandleFunc(route, mockHandler).Methods("GET")

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/test/1", nil)

	if err != nil {
		t.Errorf("error while building req: %v", err)
	}

	router.ServeHTTP(rr, req)
}

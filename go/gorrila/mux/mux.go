package mux

import (
	"net/http"

	"github.com/google/sqlcommenter/go/core"
	httpnet "github.com/google/sqlcommenter/go/net/http"
	"github.com/gorilla/mux"
)

func SQLCommenterMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		pathTemplate, err := route.GetPathTemplate()
		if err != nil {
			pathTemplate = ""
		}

		ctx := core.ContextInject(r.Context(), httpnet.NewHTTPRequestTags("gorrila/mux", pathTemplate, core.GetFunctionName(route.GetHandler())))
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

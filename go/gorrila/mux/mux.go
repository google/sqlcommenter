package mux

import (
	"fmt"
	"net/http"

	"github.com/google/sqlcommenter/go/core"
	httpnet "github.com/google/sqlcommenter/go/net/http"
	"github.com/gorilla/mux"
)

func SQLCommenterMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		muxRoute := mux.CurrentRoute(r)
		path, err := muxRoute.GetPathTemplate()
		if err != nil {
			path = ""
		}

		route := fmt.Sprintf("%s--%s", r.Method, path)
		ctx := core.ContextInject(r.Context(), httpnet.NewHTTPRequestTags("gorrila/mux", route, core.GetFunctionName(muxRoute.GetHandler())))
		h.ServeHTTP(w, r.WithContext(ctx))
	})
}

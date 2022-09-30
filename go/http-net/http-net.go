package httpnet

import (
	"context"
	"net/http"

	"google.com/sqlcommenter/core"
)

type HttpNet struct {
	r    *http.Request
	next any
}

func NewHttpNet(r *http.Request, next any) *HttpNet {
	return &HttpNet{r, next}
}

func (h *HttpNet) Route() string {
	return h.r.URL.Path
}

func (h *HttpNet) Action() string {
	return core.GetFunctionName(h.next)
}

func (h *HttpNet) Framework() string {
	return "net/http"
}

func (h *HttpNet) AddTags(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, core.Route, h.Route())
	ctx = context.WithValue(ctx, core.Action, h.Action())
	ctx = context.WithValue(ctx, core.Framework, h.Framework())
	return ctx
}

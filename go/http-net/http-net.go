// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

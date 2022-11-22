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

package core

import (
	"context"
	"fmt"
	"net/url"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"go.opentelemetry.io/otel/propagation"
)

const (
	Route       string = "route"
	Controller  string = "controller"
	Action      string = "action"
	Framework   string = "framework"
	Driver      string = "db_driver"
	Traceparent string = "traceparent"
	Application string = "application"
)

type CommenterOptions struct {
	EnableDBDriver    bool
	EnableRoute       bool
	EnableFramework   bool
	EnableController  bool
	EnableAction      bool
	EnableTraceparent bool
	EnableApplication bool
	Application       string
}

func encodeURL(k string) string {
	return url.QueryEscape(string(k))
}

func GetFunctionName(i interface{}) string {
	if i == nil {
		return ""
	}
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func ConvertMapToComment(tags map[string]string) string {
	var sb strings.Builder
	i, sz := 0, len(tags)

	//sort by keys
	sortedKeys := make([]string, 0, len(tags))
	for k := range tags {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		if i == sz-1 {
			sb.WriteString(fmt.Sprintf("%s=%v", encodeURL(key), encodeURL(tags[key])))
		} else {
			sb.WriteString(fmt.Sprintf("%s=%v,", encodeURL(key), encodeURL(tags[key])))
		}
		i++
	}
	return sb.String()
}

func ExtractTraceparent(ctx context.Context) propagation.MapCarrier {
	// Serialize the context into carrier
	textMapPropogator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	carrier := propagation.MapCarrier{}
	textMapPropogator.Inject(ctx, carrier)
	return carrier
}

type RequestTags interface {
	Route() string
	Action() string
	Framework() string
}

func ContextInject(ctx context.Context, h RequestTags) context.Context {
	ctx = context.WithValue(ctx, Route, h.Route())
	ctx = context.WithValue(ctx, Action, h.Action())
	ctx = context.WithValue(ctx, Framework, h.Framework())
	return ctx
}

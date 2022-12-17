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
)

// Constants used as key string for tags.
// It is not necessary that all SQLCommenter frameworks/ORMs will contain all these keys i.e.
// it is on best-effort basis.
const (
	Route       string = "route"
	Controller         = "controller"
	Action             = "action"
	Framework          = "framework"
	Driver             = "db_driver"
	Traceparent        = "traceparent"
	Application        = "application"
)

// CommenterConfig contains configurations for SQLCommenter library.
// We can enable and disable certain tags by enabling these configurations.
type CommenterConfig struct {
	EnableDBDriver    bool
	EnableRoute       bool
	EnableFramework   bool
	EnableController  bool
	EnableAction      bool
	EnableTraceparent bool
	EnableApplication bool
}

// StaticTags are few tags that can be set by the application and will be constant
// for every API call.
type StaticTags struct {
	Application string
	DriverName  string
}

// CommenterOptions contains all options regarding SQLCommenter library.
// This includes the configurations as well as any static tags.
type CommenterOptions struct {
	Config CommenterConfig
	Tags   StaticTags
}

func encodeURL(k string) string {
	return url.QueryEscape(k)
}

// GetFunctionName returns the name of the function passed.
func GetFunctionName(i any) string {
	if i == nil {
		return ""
	}
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

// ConvertMapToComment returns a comment string given a map of key-value pairs of tags.
// There are few steps involved here:
//   - Sorting the tags by key string
//   - url encoding the key value pairs
//   - Formatting the key value pairs as "key1=value1,key2=value2" format.
func ConvertMapToComment(tags map[string]string) string {
	var sb strings.Builder
	i, sz := 0, len(tags)

	// sort by keys
	sortedKeys := make([]string, 0, len(tags))
	for k := range tags {
		sortedKeys = append(sortedKeys, k)
	}
	sort.Strings(sortedKeys)

	for _, key := range sortedKeys {
		if i == sz-1 {
			sb.WriteString(fmt.Sprintf("%s='%s'", encodeURL(key), encodeURL(tags[key])))
		} else {
			sb.WriteString(fmt.Sprintf("%s='%s',", encodeURL(key), encodeURL(tags[key])))
		}
		i++
	}
	return sb.String()
}

// ExtractTraceparent extracts the traceparent field using OpenTelemetry library.
func ExtractTraceparent(ctx context.Context) propagation.MapCarrier {
	// Serialize the context into carrier
	textMapPropogator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	carrier := propagation.MapCarrier{}
	textMapPropogator.Inject(ctx, carrier)
	return carrier
}

// RequestTagsProvider adds a basic interface for other libraries like gorilla/mux to implement.
type RequestTagsProvider interface {
	Route() string
	Action() string
	Framework() string
}

// ContextInject injects the tags key-value pairs into context,
// which can be later passed into drivers/ORMs to finally inject them into SQL queries.
func ContextInject(ctx context.Context, h RequestTagsProvider) context.Context {
	ctx = context.WithValue(ctx, Route, h.Route())
	ctx = context.WithValue(ctx, Action, h.Action())
	ctx = context.WithValue(ctx, Framework, h.Framework())
	return ctx
}

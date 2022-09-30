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
	Driver      string = "driver"
	Traceparent string = "traceparent"
)

type CommenterOptions struct {
	EnableDBDriver    bool
	EnableRoute       bool
	EnableFramework   bool
	EnableController  bool
	EnableAction      bool
	EnableTraceparent bool
}

func encodeURL(k string) string {
	return url.QueryEscape(string(k))
}

func GetFunctionName(i interface{}) string {
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
	propgator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	carrier := propagation.MapCarrier{}
	propgator.Inject(ctx, carrier)
	return carrier
}

type RequestTagger interface {
	Route() string
	Action() string
	Framework() string
	GetContext() context.Context
}

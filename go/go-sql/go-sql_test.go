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

package gosql

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.com/sqlcommenter/core"
	httpnet "google.com/sqlcommenter/http-net"
)

var engine, connectionParams = "mysql", "root:root@/gotest"

func TestDisabled(t *testing.T) {
	mockDB, _, err := sqlmock.New()

	db := DB{DB: mockDB, options: core.CommenterOptions{}}
	if err != nil {
		t.Fatalf("MockSQL failed with unexpected error: %s", err)
	}

	query := "SELECT 2"
	if got, want := db.withComment(context.Background(), query), query; got != want {
		t.Errorf("db.withComment(context.Background(), %q) = %q, want = %q", query, got, want)
	}
}

func TestHTTP_Net(t *testing.T) {
	mockDB, _, err := sqlmock.New()

	db := DB{DB: mockDB, options: core.CommenterOptions{EnableDBDriver: true, EnableRoute: true, EnableFramework: true}}
	if err != nil {
		t.Fatalf("MockSQL failed with unexpected error: %s", err)
	}

	r, _ := http.NewRequest("GET", "hello/1", nil)
	ctx := httpnet.NewHttpNet(r, context.Background()).AddTags(r.Context())

	got := db.withComment(ctx, "Select 1")
	want := "Select 1/*driver=database%2Fsql,framework=net%2Fhttp,route=hello%2F1*/"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestQueryWithSemicolon(t *testing.T) {
	mockDB, _, err := sqlmock.New()

	db := DB{DB: mockDB, options: core.CommenterOptions{EnableDBDriver: true}}
	if err != nil {
		t.Fatalf("MockSQL failed with unexpected error: %s", err)
	}
	got := db.withComment(context.Background(), "Select 1;")
	want := "Select 1/*driver=database%2Fsql*/;"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestOtelIntegration(t *testing.T) {
	mockDB, _, err := sqlmock.New()

	db := DB{DB: mockDB, options: core.CommenterOptions{EnableTraceparent: true}}
	if err != nil {
		t.Fatalf("MockSQL failed with unexpected error: %s", err)
	}

	exp, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	bsp := sdktrace.NewSimpleSpanProcessor(exp) // You should use batch span processor in prod
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
	)

	ctx, _ := tp.Tracer("").Start(context.Background(), "parent-span-name")

	got := db.withComment(ctx, "Select 1;")
	r, _ := regexp.Compile("Select 1/\\*traceparent=\\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\\d{1,2}\\*/;")

	if !r.MatchString(got) {
		t.Errorf("got %q", got)
	}
}

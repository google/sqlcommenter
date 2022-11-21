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

package sql

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/sqlcommenter/go/core"
	httpnet "github.com/google/sqlcommenter/go/net/http"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestDisabled(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("MockSQL failed with unexpected error: %s", err)
	}
	db := DB{DB: mockDB, driverName: "mocksql", options: core.CommenterOptions{}}
	query := "SELECT 2"
	if got, want := db.withComment(context.Background(), query), query; got != want {
		t.Errorf("db.withComment(context.Background(), %q) = %q, want = %q", query, got, want)
	}
}

func TestHTTP_Net(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("MockSQL failed with unexpected error: %s", err)
	}

	db := DB{DB: mockDB, driverName: "mocksql", options: core.CommenterOptions{EnableDBDriver: true, EnableRoute: true, EnableFramework: true, EnableApplication: true, Application: "app"}, application: "app"}
	r, err := http.NewRequest("GET", "hello/1", nil)
	if err != nil {
		t.Errorf("http.NewRequest('GET', 'hello/1', nil) returned unexpected error: %v", err)
	}

	ctx := core.ContextInject(r.Context(), httpnet.NewHTTPRequestExtractor(r, nil))
	got := db.withComment(ctx, "Select 1")
	want := "Select 1/*application=app,db_driver=database%2Fsql%3Amocksql,framework=net%2Fhttp,route=hello%2F1*/"
	if got != want {
		t.Errorf("db.withComment(ctx, 'Select 1') got %q, wanted %q", got, want)
	}
}

func TestQueryWithSemicolon(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("MockSQL failed with unexpected error: %s", err)
	}

	db := DB{DB: mockDB, driverName: "mocksql", options: core.CommenterOptions{EnableDBDriver: true}}
	got := db.withComment(context.Background(), "Select 1;")
	want := "Select 1/*db_driver=database%2Fsql%3Amocksql*/;"
	if got != want {
		t.Errorf("db.withComment(context.Background(), 'Select 1;') got %q, wanted %q", got, want)
	}
}

func TestOtelIntegration(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("MockSQL failed with unexpected error: %s", err)
	}

	db := DB{DB: mockDB, driverName: "mocksql", options: core.CommenterOptions{EnableTraceparent: true}}
	exp, _ := stdouttrace.New(stdouttrace.WithPrettyPrint())
	bsp := sdktrace.NewSimpleSpanProcessor(exp) // You should use batch span processor in prod
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithSpanProcessor(bsp),
	)
	ctx, _ := tp.Tracer("").Start(context.Background(), "parent-span-name")

	got := db.withComment(ctx, "Select 1;")
	wantRegex := "Select 1/\\*traceparent=\\d{1,2}-[a-zA-Z0-9_]{32}-[a-zA-Z0-9_]{16}-\\d{1,2}\\*/;"
	r, err := regexp.Compile(wantRegex)
	if err != nil {
		t.Errorf("regex.Compile() failed with error: %v", err)
	}

	if !r.MatchString(got) {
		t.Errorf("%q does not match the given regex %q", got, wantRegex)
	}
}

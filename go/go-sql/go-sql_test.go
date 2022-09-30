package gosql

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var engine, connectionParams = "mysql", "root:root@/gotest"

func TestDisabled(t *testing.T) {
	mockDB, _, err := sqlmock.New()

	db := DB{DB: mockDB, options: CommenterOptions{}}
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

	db := DB{DB: mockDB, options: CommenterOptions{EnableDBDriver: true, EnableRoute: true, EnableFramework: true}}
	if err != nil {
		t.Fatalf("MockSQL failed with unexpected error: %s", err)
	}

	r, _ := http.NewRequest("GET", "hello/1", nil)
	ctx := AddHttpRouterTags(r, context.Background())

	got := db.withComment(ctx, "Select 1")
	want := "Select 1/*driver=database%2Fsql,framework=net%2Fhttp,route=hello%2F1*/"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestQueryWithSemicolon(t *testing.T) {
	mockDB, _, err := sqlmock.New()

	db := DB{DB: mockDB, options: CommenterOptions{EnableDBDriver: true}}
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

	db := DB{DB: mockDB, options: CommenterOptions{EnableTraceparent: true}}
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

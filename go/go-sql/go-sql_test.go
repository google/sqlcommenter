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

func TestDisabled(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("MockSQL failed with unexpected error: %s", err)
	}
	db := DB{DB: mockDB, options: core.CommenterOptions{}}
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

	db := DB{DB: mockDB, options: core.CommenterOptions{EnableDBDriver: true, EnableRoute: true, EnableFramework: true}}
	r, err := http.NewRequest("GET", "hello/1", nil)
	if err != nil {
		t.Errorf("http.NewRequest('GET', 'hello/1', nil) returned unexpected error: %v", err)
	}

	ctx := core.ContextInject(r.Context(), httpnet.NewHTTPRequestExtractor(r, nil))
	got := db.withComment(ctx, "Select 1")
	want := "Select 1/*driver=database%2Fsql,framework=net%2Fhttp,route=hello%2F1*/"
	if got != want {
		t.Errorf("db.withComment(ctx, 'Select 1') got %q, wanted %q", got, want)
	}
}

func TestQueryWithSemicolon(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("MockSQL failed with unexpected error: %s", err)
	}

	db := DB{DB: mockDB, options: core.CommenterOptions{EnableDBDriver: true}}
	got := db.withComment(context.Background(), "Select 1;")
	want := "Select 1/*driver=database%2Fsql*/;"
	if got != want {
		t.Errorf("db.withComment(context.Background(), 'Select 1;') got %q, wanted %q", got, want)
	}
}

func TestOtelIntegration(t *testing.T) {
	mockDB, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("MockSQL failed with unexpected error: %s", err)
	}

	db := DB{DB: mockDB, options: core.CommenterOptions{EnableTraceparent: true}}
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

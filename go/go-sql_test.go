package sqlcommenter

import (
	"context"
	"net/http"
	"regexp"
	"testing"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func TestDisabled(t *testing.T) {

	db, _ := Open("mysql", "root:root@/gotest", CommenterOptions{})

	got := db.withComment(context.Background(), "Select 1")
	want := "Select 1"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestHTTP_Net(t *testing.T) {

	db, _ := Open("mysql", "root:root@/gotest", CommenterOptions{EnableDBDriver: true, EnableRoute: true, EnableFramework: true})

	r, _ := http.NewRequest("GET", "hello/1", nil)
	ctx := AddHttpRouterTags(r, context.Background())

	got := db.withComment(ctx, "Select 1")
	wantOrder1 := "Select 1/*driver=database%2Fsql,framework=net%2Fhttp,route=hello%2F1*/"
	wantOrder2 := "Select 1/*framework=net%2Fhttp,route=hello%2F1,driver=database%2Fsql*/"
	//wantOrder3 := "Select 1/*driver=database%2Fsql,framework=net%2Fhttp,route=hello%2F1*/"

	if (got != wantOrder1) && (got != wantOrder2) {
		t.Errorf("got %q, wanted1 %q or wanted2 %q", got, wantOrder1, wantOrder2)
	}
}

func TestQueryWithSemicolon(t *testing.T) {

	db, _ := Open("mysql", "root:root@/gotest", CommenterOptions{EnableDBDriver: true})

	got := db.withComment(context.Background(), "Select 1;")
	want := "Select 1/*driver=database%2Fsql*/;"

	if got != want {
		t.Errorf("got %q, wanted %q", got, want)
	}
}

func TestOtelIntegration(t *testing.T) {

	db, _ := Open("mysql", "root:root@/gotest", CommenterOptions{EnableTraceparent: true})

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

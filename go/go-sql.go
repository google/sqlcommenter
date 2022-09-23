package sqlcommenter

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strings"
)

type DB struct {
	*sql.DB
	CommenterOptions CommenterOptions
}

type CommenterOptions struct {
	DBDriver   bool
	Route      bool
	Framework  bool
	Controller bool
	Action     bool
}

func Open(driverName string, dataSourceName string, commenterOptions CommenterOptions) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	return &DB{DB: db, CommenterOptions: commenterOptions}, err
}

// ***** Query Functions *****

func (db *DB) Query(query string, args ...any) (*sql.Rows, error) {
	return db.DB.Query(db.withComment(context.Background(), query), args...)
}

func (db *DB) QueryRow(query string, args ...interface{}) *sql.Row {
	return db.DB.QueryRow(db.withComment(context.Background(), query), args...)
}

func (db *DB) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	return db.DB.QueryContext(ctx, db.withComment(ctx, query), args...)
}

func (db *DB) Exec(query string, args ...any) (sql.Result, error) {
	return db.DB.Exec(db.withComment(context.Background(), query), args...)
}

func (db *DB) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	return db.DB.ExecContext(ctx, db.withComment(ctx, query), args...)
}

func (db *DB) Prepare(query string) (*sql.Stmt, error) {
	return db.DB.Prepare(db.withComment(context.Background(), query))
}

func (db *DB) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	return db.DB.PrepareContext(ctx, db.withComment(ctx, query))
}

// ***** Query Functions *****

// ***** Framework Functions *****

func AddHttpRouterTags(r *http.Request, n any, p any) context.Context { // any type is set because we need to refrain from importing http-router package
	ctx := context.Background()
	ctx = context.WithValue(ctx, "route", r.URL.Path)
	ctx = context.WithValue(ctx, "action", getFunctionName(n))
	ctx = context.WithValue(ctx, "framework", "github.com/julienschmidt/httprouter")
	return ctx
}

// ***** Framework Functions *****

// ***** Commenter Functions *****

func (db *DB) withComment(ctx context.Context, query string) string {

	var finalCommentsMap = map[string]string{}
	var finalCommentsStr string = ""
	query = strings.TrimSpace(query)

	// Sorted alphabetically
	if db.CommenterOptions.Action && (ctx.Value("action") != nil) {
		finalCommentsMap["action"] = ctx.Value("action").(string)
	}

	if db.CommenterOptions.DBDriver {
		finalCommentsMap["driver"] = "go/sql"
	}

	if db.CommenterOptions.Framework && (ctx.Value("framework") != nil) {
		finalCommentsMap["framework"] = ctx.Value("framework").(string)
	}

	if db.CommenterOptions.Route && (ctx.Value("route") != nil) {
		finalCommentsMap["route"] = ctx.Value("route").(string)
	}

	if len(finalCommentsMap) > 0 { // Converts comments map to string and appends it to query
		finalCommentsStr = fmt.Sprintf("/*%s*/", convertMapToComment(finalCommentsMap))
		fmt.Println(finalCommentsStr)
	}

	if query[len(query)-1:] == ";" {
		return fmt.Sprintf("%s%s;", strings.TrimSuffix(query, ";"), finalCommentsStr)
	}
	return fmt.Sprintf("%s%s", query, finalCommentsStr)

}

// ***** Commenter Functions *****

// ***** Util Functions *****

func encodeURL(k string) string {
	return url.QueryEscape(string(k))
}

func getFunctionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func convertMapToComment(tags map[string]string) string {
	var sb strings.Builder
	i, sz := 0, len(tags)
	for key, val := range tags {
		if i == sz-1 {
			sb.WriteString(fmt.Sprintf("%s=%v", encodeURL(key), encodeURL(val)))
		} else {
			sb.WriteString(fmt.Sprintf("%s=%v,", encodeURL(key), encodeURL(val)))
		}
		i++
	}
	return sb.String()
}

// ***** Util Functions *****

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
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"sort"
	"strings"

	"go.opentelemetry.io/otel/propagation"
)

const (
	route       string = "route"
	controller  string = "controller"
	action      string = "action"
	framework   string = "framework"
	driver      string = "driver"
	traceparent string = "traceparent"
)

type DB struct {
	*sql.DB
	options CommenterOptions
}

type CommenterOptions struct {
	EnableDBDriver    bool
	EnableRoute       bool
	EnableFramework   bool
	EnableController  bool
	EnableAction      bool
	EnableTraceparent bool
}

func Open(driverName string, dataSourceName string, options CommenterOptions) (*DB, error) {
	db, err := sql.Open(driverName, dataSourceName)
	return &DB{DB: db, options: options}, err
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

func AddHttpRouterTags(r *http.Request, next any) context.Context { // any type is set because we need to refrain from importing http-router package
	ctx := context.Background()
	ctx = context.WithValue(ctx, route, r.URL.Path)
	ctx = context.WithValue(ctx, action, getFunctionName(next))
	ctx = context.WithValue(ctx, framework, "net/http")
	return ctx
}

// ***** Framework Functions *****

// ***** Commenter Functions *****

func (db *DB) withComment(ctx context.Context, query string) string {
	var commentsMap = map[string]string{}
	query = strings.TrimSpace(query)

	// Sorted alphabetically
	if db.options.EnableAction && (ctx.Value(action) != nil) {
		commentsMap[action] = ctx.Value(action).(string)
	}

	// `driver` information should not be coming from framework.
	// So, explicitly adding that here.
	if db.options.EnableDBDriver {
		commentsMap[driver] = "database/sql"
	}

	if db.options.EnableFramework && (ctx.Value(framework) != nil) {
		commentsMap[framework] = ctx.Value(framework).(string)
	}

	if db.options.EnableRoute && (ctx.Value(route) != nil) {
		commentsMap[route] = ctx.Value(route).(string)
	}

	if db.options.EnableTraceparent {
		carrier := extractTraceparent(ctx)
		if val, ok := carrier["traceparent"]; ok {
			commentsMap[traceparent] = val
		}
	}

	var commentsString string = ""
	if len(commentsMap) > 0 { // Converts comments map to string and appends it to query
		commentsString = fmt.Sprintf("/*%s*/", convertMapToComment(commentsMap))
	}

	// A semicolon at the end of the SQL statement means the query ends there.
	// We need to insert the comment before that to be considered as part of the SQL statemtent.
	if query[len(query)-1:] == ";" {
		return fmt.Sprintf("%s%s;", strings.TrimSuffix(query, ";"), commentsString)
	}
	return fmt.Sprintf("%s%s", query, commentsString)
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

func extractTraceparent(ctx context.Context) propagation.MapCarrier {
	// Serialize the context into carrier
	propgator := propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{})
	carrier := propagation.MapCarrier{}
	propgator.Inject(ctx, carrier)
	return carrier
}

// ***** Util Functions *****

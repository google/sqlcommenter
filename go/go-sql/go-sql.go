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
	"strings"

	"google.com/sqlcommenter/core"
)

type DB struct {
	*sql.DB
	options core.CommenterOptions
}

func Open(driverName string, dataSourceName string, options core.CommenterOptions) (*DB, error) {
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

// ***** Commenter Functions *****

func (db *DB) withComment(ctx context.Context, query string) string {
	var commentsMap = map[string]string{}
	query = strings.TrimSpace(query)

	// Sorted alphabetically
	if db.options.EnableAction && (ctx.Value(core.Action) != nil) {
		commentsMap[core.Action] = ctx.Value(core.Action).(string)
	}

	// `driver` information should not be coming from framework.
	// So, explicitly adding that here.
	if db.options.EnableDBDriver {
		commentsMap[core.Driver] = "database/sql"
	}

	if db.options.EnableFramework && (ctx.Value(core.Framework) != nil) {
		commentsMap[core.Framework] = ctx.Value(core.Framework).(string)
	}

	if db.options.EnableRoute && (ctx.Value(core.Route) != nil) {
		commentsMap[core.Route] = ctx.Value(core.Route).(string)
	}

	if db.options.EnableTraceparent {
		carrier := core.ExtractTraceparent(ctx)
		if val, ok := carrier["traceparent"]; ok {
			commentsMap[core.Traceparent] = val
		}
	}

	var commentsString string = ""
	if len(commentsMap) > 0 { // Converts comments map to string and appends it to query
		commentsString = fmt.Sprintf("/*%s*/", core.ConvertMapToComment(commentsMap))
	}

	// A semicolon at the end of the SQL statement means the query ends there.
	// We need to insert the comment before that to be considered as part of the SQL statemtent.
	if query[len(query)-1:] == ";" {
		return fmt.Sprintf("%s%s;", strings.TrimSuffix(query, ";"), commentsString)
	}
	return fmt.Sprintf("%s%s", query, commentsString)
}

// ***** Commenter Functions *****

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
	"database/sql/driver"
	"fmt"
	"runtime/debug"
	"strings"

	"github.com/google/sqlcommenter/go/core"
)

var attemptedToAutosetApplication = false

type sqlCommenterConn struct {
	driver.Conn
	options core.CommenterOptions
}

func newSQLCommenterConn(conn driver.Conn, options core.CommenterOptions) *sqlCommenterConn {
	return &sqlCommenterConn{
		Conn:    conn,
		options: options,
	}
}

func (s *sqlCommenterConn) Query(query string, args []driver.Value) (driver.Rows, error) {
	queryer, ok := s.Conn.(driver.Queryer)
	if !ok {
		return nil, driver.ErrSkip
	}
	ctx := context.Background()
	commentedQuery := s.withComment(ctx, query)
	return queryer.Query(commentedQuery, args)
}

func (s *sqlCommenterConn) QueryContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Rows, error) {
	queryer, ok := s.Conn.(driver.QueryerContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	commentedQuery := s.withComment(ctx, query)
	return queryer.QueryContext(ctx, commentedQuery, args)
}

func (s *sqlCommenterConn) Exec(query string, args []driver.Value) (driver.Result, error) {
	execor, ok := s.Conn.(driver.Execer)
	if !ok {
		return nil, driver.ErrSkip
	}
	ctx := context.Background()
	commentedQuery := s.withComment(ctx, query)
	return execor.Exec(commentedQuery, args)
}

func (s *sqlCommenterConn) ExecContext(ctx context.Context, query string, args []driver.NamedValue) (driver.Result, error) {
	execor, ok := s.Conn.(driver.ExecerContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	commentedQuery := s.withComment(ctx, query)
	return execor.ExecContext(ctx, commentedQuery, args)
}

func (s *sqlCommenterConn) PrepareContext(ctx context.Context, query string) (stmt driver.Stmt, err error) {
	preparer, ok := s.Conn.(driver.ConnPrepareContext)
	if !ok {
		return nil, driver.ErrSkip
	}
	commentedQuery := s.withComment(ctx, query)
	return preparer.PrepareContext(ctx, commentedQuery)
}

func (s *sqlCommenterConn) Raw() driver.Conn {
	return s.Conn
}

// ***** Commenter Functions *****

func (conn *sqlCommenterConn) withComment(ctx context.Context, query string) string {
	var commentsMap = map[string]string{}
	query = strings.TrimSpace(query)
	config := conn.options.Config

	// Sorted alphabetically
	if config.EnableAction && (ctx.Value(core.Action) != nil) {
		commentsMap[core.Action] = ctx.Value(core.Action).(string)
	}

	if config.EnableController && (ctx.Value(core.Controller) != nil) {
		commentsMap[core.Controller] = ctx.Value(core.Controller).(string)
	}

	// `driver` information should not be coming from framework.
	// So, explicitly adding that here.
	if config.EnableDBDriver {
		commentsMap[core.Driver] = fmt.Sprintf("database/sql:%s", conn.options.Tags.DriverName)
	}

	if config.EnableFramework && (ctx.Value(core.Framework) != nil) {
		commentsMap[core.Framework] = ctx.Value(core.Framework).(string)
	}

	if config.EnableRoute && (ctx.Value(core.Route) != nil) {
		commentsMap[core.Route] = ctx.Value(core.Route).(string)
	}

	if config.EnableTraceparent {
		carrier := core.ExtractTraceparent(ctx)
		if val, ok := carrier["traceparent"]; ok {
			commentsMap[core.Traceparent] = val
		}
	}

	if config.EnableApplication {
		if !attemptedToAutosetApplication && conn.options.Tags.Application == "" {
			attemptedToAutosetApplication = true
			bi, ok := debug.ReadBuildInfo()
			if ok {
				conn.options.Tags.Application = bi.Path
			}
		}
		if conn.options.Tags.Application != "" {
			commentsMap[core.Application] = conn.options.Tags.Application
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

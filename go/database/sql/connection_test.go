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
	"testing"

	"github.com/google/sqlcommenter/go/core"
)

func TestWithComment_NoContext(t *testing.T) {
	testBasicConn := &mockConn{}
	testCases := []struct {
		desc             string
		commenterOptions core.CommenterOptions
		query            string
		wantQuery        string
	}{
		{
			desc:             "empty commenter options",
			commenterOptions: core.CommenterOptions{},
			query:            "SELECT 1;",
			wantQuery:        "SELECT 1;",
		},
		{
			desc: "only enable DBDriver",
			commenterOptions: core.CommenterOptions{
				Config: core.CommenterConfig{EnableDBDriver: true},
			},
			query:     "SELECT 1;",
			wantQuery: "SELECT 1/*db_driver=database%2Fsql%3A*/;",
		},
		{
			desc: "enable DBDriver and pass static tag driver name",
			commenterOptions: core.CommenterOptions{
				Config: core.CommenterConfig{EnableDBDriver: true},
				Tags:   core.StaticTags{DriverName: "postgres"},
			},
			query:     "SELECT 1;",
			wantQuery: "SELECT 1/*db_driver=database%2Fsql%3Apostgres*/;",
		},
		{
			desc: "enable DBDriver and pass all static tags",
			commenterOptions: core.CommenterOptions{
				Config: core.CommenterConfig{EnableDBDriver: true},
				Tags:   core.StaticTags{DriverName: "postgres", Application: "app-1"},
			},
			query:     "SELECT 1;",
			wantQuery: "SELECT 1/*db_driver=database%2Fsql%3Apostgres*/;",
		},
		{
			desc: "enable other tags and pass all static tags",
			commenterOptions: core.CommenterOptions{
				Config: core.CommenterConfig{EnableDBDriver: true, EnableApplication: true, EnableFramework: true},
				Tags:   core.StaticTags{DriverName: "postgres", Application: "app-1"},
			},
			query:     "SELECT 1;",
			wantQuery: "SELECT 1/*application=app-1,db_driver=database%2Fsql%3Apostgres*/;",
		},
	}
	for _, tc := range testCases {
		testConn := newSQLCommenterConn(testBasicConn, tc.commenterOptions)
		ctx := context.Background()
		if got, want := testConn.withComment(ctx, tc.query), tc.wantQuery; got != want {
			t.Errorf("testConn.withComment(ctx, %q) = %q, want = %q", tc.query, got, want)
		}
	}
}

func TestWithComment_WithContext(t *testing.T) {
	testBasicConn := &mockConn{}
	testCases := []struct {
		desc             string
		commenterOptions core.CommenterOptions
		ctx              context.Context
		query            string
		wantQuery        string
	}{
		{
			desc:             "empty commenter options",
			commenterOptions: core.CommenterOptions{},
			ctx: getContextWithKeyValue(
				map[string]string{
					"route":     "listData",
					"framework": "custom-golang",
				},
			),
			query:     "SELECT 1;",
			wantQuery: "SELECT 1;",
		},
		{
			desc: "only all options but context has few tags",
			commenterOptions: core.CommenterOptions{
				Config: core.CommenterConfig{
					EnableDBDriver:    true,
					EnableRoute:       true,
					EnableFramework:   true,
					EnableController:  true,
					EnableAction:      true,
					EnableTraceparent: true,
					EnableApplication: true,
				},
				Tags: core.StaticTags{DriverName: "postgres", Application: "app-1"},
			},
			ctx: getContextWithKeyValue(
				map[string]string{
					"route":     "listData",
					"framework": "custom-golang",
				},
			),
			query:     "SELECT 1;",
			wantQuery: "SELECT 1/*application=app-1,db_driver=database%2Fsql%3Apostgres,framework=custom-golang,route=listData*/;",
		},
		{
			desc: "only all options but context contains all tags",
			commenterOptions: core.CommenterOptions{
				Config: core.CommenterConfig{
					EnableDBDriver:    true,
					EnableRoute:       true,
					EnableFramework:   true,
					EnableController:  true,
					EnableAction:      true,
					EnableTraceparent: true,
					EnableApplication: true,
				},
				Tags: core.StaticTags{DriverName: "postgres", Application: "app-1"},
			},
			ctx: getContextWithKeyValue(
				map[string]string{
					"route":      "listData",
					"framework":  "custom-golang",
					"controller": "custom-controller",
					"action":     "any action",
				},
			),
			query:     "SELECT 1;",
			wantQuery: "SELECT 1/*action=any+action,application=app-1,db_driver=database%2Fsql%3Apostgres,framework=custom-golang,route=listData*/;",
		},
	}
	for _, tc := range testCases {
		testConn := newSQLCommenterConn(testBasicConn, tc.commenterOptions)
		if got, want := testConn.withComment(tc.ctx, tc.query), tc.wantQuery; got != want {
			t.Errorf("testConn.withComment(ctx, %q) = %q, want = %q", tc.query, got, want)
		}
	}
}

func getContextWithKeyValue(vals map[string]string) context.Context {
	ctx := context.Background()
	for k, v := range vals {
		ctx = context.WithValue(ctx, k, v)
	}
	return ctx
}

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
	"database/sql/driver"
)

type mockConn struct {
	prepareStmt driver.Stmt
	prepareErr error

	closeErr error

	beginTx driver.Tx
	beginErr error
}

func (c *mockConn) Prepare(query string) (driver.Stmt, error) {
	if c.prepareErr != nil {
		return nil, c.prepareErr
	}
	return c.prepareStmt, nil
}

func (c *mockConn) Close() error {
	return c.closeErr
}

func (c *mockConn) Begin() (driver.Tx, error) {
	if c.beginErr != nil {
		return nil, c.beginErr
	}
	return c.beginTx, nil 
}

type mockDriver struct {
	conn driver.Conn
	openError error
}

func (d *mockDriver) Open(name string) (driver.Conn, error) {
	if d.openError != nil {
		return nil, d.openError
	}
	return d.conn, nil
}


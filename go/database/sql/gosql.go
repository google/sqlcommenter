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
	"database/sql"
	"database/sql/driver"

	"github.com/google/sqlcommenter/go/core"
)

var (
	_ driver.Driver        = (*sqlCommenterDriver)(nil)
	_ driver.DriverContext = (*sqlCommenterDriver)(nil)
	_ driver.Connector     = (*sqlCommenterConnector)(nil)
)

// SQLCommenterDriver returns a driver object that contains SQLCommenter drivers.
type sqlCommenterDriver struct {
	driver  driver.Driver
	options core.CommenterOptions
}

func newSQLCommenterDriver(dri driver.Driver, options core.CommenterOptions) *sqlCommenterDriver {
	return &sqlCommenterDriver{driver: dri, options: options}
}

func (d *sqlCommenterDriver) Open(name string) (driver.Conn, error) {
	rawConn, err := d.driver.Open(name)
	if err != nil {
		return nil, err
	}
	return newConn(rawConn, d.options), nil
}

func (d *sqlCommenterDriver) OpenConnector(name string) (driver.Connector, error) {
	rawConnector, err := d.driver.(driver.DriverContext).OpenConnector(name)
	if err != nil {
		return nil, err
	}
	return newConnector(rawConnector, d, d.options), err
}

type sqlCommenterConnector struct {
	driver.Connector
	driver  *sqlCommenterDriver
	options core.CommenterOptions
}

func newConnector(connector driver.Connector, driver *sqlCommenterDriver, options core.CommenterOptions) *sqlCommenterConnector {
	return &sqlCommenterConnector{
		Connector: connector,
		driver:    driver,
		options:   options,
	}
}

func (c *sqlCommenterConnector) Connect(ctx context.Context) (connection driver.Conn, err error) {
	connection, err = c.Connector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return newConn(connection, c.options), nil
}

func (c *sqlCommenterConnector) Driver() driver.Driver {
	return c.driver
}

type dsnConnector struct {
	dsn    string
	driver driver.Driver
}

func (t dsnConnector) Connect(_ context.Context) (driver.Conn, error) {
	return t.driver.Open(t.dsn)
}

func (t dsnConnector) Driver() driver.Driver {
	return t.driver
}

// Open is a wrapper over sql.Open with OTel instrumentation.
func Open(driverName, dataSourceName string, options core.CommenterOptions) (*sql.DB, error) {
	// Retrieve the driver implementation we need to wrap with instrumentation
	db, err := sql.Open(driverName, "")
	if err != nil {
		return nil, err
	}
	d := db.Driver()
	if err = db.Close(); err != nil {
		return nil, err
	}

	options.Tags.DriverName = driverName
	sqlCommenterDriver := newSQLCommenterDriver(d, options)

	if _, ok := d.(driver.DriverContext); ok {
		connector, err := sqlCommenterDriver.OpenConnector(dataSourceName)
		if err != nil {
			return nil, err
		}
		return sql.OpenDB(connector), nil
	}

	return sql.OpenDB(dsnConnector{dsn: dataSourceName, driver: sqlCommenterDriver}), nil
}

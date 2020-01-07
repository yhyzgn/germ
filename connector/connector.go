// Copyright 2019 yhyzgn germ
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// author : 颜洪毅
// e-mail : yhyzgn@gmail.com
// time   : 2020-01-06 15:53
// version: 1.0.0
// desc   : 

package connector

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/yhyzgn/germ/external"
	"github.com/yhyzgn/germ/util"
)

type Connector struct {
	dialect external.Dialect
	source  *external.DataSource
	DB      *sql.DB
}

func Connect(dlt external.Dialect, src *external.DataSource) (*Connector, error) {
	cn := &Connector{
		dialect: dlt,
		source:  src,
	}
	db, err := sql.Open(cn.dialect.Name(), cn.url())
	if err != nil {
		return nil, err
	}
	if cn.source.MaxOpen == 0 {
		cn.source.MaxOpen = 20
	}
	if cn.source.MaxIdle == 0 {
		cn.source.MaxIdle = 10
	}
	db.SetMaxOpenConns(cn.source.MaxOpen)
	db.SetMaxIdleConns(cn.source.MaxIdle)
	cn.DB = db

	// 当前方言
	Current = cn

	return cn, nil
}

func (cn *Connector) Ping() error {
	return cn.DB.Ping()
}

func (cn *Connector) Close() error {
	if cn.DB == nil {
		return errors.New("Can not close a closed connection.")
	}
	return cn.DB.Close()
}

func (cn *Connector) Dialect() external.Dialect {
	return cn.dialect
}

func (cn *Connector) DataSource() *external.DataSource {
	return cn.source
}

func (cn *Connector) Query(sql string, args ...interface{}) (*sql.Rows, error) {
	//logger.SQL(sql, args...)
	return cn.DB.Query(sql, args...)
}

func (cn *Connector) Exec(sql string, args ...interface{}) (sql.Result, error) {
	return cn.DB.Exec(sql, args...)
}

func (cn *Connector) url() string {
	return fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?%v", cn.source.Username, cn.source.Password, cn.source.Host, cn.source.Port, cn.source.DBName, util.JoinMapParams(cn.source.Params))
}

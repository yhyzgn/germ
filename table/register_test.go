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
// time   : 2020-01-07 9:35
// version: 1.0.0
// desc   : 

package table

import (
	"github.com/yhyzgn/germ/connector"
	"github.com/yhyzgn/germ/dialect/mysql"
	"github.com/yhyzgn/germ/external"
	"github.com/yhyzgn/germ/logger"
	"testing"
	"time"
)

type Test struct {
	ID         int64     `germ:"column:id;primary;"`
	Name       string    `germ:"column:name;default:null"`
	Age        int       `germ:"type:int;"`
	CreateTime time.Time `germ:"column:create_time"`
	Common
}

func (Test) TableName() string {
	return "test"
}

func TestRegister(t *testing.T) {
	conn, err := connector.Connect(&mysql.Dialect{}, &external.DataSource{
		Host:     "localhost",
		Port:     3306,
		Username: "root",
		Password: "root",
		DBName:   "germ",
		MaxOpen:  20,
		MaxIdle:  10,
		Params: map[string]interface{}{
			"charset": "utf8",
			"useSSL":  false,
		},
	})
	if err != nil {
		logger.Fatal(err)
		return
	}
	logger.Info(conn)

	test := &Test{}
	err = Register(test)
	if err != nil {
		logger.Fatal(err)
		return
	}
}

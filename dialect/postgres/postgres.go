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
// time   : 2020-01-09 18:18
// version: 1.0.0
// desc   : 

package postgres

import (
	"database/sql"
	"github.com/yhyzgn/germ/dialect"
	"github.com/yhyzgn/germ/errors"
	"github.com/yhyzgn/germ/external"
)

type Postgres struct {
	dialect.Adapter
}

func (*Postgres) Name() string {
	return "postgres"
}

func (*Postgres) Insert(executor external.Executor, command *external.SQLCommand) (sql.Result, error) {
	var (
		res rawResult
		err error
	)
	if scanErr := executor.QueryRow(command.SQL, command.Args...).Scan(&res.ID); scanErr == sql.ErrNoRows {
		res.Err = scanErr
	} else {
		err = scanErr
	}
	return res, err
}

func (*Postgres) Update(executor external.Executor, command *external.SQLCommand) (sql.Result, error) {
	result, err := executor.Exec(command.SQL, command.Args...)
	if err != nil {
		return nil, err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return nil, err
	}
	if affected == 0 {
		return nil, errors.UpdateError{}
	}
	return result, nil
}

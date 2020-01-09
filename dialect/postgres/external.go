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
// time   : 2020-01-09 18:22
// version: 1.0.0
// desc   : 

package postgres

import "github.com/yhyzgn/germ/errors"

type rawResult struct {
	ID  int64
	Err error
}

func (t rawResult) LastInsertId() (int64, error) {
	return t.ID, t.Err
}

func (t rawResult) RowsAffected() (int64, error) {
	return 0, errors.NewAffectedError(t.Err.Error())
}

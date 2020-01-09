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
// time   : 2020-01-06 15:44
// version: 1.0.0
// desc   : 

package external

import (
	"database/sql"
	"reflect"
)

type Dialect interface {
	Name() string

	Quote(key string) string

	TypeToSQLType(tp reflect.Type) string

	// 有些数据库不支持返回最后插入id，如：postgres，此时需要手动查询
	Insert(executor Executor, command *SQLCommand) (sql.Result, error)

	// 有些数据库根据 RowsAffected 来判断是否更新成功，
	// 而有些不能，如：MySQL，只要数据无变化，RowsAffected就是0，此时只能通过SQL语句执行是否成功来判断
	Update(executor Executor, command *SQLCommand) (sql.Result, error)

	HasTable(tableName string) *SQLCommand

	TableColumns(tableName string) *SQLCommand

	CreateTable(model *Model, dialect Dialect) []*SQLCommand

	ModifyTable(model *Model, dialect Dialect) []*SQLCommand

	DropTable(model *Model, dialect Dialect) *SQLCommand
}

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
// time   : 2020-01-07 16:11
// version: 1.0.0
// desc   : 

package table

import (
	"fmt"
	"github.com/yhyzgn/germ/external"
	"github.com/yhyzgn/germ/external/table/primary/strategy"
	"github.com/yhyzgn/germ/logger"
	"reflect"
)

func Exist(tableName string) *external.SQLCommand {
	return dialect.TableExistSQL(source.DBName, tableName)
}

func Columns(tableName string) *external.SQLCommand {
	return dialect.ColumnsOfTableSQL(source.DBName, tableName)
}

func Create(tableName string) *external.SQLCommand {
	info, ok := cacheTableName[tableName]
	if !ok {
		logger.ErrorF("The table [%v] has not been registered.")
		return nil
	}
	if info.Struct.Fields == nil || len(info.Struct.Fields) == 0 {
		return nil
	}

	var (
		primaryKey string
	)

	ln := len(info.Struct.Fields)

	cmd := external.NewCommand(fmt.Sprintf("CREATE TABLE %v.%v (", dialect.Quote(source.DBName), dialect.Quote(tableName)))
	for idx, field := range info.Struct.Fields {
		fmt.Println(field)
		cmd.LineTab(dialect.Quote(field.Column), 1).Append(field.SQLType)
		if field.NotNull {
			cmd.Append("NOT")
		}
		cmd.Append("NULL")

		// 主键
		if field.IsPrimary {
			primaryKey = field.Column

			stg := info.Strategy
			if stg != nil && reflect.TypeOf(stg).Elem() == reflect.TypeOf(strategy.AutoIncrement{}) {
				// AutoIncrement
				cmd.Append("AUTO_INCREMENT")
			}
		} else if field.Default != nil {
			cmd.Append("DEFAULT").Append(field.Default.(string))
		}

		if field.Comment != "" {
			cmd.Append("COMMENT").Append(fmt.Sprintf("'%v'", field.Comment))
		}
		if idx < ln-1 {
			cmd.Link(",")
		}

	}
	if primaryKey != "" {
		cmd.Link(",")
	}
	cmd.LineTab(fmt.Sprintf("PRIMARY KEY (%v) USING BTREE", dialect.Quote(primaryKey)), 1)
	cmd.Line(")")
	return cmd
}

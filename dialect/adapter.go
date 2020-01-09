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
// time   : 2020-01-06 16:08
// version: 1.0.0
// desc   : 

package dialect

import (
	"database/sql"
	"fmt"
	"github.com/yhyzgn/germ/external"
	"github.com/yhyzgn/germ/external/table/primary/strategy"
	"reflect"
)

// Adapter 数据库方言适配器
//
// 默认为 MySQL 驱动
//
// 其他方言则通过继承该适配器并重写方法实现
type Adapter struct {
}

func (*Adapter) Name() string {
	return "mysql"
}

func (*Adapter) Quote(key string) string {
	return fmt.Sprintf("`%s`", key)
}

func (a *Adapter) TypeToSQLType(tp reflect.Type) string {
	switch tp.Kind() {
	case reflect.Ptr:
		return a.TypeToSQLType(tp.Elem())
	case reflect.Int8, reflect.Uint8:
		// 1 字节/8 位bit
		// 默认长度为：6
		return "TINYINT"
	case reflect.Int16, reflect.Uint16:
		// 2 字节/16 位bit
		// 默认长度为：9
		return "SMALLINT"
	case reflect.Int, reflect.Uint, reflect.Int32, reflect.Uint32:
		// 2 字节/16 位bit
		// 默认长度为：11
		return "INT"
	case reflect.Int64, reflect.Uint64:
		// 8 字节/64 位bit
		// 默认长度为：20
		return "BIGINT"
	case reflect.String:
		return "VARCHAR(255)"
	case reflect.Bool:
		// 0 -> false
		// 1 -> true
		return "TINYINT(1)"
	case reflect.Float32:
		return "FLOAT"
	case reflect.Float64:
		return "DOUBLE"
	case reflect.Struct:
		switch tp.String() {
		case "time.Time":
			return "DATETIME"
		}
	}
	if _, ok := reflect.New(tp).Elem().Interface().([]byte); ok {
		return "VARCHAR(255)"
	}
	return ""
}

func (*Adapter) Insert(executor external.Executor, command *external.SQLCommand) (sql.Result, error) {
	return executor.Exec(command.SQL, command.Args...)
}

func (*Adapter) Update(executor external.Executor, command *external.SQLCommand) (sql.Result, error) {
	return executor.Exec(command.SQL, command.Args...)
}

func (*Adapter) HasTable(tableName string) *external.SQLCommand {
	return external.
		NewCommand("SELECT").
		LineTab("COUNT(1)", 1).
		Line("FROM").
		LineTab("INFORMATION_SCHEMA.TABLES", 1).
		Line("WHERE").
		LineTab("TABLE_SCHEMA = (SELECT DATABASE()) AND", 1).
		LineTab("TABLE_NAME = ?", 1).
		Arguments(tableName)
}

func (*Adapter) TableColumns(tableName string) *external.SQLCommand {
	return external.
		NewCommand("SELECT").
		LineTab("*", 1).
		Line("FROM").
		LineTab("information_schema.COLUMNS", 1).
		Line("WHERE").
		LineTab("TABLE_SCHEMA = (SELECT DATABASE()) AND", 1).
		LineTab("TABLE_NAME = ?", 1).
		Line("ORDER BY").
		LineTab("ORDINAL_POSITION ASC", 1).
		Arguments(tableName)
}

func (*Adapter) CreateTable(model *external.Model, dialect external.Dialect) []*external.SQLCommand {
	var (
		primaryKey string
	)

	result := make([]*external.SQLCommand, 0)
	indexes := make([]*external.Index, 0)

	ln := len(model.Fields)

	cmd := external.NewCommand(fmt.Sprintf("CREATE TABLE %v (", dialect.Quote(model.TableName)))
	for idx, field := range model.Fields {
		cmd.LineTab(dialect.Quote(field.Column), 1).Append(field.SQLType)
		if field.NotNull {
			cmd.Append("NOT")
		}
		cmd.Append("NULL")

		// 主键
		if field.IsPrimary {
			primaryKey = field.Column

			stg := model.Strategy
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

		// 记录 索引
		if field.Indexes != nil && len(field.Indexes) > 0 {
			indexes = append(indexes, field.Indexes...)
		}
	}
	if primaryKey != "" {
		cmd.Link(",")
		cmd.LineTab(fmt.Sprintf("PRIMARY KEY (%v)", dialect.Quote(primaryKey)), 1)
	}

	// 添加索引
	for _, index := range indexes {
		cmd.Link(",").LineTab("", 1)
		switch index.Type {
		case external.IndexUnique:
			cmd.Link("UNIQUE ")
			break
		case external.IndexFullText:
			cmd.Link("FULLTEXT ")
			break
		case external.IndexSpatial:
			cmd.Link("SPATIAL ")
			break
		}
		cmd.Link(fmt.Sprintf("INDEX %v (%v)", dialect.Quote(index.Name), dialect.Quote(index.Column)))
	}

	cmd.Line(")")

	// 基本建表命令
	result = append(result, cmd)

	return result
}

func (*Adapter) ModifyTable(model *external.Model, dialect external.Dialect) []*external.SQLCommand {
	result := make([]*external.SQLCommand, 0)
	// TODO modify
	return result
}

func (*Adapter) DropTable(model *external.Model, dialect external.Dialect) *external.SQLCommand {
	return external.NewCommand(fmt.Sprintf("DROP TABLE %v", dialect.Quote(model.TableName)))
}

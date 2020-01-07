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
// time   : 2020-01-06 16:58
// version: 1.0.0
// desc   : 

package table

import (
	"errors"
	"fmt"
	"github.com/yhyzgn/germ/connector"
	"github.com/yhyzgn/germ/external"
	"github.com/yhyzgn/germ/util"
	"reflect"
	"strings"
	"sync"
)

var (
	once            sync.Once
	cacheTableName  map[string]Info   // 表名：类信息
	cacheStructName map[string]string // 类名：表名
)

func init() {
	once.Do(func() {
		cacheTableName = make(map[string]Info)
		cacheStructName = make(map[string]string)
	})
}

func Register(tl Table) (err error) {
	if tl == nil {
		err = errors.New("The table can not be nil pointer reference.")
		return
	}

	if connector.Current == nil {
		err = errors.New("Must connect to sql driver at first.")
		return
	}

	hasPrimary := false

	// 类基本信息
	tp := reflect.TypeOf(tl)
	elmTp := util.GetEleType(tp)

	tableName := tl.TableName()
	if tableName == "" {
		tableName = elmTp.Name()
	}
	if _, ok := cacheTableName[tableName]; ok {
		err = fmt.Errorf("The struct [%v] has been regiestered.", elmTp.String())
		return
	}

	// 字段们
	fields := make([]Field, 0)

	count := elmTp.NumField()
	var (
		fld   reflect.StructField
		field Field
		tags  map[string]string
	)
	for i := 0; i < count; i++ {
		fld = elmTp.Field(i)

		// 字段信息
		field = Field{
			Name:    fld.Name,
			Type:    fld.Type,
			ELmType: util.GetEleType(fld.Type),
		}

		// 标签信息
		tags = util.GetTagMap(external.TagGerm, fld.Tag)
		if tags == nil {
			field.Ignored = true
		} else {
			// 列名
			column, ok := tags[external.KeyColumn]
			if !ok || column == "" {
				column = field.Name
			}
			field.Column = column

			// 是否是主键
			_, isPrimary := tags[external.KeyIsPrimary]
			if isPrimary {
				if hasPrimary {
					err = fmt.Errorf("Only allow one primary key in one table [%v].", fld.Name)
					return
				}
				field.IsPrimary = true
				field.Nullable = false
				hasPrimary = true
			}

			// 字段类型
			sqlType, ok := tags[external.KeySQLType]
			if ok && sqlType != "" {
				field.SQLType = strings.ToUpper(sqlType)
			} else {
				// 根据类中字段类型自动推导数据库类型
				sqlType = connector.Current.Dialect().TypeToSQLType(field.ELmType)
				if sqlType == "" {
					err = fmt.Errorf("Unknown SQL type for field [%v] in struct [%v].", fld.Name, elmTp.String())
					return
				}
				field.SQLType = sqlType
			}

			// 索引
			index, ok := tags[external.KeyIndex]
			if ok {
				field.Index = index
			}

			// 是否可空
			_, nullable := tags[external.KeyNullable]
			// 主键不能为空
			if nullable && field.IsPrimary {
				err = fmt.Errorf("The primary key field [%v] can not be null in struct [%v].", fld.Name, elmTp.String())
				return
			}
			field.Nullable = nullable

			// 默认值
			dft, ok := tags[external.KeyDefault]
			if ok {
				// 自动设置是否可空
				if strings.ToUpper(dft) == "NULL" {
					// 主键不能为空
					if field.IsPrimary {
						err = fmt.Errorf("The primary key field [%v] can not be null in struct [%v].", fld.Name, elmTp.String())
						return
					}
					field.Nullable = true
				}
				field.Default = dft
			}

			// 注释
			comment, ok := tags[external.KeyComment]
			if ok {
				field.Comment = comment
			}

			fields = append(fields, field)
		}

		// 表名：类信息
		cacheTableName[tableName] = Info{
			Name: tableName,
			Struct: Struct{
				Type:    tp,
				ELmType: elmTp,
				Fields:  fields,
			},
		}
		// 类名 : 表名
		cacheStructName[elmTp.String()] = tableName
	}
	return
}

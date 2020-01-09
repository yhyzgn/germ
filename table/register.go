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
	cacheTableName  map[string]*external.Model // 表名：类信息
	cacheStructName map[string]string          // 类名：表名
	dialect         external.Dialect
	source          *external.DataSource
)

func init() {
	once.Do(func() {
		cacheTableName = make(map[string]*external.Model)
		cacheStructName = make(map[string]string)
	})
}

func Register(tl external.Table) (err error) {
	if tl == nil {
		err = errors.New("The table can not be nil pointer reference.")
		return
	}

	if connector.Current == nil {
		err = errors.New("Must connect to sql connector at first.")
		return
	}
	dialect = connector.Current.Dialect()
	source = connector.Current.DataSource()

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
	fields := make([]*external.Field, 0)

	count := elmTp.NumField()
	var (
		fld   reflect.StructField
		field *external.Field
		tags  map[string]string
	)
	for i := 0; i < count; i++ {
		fld = elmTp.Field(i)

		// 标签信息
		tags = util.GetTagMap(external.TagGerm, fld.Tag)
		if tags != nil {
			// 字段信息
			field = &external.Field{
				Name:    fld.Name,
				Type:    fld.Type,
				ELmType: util.GetEleType(fld.Type),
				Indexes: make([]*external.Index, 0),
			}

			// 列名
			column, ok := tags[external.KeyColumn]
			if !ok || column == "" {
				column = field.Name
			}
			field.Column = column

			// 是否是主键
			_, isPrimary := tags[external.KeyPrimary]
			if isPrimary {
				if hasPrimary {
					err = fmt.Errorf("Only allow one primary key in one table [%v].", fld.Name)
					return
				}
				field.IsPrimary = true
				field.NotNull = true
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

			// 是否不可空
			if !field.IsPrimary {
				_, notnull := tags[external.KeyNotNull]
				field.NotNull = notnull
			}

			// 默认值
			dft, ok := tags[external.KeyDefault]
			if ok {
				if strings.ToUpper(dft) == "NULL" {
					// 主键不能为空
					if field.IsPrimary {
						err = fmt.Errorf("The primary key field [%v] can not be null in struct [%v], but set default value is null.", fld.Name, elmTp.String())
						return
					}
					// 冲突
					if field.NotNull {
						err = fmt.Errorf("The field [%v] can not be null in struct [%v], but set default value is null.", fld.Name, elmTp.String())
						return
					}
				}
				field.Default = dft
			}

			// 注释
			comment, ok := tags[external.KeyComment]
			if ok {
				field.Comment = comment
			}

			// 普通索引
			index(field, tags, external.KeyIndex)

			// 唯一索引
			index(field, tags, external.KeyUnique)

			// 全文索引
			index(field, tags, external.KeyFullText)

			// 空间索引
			index(field, tags, external.KeySpatial)

			fields = append(fields, field)
		}
	}

	// 表名：类信息
	cacheTableName[tableName] = &external.Model{
		TableName: tableName,
		Struct: external.Struct{
			Type:    tp,
			ELmType: elmTp,
		},
		Strategy: tl.PrimaryStrategy(),
		Fields:   fields,
	}
	// 类名 : 表名
	cacheStructName[elmTp.String()] = tableName

	// 检查数据库表信息
	CheckTable(tableName)
	return
}

func index(field *external.Field, tags map[string]string, keyName string) {
	// 空间索引
	name, ok := tags[keyName]
	if ok {
		if name == "" {
			// 默认为列名
			name = field.Column
		}

		index := &external.Index{
			Name:   name,
			Column: field.Column,
		}

		switch keyName {
		case external.KeyUnique:
			index.Type = external.IndexUnique
			break
		case external.KeyFullText:
			index.Type = external.IndexFullText
			break
		case external.KeySpatial:
			index.Type = external.IndexSpatial
			break
		default:
			index.Type = external.IndexNormal
			break
		}

		field.Indexes = append(field.Indexes, index)
	}
}

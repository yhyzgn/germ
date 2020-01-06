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
	"github.com/yhyzgn/germ/util"
	"reflect"
	"sync"
)

var (
	once  sync.Once
	cache map[string]Info
)

func init() {
	once.Do(func() {
		cache = make(map[string]Info)
	})
}

func Register(tl Table) {
	if tl == nil {
		return
	}

	// 类基本信息
	tp := reflect.TypeOf(tl)
	elmTp := util.GetEleType(tp)

	// 字段们
	fields := make([]Field, 0)

	count := elmTp.NumField()
	var (
		fld   reflect.StructField
		field Field
	)
	for i := 0; i < count; i++ {
		fld = elmTp.Field(i)

		// 字段信息
		field.Name = fld.Name
		field.Type = fld.Type
		field.ELmType = util.GetEleType(fld.Type)

		// 标签信息
	}

	// 缓存
	cache[tl.Name()] = Info{
		Name: tl.Name(),
		Struct: Struct{
			Type:    tp,
			ELmType: elmTp,
			Fields:  fields,
		},
	}
}

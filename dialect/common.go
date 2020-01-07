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
	"fmt"
	"reflect"
)

type Common struct {
}

func (*Common) Name() string {
	return ""
}

func (*Common) Quote(key string) string {
	return fmt.Sprintf(`"%s"`, key)
}

func (*Common) TypeToSQLType(tp reflect.Type) string {
	switch tp.Kind() {
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
	return ""
}

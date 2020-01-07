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
// time   : 2020-01-06 18:10
// version: 1.0.0
// desc   : 

package util

import (
	"reflect"
	"strings"
)

func GetEleType(tp reflect.Type) reflect.Type {
	if tp.Kind() == reflect.Ptr || tp.Kind() == reflect.Slice {
		return tp.Elem()
	}
	return tp
}

func GetTag(name string, tag reflect.StructTag) (value string, ok bool) {
	value, ok = tag.Lookup(name)
	return
}

func GetTagMap(name string, tag reflect.StructTag) map[string]string {
	value, ok := GetTag(name, tag)
	if ok {
		// 去除所有空格
		value = strings.ReplaceAll(value, " ", "")
		tags := strings.Split(value, ";")
		if len(tags) > 0 {
			result := make(map[string]string)
			for _, tag := range tags {
				if tag == "" {
					continue
				}
				if strings.Contains(tag, ":") {
					temp := strings.Split(tag, ":")
					result[temp[0]] = temp[1]
				} else {
					result[tag] = ""
				}
			}
			return result
		}
	}
	return nil
}
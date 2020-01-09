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
// time   : 2020-01-07 16:15
// version: 1.0.0
// desc   : 

package table

import (
	"fmt"
	"github.com/yhyzgn/germ/connector"
	"github.com/yhyzgn/germ/logger"
)

func CheckTable(tableName string) {
	hasCMD := HasTable(tableName)

	rows, err := connector.Current.Query(hasCMD.SQL, hasCMD.Args...)
	if err != nil {
		logger.Error(err)
		return
	}
	defer rows.Close()

	// 表是否存在
	if exist := rows.Next(); exist {
		var count int
		err = rows.Scan(&count)
		if err == nil && count == 1 {
			// 已存在，检查结构是否同步
			logger.Info(tableName + "表已存在")
		} else {
			// 不存在，生成表
			create := Create(tableName)
			fmt.Println(create[0].SQL)
			//logger.SQL(create.SQL)
			//
			//res, err := connector.Current.Exec(create.SQL)
			//if err != nil {
			//	logger.Error(err)
			//	return
			//}
			//logger.Info(res)
		}
	}
}

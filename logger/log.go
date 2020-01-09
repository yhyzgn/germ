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
// time   : 2020-01-07 10:13
// version: 1.0.0
// desc   : 

package logger

import (
	"os"
	"strings"
	"sync"
)

var (
	once    sync.Once
	showSQL bool
	lgr     Logger
)

func init() {
	once.Do(func() {
		if lgr == nil {
			lgr = logger{}
			showSQL = true
		}
	})
}

func SQL(sql string, args ...interface{}) {
	if showSQL {
		lgr.InfoF("\n"+strings.ReplaceAll(sql, "?", "%v"), args...)
	}
}

func Info(info interface{}) {
	lgr.Info(info)
}

func InfoF(format string, args ...interface{}) {
	lgr.InfoF(format, args...)
}

func Error(err interface{}) {
	lgr.Error(err)
}

func ErrorF(format string, args ...interface{}) {
	lgr.ErrorF(format, args...)
}

func Fatal(fatal interface{}) {
	lgr.Error(fatal)
	os.Exit(1)
}

func FatalF(format string, fts ...interface{}) {
	lgr.ErrorF(format, fts...)
	os.Exit(1)
}

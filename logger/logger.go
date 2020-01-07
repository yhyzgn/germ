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
// time   : 2020-01-07 10:17
// version: 1.0.0
// desc   : 

package logger

import "log"

type Logger interface {
	Info(info interface{})
	InfoF(format string, args ...interface{})
	Error(err interface{})
	ErrorF(format string, args ...interface{})
}

type common struct {
}

func (common) Info(info interface{}) {
	log.Println(info)
}

func (common) InfoF(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func (common) Error(err interface{}) {
	log.Println(err)
}

func (common) ErrorF(format string, args ...interface{}) {
	log.Printf(format, args...)
}
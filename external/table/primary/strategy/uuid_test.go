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
// time   : 2020-01-06 17:48
// version: 1.0.0
// desc   : 

package strategy

import (
	"fmt"
	"testing"
)

func TestUUID_Auto(t *testing.T) {
	uid := &UUID{}

	fmt.Println(uid.Auto())
	fmt.Println(uid.Auto())
	fmt.Println(uid.Auto())
	fmt.Println(uid.Auto())
}

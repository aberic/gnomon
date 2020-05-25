/*
 * Copyright (c) 2020. Aberic - All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 * http://www.apache.org/licenses/LICENSE-2.0
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package grope

import (
	"net/url"
	"testing"
)

func TestURLParseQuery(t *testing.T) {
	// 参数 encode = %E5%8F%82%E6%95%B0
	urlStr := "c=c&%E5%8F%82%E6%95%B0=%E5%8F%82%E6%95%B0&e=e"
	values, err := url.ParseQuery(urlStr)
	if nil != err {
		t.Error(err)
	}
	t.Log(values.Encode())
	t.Log(values.Get("c"), values.Get("参数"), values.Get("e"))
}

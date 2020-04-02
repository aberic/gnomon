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

package gnomon

import (
	"github.com/aberic/gnomon/grope"
	"log"
	"net/http"
)

// GHTTPCommon Http工具
type GHTTPCommon struct{}

// NewHttpServe 新建一个Http服务
func (ghc *GHTTPCommon) NewHttpServe() *grope.GHttpServe {
	return grope.NewGHttpServe()
}

// ListenAndServe 启动监听
//
// Addr 期望监听的端口号，如“:8080”
func (ghc *GHTTPCommon) ListenAndServe(Addr string, gs *grope.GHttpServe) {
	err := http.ListenAndServe(Addr, gs) //设置监听的端口
	if err != nil {
		log.Panic("ListenAndServe: ", err)
	}
}

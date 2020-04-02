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
	"encoding/json"
	"net/http"
	"strings"
)

// NewGHttpServe 新建一个Http服务
func NewGHttpServe() *GHttpServe {
	return &GHttpServe{routerMap: map[string]*GHttpRouter{}}
}

type GHttpServe struct {
	routerMap map[string]*GHttpRouter
}

// Group 设置路由根路径
//
// pattern 路由根路径，如“/test”
func (ghs *GHttpServe) Group(pattern string) *GHttpRouter {
	ghr := &GHttpRouter{methodMap: map[string]*router{}}
	ghs.routerMap[pattern] = ghr
	return ghr
}

func (ghs *GHttpServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ghs.doMethod(w, r)
}

// doMethod 处理请求具体方法
func (ghs *GHttpServe) doMethod(w http.ResponseWriter, r *http.Request) {
	var (
		groupPattern string // 项目根路径
		patterned    string // 处理后的url
	)
	// 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
	pattern := r.URL.String()
	ps := strings.Split(pattern, "/")
	groupPattern = strings.Join([]string{"/", ps[1]}, "")
	ghr := ghs.routerMap[groupPattern]
	for position, param := range ps {
		if position == 0 || position == 1 {
			continue
		}
		patterned = strings.Join([]string{patterned, param}, "/")
		if router, exist := ghr.methodMap[r.Method]; exist { // 判断router中是否存在当前请求方法
			if route, ok := router.routes[patterned]; ok { // 判断当前url是否存在route中
				if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions {
					if err := json.NewDecoder(r.Body).Decode(route.model); nil != err {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					}
				}
				var (
					offset   = 0
					paramMap = map[string]string{}
				)
				for index, p := range ps {
					if index > position {
						paramMap[route.paramMap[offset]] = p
						offset++
					}
				}
				if respModel, custom := route.handler(w, r, route.model, paramMap); !custom {
					if bytes, err := json.Marshal(respModel); nil != err {
						http.Error(w, err.Error(), http.StatusInternalServerError)
						return
					} else {
						if _, err := w.Write(bytes); nil != err {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}
						return
					}
				}
				return
			}
		}
	}
	http.NotFound(w, r)
}

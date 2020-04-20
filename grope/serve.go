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
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
)

var (
	patterns    []string
	patternLock sync.Mutex
)

// NewGHttpServe 新建一个Http服务
func NewGHttpServe(filters ...Filter) *GHttpServe {
	return &GHttpServe{filters: filters, routerMap: map[string]*GHttpRouter{}}
}

type GHttpServe struct {
	filters   []Filter // 过滤器/拦截器数组
	routerMap map[string]*GHttpRouter
}

// Group 设置路由根路径
//
// pattern 路由根路径，如“/test”
//
// filters 待实现拦截器/过滤器方法数组
func (ghs *GHttpServe) Group(pattern string, filters ...Filter) *GHttpRouter {
	filters = append(ghs.filters, filters...)
	ghr := &GHttpRouter{groupPattern: pattern, methodMap: map[string]*router{}, filters: filters}
	ghs.routerMap[pattern] = ghr
	return ghr
}

func (ghs *GHttpServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ghs.doServe(w, r)
}

// doMethod 处理请求具体方法
func (ghs *GHttpServe) doServe(w http.ResponseWriter, r *http.Request) {
	var (
		// 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
		pattern      = r.URL.String()
		ps           = strings.Split(pattern, "/")[1:]
		offset       = 0
		ghr          *GHttpRouter
		groupPattern string // 项目根路径
		patterned    string // 处理后的url
	)
	for position, p := range ps {
		var exist bool
		groupPattern = strings.Join([]string{groupPattern, "/", p}, "")
		if ghr, exist = ghs.routerMap[groupPattern]; exist {
			offset = position
			break
		}
	}
	if nil == ghr {
		http.NotFound(w, r)
		return
	}
	for position, param := range ps {
		if position <= offset {
			continue
		}
		patterned = strings.Join([]string{patterned, param}, "/")
		if router, exist := ghr.methodMap[r.Method]; exist { // 判断router中是否存在当前请求方法
			if route, ok := router.routes[patterned]; ok { // 判断当前url是否存在route中
				for _, filter := range route.filters { // 过滤无效请求
					if custom, code, err := filter(w, r); nil != err {
						if custom {
							return
						}
						http.Error(w, err.Error(), code)
						return
					}
				}
				if err := ghs.parseReqMethod(r, route); nil != err {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
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
				if err := ghs.parseHandler(w, r, route, paramMap); nil != err {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				return
			}
		}
	}
	http.NotFound(w, r)
}

// parseReqMethod 解析请求方法
func (ghs *GHttpServe) parseReqMethod(r *http.Request, route *route) error {
	if r.Method != http.MethodGet && r.Method != http.MethodHead && r.Method != http.MethodOptions && r.Method != http.MethodDelete {
		if err := ghs.parseReqModel(r, route); nil != err {
			return err
		}
	}
	return nil
}

// parseReqModel 解析请求参数
func (ghs *GHttpServe) parseReqModel(r *http.Request, route *route) error {
	contentType := r.Header.Get("Content-Type")
	switch contentType {
	case "application/json":
		return json.NewDecoder(r.Body).Decode(route.model)
	case "application/x-www-form-urlencoded":
		if err := r.ParseForm(); nil != err { //解析参数，默认是不会解析的
			return err
		}
		filedMap := make(map[string]interface{})
		for k, v := range r.Form {
			filedMap[k] = strings.Join(v, "")
		}
		route.model = filedMap
	default:
		if strings.Contains(contentType, "multipart/form-data") {
			if reader, err := r.MultipartReader(); nil != err {
				return err
			} else {
				filedMap := make(map[string]interface{})
				for {
					part, err := reader.NextPart()
					if err == io.EOF {
						break
					}
					if part == nil {
						return err
					}
					if part.FileName() == "" { // this is FormData
						if data, err := ioutil.ReadAll(part); nil != err {
							return err
						} else {
							filedMap[part.FormName()] = string(data)
						}
					} else { // This is FileData
						if bytes, err := ioutil.ReadAll(part); nil != err {
							return err
						} else {
							filedMap[part.FormName()] = &FormFile{
								FileName: part.FileName(),
								Data:     bytes,
							}
						}
					}
					func() { _ = part.Close() }()
				}
				route.model = filedMap
			}
		}
	}
	return nil
}

type FormFile struct {
	FileName string // file name
	Data     []byte // file bytes content
}

// parseHandler 解析请求处理方法
func (ghs *GHttpServe) parseHandler(w http.ResponseWriter, r *http.Request, route *route, paramMap map[string]string) error {
	if respModel, custom := route.handler(w, r, route.model, paramMap); !custom {
		if bytes, err := json.Marshal(respModel); nil != err {
			return err
		} else if _, err := w.Write(bytes); nil != err {
			return err
		}
	}
	return nil
}

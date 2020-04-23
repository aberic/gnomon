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
	"github.com/aberic/gnomon/grope/log"
	"net/http"
	"strings"
	"sync"
)

var (
	patternMap  map[string][]string // method下所属url集合
	patternLock sync.Mutex
)

// NewGHttpServe 新建一个Http服务
func NewGHttpServe(filters ...Filter) *GHttpServe {
	patternMap = map[string][]string{}
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
		offset       = 0
		ghr          *GHttpRouter
		groupPattern string // 项目根路径
		patterned    string // 处理后的url
		ctx          = &Context{writer: w, request: r}
		rtr          *router
		rt           *route
		exist        bool
	)
	pattern, paramMap := ghs.parseUrlParams(r)
	ps := strings.Split(pattern, "/")[1:]
	for position, p := range ps {
		groupPattern = strings.Join([]string{groupPattern, "/", p}, "")
		if ghrNow, exist := ghs.routerMap[groupPattern]; exist {
			ghr = ghrNow
			offset = position
			break
		}
	}
	if nil == ghr {
		http.NotFound(w, r)
		return
	}
	if rtr, exist = ghr.methodMap[r.Method]; exist { // 判断router中是否存在当前请求方法
		var pos int
		for position, param := range ps {
			if position <= offset {
				continue
			}
			patterned = strings.Join([]string{patterned, param}, "/")
			if route, ok := rtr.routes[patterned]; ok { // 判断当前url是否存在route中
				rt = route
				pos = position
				continue
			}
		}
		if nil != rt {
			for _, filter := range rt.filters { // 过滤无效请求
				filter(ctx)
				if ctx.responded {
					return
				}
			}
			var (
				offset   = 0
				valueMap = map[string]string{}
			)
			for index, p := range ps {
				if index > pos {
					valueMap[rt.valueMap[offset]] = p
					offset++
				}
			}
			ctx.valueMap = valueMap
			ctx.paramMap = paramMap
			ghs.parseHandler(ctx, rt)
			return
		}
	}
	http.NotFound(w, r)
}

func (ghs *GHttpServe) parseUrlParams(r *http.Request) (pattern string, paramMap map[string]string) {
	var (
		// 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g?name=hello&pass=work”方式进行访问
		p  = ghs.singleSeparator(r.URL.String())
		ps = strings.Split(p, "?")
	)
	if len(ps) != 2 {
		pattern = p
	} else {
		paramMap = ghs.execUrlParams(ps[1])
		if len(paramMap) == 0 {
			pattern = p
		} else {
			pattern = ps[0]
		}
	}
	return
}

func (ghs *GHttpServe) execUrlParams(paramStr string) map[string]string {
	paramMap := map[string]string{}
	paramPair := strings.Split(paramStr, "&")
	for _, pair := range paramPair {
		valuePair := strings.Split(pair, "=")
		if len(valuePair) == 1 {
			return map[string]string{}
		} else {
			paramMap[valuePair[0]] = valuePair[1]
		}
	}
	return paramMap
}

// parseHandler 解析请求处理方法
func (ghs *GHttpServe) parseHandler(ctx *Context, route *route) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("parseHandler", log.Field("error", err))
			ctx.Status(http.StatusInternalServerError)
		}
	}()
	route.handler(ctx)
}

// singleSeparator 将字符串内所有连续/替换为单个/
func (ghs *GHttpServe) singleSeparator(res string) string {
	for skip := false; !skip; {
		resNew := strings.Replace(res, "//", "/", -1)
		if res == resNew {
			skip = true
		}
		res = resNew
	}
	return res
}

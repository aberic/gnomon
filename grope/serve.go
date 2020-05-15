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
	"github.com/aberic/gnomon/grope/tune"
	"github.com/aberic/gnomon/log"
	"net/http"
	"net/url"
	"strings"
)

// newGHttpServe 新建一个Http服务
func newGHttpServe(filters ...Filter) *GHttpServe {
	nodal := newNode(filters...)
	return &GHttpServe{nodal: nodal}
}

// GHttpServe Http服务
type GHttpServe struct {
	nodal *node
}

// Group 设置路由根路径
//
// pattern 路由根路径，如“/test”
//
// filters 待实现拦截器/过滤器方法数组
func (ghs *GHttpServe) Group(pattern string, filters ...Filter) *GHttpRouter {
	ghs.nodal.add(pattern, "", nil, nil, filters...)
	ghr := &GHttpRouter{pattern: pattern, nodal: ghs.nodal}
	return ghr
}

func (ghs *GHttpServe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ghs.doServe(w, r)
}

// doMethod 处理请求具体方法
func (ghs *GHttpServe) doServe(w http.ResponseWriter, r *http.Request) {
	var ctx = &Context{writer: w, request: r, valueMap: map[string]string{}}
	pattern, paramMap := ghs.parseURLParams(r)
	ctx.paramMap = paramMap
	n := ghs.nodal.fetch(pattern, r.Method)
	if nil == n {
		http.NotFound(w, r)
		return
	} else if nil != n.error {
		w.Header().Set("Content-Type", tune.ContentTypeJSON)
		w.WriteHeader(http.StatusOK)
		bytes, _ := json.Marshal(n.error)
		_, _ = w.Write(bytes)
		return
	}
	psURLReq := strings.Split(pattern, "/")[1:]
	psURLLocal := strings.Split(n.pattern, "/")[1:]
	for index, p := range psURLLocal {
		if p[0] == ':' {
			ctx.valueMap[p[1:]] = psURLReq[index]
		}
	}
	ghs.execRoute(ctx, n)
}

// execRoute 处理请求逻辑
func (ghs *GHttpServe) execRoute(ctx *Context, nodal *node) {
	for _, filter := range nodal.filters { // 过滤无效请求
		filter(ctx)
		if ctx.responded {
			return
		}
	}
	ghs.parseHandler(ctx, nodal)
}

func (ghs *GHttpServe) parseURLParams(r *http.Request) (pattern string, paramMap map[string]string) {
	var (
		// 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g?name=hello&pass=work”方式进行访问
		p  = ghs.singleSeparator(r.URL.String())
		ps = strings.Split(p, "?")
	)
	if len(ps) != 2 {
		pattern = p
	} else {
		paramMap = ghs.execURLParams(ps[1])
		if len(paramMap) == 0 {
			pattern = p
		} else {
			pattern = ps[0]
		}
	}
	return
}

func (ghs *GHttpServe) execURLParams(paramStr string) map[string]string {
	paramMap := map[string]string{}
	values, err := url.ParseQuery(paramStr)
	if nil != err {
		return paramMap
	}
	paramPair := strings.Split(paramStr, "&")
	for _, pair := range paramPair {
		valuePair := strings.Split(pair, "=")
		if len(valuePair) != 1 {
			paramMap[valuePair[0]] = values.Get(valuePair[0])
		} else {
			return map[string]string{}
		}
	}
	return paramMap
}

// parseHandler 解析请求处理方法
func (ghs *GHttpServe) parseHandler(ctx *Context, nodal *node) {
	defer func() {
		if err := recover(); err != nil {
			log.Error("parseHandler", log.Field("error", err))
			ctx.Status(http.StatusInternalServerError)
		}
	}()
	nodal.handler(ctx)
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

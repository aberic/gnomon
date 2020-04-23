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
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

// Handler 待实现接收请求方法
//
// ctx 请求处理上下文结构
type Handler func(ctx *Context)

// Filter 过滤器/拦截器处理
//
// ctx 请求处理上下文结构
type Filter func(ctx *Context)

type router struct {
	routes map[string]*route
}

// route 路由子项目结构
type route struct {
	handler  Handler        // 待实现接收请求方法
	valueMap map[int]string // url泛型下标对应字符串集合
	filters  []Filter       // 过滤器/拦截器数组
}

// GHttpRouter Http服务路由结构
type GHttpRouter struct {
	groupPattern string             // 跟路由名
	methodMap    map[string]*router // 路由子项目集
	filters      []Filter           // 过滤器/拦截器数组
	lock         sync.RWMutex
}

func (ghr *GHttpRouter) repo(method, pattern string, handler Handler, filters ...Filter) {
	if pattern[0] != '/' {
		panic("path must begin with '/'")
	}
	var (
		patterned string
		valueMap  map[int]string
		rtr       *router
		exist     bool
		err       error
	)
	if patterned, valueMap, err = ghr.execUrl(pattern); nil != err {
		panic(err.Error())
	}
	ghr.lock.Lock()
	if rtr, exist = ghr.methodMap[method]; !exist {
		rtr = &router{map[string]*route{}}
		ghr.methodMap[method] = rtr
	}
	filters = append(ghr.filters, filters...)
	if _, exist := rtr.routes[patterned]; exist {
		panic(fmt.Sprintf("already have the same url, with method:%s and patterned:%s in group with pattern:%s",
			method, patterned, ghr.groupPattern))
	}
	rtr.routes[patterned] = &route{handler: handler, valueMap: valueMap, filters: filters}
	ghr.lock.Unlock()

	assemblyPattern := strings.Join([]string{ghr.groupPattern, patterned}, "")
	defer patternLock.Unlock()
	patternLock.Lock()
	if patterns, exist := patternMap[method]; exist {
		for _, p := range patterns {
			if p == assemblyPattern {
				panic(fmt.Sprintf("already have the same url, with assemblyPattern:%s and p:%s", assemblyPattern, p))
			}
		}
		patternMap[method] = append(patternMap[method], assemblyPattern)
	} else {
		patternMap[method] = []string{assemblyPattern}
	}
}

// execUrl 特殊处理Url
//
// pattern 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
//
// return patterned 处理后的url
//
// return valueMap url泛型下标对应字符串集合
//
// return err 处理错误内容
func (ghr *GHttpRouter) execUrl(pattern string) (patterned string, valueMap map[int]string, err error) {
	patterned = ""
	valueMap = map[int]string{}
	ps := strings.Split(pattern, "/")[1:]
	index := 0
	for _, param := range ps {
		if strings.HasPrefix(param, ":") {
			valueMap[index] = strings.Split(param, ":")[1]
			index++
		} else {
			if index > 0 {
				err = errors.New("custom url must continue until the end")
				return
			}
			patterned = strings.Join([]string{patterned, param}, "/")
		}
	}
	return
}

// Get 发起一个 Get 请求接收项目
//
// GET请求会显示请求指定的资源。一般来说GET方法应该只用于数据的读取，而不应当用于会产生副作用的非幂等的操作中。
// GET会方法请求指定的页面信息，并返回响应主体，GET被认为是不安全的方法，因为GET方法会被网络蜘蛛等任意的访问。
//
// pattern 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
//
// model 期望接收的结构，如“&Test{}”，最终在Handler方法中得以调用
//
// handler 待实现接收请求方法
//
// filters 待实现拦截器/过滤器方法数组
func (ghr *GHttpRouter) Get(pattern string, handler Handler, filters ...Filter) {
	go ghr.repo(http.MethodGet, pattern, handler, filters...)
}

// Head 发起一个 Head 请求接收项目
//
// HEAD方法与GET方法一样，都是向服务器发出指定资源的请求。
// 但是，服务器在响应HEAD请求时不会回传资源的内容部分，即：响应主体。
// 这样，我们可以不传输全部内容的情况下，就可以获取服务器的响应头信息。HEAD方法常被用于客户端查看服务器的性能。
//
// pattern 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
//
// model 期望接收的结构，如“&Test{}”，最终在Handler方法中得以调用
//
// handler 待实现接收请求方法
//
// filters 待实现拦截器/过滤器方法数组
func (ghr *GHttpRouter) Head(pattern string, handler Handler, filters ...Filter) {
	go ghr.repo(http.MethodHead, pattern, handler, filters...)
}

// Post 发起一个 Post 请求接收项目
//
// POST请求会 向指定资源提交数据，请求服务器进行处理，如：表单数据提交、文件上传等，请求数据会被包含在请求体中。
// POST方法是非幂等的方法，因为这个请求可能会创建新的资源或/和修改现有资源。
//
// pattern 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
//
// model 期望接收的结构，如“&Test{}”，最终在Handler方法中得以调用
//
// handler 待实现接收请求方法
//
// filters 待实现拦截器/过滤器方法数组
func (ghr *GHttpRouter) Post(pattern string, handler Handler, filters ...Filter) {
	go ghr.repo(http.MethodPost, pattern, handler, filters...)
}

// Put 发起一个 Put 请求接收项目
//
// PUT请求会身向指定资源位置上传其最新内容，PUT方法是幂等的方法。通过该方法客户端可以将指定资源的最新数据传送给服务器取代指定的资源的内容。
//
// pattern 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
//
// model 期望接收的结构，如“&Test{}”，最终在Handler方法中得以调用
//
// handler 待实现接收请求方法
//
// filters 待实现拦截器/过滤器方法数组
func (ghr *GHttpRouter) Put(pattern string, handler Handler, filters ...Filter) {
	go ghr.repo(http.MethodPut, pattern, handler, filters...)
}

// Patch 发起一个 Patch 请求接收项目
//
// PATCH方法出现的较晚，它在2010年的RFC 5789标准中被定义。PATCH请求与PUT请求类似，同样用于资源的更新。二者有以下两点不同：
// PATCH一般用于资源的部分更新，而PUT一般用于资源的整体更新。
// 当资源不存在时，PATCH会创建一个新的资源，而PUT只会对已在资源进行更新。
//
// pattern 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
//
// model 期望接收的结构，如“&Test{}”，最终在Handler方法中得以调用
//
// handler 待实现接收请求方法
//
// filters 待实现拦截器/过滤器方法数组
func (ghr *GHttpRouter) Patch(pattern string, handler Handler, filters ...Filter) {
	go ghr.repo(http.MethodPatch, pattern, handler, filters...)
}

// Delete 发起一个 Delete 请求接收项目
//
// DELETE请求用于请求服务器删除所请求URI（统一资源标识符，Uniform Resource Identifier）所标识的资源。
// DELETE请求后指定资源会被删除，DELETE方法也是幂等的。
//
// pattern 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
//
// model 期望接收的结构，如“&Test{}”，最终在Handler方法中得以调用
//
// handler 待实现接收请求方法
//
// filters 待实现拦截器/过滤器方法数组
func (ghr *GHttpRouter) Delete(pattern string, handler Handler, filters ...Filter) {
	go ghr.repo(http.MethodDelete, pattern, handler, filters...)
}

// Connect 发起一个 Connect 请求接收项目
//
// CONNECT方法是HTTP/1.1协议预留的，能够将连接改为管道方式的代理服务器。通常用于SSL加密服务器的链接与非加密的HTTP代理服务器的通信。
//
// pattern 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
//
// model 期望接收的结构，如“&Test{}”，最终在Handler方法中得以调用
//
// handler 待实现接收请求方法
//
// filters 待实现拦截器/过滤器方法数组
func (ghr *GHttpRouter) Connect(pattern string, handler Handler, filters ...Filter) {
	go ghr.repo(http.MethodConnect, pattern, handler, filters...)
}

// Options 发起一个 Options 请求接收项目
//
// OPTIONS请求与HEAD类似，一般也是用于客户端查看服务器的性能。
// 这个方法会请求服务器返回该资源所支持的所有HTTP请求方法，该方法会用’*’来代替资源名称，向服务器发送OPTIONS请求，可以测试服务器功能是否正常。
// JavaScript的XMLHttpRequest对象进行CORS跨域资源共享时，就是使用OPTIONS方法发送嗅探请求，以判断是否有对指定资源的访问权限。
//
// pattern 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
//
// model 期望接收的结构，如“&Test{}”，最终在Handler方法中得以调用
//
// handler 待实现接收请求方法
//
// filters 待实现拦截器/过滤器方法数组
func (ghr *GHttpRouter) Options(pattern string, handler Handler, filters ...Filter) {
	go ghr.repo(http.MethodOptions, pattern, handler, filters...)
}

// Trace 发起一个 Trace 请求接收项目
//
// TRACE请求服务器回显其收到的请求信息，该方法主要用于HTTP请求的测试或诊断。
//
// pattern 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
//
// model 期望接收的结构，如“&Test{}”，最终在Handler方法中得以调用
//
// handler 待实现接收请求方法
//
// filters 待实现拦截器/过滤器方法数组
func (ghr *GHttpRouter) Trace(pattern string, handler Handler, filters ...Filter) {
	go ghr.repo(http.MethodTrace, pattern, handler, filters...)
}

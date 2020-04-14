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
	"net/http"
	"strings"
	"sync"
)

// Handler 待实现接收请求方法
//
// writer 原生 net/http 结构
//
// request 原生 net/http 结构
//
// reqModel 如果需要，这里是期望接收的结构，如Test，通过reqModel.(*Test)转换
//
// paramMap 如果需要，这里是与url请求中对应的参数集合，如“/demo/:id/”，则通过 paramMap[id] 获取url中的值
//
// return respModel 自行返回的结构
//
// return custom 是否自行处理返回结果，如果自行处理，则本次返回不再执行默认操作
type Handler func(writer http.ResponseWriter, request *http.Request, reqModel interface{}, paramMap map[string]string) (respModel interface{}, custom bool)

type router struct {
	routes map[string]*route
}

// route 路由子项目结构
type route struct {
	model    interface{}    // 期望接收的结构，如“&Test{}”，最终在Handler方法中得以调用
	handler  Handler        // 待实现接收请求方法
	paramMap map[int]string // url泛型下标对应字符串集合
}

// GHttpRouter Http服务路由结构
type GHttpRouter struct {
	methodMap map[string]*router // 路由子项目集
	lock      sync.RWMutex
}

// execUrl 特殊处理Url
//
// pattern 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
//
// return patterned 处理后的url
//
// return paramMap url泛型下标对应字符串集合
//
// return err 处理错误内容
func (ghr *GHttpRouter) execUrl(pattern string) (patterned string, paramMap map[int]string, err error) {
	patterned = ""
	paramMap = map[int]string{}
	ps := strings.Split(pattern, "/")
	index := 0
	for position, param := range ps {
		if position == 0 {
			continue
		}
		if strings.HasPrefix(param, ":") {
			paramMap[index] = strings.Split(param, ":")[1]
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

// repo 发起一个接收项目，请求方法复用 net/http “http.MethodGet”等
//
// 安全性与幂等性
//
// 关于HTTP请求采用的这些个方法，具有两个基本的特性，即“安全性”和“幂等性”。
//
// 对于7种HTTP方法，GET、HEAD和OPTIONS均被认为是安全的方法，因为它们旨在实现对数据的获取，并不具有“边界效应”。
//
// 至于其它4个HTTP方法，由于它们会导致服务端资源的变化，所以被认为是不安全的方法。
//
// 幂等性（Idempotent）是一个数学上的概念，在这里表示发送一次和多次请求引起的边界效应是一致的。
// 在网速不够快的情况下，客户端发送一个请求后不能立即得到响应，由于不能确定是否请求是否被成功提交，所以它有可能会再次发送另一个相同的请求，幂等性决定了第二个请求是否有效。
//
// 3种安全的HTTP方法（GET、HEAD和OPTIONS）均是幂等方法。由于DELETE和PATCH请求操作的是现有的某个资源，所以它们是幂等方法。
// 对于PUT请求，只有在对应资源不存在的情况下服务器才会进行添加操作，否则只作修改操作，所以它也是幂等方法。
// 至于最后一种POST，由于它总是进行添加操作，如果服务器接收到两次相同的POST操作，将导致两个相同的资源被创建，所以这是一个非幂等的方法。
//
// 当我们在设计Web API的时候，应该尽量根据请求HTTP方法的幂等型来决定处理的逻辑。
// 由于PUT是一个幂等方法，所以携带相同资源的PUT请求不应该引起资源的状态变化，如果我们在资源上附加一个自增长的计数器表示被修改的次数，这实际上就破坏了幂等型。
//
// 无状态性
//
// restful只要维护资源的状态，而不需要维护客户端的状态。
// 对于它来说，每次请求都是全新的，它只需要针对本次请求作相应的操作，不需要将本次请求的相关信息记录下来以便用于后续来自相同客户端请求的处理。
// 对于上面所述restful的这些个特性，它们都是要求为了满足这些特征做点什么，唯有这个无状态却是要求不要做什么，因为HTTP本身就是无状态的。
//
// 举个例子，一个网页通过调用Web API分页获取符合查询条件的记录。
// 一般情况下，页面导航均具有“上一页”和“下一页”链接用于呈现当前页的前一页和后一页的记录。那么现在有两种实现方式返回上下页的记录。
//
// Web API不仅仅会定义根据具体页码的数据查询定义相关的操作，还会针对“上一页”和“下一页”这样的请求定义单独的操作。
// 它自身会根据客户端的Session ID对每次数据返回的页面在本地进行保存，以便能够知道上一页和下一页具体是哪一页。
//
// Web API只会定义根据具体页码的数据查询定义相关的操作，当前返回数据的页码由客户端来维护。
//
// 第一种貌似很“智能”，其实就是一种画蛇添足的作法，因为它破坏了Web API的无状态性。
// 设计无状态的Web API不仅仅使Web API自身显得简单而精炼，还因减除了针对客户端的“亲和度（Affinty）”使我们可以有效地实施负载均衡，
// 因为只有这样集群中的每一台服务器对于每个客户端才是等效的。
//
// pattern 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
//
// model 期望接收的结构，如“&Test{}”，最终在Handler方法中得以调用
//
// handler 待实现接收请求方法
func (ghr *GHttpRouter) repo(method, pattern string, model interface{}, handler Handler) {
	if pattern[0] != '/' {
		panic("path must begin with '/'")
	}
	var (
		patterned string
		paramMap  map[int]string
		rtr       *router
		exist     bool
		err       error
	)
	if patterned, paramMap, err = ghr.execUrl(pattern); nil != err {
		panic(err.Error())
	}
	defer ghr.lock.Unlock()
	ghr.lock.Lock()
	if rtr, exist = ghr.methodMap[method]; !exist {
		rtr = &router{map[string]*route{}}
		ghr.methodMap[method] = rtr
	}
	if _, exist := rtr.routes[patterned]; exist {
		panic("already have the same url")
	}
	rtr.routes[patterned] = &route{model: model, handler: handler, paramMap: paramMap}
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
func (ghr *GHttpRouter) Get(pattern string, model interface{}, handler Handler) {
	go ghr.repo(http.MethodGet, pattern, model, handler)
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
func (ghr *GHttpRouter) Head(pattern string, model interface{}, handler Handler) {
	go ghr.repo(http.MethodHead, pattern, model, handler)
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
func (ghr *GHttpRouter) Post(pattern string, model interface{}, handler Handler) {
	go ghr.repo(http.MethodPost, pattern, model, handler)
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
func (ghr *GHttpRouter) Put(pattern string, model interface{}, handler Handler) {
	go ghr.repo(http.MethodPut, pattern, model, handler)
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
func (ghr *GHttpRouter) Patch(pattern string, model interface{}, handler Handler) {
	go ghr.repo(http.MethodPatch, pattern, model, handler)
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
func (ghr *GHttpRouter) Delete(pattern string, model interface{}, handler Handler) {
	go ghr.repo(http.MethodDelete, pattern, model, handler)
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
func (ghr *GHttpRouter) Connect(pattern string, model interface{}, handler Handler) {
	go ghr.repo(http.MethodConnect, pattern, model, handler)
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
func (ghr *GHttpRouter) Options(pattern string, model interface{}, handler Handler) {
	go ghr.repo(http.MethodOptions, pattern, model, handler)
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
func (ghr *GHttpRouter) Trace(pattern string, model interface{}, handler Handler) {
	go ghr.repo(http.MethodTrace, pattern, model, handler)
}

// GetForm 发起一个 PostForm 请求接收项目
//
// POST请求会 向指定资源提交数据，请求服务器进行处理，如：表单数据提交、文件上传等，请求数据会被包含在请求体中。
// POST方法是非幂等的方法，因为这个请求可能会创建新的资源或/和修改现有资源。
//
// pattern 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
//
// fieldMap 如果需要，这里是期望接收的kv集合
//
// handler 待实现接收请求方法
func (ghr *GHttpRouter) GetForm(pattern string, fieldMap map[string]string, handler Handler) {
	go ghr.repo(http.MethodGet, pattern, fieldMap, handler)
}

// PostForm 发起一个 PostForm 请求接收项目
//
// POST请求会 向指定资源提交数据，请求服务器进行处理，如：表单数据提交、文件上传等，请求数据会被包含在请求体中。
// POST方法是非幂等的方法，因为这个请求可能会创建新的资源或/和修改现有资源。
//
// pattern 项目路径，如“/demo/:id/:name”，与路由根路径相结合，最终会通过类似“http://127.0.0.1:8080/test/demo/1/g”方式进行访问
//
// fieldMap 如果需要，这里是期望接收的kv集合
//
// handler 待实现接收请求方法
func (ghr *GHttpRouter) PostForm(pattern string, fieldMap map[string]interface{}, handler Handler) {
	go ghr.repo(http.MethodPost, pattern, fieldMap, handler)
}

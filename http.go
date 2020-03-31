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
	"encoding/json"
	"net/http"
)

// HTTPCommon Http工具
type HTTPCommon struct{}

type Handler func(writer http.ResponseWriter, request *http.Request, reqModel interface{}) (respModel interface{})

func (hc *HTTPCommon) NewHttpServe() *HttpServe {
	mux := http.NewServeMux()
	return &HttpServe{mux: mux, routers: []*HttpRouter{}}
}

func (hc *HTTPCommon) NewHttpRouter(hs *HttpServe) *HttpRouter {
	hr := &HttpRouter{routes: []*route{}}
	hs.routers = append(hs.routers, hr)
	return hr
}

type HttpServe struct {
	mux     *http.ServeMux
	routers []*HttpRouter
}

func (hs *HttpServe) route() {
	for _, router := range hs.routers {
		for _, r := range router.routes {
			r.handler(hs.mux)
		}
	}
}

func (hs *HttpServe) ListenAndServe(Addr string) error {
	hs.route()
	server := &http.Server{
		Addr:    Addr,
		Handler: hs.mux,
	}
	return server.ListenAndServe()
}

func (hs *HttpServe) Listen(server *http.Server) error {
	hs.route()
	server.Handler = hs.mux
	return server.ListenAndServe()
}

type HttpRouter struct {
	pattern string
	routes  []*route
}

func (hr *HttpRouter) Group(pattern string) {
	hr.pattern = pattern
}

func (hr *HttpRouter) Post(pattern string, model interface{}, handler Handler) {
	hr.routes = append(hr.routes, &route{
		Method:  http.MethodPost,
		Pattern: String().StringBuilder(hr.pattern, pattern),
		Model:   model,
		Handler: handler,
	})
}

type route struct {
	Method  string
	Pattern string
	Model   interface{}
	Handler Handler
}

func (r *route) handler(mux *http.ServeMux) {
	mux.HandleFunc(r.Pattern, func(writer http.ResponseWriter, request *http.Request) {
		if r.Method == request.Method {
			//r.HandlerFunc(writer, request)
			if err := json.NewDecoder(request.Body).Decode(r.Model); nil != err {
				_, _ = writer.Write([]byte(err.Error()))
			}
			if bytes, err := json.Marshal(r.Handler(writer, request, r.Model)); nil != err {
				_, _ = writer.Write([]byte(err.Error()))
			} else {
				_, _ = writer.Write(bytes)
			}
		} else {
			_, _ = writer.Write([]byte("method error"))
		}
	})
}

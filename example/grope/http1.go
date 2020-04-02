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

package main

import (
	"github.com/aberic/gnomon"
	"github.com/aberic/gnomon/grope"
	"net/http"
)

type TestOne struct {
	One   string `json:"one"`
	Ones  bool   `json:"ones"`
	OneGo int    `json:"one_go"`
}

type TestTwo struct {
	Two   string `json:"two"`
	Twos  bool   `json:"twos"`
	TwoGo int    `json:"two_go"`
}

func main() {
	httpServe := gnomon.Grope().NewHttpServe()
	router1(httpServe)
	router2(httpServe)
	gnomon.Grope().ListenAndServe(":8888", httpServe)
}

func router1(hs *grope.GHttpServe) {
	// 仓库相关路由设置
	route := hs.Group("/one")
	route.Post("/test1", &TestOne{}, one1)
	route.Post("/test2/:a/:b", &TestOne{}, one2)
}

func one1(_ http.ResponseWriter, r *http.Request, reqModel interface{}, _ map[string]string) (respModel interface{}, custom bool) {
	ones := reqModel.(*TestOne)
	gnomon.Log().Info("one", gnomon.Log().Field("one", &ones), gnomon.Log().Field("url", r.URL.String()))
	return &TestTwo{
		Two:   "2",
		Twos:  false,
		TwoGo: 2,
	}, false
}

func one2(_ http.ResponseWriter, r *http.Request, reqModel interface{}, paramMaps map[string]string) (respModel interface{}, custom bool) {
	ones := reqModel.(*TestOne)
	gnomon.Log().Info("one", gnomon.Log().Field("one", &ones), gnomon.Log().Field("url", r.URL.String()),
		gnomon.Log().Field("a", paramMaps["a"]), gnomon.Log().Field("b", paramMaps["b"]))
	return &TestTwo{
		Two:   "2",
		Twos:  false,
		TwoGo: 2,
	}, false
}

func router2(hs *grope.GHttpServe) {
	// 仓库相关路由设置
	route := hs.Group("/two")
	route.Post("/test1", &TestTwo{}, two1)
	route.Get("/test2/:id/:name/:pass", nil, two2)
}

func two1(_ http.ResponseWriter, r *http.Request, reqModel interface{}, _ map[string]string) (respModel interface{}, custom bool) {
	twos := reqModel.(*TestTwo)
	gnomon.Log().Info("two", gnomon.Log().Field("two", &twos), gnomon.Log().Field("url", r.URL.String()))
	return &TestOne{
		One:   "1",
		Ones:  true,
		OneGo: 1,
	}, false
}

func two2(w http.ResponseWriter, r *http.Request, _ interface{}, paramMaps map[string]string) (respModel interface{}, custom bool) {
	gnomon.Log().Info("one", gnomon.Log().Field("url", r.URL.String()),
		gnomon.Log().Field("id", paramMaps["id"]), gnomon.Log().Field("name", paramMaps["name"]),
		gnomon.Log().Field("pass", paramMaps["pass"]))
	_, _ = w.Write([]byte("custom true"))
	return nil, true
}

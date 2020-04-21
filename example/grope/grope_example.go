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
	"errors"
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
	//gnomon.Grope().ListenAndServeTLS(
	//	httpServe,
	//	":8888",
	//	"./example/ca/server/rootCA.crt",
	//	"./example/ca/server/rootCA.key",
	//	"./example/ca/client/rootCA.crt")
}

func doFilter1(w http.ResponseWriter, r *http.Request) (bool, int, error) {
	if r.Header.Get("name") == "name1" {
		return false, 0, nil
	} else if r.Header.Get("name") == "name2" {
		return false, http.StatusBadGateway, errors.New("filter name error")
	}
	_, _ = w.Write([]byte("filter custom name error"))
	return true, http.StatusBadGateway, errors.New("filter name error")
}

func doFilter2(w http.ResponseWriter, r *http.Request) (bool, int, error) {
	if r.Header.Get("pass") == "pass1" {
		return false, 0, nil
	} else if r.Header.Get("pass") == "pass2" {
		return false, http.StatusBadGateway, errors.New("filter pass error")
	}
	_, _ = w.Write([]byte("filter custom pass error"))
	return true, http.StatusBadGateway, errors.New("filter pass error")
}

func router1(hs *grope.GHttpServe) {
	// 仓库相关路由设置
	route := hs.Group("/one/get")
	route.Post("/test1", &TestOne{}, one1)
	route.Put("/test1", &TestOne{}, one1)
	route.Post("/test2/:a/:b", &TestOne{}, one2)
	route.PostForm("/test3/:a/:b", map[string]interface{}{}, one3)
	route.PostForm("/test4/:a/:b", map[string]interface{}{}, one4)
	route.PostForm("/test5/:a/:b", map[string]interface{}{}, one5)
	route.Get("/test6", &TestOne{}, one1, doFilter2)
}

func one1(_ http.ResponseWriter, r *http.Request, reqModel interface{}, _ map[string]string, paramMap map[string]string) (respModel interface{}, custom bool) {
	ones := reqModel.(*TestOne)
	gnomon.Log().Info("one", gnomon.Log().Field("one", &ones), gnomon.Log().Field("url", r.URL.String()), gnomon.Log().Field("paramMap", paramMap))
	return &TestTwo{
		Two:   "2",
		Twos:  false,
		TwoGo: 2,
	}, false
}

func one2(_ http.ResponseWriter, r *http.Request, reqModel interface{}, valueMaps map[string]string, _ map[string]string) (respModel interface{}, custom bool) {
	ones := reqModel.(*TestOne)
	gnomon.Log().Info("one", gnomon.Log().Field("one", &ones), gnomon.Log().Field("url", r.URL.String()),
		gnomon.Log().Field("a", valueMaps["a"]), gnomon.Log().Field("b", valueMaps["b"]))
	return &TestTwo{
		Two:   "2",
		Twos:  false,
		TwoGo: 2,
	}, false
}

func one3(_ http.ResponseWriter, r *http.Request, reqModel interface{}, valueMaps map[string]string, _ map[string]string) (respModel interface{}, custom bool) {
	ones := reqModel.(map[string]interface{})
	gnomon.Log().Info("one", gnomon.Log().Field("one", &ones), gnomon.Log().Field("url", r.URL.String()),
		gnomon.Log().Field("a", valueMaps["a"]), gnomon.Log().Field("b", valueMaps["b"]))
	return &TestTwo{
		Two:   "2",
		Twos:  false,
		TwoGo: 2,
	}, false
}

func one4(_ http.ResponseWriter, r *http.Request, reqModel interface{}, valueMaps map[string]string, _ map[string]string) (respModel interface{}, custom bool) {
	ones := reqModel.(map[string]interface{})
	gnomon.Log().Info("one", gnomon.Log().Field("u", ones["u"]), gnomon.Log().Field("v", ones["v"]), gnomon.Log().Field("url", r.URL.String()),
		gnomon.Log().Field("a", valueMaps["a"]), gnomon.Log().Field("b", valueMaps["b"]))

	file := ones["file1"].(*grope.FormFile)
	gnomon.File().Append("tmp/httpFileTest/"+file.FileName, file.Data, true)

	return &TestTwo{
		Two:   "2",
		Twos:  false,
		TwoGo: 2,
	}, false
}

func one5(_ http.ResponseWriter, r *http.Request, reqModel interface{}, valueMaps map[string]string, _ map[string]string) (respModel interface{}, custom bool) {
	ones := reqModel.(map[string]interface{})
	gnomon.Log().Info("one", gnomon.Log().Field("u", ones["u"]), gnomon.Log().Field("v", ones["v"]), gnomon.Log().Field("url", r.URL.String()),
		gnomon.Log().Field("aaa", valueMaps["aaa"]), gnomon.Log().Field("bbb", valueMaps["bbb"]))

	file1 := ones["wk"].(*grope.FormFile)
	file2 := ones["kw"].(*grope.FormFile)
	gnomon.File().Append("tmp/httpFileTest/"+file1.FileName, file1.Data, true)
	gnomon.File().Append("tmp/httpFileTest/"+file2.FileName, file2.Data, true)

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

func two1(_ http.ResponseWriter, r *http.Request, reqModel interface{}, _ map[string]string, _ map[string]string) (respModel interface{}, custom bool) {
	twos := reqModel.(*TestTwo)
	gnomon.Log().Info("two", gnomon.Log().Field("two", &twos), gnomon.Log().Field("url", r.URL.String()))
	return &TestOne{
		One:   "1",
		Ones:  true,
		OneGo: 1,
	}, false
}

func two2(w http.ResponseWriter, r *http.Request, _ interface{}, valueMaps map[string]string, _ map[string]string) (respModel interface{}, custom bool) {
	gnomon.Log().Info("one", gnomon.Log().Field("url", r.URL.String()),
		gnomon.Log().Field("id", valueMaps["id"]), gnomon.Log().Field("name", valueMaps["name"]),
		gnomon.Log().Field("pass", valueMaps["pass"]))
	_, _ = w.Write([]byte("custom true"))
	return nil, true
}

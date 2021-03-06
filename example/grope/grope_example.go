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
	"github.com/aberic/gnomon/log"
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
	httpServe := grope.NewHTTPServe()
	router1(httpServe)
	router2(httpServe)
	grope.ListenAndServe(":8888", httpServe)
	//grope.ListenAndServeTLS(
	//	httpServe,
	//	":8888",
	//	"/Users/aberic/Downloads/test1/node.org.cert.pem",
	//	"/Users/aberic/Downloads/test1/node.key.pem",
	//	true,
	//	"/Users/aberic/Downloads/test1/org.root.cert.pem")
}

func doFilter1(ctx *grope.Context) {
	if ctx.HeaderGet("name") != "name" {
		log.Info("doFilter1", log.Field("resp", ctx.ResponseText(http.StatusForbidden, "filter name")))
	}
}

func doFilter2(ctx *grope.Context) {
	filterStr := ctx.HeaderGet("test")
	if filterStr != "pass" {
		log.Info("doFilter2", log.Field("resp", ctx.ResponseText(http.StatusForbidden, "filter pass")))
	}
}

func doFilter3(ctx *grope.Context) {
	filterStr := ctx.HeaderGet("test")
	if filterStr != "test" {
		log.Info("doFilter3", log.Field("resp", ctx.ResponseText(http.StatusForbidden, "filter test")))
	}
}

func router1(hs *grope.GHttpServe) {
	// 仓库相关路由设置
	route := hs.Group("/one/test")
	route.Post("/test1", one1)
	route.Put("/test1", one2)
	route.Post("/test2/:a/:b", one2)
	route.Post("/test3/:a/:b", one3)
	route.Post("/test4/:a/:b", one4)
	route.Post("/test5/:a/:b", one5)
	route.Put("/test6/ok", one1)
	route.Put("/test6/ok/no", one6)
	route.Put("/test6/:a/:b", one6)
}

func one1(ctx *grope.Context) {
	//oness := nilOne()
	//oness.One = "1"
	ones := &TestOne{}
	_ = ctx.ReceiveJSON(ones)
	log.Info("one", log.Field("one", &ones),
		log.Field("url", ctx.Request().URL.String()), log.Field("paramMap", ctx.Params()))
	log.Info("one1", log.Field("resp", ctx.ResponseJSON(http.StatusOK, &TestTwo{
		Two:   "1",
		Twos:  false,
		TwoGo: 1,
	})))
}

func one2(ctx *grope.Context) {
	ones := &TestOne{}
	_ = ctx.ReceiveJSON(ones)
	log.Info("one", log.Field("one", &ones),
		log.Field("url", ctx.Request().URL.String()),
		log.Field("a", ctx.Value("a")), log.Field("b", ctx.Value("b")),
		log.Field("c", ctx.Param("c")), log.Field("d", ctx.Param("d")))
	log.Info("one1", log.Field("resp", ctx.ResponseJSON(http.StatusOK, &TestTwo{
		Two:   "2",
		Twos:  false,
		TwoGo: 2,
	})))
}

func one3(ctx *grope.Context) {
	ones, err := ctx.ReceiveForm()
	if nil != err {
		log.Error("one3", log.Err(err))
	}
	log.Info("one", log.Field("one", &ones),
		log.Field("url", ctx.Request().URL.String()),
		log.Field("a", ctx.Value("a")), log.Field("b", ctx.Value("b")),
		log.Field("xxx", ctx.Param("xxx")), log.Field("yyy", ctx.Param("yyy")))
	log.Info("one1", log.Field("resp", ctx.ResponseJSON(http.StatusOK, &TestTwo{
		Two:   "2",
		Twos:  false,
		TwoGo: 2,
	})))
}

func one4(ctx *grope.Context) {
	ones, _ := ctx.ReceiveMultipartForm()
	log.Info("one", log.Field("u", ones.Params["u"]), log.Field("v", ones.Params["v"]),
		log.Field("url", ctx.Request().URL.String()),
		log.Field("a", ctx.Values()["a"]), log.Field("b", ctx.Values()["b"]))
	file := ones.Files["file1"][0]
	_, _ = gnomon.FileAppend("tmp/httpFileTest/"+file.FileName, file.Data, true)
	log.Info("one1", log.Field("resp", ctx.ResponseJSON(http.StatusOK, &TestTwo{
		Two:   "2",
		Twos:  false,
		TwoGo: 2,
	})))
}

func one5(ctx *grope.Context) {
	ones, _ := ctx.ReceiveMultipartForm()
	log.Info("one", log.Field("u", ones.Params["u"]), log.Field("v", ones.Params["v"]),
		log.Field("url", ctx.Request().URL.String()),
		log.Field("a", ctx.Values()["a"]), log.Field("b", ctx.Values()["b"]))

	file1 := ones.Files["wk"][0]
	file2 := ones.Files["kw"][0]
	_, _ = gnomon.FileAppend("tmp/httpFileTest/"+file1.FileName, file1.Data, true)
	_, _ = gnomon.FileAppend("tmp/httpFileTest/"+file2.FileName, file2.Data, true)
	log.Info("one1", log.Field("resp", ctx.ResponseJSON(http.StatusOK, &TestTwo{
		Two:   "2",
		Twos:  false,
		TwoGo: 2,
	})))
}

func one6(ctx *grope.Context) {
	ones := &TestOne{}
	_ = ctx.ReceiveJSON(ones)
	log.Info("one", log.Field("one", &ones),
		log.Field("url", ctx.Request().URL.String()), log.Field("a", ctx.Values()["a"]))
	log.Info("one6", log.Field("resp", ctx.ResponseJSON(http.StatusOK, &TestTwo{
		Two:   "22",
		Twos:  false,
		TwoGo: 22,
	})))
}

func router2(hs *grope.GHttpServe) {
	// 仓库相关路由设置
	route := hs.Group("/two")
	route.Post("/test1", two1)
	route.Get("/test2/:id/:name/:pass/:test", two2)
	route.Get("/test2/test", two3)
	route.Get("/test3", two3)
	route.Get("/test4", two4)
}

func two1(ctx *grope.Context) {
	twos := &TestTwo{}
	_ = ctx.ReceiveJSON(twos)
	log.Info("two", log.Field("two", &twos), log.Field("url", ctx.Request().URL.String()))
	log.Info("one1", log.Field("resp", ctx.ResponseJSON(http.StatusOK, &TestOne{
		One:   "1",
		Ones:  true,
		OneGo: 1,
	})))
}

func two2(ctx *grope.Context) {
	log.Info("one", log.Field("url", ctx.Request().URL.String()),
		log.Field("id", ctx.Values()["id"]), log.Field("name", ctx.Values()["name"]),
		log.Field("pass", ctx.Values()["pass"]),
		log.Field("ok", ctx.HeaderGet("ok")), log.Field("no", ctx.HeaderGet("no")))
	log.Info("one1", log.Field("resp", ctx.ResponseJSON(http.StatusOK, &TestOne{
		One:   "8888",
		Ones:  true,
		OneGo: 1,
	})))
}

func two3(ctx *grope.Context) {
	twos := &TestTwo{}
	_ = ctx.ReceiveJSON(twos)
	log.Info("two", log.Field("two", &twos), log.Field("url", ctx.Request().URL.String()))
	log.Info("one1", log.Field("resp", ctx.ResponseFile(http.StatusOK, "1.sql", "tmp/httpFileTest/1.sql")))
}

func two4(ctx *grope.Context) {
	twos := &TestTwo{}
	_ = ctx.ReceiveJSON(twos)
	log.Info("two", log.Field("two", &twos), log.Field("url", ctx.Request().URL.String()))
	ctx.Distribution("http://www.baidu.com", func(err error) {
		log.Error("two4", log.Err(err))
	})
}

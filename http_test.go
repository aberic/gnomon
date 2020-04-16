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
	"io/ioutil"
	"testing"
)

type TestOne struct {
	One   string `json:"one"`
	Ones  bool   `json:"ones"`
	OneGo int    `json:"one_go"`
}

func TestHttpClientCommon_Get(t *testing.T) {
	if resp, err := HttpClient().Get("http://localhost:8888/two/test2/0/hello/word"); nil != err {
		t.Error(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Error("unable to read response body:", err.Error())
		} else {
			t.Log(string(bytes))
		}
	}
}

func TestHttpClientCommon_GetTLS(t *testing.T) {
	if resp, err := HttpClient().GetTLS("https://localhost:8888/two/test2/0/hello/word", &HttpTLSConfig{
		CACrtFilePath:      "./example/ca/server/rootCA.crt",
		CertFilePath:       "./example/ca/client/rootCA.crt",
		KeyFilePath:        "./example/ca/client/rootCA.key",
		InsecureSkipVerify: false,
	}); nil != err {
		t.Error(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Error("unable to read response body:", err.Error())
		} else {
			t.Log(string(bytes))
		}
	}
}

func TestHttpClientCommon_Post(t *testing.T) {
	if resp, err := HttpClient().Post("http://localhost:8888/one/test1", &TestOne{
		One:   "1",
		Ones:  true,
		OneGo: 1,
	}); nil != err {
		t.Error(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Error("unable to read response body:", err.Error())
		} else {
			t.Log(string(bytes))
		}
	}
}

func TestHttpClientCommon_PostTLS(t *testing.T) {
	if resp, err := HttpClient().PostTLS("https://localhost:8888/one/test1", &TestOne{
		One:   "1",
		Ones:  true,
		OneGo: 1,
	}, &HttpTLSConfig{
		CACrtFilePath:      "./example/ca/server/rootCA.crt",
		CertFilePath:       "./example/ca/client/rootCA.crt",
		KeyFilePath:        "./example/ca/client/rootCA.key",
		InsecureSkipVerify: false,
	}); nil != err {
		t.Error(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Error("unable to read response body:", err.Error())
		} else {
			t.Log(string(bytes))
		}
	}
}

func TestHttpClientCommon_PostForm(t *testing.T) {
	paramMap := map[string]string{}
	paramMap["xxx"] = "111"
	paramMap["yyy"] = "222"
	if resp, err := HttpClient().PostForm("http://localhost:8888/one/test3/x/y", paramMap, nil); nil != err {
		t.Error(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Error("unable to read response body:", err.Error())
		} else {
			t.Log(string(bytes))
		}
	}
}

func TestHttpClientCommon_PostForm1(t *testing.T) {
	paramMap := map[string]string{}
	paramMap["aaa"] = "333"
	paramMap["bbb"] = "444"
	fileMap := map[string]string{}
	fileMap["wk"] = "/Users/aberic/Downloads/plantuml4idea.zip"
	fileMap["kw"] = "/Users/aberic/Documents/1400115281_report_pb.dump"
	if resp, err := HttpClient().PostForm("http://localhost:8888/one/test5/a/b", paramMap, fileMap); nil != err {
		t.Error(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Error("unable to read response body:", err.Error())
		} else {
			t.Log(string(bytes))
		}
	}
}

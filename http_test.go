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
	"context"
	"fmt"
	"gotest.tools/assert"
	"io/ioutil"
	"net"
	"testing"
)

type TestOne struct {
	One   string `json:"one"`
	Ones  bool   `json:"ones"`
	OneGo int    `json:"one_go"`
}

func TestHttpClientCommon_Get(t *testing.T) {
	if resp, err := HTTPGet("http://localhost:8888/two/test2/0/hello/word/test"); nil != err {
		t.Skip(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Skip("unable to read response body:", err.Error())
		} else {
			t.Log(string(bytes))
		}
	}
}

func TestHttpClientCommon_GetTLS(t *testing.T) {
	if resp, err := HTTPGetTLS("https://localhost:8888/two/test2/0/hello/word", &HTTPTLSConfig{
		RootCrtFilePath:    "/Users/aberic/Downloads/test1/org.root.cert.pem",
		CertFilePath:       "/Users/aberic/Downloads/test1/client.org.cert.pem",
		KeyFilePath:        "/Users/aberic/Downloads/test1/client.key.pem",
		InsecureSkipVerify: false,
	}); nil != err {
		t.Skip(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Skip("unable to read response body:", err.Error())
		} else {
			t.Log(string(bytes))
		}
	}
}

func TestHttpClientCommon_GetTLS1(t *testing.T) {
	if resp, err := HTTPGetTLS("https://localhost:8888/two/test2/0/hello/pass/test", &HTTPTLSConfig{
		RootCrtFilePath:    "/Users/aberic/Downloads/test1/org.root.cert.pem",
		CertFilePath:       "/Users/aberic/Downloads/test1/client.org.cert.pem",
		KeyFilePath:        "/Users/aberic/Downloads/test1/client.key.pem",
		InsecureSkipVerify: false,
	}); nil != err {
		t.Skip(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Skip("unable to read response body:", err.Error())
		} else {
			t.Log(string(bytes))
		}
	}
}

func TestHttpClientCommon_GetTLSProxy(t *testing.T) {
	if resp, err := HttpGetTLSProxy("https://localhost:8888", "https://localhost:8888/two/test2/0/hello/pass/test", &HTTPTLSConfig{
		RootCrtFilePath:    "/Users/aberic/Downloads/test1/org.root.cert.pem",
		CertFilePath:       "/Users/aberic/Downloads/test1/client.org.cert.pem",
		KeyFilePath:        "/Users/aberic/Downloads/test1/client.key.pem",
		InsecureSkipVerify: false,
	}); nil != err {
		t.Skip(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Skip("unable to read response body:", err.Error())
		} else {
			t.Log(string(bytes))
		}
	}
}

func TestHttpClientCommon_GetTLSDial(t *testing.T) {
	resolver := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, _, _ string) (net.Conn, error) {
			d := net.Dialer{}
			return d.DialContext(ctx, "tcp", "127.0.0.1:8888")
		},
	}
	net.DefaultResolver = resolver
	ctx, cancelfunc := context.WithTimeout(context.Background(), 30)
	defer cancelfunc()
	_, err := resolver.LookupHost(ctx, "www.baidu.com")

	fmt.Printf("Error is: %s\n", err)
}

func TestHttpClientCommon_Post(t *testing.T) {
	if resp, err := HTTPPostJSON("http://localhost:8888/one/test/test1?test=he&go=to", &TestOne{
		One:   "1",
		Ones:  true,
		OneGo: 1,
	}); nil != err {
		t.Skip(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Skip("unable to read response body:", err.Error())
		} else {
			t.Log(string(bytes))
		}
	}
}

func TestHttpClientCommon_PostTLS(t *testing.T) {
	if resp, err := HTTPPostJSONTLS("https://localhost:8888/one/test/test1?test=he&go=to", &TestOne{
		One:   "1",
		Ones:  true,
		OneGo: 1,
	}, &HTTPTLSConfig{
		RootCrtFilePath: "./example/ca/server/rootCA.crt",
		CertFilePath:    "./example/ca/client/rootCA.crt",
		KeyFilePath:     "./example/ca/client/rootCA.key",
		//RootCrtFilePath:    "tmp/example/ca/fabric/pksc1/2048/rootCA.crt",
		//CertFilePath:       "tmp/example/ca/fabric/pksc1/2048/rootCA.crt",
		//KeyFilePath:        "tmp/example/ca/fabric/pksc1/2048/rootCA.key",
		InsecureSkipVerify: false,
	}); nil != err {
		t.Skip(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Skip("unable to read response body:", err.Error())
		} else {
			t.Log(string(bytes))
		}
	}
}

func TestHttpClientCommon_PostTLSBytes(t *testing.T) {
	var (
		rootCrtBytes, keyBytes, certBytes []byte
		err                               error
	)
	rootCrtBytes, err = ioutil.ReadFile("./example/ca/server/rootCA.crt")
	assert.NilError(t, err)
	keyBytes, err = ioutil.ReadFile("./example/ca/client/rootCA.key")
	assert.NilError(t, err)
	certBytes, err = ioutil.ReadFile("./example/ca/client/rootCA.crt")
	assert.NilError(t, err)
	if resp, err := HTTPPostJSONTLSBytes("https://localhost:8888/one/test/test1?test=he&go=to", &TestOne{
		One:   "1",
		Ones:  true,
		OneGo: 1,
	}, &HTTPTLSBytesConfig{
		RootCrtBytes:       rootCrtBytes,
		KeyBytes:           keyBytes,
		CertBytes:          certBytes,
		InsecureSkipVerify: false,
	}); nil != err {
		t.Skip(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Skip("unable to read response body:", err.Error())
		} else {
			t.Log(string(bytes))
		}
	}
}

func TestHttpClientCommon_PostForm(t *testing.T) {
	paramMap := map[string]string{}
	paramMap["xxx"] = "111"
	paramMap["yyy"] = "222"
	if resp, err := HTTPPostForm("http://localhost:8888/one/test3/x/y", paramMap); nil != err {
		t.Skip(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Skip("unable to read response body:", err.Error())
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
	if resp, err := HTTPPostFormMultipart("http://localhost:8888/one/test5/a/b", paramMap, fileMap); nil != err {
		t.Skip(err)
	} else {
		defer func() { _ = resp.Body.Close() }()
		if bytes, err := ioutil.ReadAll(resp.Body); err != nil {
			t.Skip("unable to read response body:", err.Error())
		} else {
			t.Log(string(bytes))
		}
	}
}

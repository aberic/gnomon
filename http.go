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
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// HttpClientCommon http客户端工具
type HttpClientCommon struct{}

func (hc *HttpClientCommon) Get(url string) (resp *http.Response, err error) {
	return hc.GetTLS(url, "", "", "")
}

// Post
//
// content-type=application/json
func (hc *HttpClientCommon) Post(url string, model interface{}) (resp *http.Response, err error) {
	return hc.PostTLS(url, model, "", "", "")
}

// Put
//
// content-type=application/json
func (hc *HttpClientCommon) Put(url string, model interface{}) (resp *http.Response, err error) {
	return hc.PutTLS(url, model, "", "", "")
}

// Patch
//
// content-type=application/json
func (hc *HttpClientCommon) Patch(url string, model interface{}) (resp *http.Response, err error) {
	return hc.PatchTLS(url, model, "", "", "")
}

func (hc *HttpClientCommon) Delete(url string) (resp *http.Response, err error) {
	return hc.DeleteTLS(url, "", "", "")
}

func (hc *HttpClientCommon) Do(req *http.Request) (resp *http.Response, err error) {
	return hc.DoTLS(req, "", "", "")
}

func (hc *HttpClientCommon) GetTLS(url string, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return hc.requestTLS(http.MethodGet, url, nil, caCrtFilePath, certFilePath, keyFilePath)
}

// PostTLS
//
// content-type=application/json
func (hc *HttpClientCommon) PostTLS(url string, model interface{}, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return hc.requestJson(http.MethodPost, url, model, caCrtFilePath, certFilePath, keyFilePath)
}

// PutTLS
//
// content-type=application/json
func (hc *HttpClientCommon) PutTLS(url string, model interface{}, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return hc.requestJson(http.MethodPut, url, model, caCrtFilePath, certFilePath, keyFilePath)
}

// PatchTLS
//
// content-type=application/json
func (hc *HttpClientCommon) PatchTLS(url string, model interface{}, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return hc.requestJson(http.MethodPatch, url, model, caCrtFilePath, certFilePath, keyFilePath)
}

func (hc *HttpClientCommon) DeleteTLS(url string, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return hc.requestTLS(http.MethodDelete, url, nil, caCrtFilePath, certFilePath, keyFilePath)
}

func (hc *HttpClientCommon) DoTLS(req *http.Request, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return hc.requestTLSDo(req, caCrtFilePath, certFilePath, keyFilePath)
}

// PostForm
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HttpClientCommon) PostForm(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return hc.PostFormTLS(url, paramMap, fileMap, "", "", "")
}

// PutForm
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HttpClientCommon) PutForm(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return hc.PutFormTLS(url, paramMap, fileMap, "", "", "")
}

// PatchForm
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HttpClientCommon) PatchForm(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return hc.PatchFormTLS(url, paramMap, fileMap, "", "", "")
}

// PostFormTLS
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HttpClientCommon) PostFormTLS(url string, paramMap map[string]string, fileMap map[string]string, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return hc.requestForm(http.MethodPost, url, paramMap, fileMap, caCrtFilePath, certFilePath, keyFilePath)
}

// PutFormTLS
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HttpClientCommon) PutFormTLS(url string, paramMap map[string]string, fileMap map[string]string, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return hc.requestForm(http.MethodPut, url, paramMap, fileMap, caCrtFilePath, certFilePath, keyFilePath)
}

// PatchFormTLS
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HttpClientCommon) PatchFormTLS(url string, paramMap map[string]string, fileMap map[string]string, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return hc.requestForm(http.MethodPatch, url, paramMap, fileMap, caCrtFilePath, certFilePath, keyFilePath)
}

func (hc *HttpClientCommon) requestJson(method, url string, model interface{}, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	var (
		data []byte
		req  *http.Request
	)
	if data, err = json.Marshal(model); err != nil {
		return nil, err
	}
	if req, err = http.NewRequest(method, url, bytes.NewReader(data)); nil != err {
		return
	}
	req.Header.Set("content-type", "application/json")
	return hc.requestTLSDo(req, caCrtFilePath, certFilePath, keyFilePath)
}

// requestForm
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HttpClientCommon) requestForm(method, url string, paramMap map[string]string, fileMap map[string]string, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	var (
		req        *http.Request
		bodyBuffer = &bytes.Buffer{}
		bodyWriter = multipart.NewWriter(bodyBuffer)
	)
	for key, value := range paramMap {
		if err = bodyWriter.WriteField(key, value); nil != err {
			return nil, err
		}
	}
	for key, value := range fileMap {
		var (
			fileSplitArr []string
			fileWriter   io.Writer
			file         *os.File
		)
		fileSplitArr = strings.Split(value, string(filepath.Separator))
		if fileWriter, err = bodyWriter.CreateFormFile(key, fileSplitArr[len(fileSplitArr)-1]); nil != err {
			return nil, err
		}
		if file, err = os.Open(value); nil != err {
			return nil, err
		}
		if _, err = io.Copy(fileWriter, file); nil != err {
			return nil, err
		}
		func() { _ = file.Close() }()
	}
	contentType := bodyWriter.FormDataContentType()
	if err = bodyWriter.Close(); nil != err {
		return nil, err
	}
	if req, err = http.NewRequest(method, url, bodyBuffer); nil != err {
		return
	}
	req.Header.Set("Content-Type", contentType)
	return hc.requestTLSDo(req, caCrtFilePath, certFilePath, keyFilePath)
}

func (hc *HttpClientCommon) requestTLS(method, url string, body io.Reader, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	var req *http.Request
	if req, err = http.NewRequest(method, url, body); nil != err {
		return
	}
	return hc.requestTLSDo(req, caCrtFilePath, certFilePath, keyFilePath)
}

func (hc *HttpClientCommon) requestTLSDo(req *http.Request, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	var (
		tlsClient    *http.Client
		tlsClientKey = CryptoHash().MD516(String().StringBuilder(caCrtFilePath, certFilePath, keyFilePath))
	)
	if tlsClient, err = getTLSClient(tlsClientKey, caCrtFilePath, certFilePath, keyFilePath); nil != err {
		return
	}
	return tlsClient.Do(req)
}

var (
	clients    = map[string]*http.Client{}
	clientLock sync.Mutex
)

func getTLSClient(tlsClientKey, caCrtFilePath, certFilePath, keyFilePath string) (*http.Client, error) {
	if tlsClient, exist := clients[tlsClientKey]; exist {
		return tlsClient, nil
	}
	defer clientLock.Unlock()
	clientLock.Lock()
	if tlsClient, exist := clients[tlsClientKey]; exist {
		return tlsClient, nil
	}
	var (
		tlsClient = &http.Client{}
		err       error
	)
	if tlsClient.Transport, err = getTlsTransport(caCrtFilePath, certFilePath, keyFilePath); nil != err {
		return nil, err
	}
	clients[tlsClientKey] = tlsClient
	return tlsClient, nil
}

func getTlsTransport(caCrtFilePath, certFilePath, keyFilePath string) (transport *http.Transport, err error) {
	var (
		pool               = x509.NewCertPool()
		caCrtBytes         []byte
		cert               tls.Certificate
		insecureSkipVerify = false // 是否验证服务端整数，即双向认证
	)
	if String().IsNotEmpty(caCrtFilePath) {
		// 对方验证我方整数合法性
		if caCrtBytes, err = ioutil.ReadFile(caCrtFilePath); nil != err {
			return
		}
		pool.AppendCertsFromPEM(caCrtBytes)
	}
	if String().IsNotEmpty(certFilePath) && String().IsNotEmpty(keyFilePath) {
		insecureSkipVerify = true
		// 我方验证对方整数合法性
		if cert, err = tls.LoadX509KeyPair(certFilePath, keyFilePath); nil != err {
			return
		}
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: insecureSkipVerify,
				Certificates:       []tls.Certificate{cert},
			},
		}
	} else {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: insecureSkipVerify,
			},
		}
	}
	return
}

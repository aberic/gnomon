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

// HTTPClientCommon http客户端工具
type HTTPClientCommon struct{}

// Get get 请求
func (hc *HTTPClientCommon) Get(url string) (resp *http.Response, err error) {
	return hc.GetTLS(url, &HTTPTLSConfig{})
}

// Post post 请求
//
// content-type=application/json
func (hc *HTTPClientCommon) Post(url string, model interface{}) (resp *http.Response, err error) {
	return hc.PostTLS(url, model, &HTTPTLSConfig{})
}

// Put put 请求
//
// content-type=application/json
func (hc *HTTPClientCommon) Put(url string, model interface{}) (resp *http.Response, err error) {
	return hc.PutTLS(url, model, &HTTPTLSConfig{})
}

// Patch patch 请求
//
// content-type=application/json
func (hc *HTTPClientCommon) Patch(url string, model interface{}) (resp *http.Response, err error) {
	return hc.PatchTLS(url, model, &HTTPTLSConfig{})
}

// Delete delete 请求
func (hc *HTTPClientCommon) Delete(url string) (resp *http.Response, err error) {
	return hc.DeleteTLS(url, &HTTPTLSConfig{})
}

// Do 自定义请求处理
func (hc *HTTPClientCommon) Do(req *http.Request) (resp *http.Response, err error) {
	return hc.DoTLS(req, &HTTPTLSConfig{})
}

// GetTLS get tls 请求
func (hc *HTTPClientCommon) GetTLS(url string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return hc.requestTLS(http.MethodGet, url, nil, tlsConfig)
}

// PostTLS post tls 请求
//
// content-type=application/json
func (hc *HTTPClientCommon) PostTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return hc.requestJSON(http.MethodPost, url, model, tlsConfig)
}

// PutTLS put tls 请求
//
// content-type=application/json
func (hc *HTTPClientCommon) PutTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return hc.requestJSON(http.MethodPut, url, model, tlsConfig)
}

// PatchTLS patch tls 请求
//
// content-type=application/json
func (hc *HTTPClientCommon) PatchTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return hc.requestJSON(http.MethodPatch, url, model, tlsConfig)
}

// DeleteTLS delete tls 请求
func (hc *HTTPClientCommon) DeleteTLS(url string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return hc.requestTLS(http.MethodDelete, url, nil, tlsConfig)
}

// DoTLS 处理 tls 请求
func (hc *HTTPClientCommon) DoTLS(req *http.Request, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return hc.requestTLSDo(req, tlsConfig)
}

// PostForm post 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HTTPClientCommon) PostForm(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return hc.PostFormTLS(url, paramMap, fileMap, &HTTPTLSConfig{})
}

// PutForm put 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HTTPClientCommon) PutForm(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return hc.PutFormTLS(url, paramMap, fileMap, &HTTPTLSConfig{})
}

// PatchForm patch 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HTTPClientCommon) PatchForm(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return hc.PatchFormTLS(url, paramMap, fileMap, &HTTPTLSConfig{})
}

// PostFormTLS post tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HTTPClientCommon) PostFormTLS(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return hc.requestForm(http.MethodPost, url, paramMap, fileMap, tlsConfig)
}

// PutFormTLS put tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HTTPClientCommon) PutFormTLS(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return hc.requestForm(http.MethodPut, url, paramMap, fileMap, tlsConfig)
}

// PatchFormTLS patch tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HTTPClientCommon) PatchFormTLS(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return hc.requestForm(http.MethodPatch, url, paramMap, fileMap, tlsConfig)
}

// requestJSON
//
// model 结构体
func (hc *HTTPClientCommon) requestJSON(method, url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
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
	return hc.requestTLSDo(req, tlsConfig)
}

// requestForm
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HTTPClientCommon) requestForm(method, url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
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
	return hc.requestTLSDo(req, tlsConfig)
}

func (hc *HTTPClientCommon) requestTLS(method, url string, body io.Reader, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	var req *http.Request
	if req, err = http.NewRequest(method, url, body); nil != err {
		return
	}
	return hc.requestTLSDo(req, tlsConfig)
}

func (hc *HTTPClientCommon) requestTLSDo(req *http.Request, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	var (
		tlsClient    *http.Client
		tlsClientKey string
	)
	if nil == tlsConfig {
		tlsClientKey = ""
	} else {
		tlsClientKey = CryptoHash().MD516(String().StringBuilder(tlsConfig.CACrtFilePath, tlsConfig.CertFilePath, tlsConfig.KeyFilePath))
	}
	if tlsClient, err = getTLSClient(tlsClientKey, tlsConfig); nil != err {
		return
	}
	return tlsClient.Do(req)
}

var (
	clients    = map[string]*http.Client{}
	clientLock sync.Mutex
)

func getTLSClient(tlsClientKey string, tlsConfig *HTTPTLSConfig) (*http.Client, error) {
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
	if tlsClient.Transport, err = getTLSTransport(tlsConfig); nil != err {
		return nil, err
	}
	clients[tlsClientKey] = tlsClient
	return tlsClient, nil
}

func getTLSTransport(tlsConfig *HTTPTLSConfig) (transport *http.Transport, err error) {
	var (
		pool       = x509.NewCertPool()
		caCrtBytes []byte
		cert       tls.Certificate
	)
	if String().IsNotEmpty(tlsConfig.CACrtFilePath) {
		// 用于我方验证对方证书合法性
		if caCrtBytes, err = ioutil.ReadFile(tlsConfig.CACrtFilePath); nil != err {
			return
		}
		pool.AppendCertsFromPEM(caCrtBytes)
	} else {
		tlsConfig.InsecureSkipVerify = false
	}
	if String().IsNotEmpty(tlsConfig.CertFilePath) && String().IsNotEmpty(tlsConfig.KeyFilePath) {
		// 用于对方验证我方证书合法性
		if cert, err = tls.LoadX509KeyPair(tlsConfig.CertFilePath, tlsConfig.KeyFilePath); nil != err {
			return
		}
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: tlsConfig.InsecureSkipVerify,
				Certificates:       []tls.Certificate{cert},
			},
		}
	} else {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: tlsConfig.InsecureSkipVerify,
			},
		}
	}
	return
}

// HTTPTLSConfig http tls 请求配置
type HTTPTLSConfig struct {
	CACrtFilePath      string // 服务端根证书，用于我方验证对方证书合法性
	CertFilePath       string // 服务端签发的子证书，用于对方验证我方证书合法性
	KeyFilePath        string // 客户端私钥，用于对方验证我方证书合法性
	InsecureSkipVerify bool   // 是否验证服务端证书，即双向认证
}

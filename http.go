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
	return hc.GetTLS(url, &HttpTLSConfig{})
}

// Post
//
// content-type=application/json
func (hc *HttpClientCommon) Post(url string, model interface{}) (resp *http.Response, err error) {
	return hc.PostTLS(url, model, &HttpTLSConfig{})
}

// Put
//
// content-type=application/json
func (hc *HttpClientCommon) Put(url string, model interface{}) (resp *http.Response, err error) {
	return hc.PutTLS(url, model, &HttpTLSConfig{})
}

// Patch
//
// content-type=application/json
func (hc *HttpClientCommon) Patch(url string, model interface{}) (resp *http.Response, err error) {
	return hc.PatchTLS(url, model, &HttpTLSConfig{})
}

func (hc *HttpClientCommon) Delete(url string) (resp *http.Response, err error) {
	return hc.DeleteTLS(url, &HttpTLSConfig{})
}

func (hc *HttpClientCommon) Do(req *http.Request) (resp *http.Response, err error) {
	return hc.DoTLS(req, &HttpTLSConfig{})
}

func (hc *HttpClientCommon) GetTLS(url string, tlsConfig *HttpTLSConfig) (resp *http.Response, err error) {
	return hc.requestTLS(http.MethodGet, url, nil, tlsConfig)
}

// PostTLS
//
// content-type=application/json
func (hc *HttpClientCommon) PostTLS(url string, model interface{}, tlsConfig *HttpTLSConfig) (resp *http.Response, err error) {
	return hc.requestJson(http.MethodPost, url, model, tlsConfig)
}

// PutTLS
//
// content-type=application/json
func (hc *HttpClientCommon) PutTLS(url string, model interface{}, tlsConfig *HttpTLSConfig) (resp *http.Response, err error) {
	return hc.requestJson(http.MethodPut, url, model, tlsConfig)
}

// PatchTLS
//
// content-type=application/json
func (hc *HttpClientCommon) PatchTLS(url string, model interface{}, tlsConfig *HttpTLSConfig) (resp *http.Response, err error) {
	return hc.requestJson(http.MethodPatch, url, model, tlsConfig)
}

func (hc *HttpClientCommon) DeleteTLS(url string, tlsConfig *HttpTLSConfig) (resp *http.Response, err error) {
	return hc.requestTLS(http.MethodDelete, url, nil, tlsConfig)
}

func (hc *HttpClientCommon) DoTLS(req *http.Request, tlsConfig *HttpTLSConfig) (resp *http.Response, err error) {
	return hc.requestTLSDo(req, tlsConfig)
}

// PostForm
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HttpClientCommon) PostForm(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return hc.PostFormTLS(url, paramMap, fileMap, &HttpTLSConfig{})
}

// PutForm
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HttpClientCommon) PutForm(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return hc.PutFormTLS(url, paramMap, fileMap, &HttpTLSConfig{})
}

// PatchForm
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HttpClientCommon) PatchForm(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return hc.PatchFormTLS(url, paramMap, fileMap, &HttpTLSConfig{})
}

// PostFormTLS
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HttpClientCommon) PostFormTLS(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HttpTLSConfig) (resp *http.Response, err error) {
	return hc.requestForm(http.MethodPost, url, paramMap, fileMap, tlsConfig)
}

// PutFormTLS
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HttpClientCommon) PutFormTLS(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HttpTLSConfig) (resp *http.Response, err error) {
	return hc.requestForm(http.MethodPut, url, paramMap, fileMap, tlsConfig)
}

// PatchFormTLS
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func (hc *HttpClientCommon) PatchFormTLS(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HttpTLSConfig) (resp *http.Response, err error) {
	return hc.requestForm(http.MethodPatch, url, paramMap, fileMap, tlsConfig)
}

// requestJson
//
// model 结构体
func (hc *HttpClientCommon) requestJson(method, url string, model interface{}, tlsConfig *HttpTLSConfig) (resp *http.Response, err error) {
	var (
		data []byte
		req  *http.Request
	)
	switch model.(type) {
	case []byte:
		data = model.([]byte)
	default:
		if data, err = json.Marshal(model); err != nil {
			return nil, err
		}
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
func (hc *HttpClientCommon) requestForm(method, url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HttpTLSConfig) (resp *http.Response, err error) {
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

func (hc *HttpClientCommon) requestTLS(method, url string, body io.Reader, tlsConfig *HttpTLSConfig) (resp *http.Response, err error) {
	var req *http.Request
	if req, err = http.NewRequest(method, url, body); nil != err {
		return
	}
	return hc.requestTLSDo(req, tlsConfig)
}

func (hc *HttpClientCommon) requestTLSDo(req *http.Request, tlsConfig *HttpTLSConfig) (resp *http.Response, err error) {
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

func getTLSClient(tlsClientKey string, tlsConfig *HttpTLSConfig) (*http.Client, error) {
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
	if tlsClient.Transport, err = getTlsTransport(tlsConfig); nil != err {
		return nil, err
	}
	clients[tlsClientKey] = tlsClient
	return tlsClient, nil
}

func getTlsTransport(tlsConfig *HttpTLSConfig) (transport *http.Transport, err error) {
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

// HttpTLSConfig http tls 请求配置
type HttpTLSConfig struct {
	CACrtFilePath      string // 服务端根证书，用于我方验证对方证书合法性
	CertFilePath       string // 服务端签发的子证书，用于对方验证我方证书合法性
	KeyFilePath        string // 客户端私钥，用于对方验证我方证书合法性
	InsecureSkipVerify bool   // 是否验证服务端证书，即双向认证
}

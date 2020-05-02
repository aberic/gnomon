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
	"encoding/xml"
	"github.com/golang/protobuf/proto"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// HTTPGet get 请求
func HTTPGet(url string) (resp *http.Response, err error) {
	return HTTPGetTLS(url, &HTTPTLSConfig{})
}

// HTTPPostJSON post 请求
//
// content-type=application/json
func HTTPPostJSON(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPPostJSONTLS(url, model, &HTTPTLSConfig{})
}

// HTTPPutJSON put 请求
//
// content-type=application/json
func HTTPPutJSON(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPPutJSONTLS(url, model, &HTTPTLSConfig{})
}

// HTTPPatchJSON patch 请求
//
// content-type=application/json
func HTTPPatchJSON(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPPatchJSONTLS(url, model, &HTTPTLSConfig{})
}

// HTTPPostXML post 请求
//
// content-type=application/xml
func HTTPPostXML(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPPostXMLTLS(url, model, &HTTPTLSConfig{})
}

// HTTPPutXML put 请求
//
// content-type=application/xml
func HTTPPutXML(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPPutXMLTLS(url, model, &HTTPTLSConfig{})
}

// HTTPPatchXML patch 请求
//
// content-type=application/xml
func HTTPPatchXML(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPPatchXMLTLS(url, model, &HTTPTLSConfig{})
}

// HTTPPostYaml post 请求
//
// content-type=application/x-yaml
func HTTPPostYaml(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPPostYamlTLS(url, model, &HTTPTLSConfig{})
}

// HTTPPutYaml put 请求
//
// content-type=application/x-yaml
func HTTPPutYaml(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPPutYamlTLS(url, model, &HTTPTLSConfig{})
}

// HTTPPatchYaml patch 请求
//
// content-type=application/x-yaml
func HTTPPatchYaml(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPPatchYamlTLS(url, model, &HTTPTLSConfig{})
}

// HTTPPostMsgPack post 请求
//
// content-type=application/x-msgpack
func HTTPPostMsgPack(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPPostMsgPackTLS(url, model, &HTTPTLSConfig{})
}

// HTTPPutMsgPack put 请求
//
// content-type=application/x-msgpack
func HTTPPutMsgPack(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPPutMsgPackTLS(url, model, &HTTPTLSConfig{})
}

// HTTPPatchMsgPack patch 请求
//
// content-type=application/x-msgpack
func HTTPPatchMsgPack(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPPatchMsgPackTLS(url, model, &HTTPTLSConfig{})
}

// HTTPPostProtoBuf post 请求
//
// content-type=application/x-protobuf
func HTTPPostProtoBuf(url string, pm proto.Message) (resp *http.Response, err error) {
	return HTTPPostProtoBufTLS(url, pm, &HTTPTLSConfig{})
}

// HTTPPutProtoBuf put 请求
//
// content-type=application/x-protobuf
func HTTPPutProtoBuf(url string, pm proto.Message) (resp *http.Response, err error) {
	return HTTPPutProtoBufTLS(url, pm, &HTTPTLSConfig{})
}

// HTTPPatchProtoBuf patch 请求
//
// content-type=application/x-protobuf
func HTTPPatchProtoBuf(url string, pm proto.Message) (resp *http.Response, err error) {
	return HTTPPatchProtoBufTLS(url, pm, &HTTPTLSConfig{})
}

// HTTPDelete delete 请求
func HTTPDelete(url string) (resp *http.Response, err error) {
	return HTTPDeleteTLS(url, &HTTPTLSConfig{})
}

// HTTPDo 自定义请求处理
func HTTPDo(req *http.Request) (resp *http.Response, err error) {
	return HTTPDoTLS(req, &HTTPTLSConfig{})
}

// HTTPGetTLS get tls 请求
func HTTPGetTLS(url string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestTLS(http.MethodGet, url, nil, tlsConfig)
}

// HTTPPostJSONTLS post tls 请求
//
// content-type=application/json
func HTTPPostJSONTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestJSON(http.MethodPost, url, model, tlsConfig)
}

// HTTPPutJSONTLS put tls 请求
//
// content-type=application/json
func HTTPPutJSONTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestJSON(http.MethodPut, url, model, tlsConfig)
}

// HTTPPatchJSONTLS patch tls 请求
//
// content-type=application/json
func HTTPPatchJSONTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestJSON(http.MethodPatch, url, model, tlsConfig)
}

// HTTPPostXMLTLS post tls 请求
//
// content-type=application/xml
func HTTPPostXMLTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestXML(http.MethodPost, url, model, tlsConfig)
}

// HTTPPutXMLTLS put tls 请求
//
// content-type=application/xml
func HTTPPutXMLTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestXML(http.MethodPut, url, model, tlsConfig)
}

// HTTPPatchXMLTLS patch tls 请求
//
// content-type=application/xml
func HTTPPatchXMLTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestXML(http.MethodPatch, url, model, tlsConfig)
}

// HTTPPostYamlTLS post tls 请求
//
// content-type=application/x-yaml
func HTTPPostYamlTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestYaml(http.MethodPost, url, model, tlsConfig)
}

// HTTPPutYamlTLS put tls 请求
//
// content-type=application/x-yaml
func HTTPPutYamlTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestYaml(http.MethodPut, url, model, tlsConfig)
}

// HTTPPatchYamlTLS patch tls 请求
//
// content-type=application/x-yaml
func HTTPPatchYamlTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestYaml(http.MethodPatch, url, model, tlsConfig)
}

// HTTPPostMsgPackTLS post tls 请求
//
// content-type=application/x-msgpack
func HTTPPostMsgPackTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestMsgPack(http.MethodPost, url, model, tlsConfig)
}

// HTTPPutMsgPackTLS put tls 请求
//
// content-type=application/x-msgpack
func HTTPPutMsgPackTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestMsgPack(http.MethodPut, url, model, tlsConfig)
}

// HTTPPatchMsgPackTLS patch tls 请求
//
// content-type=application/x-msgpack
func HTTPPatchMsgPackTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestMsgPack(http.MethodPatch, url, model, tlsConfig)
}

// HTTPPostProtoBufTLS post tls 请求
//
// content-type=application/x-protobuf
func HTTPPostProtoBufTLS(url string, pm proto.Message, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestProtoBuf(http.MethodPost, url, pm, tlsConfig)
}

// HTTPPutProtoBufTLS put tls 请求
//
// content-type=application/x-protobuf
func HTTPPutProtoBufTLS(url string, pm proto.Message, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestProtoBuf(http.MethodPut, url, pm, tlsConfig)
}

// HTTPPatchProtoBufTLS patch tls 请求
//
// content-type=application/x-protobuf
func HTTPPatchProtoBufTLS(url string, pm proto.Message, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestProtoBuf(http.MethodPatch, url, pm, tlsConfig)
}

// HTTPDeleteTLS delete tls 请求
func HTTPDeleteTLS(url string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestTLS(http.MethodDelete, url, nil, tlsConfig)
}

// HTTPDoTLS 处理 tls 请求
func HTTPDoTLS(req *http.Request, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestTLSDo(req, tlsConfig)
}

// HTTPPostForm post 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPostForm(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return HTTPPostFormTLS(url, paramMap, fileMap, &HTTPTLSConfig{})
}

// HTTPPutForm put 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPutForm(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return HTTPPutFormTLS(url, paramMap, fileMap, &HTTPTLSConfig{})
}

// HTTPPatchForm patch 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPatchForm(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return HTTPPatchFormTLS(url, paramMap, fileMap, &HTTPTLSConfig{})
}

// HTTPPostFormTLS post tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPostFormTLS(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestForm(http.MethodPost, url, paramMap, fileMap, tlsConfig)
}

// HTTPPutFormTLS put tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPutFormTLS(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestForm(http.MethodPut, url, paramMap, fileMap, tlsConfig)
}

// HTTPPatchFormTLS patch tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPatchFormTLS(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestForm(http.MethodPatch, url, paramMap, fileMap, tlsConfig)
}

// httpRequestJSON json 请求
//
// model 结构体
func httpRequestJSON(method, url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
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
	return httpRequestTLSDo(req, tlsConfig)
}

// httpRequestXML xml 请求
//
// model 结构体
func httpRequestXML(method, url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	var (
		data []byte
		req  *http.Request
	)
	if data, err = xml.Marshal(model); err != nil {
		return nil, err
	}
	if req, err = http.NewRequest(method, url, bytes.NewReader(data)); nil != err {
		return
	}
	req.Header.Set("content-type", "application/xml")
	return httpRequestTLSDo(req, tlsConfig)
}

// httpRequestYaml yaml 请求
//
// model 结构体
func httpRequestYaml(method, url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	var (
		data []byte
		req  *http.Request
	)
	if data, err = yaml.Marshal(model); err != nil {
		return nil, err
	}
	if req, err = http.NewRequest(method, url, bytes.NewReader(data)); nil != err {
		return
	}
	req.Header.Set("content-type", "application/x-yaml")
	return httpRequestTLSDo(req, tlsConfig)
}

// httpRequestMsgPack msgpack 请求
//
// model 结构体
func httpRequestMsgPack(method, url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	var (
		data []byte
		req  *http.Request
	)
	if data, err = msgpack.Marshal(model); err != nil {
		return nil, err
	}
	if req, err = http.NewRequest(method, url, bytes.NewReader(data)); nil != err {
		return
	}
	req.Header.Set("content-type", "application/x-msgpack")
	return httpRequestTLSDo(req, tlsConfig)
}

// httpRequestProtoBuf protobuf 请求
//
// model 结构体
func httpRequestProtoBuf(method, url string, pm proto.Message, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	var (
		data []byte
		req  *http.Request
	)
	if data, err = proto.Marshal(pm); err != nil {
		return nil, err
	}
	if req, err = http.NewRequest(method, url, bytes.NewReader(data)); nil != err {
		return
	}
	req.Header.Set("content-type", "application/x-protobuf")
	return httpRequestTLSDo(req, tlsConfig)
}

// httpRequestForm
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func httpRequestForm(method, url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
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
	return httpRequestTLSDo(req, tlsConfig)
}

func httpRequestTLS(method, url string, body io.Reader, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	var req *http.Request
	if req, err = http.NewRequest(method, url, body); nil != err {
		return
	}
	return httpRequestTLSDo(req, tlsConfig)
}

func httpRequestTLSDo(req *http.Request, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	var (
		tlsClient    *http.Client
		tlsClientKey string
	)
	if nil == tlsConfig {
		tlsClientKey = ""
	} else {
		tlsClientKey = HashMD516(StringBuild(tlsConfig.CACrtFilePath, tlsConfig.CertFilePath, tlsConfig.KeyFilePath))
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
	if StringIsNotEmpty(tlsConfig.CACrtFilePath) {
		// 用于我方验证对方证书合法性
		if caCrtBytes, err = ioutil.ReadFile(tlsConfig.CACrtFilePath); nil != err {
			return
		}
		pool.AppendCertsFromPEM(caCrtBytes)
	} else {
		tlsConfig.InsecureSkipVerify = false
	}
	if StringIsNotEmpty(tlsConfig.CertFilePath) && StringIsNotEmpty(tlsConfig.KeyFilePath) {
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

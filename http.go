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
	"net/url"
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

// HTTPDeleteJSON delete 请求
func HTTPDeleteJSON(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPDeleteJSONTLS(url, model, &HTTPTLSConfig{})
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

// HTTPDeleteXML delete 请求
func HTTPDeleteXML(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPDeleteXMLTLS(url, model, &HTTPTLSConfig{})
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

// HTTPDeleteYaml delete 请求
func HTTPDeleteYaml(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPDeleteYamlTLS(url, model, &HTTPTLSConfig{})
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

// HTTPDeleteMsgPack delete 请求
func HTTPDeleteMsgPack(url string, model interface{}) (resp *http.Response, err error) {
	return HTTPDeleteMsgPackTLS(url, model, &HTTPTLSConfig{})
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

// HTTPDeleteProtoBuf delete 请求
func HTTPDeleteProtoBuf(url string, pm proto.Message) (resp *http.Response, err error) {
	return HTTPDeleteProtoBufTLS(url, pm, &HTTPTLSConfig{})
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
	return HTTPGetTLSBytes(url, tlsConfig.trans())
}

// HTTPGetHostTLS get tls 请求
func HTTPGetHostTLS(url, host string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPGetHostTLSBytes(url, host, tlsConfig.trans())
}

func HttpGetTLSProxy(expectURL, proxyURL string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return httpRequestTLSProxy(expectURL, proxyURL, tlsConfig.trans())
}

// HTTPPostJSONTLS post tls 请求
//
// content-type=application/json
func HTTPPostJSONTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPostJSONTLSBytes(url, model, tlsConfig.trans())
}

// HTTPPutJSONTLS put tls 请求
//
// content-type=application/json
func HTTPPutJSONTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPutJSONTLSBytes(url, model, tlsConfig.trans())
}

// HTTPPatchJSONTLS patch tls 请求
//
// content-type=application/json
func HTTPPatchJSONTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPatchJSONTLSBytes(url, model, tlsConfig.trans())
}

// HTTPDeleteJSONTLS delete tls 请求
func HTTPDeleteJSONTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPDeleteJSONTLSBytes(url, model, tlsConfig.trans())
}

// HTTPPostXMLTLS post tls 请求
//
// content-type=application/xml
func HTTPPostXMLTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPostXMLTLSBytes(url, model, tlsConfig.trans())
}

// HTTPPutXMLTLS put tls 请求
//
// content-type=application/xml
func HTTPPutXMLTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPutXMLTLSBytes(url, model, tlsConfig.trans())
}

// HTTPPatchXMLTLS patch tls 请求
//
// content-type=application/xml
func HTTPPatchXMLTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPatchXMLTLSBytes(url, model, tlsConfig.trans())
}

// HTTPDeleteXMLTLS delete tls 请求
func HTTPDeleteXMLTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPDeleteXMLTLSBytes(url, model, tlsConfig.trans())
}

// HTTPPostYamlTLS post tls 请求
//
// content-type=application/x-yaml
func HTTPPostYamlTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPostYamlTLSBytes(url, model, tlsConfig.trans())
}

// HTTPPutYamlTLS put tls 请求
//
// content-type=application/x-yaml
func HTTPPutYamlTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPutYamlTLSBytes(url, model, tlsConfig.trans())
}

// HTTPPatchYamlTLS patch tls 请求
//
// content-type=application/x-yaml
func HTTPPatchYamlTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPatchYamlTLSBytes(url, model, tlsConfig.trans())
}

// HTTPDeleteYamlTLS delete tls 请求
func HTTPDeleteYamlTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPDeleteYamlTLSBytes(url, model, tlsConfig.trans())
}

// HTTPPostMsgPackTLS post tls 请求
//
// content-type=application/x-msgpack
func HTTPPostMsgPackTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPostMsgPackTLSBytes(url, model, tlsConfig.trans())
}

// HTTPPutMsgPackTLS put tls 请求
//
// content-type=application/x-msgpack
func HTTPPutMsgPackTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPutMsgPackTLSBytes(url, model, tlsConfig.trans())
}

// HTTPPatchMsgPackTLS patch tls 请求
//
// content-type=application/x-msgpack
func HTTPPatchMsgPackTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPatchMsgPackTLSBytes(url, model, tlsConfig.trans())
}

// HTTPDeleteMsgPackTLS delete tls 请求
func HTTPDeleteMsgPackTLS(url string, model interface{}, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPDeleteMsgPackTLSBytes(url, model, tlsConfig.trans())
}

// HTTPPostProtoBufTLS post tls 请求
//
// content-type=application/x-protobuf
func HTTPPostProtoBufTLS(url string, pm proto.Message, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPostProtoBufTLSBytes(url, pm, tlsConfig.trans())
}

// HTTPPutProtoBufTLS put tls 请求
//
// content-type=application/x-protobuf
func HTTPPutProtoBufTLS(url string, pm proto.Message, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPutProtoBufTLSBytes(url, pm, tlsConfig.trans())
}

// HTTPPatchProtoBufTLS patch tls 请求
//
// content-type=application/x-protobuf
func HTTPPatchProtoBufTLS(url string, pm proto.Message, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPatchProtoBufTLSBytes(url, pm, tlsConfig.trans())
}

// HTTPDeleteProtoBufTLS delete tls 请求
func HTTPDeleteProtoBufTLS(url string, pm proto.Message, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPDeleteProtoBufTLSBytes(url, pm, tlsConfig.trans())
}

// HTTPDeleteTLS delete tls 请求
func HTTPDeleteTLS(url string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPDeleteTLSBytes(url, tlsConfig.trans())
}

// HTTPDoTLS 处理 tls 请求
func HTTPDoTLS(req *http.Request, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPDoTLSBytes(req, tlsConfig.trans())
}

// HTTPGetTLSBytes get tls 请求
func HTTPGetTLSBytes(url string, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestTLSBytes(http.MethodGet, url, "", nil, tlsConfig)
}

// HTTPGetHostTLSBytes get tls 请求
func HTTPGetHostTLSBytes(url, host string, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestTLSBytes(http.MethodGet, url, host, nil, tlsConfig)
}

func HttpGetTLSBytesProxy(expectURL, proxyURL string, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestTLSProxy(expectURL, proxyURL, tlsConfig)
}

// HTTPPostJSONTLSBytes post tls 请求
//
// content-type=application/json
func HTTPPostJSONTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestJSON(http.MethodPost, url, model, tlsConfig)
}

// HTTPPutJSONTLSBytes put tls 请求
//
// content-type=application/json
func HTTPPutJSONTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestJSON(http.MethodPut, url, model, tlsConfig)
}

// HTTPPatchJSONTLSBytes patch tls 请求
//
// content-type=application/json
func HTTPPatchJSONTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestJSON(http.MethodPatch, url, model, tlsConfig)
}

// HTTPDeleteJSONTLSBytes delete tls 请求
func HTTPDeleteJSONTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestJSON(http.MethodDelete, url, model, tlsConfig)
}

// HTTPPostXMLTLSBytes post tls 请求
//
// content-type=application/xml
func HTTPPostXMLTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestXML(http.MethodPost, url, model, tlsConfig)
}

// HTTPPutXMLTLSBytes put tls 请求
//
// content-type=application/xml
func HTTPPutXMLTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestXML(http.MethodPut, url, model, tlsConfig)
}

// HTTPPatchXMLTLSBytes patch tls 请求
//
// content-type=application/xml
func HTTPPatchXMLTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestXML(http.MethodPatch, url, model, tlsConfig)
}

// HTTPDeleteXMLTLSBytes delete tls 请求
func HTTPDeleteXMLTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestXML(http.MethodDelete, url, model, tlsConfig)
}

// HTTPPostYamlTLSBytes post tls 请求
//
// content-type=application/x-yaml
func HTTPPostYamlTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestYaml(http.MethodPost, url, model, tlsConfig)
}

// HTTPPutYamlTLSBytes put tls 请求
//
// content-type=application/x-yaml
func HTTPPutYamlTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestYaml(http.MethodPut, url, model, tlsConfig)
}

// HTTPPatchYamlTLSBytes patch tls 请求
//
// content-type=application/x-yaml
func HTTPPatchYamlTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestYaml(http.MethodPatch, url, model, tlsConfig)
}

// HTTPDeleteYamlTLSBytes delete tls 请求
func HTTPDeleteYamlTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestYaml(http.MethodDelete, url, model, tlsConfig)
}

// HTTPPostMsgPackTLSBytes post tls 请求
//
// content-type=application/x-msgpack
func HTTPPostMsgPackTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestMsgPack(http.MethodPost, url, model, tlsConfig)
}

// HTTPPutMsgPackTLSBytes put tls 请求
//
// content-type=application/x-msgpack
func HTTPPutMsgPackTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestMsgPack(http.MethodPut, url, model, tlsConfig)
}

// HTTPPatchMsgPackTLSBytes patch tls 请求
//
// content-type=application/x-msgpack
func HTTPPatchMsgPackTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestMsgPack(http.MethodPatch, url, model, tlsConfig)
}

// HTTPDeleteMsgPackTLSBytes delete tls 请求
func HTTPDeleteMsgPackTLSBytes(url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestMsgPack(http.MethodDelete, url, model, tlsConfig)
}

// HTTPPostProtoBufTLSBytes post tls 请求
//
// content-type=application/x-protobuf
func HTTPPostProtoBufTLSBytes(url string, pm proto.Message, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestProtoBuf(http.MethodPost, url, pm, tlsConfig)
}

// HTTPPutProtoBufTLSBytes put tls 请求
//
// content-type=application/x-protobuf
func HTTPPutProtoBufTLSBytes(url string, pm proto.Message, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestProtoBuf(http.MethodPut, url, pm, tlsConfig)
}

// HTTPPatchProtoBufTLSBytes patch tls 请求
//
// content-type=application/x-protobuf
func HTTPPatchProtoBufTLSBytes(url string, pm proto.Message, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestProtoBuf(http.MethodPatch, url, pm, tlsConfig)
}

// HTTPDeleteProtoBufTLSBytes delete tls 请求
//
// content-type=application/x-protobuf
func HTTPDeleteProtoBufTLSBytes(url string, pm proto.Message, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestProtoBuf(http.MethodDelete, url, pm, tlsConfig)
}

// HTTPDeleteTLSBytes delete tls 请求
func HTTPDeleteTLSBytes(url string, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestTLSBytes(http.MethodDelete, url, "", nil, tlsConfig)
}

// HTTPDoTLSBytes 处理 tls 请求
func HTTPDoTLSBytes(req *http.Request, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestTLSBytesDo(req, tlsConfig)
}

// HTTPPostForm post 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPostForm(url string, paramMap map[string]string) (resp *http.Response, err error) {
	return HTTPPostFormTLS(url, paramMap, &HTTPTLSConfig{})
}

// HTTPPutForm put 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPutForm(url string, paramMap map[string]string) (resp *http.Response, err error) {
	return HTTPPutFormTLS(url, paramMap, &HTTPTLSConfig{})
}

// HTTPPatchForm patch 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPatchForm(url string, paramMap map[string]string) (resp *http.Response, err error) {
	return HTTPPatchFormTLS(url, paramMap, &HTTPTLSConfig{})
}

// HTTPPostFormMultipart post 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPostFormMultipart(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return HTTPPostFormMultipartTLS(url, paramMap, fileMap, &HTTPTLSConfig{})
}

// HTTPPutFormMultipart put 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPutFormMultipart(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return HTTPPutFormMultipartTLS(url, paramMap, fileMap, &HTTPTLSConfig{})
}

// HTTPPatchFormMultipart patch 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPatchFormMultipart(url string, paramMap map[string]string, fileMap map[string]string) (resp *http.Response, err error) {
	return HTTPPatchFormMultipartTLS(url, paramMap, fileMap, &HTTPTLSConfig{})
}

// HTTPPostFormTLS post tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPostFormTLS(url string, paramMap map[string]string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPostFormTLSBytes(url, paramMap, tlsConfig.trans())
}

// HTTPPutFormTLS put tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPutFormTLS(url string, paramMap map[string]string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPutFormTLSBytes(url, paramMap, tlsConfig.trans())
}

// HTTPPatchFormTLS patch tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPatchFormTLS(url string, paramMap map[string]string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPatchFormTLSBytes(url, paramMap, tlsConfig.trans())
}

// HTTPPostFormMultipartTLS post tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPostFormMultipartTLS(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPostFormMultipartTLSBytes(url, paramMap, fileMap, tlsConfig.trans())
}

// HTTPPutFormMultipartTLS put tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPutFormMultipartTLS(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPutFormMultipartTLSBytes(url, paramMap, fileMap, tlsConfig.trans())
}

// HTTPPatchFormMultipartTLS patch tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPatchFormMultipartTLS(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSConfig) (resp *http.Response, err error) {
	return HTTPPatchFormMultipartTLSBytes(url, paramMap, fileMap, tlsConfig.trans())
}

// HTTPPostFormTLSBytes post tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPostFormTLSBytes(url string, paramMap map[string]string, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestForm(http.MethodPost, url, paramMap, tlsConfig)
}

// HTTPPutFormTLSBytes put tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPutFormTLSBytes(url string, paramMap map[string]string, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestForm(http.MethodPut, url, paramMap, tlsConfig)
}

// HTTPPatchFormTLSBytes patch tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPatchFormTLSBytes(url string, paramMap map[string]string, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestForm(http.MethodPatch, url, paramMap, tlsConfig)
}

// HTTPPostFormMultipartTLSBytes post tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPostFormMultipartTLSBytes(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestFormMultipart(http.MethodPost, url, paramMap, fileMap, tlsConfig)
}

// HTTPPutFormMultipartTLSBytes put tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPutFormMultipartTLSBytes(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestFormMultipart(http.MethodPut, url, paramMap, fileMap, tlsConfig)
}

// HTTPPatchFormMultipartTLSBytes patch tls 请求
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func HTTPPatchFormMultipartTLSBytes(url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	return httpRequestFormMultipart(http.MethodPatch, url, paramMap, fileMap, tlsConfig)
}

// httpRequestJSON json 请求
//
// model 结构体
func httpRequestJSON(method, url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
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
	return httpRequestTLSBytesDo(req, tlsConfig)
}

// httpRequestXML xml 请求
//
// model 结构体
func httpRequestXML(method, url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
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
	return httpRequestTLSBytesDo(req, tlsConfig)
}

// httpRequestYaml yaml 请求
//
// model 结构体
func httpRequestYaml(method, url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
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
	return httpRequestTLSBytesDo(req, tlsConfig)
}

// httpRequestMsgPack msgpack 请求
//
// model 结构体
func httpRequestMsgPack(method, url string, model interface{}, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
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
	return httpRequestTLSBytesDo(req, tlsConfig)
}

// httpRequestProtoBuf protobuf 请求
//
// model 结构体
func httpRequestProtoBuf(method, url string, pm proto.Message, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
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
	return httpRequestTLSBytesDo(req, tlsConfig)
}

// httpRequestForm
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func httpRequestForm(method, url string, paramMap map[string]string, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
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
	if err = bodyWriter.Close(); nil != err {
		return nil, err
	}
	if req, err = http.NewRequest(method, url, bodyBuffer); nil != err {
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	return httpRequestTLSBytesDo(req, tlsConfig)
}

// httpRequestFormMultipart
//
// paramMap form普通参数
//
// fileMap form附件key及附件路径
func httpRequestFormMultipart(method, url string, paramMap map[string]string, fileMap map[string]string, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
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
	if err = bodyWriter.Close(); nil != err {
		return nil, err
	}
	if req, err = http.NewRequest(method, url, bodyBuffer); nil != err {
		return
	}
	req.Header.Set("Content-Type", "multipart/form-data")
	return httpRequestTLSBytesDo(req, tlsConfig)
}

func httpRequestTLSBytes(method, url, host string, body io.Reader, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	var req *http.Request
	if req, err = http.NewRequest(method, url, body); nil != err {
		return
	}
	if StringIsNotEmpty(host) {
		req.Host = host
	}
	return httpRequestTLSBytesDo(req, tlsConfig)
}

func httpRequestTLSBytesDo(req *http.Request, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	var (
		tlsClient    *http.Client
		tlsClientKey string
	)
	if nil == tlsConfig {
		tlsClientKey = ""
	} else {
		bs := append(tlsConfig.RootCrtBytes, append(tlsConfig.KeyBytes, tlsConfig.CertBytes...)...)
		tlsClientKey = HashMD516Bytes(bs)
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

func getTLSClient(tlsClientKey string, tlsConfig *HTTPTLSBytesConfig) (*http.Client, error) {
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

func getTLSTransport(tlsConfig *HTTPTLSBytesConfig) (transport *http.Transport, err error) {
	var (
		pool = x509.NewCertPool()
		cert tls.Certificate
	)
	if nil != tlsConfig.RootCrtBytes {
		pool.AppendCertsFromPEM(tlsConfig.RootCrtBytes)
	} else {
		tlsConfig.InsecureSkipVerify = false
	}
	if nil != tlsConfig.KeyBytes && nil != tlsConfig.CertBytes {
		// 用于对方验证我方证书合法性
		if cert, err = tls.X509KeyPair(tlsConfig.CertBytes, tlsConfig.KeyBytes); nil != err {
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

// expectURL https://localhost:8888
//
// proxyURL https://www.baidu.com:443
func httpRequestTLSProxy(expectURL, proxyURL string, tlsConfig *HTTPTLSBytesConfig) (resp *http.Response, err error) {
	var (
		u         *url.URL
		transport *http.Transport
	)
	if u, err = url.Parse(expectURL); err != nil {
		panic(err)
	}
	if transport, err = getTLSTransport(tlsConfig); nil != err {
		return nil, err
	}
	transport.Proxy = http.ProxyURL(u)
	// disabled HTTP/2
	//transport.TLSNextProto = make(map[string]func(authority string, c *tls.Conn) http.RoundTripper)
	//transport = &http.Transport{
	//	Proxy: http.ProxyURL(u),
	//	// disabled HTTP/2
	//	TLSNextProto: make(map[string]func(authority string, c *tls.Conn) http.RoundTripper),
	//}
	client := &http.Client{Transport: transport}
	return client.Get(proxyURL)
}

// HTTPTLSConfig http tls 请求配置
type HTTPTLSConfig struct {
	RootCrtFilePath    string // 服务端根证书，用于我方验证对方证书合法性
	CertFilePath       string // 服务端签发的子证书，用于对方验证我方证书合法性
	KeyFilePath        string // 客户端私钥，用于对方验证我方证书合法性
	InsecureSkipVerify bool   // 是否验证服务端证书，即双向认证
}

func (htc *HTTPTLSConfig) trans() *HTTPTLSBytesConfig {
	var (
		htbc                              = &HTTPTLSBytesConfig{InsecureSkipVerify: htc.InsecureSkipVerify}
		rootCrtBytes, keyBytes, certBytes []byte
		err                               error
	)
	if rootCrtBytes, err = ioutil.ReadFile(htc.RootCrtFilePath); nil == err {
		htbc.RootCrtBytes = rootCrtBytes
	}
	if certBytes, err = ioutil.ReadFile(htc.CertFilePath); nil == err {
		htbc.CertBytes = certBytes
	}
	if keyBytes, err = ioutil.ReadFile(htc.KeyFilePath); nil == err {
		htbc.KeyBytes = keyBytes
	}
	return htbc
}

// HTTPTLSBytesConfig http tls 请求配置
type HTTPTLSBytesConfig struct {
	RootCrtBytes       []byte // 服务端根证书，用于我方验证对方证书合法性
	CertBytes          []byte // 服务端签发的子证书，用于对方验证我方证书合法性
	KeyBytes           []byte // 客户端私钥，用于对方验证我方证书合法性
	InsecureSkipVerify bool   // 是否验证服务端证书，即双向认证
}

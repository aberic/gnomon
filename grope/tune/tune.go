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

package tune

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	// ContentTypeJSON application/json
	ContentTypeJSON = "application/json"
	// ContentTypeXML application/xml
	ContentTypeXML = "application/xml"
	// ContentTypePlain text/plain
	ContentTypePlain = "text/plain"
	// ContentTypePostForm application/x-www-form-urlencoded
	ContentTypePostForm = "application/x-www-form-urlencoded"
	// ContentTypeMultipartPostForm multipart/form-data
	ContentTypeMultipartPostForm = "multipart/form-data"
	// ContentTypeYaml application/x-yaml
	ContentTypeYaml = "application/x-yaml"
	// ContentTypeMsgPack application/x-msgpack
	ContentTypeMsgPack = "application/x-msgpack"
	// ContentTypeProtoBuf application/x-protobuf
	ContentTypeProtoBuf = "application/x-protobuf"
	//ContentTypeHtml              = "text/html"
	//ContentTypeXml2              = "text/xml"
	//ContentTypeMsgPack2          = "application/msgpack"
)

var (
	// ErrContentType context type error
	ErrContentType = errors.New("context type error")
)

// ParseJSON 解析请求参数
func ParseJSON(r *http.Request, obj interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, ContentTypeJSON) {
		if err := ValidateStruct(obj); nil != err {
			return err
		}
		return json.NewDecoder(r.Body).Decode(obj)
	}
	return ErrContentType
}

// ParseXML 解析请求参数
func ParseXML(r *http.Request, obj interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, ContentTypeXML) {
		if err := ValidateStruct(obj); nil != err {
			return err
		}
		return xml.NewDecoder(r.Body).Decode(obj)
	}
	return ErrContentType
}

// ParseYaml 解析请求参数
func ParseYaml(r *http.Request, obj interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, ContentTypeYaml) {
		if err := ValidateStruct(obj); nil != err {
			return err
		}
		return yaml.NewDecoder(r.Body).Decode(obj)
	}
	return ErrContentType
}

// ParseMsgPack 解析请求参数
func ParseMsgPack(r *http.Request, obj interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, ContentTypeMsgPack) {
		if err := ValidateStruct(obj); nil != err {
			return err
		}
		return msgpack.NewDecoder(r.Body).Decode(obj)
	}
	return ErrContentType
}

// ParseProtoBuf 解析请求参数
func ParseProtoBuf(r *http.Request, pm proto.Message) error {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, ContentTypeProtoBuf) {
		bs, err := ioutil.ReadAll(r.Body)
		if nil != err {
			return err
		}
		return proto.UnmarshalMerge(bs, pm)
	}
	return ErrContentType
}

// ParseForm 解析请求参数
func ParseForm(r *http.Request) (map[string]interface{}, error) {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, ContentTypePostForm) {
		if err := r.ParseForm(); nil != err { //解析参数，默认是不会解析的
			return nil, err
		}
		filedMap := make(map[string]interface{})
		for k, v := range r.Form {
			filedMap[k] = strings.Join(v, "")
		}
		return filedMap, nil
	}
	return nil, ErrContentType
}

// ParseMultipartForm 解析请求参数
func ParseMultipartForm(r *http.Request) (map[string]interface{}, error) {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, ContentTypeMultipartPostForm) {
		reader, err := r.MultipartReader()
		if nil == err {
			filedMap := make(map[string]interface{})
			for {
				part, err := reader.NextPart()
				if err == io.EOF {
					break
				}
				if part == nil {
					return nil, err
				}
				if part.FileName() == "" { // this is FormData
					if data, err := ioutil.ReadAll(part); nil == err {
						filedMap[part.FormName()] = string(data)
					} else {
						return nil, err
					}
				} else { // This is FileData
					if bytes, err := ioutil.ReadAll(part); nil == err {
						filedMap[part.FormName()] = &FormFile{
							FileName: part.FileName(),
							Data:     bytes,
						}
					} else {
						return nil, err
					}
				}
				func() { _ = part.Close() }()
			}
			return filedMap, nil
		}
		return nil, err
	}
	return nil, ErrContentType
}

// FormFile 表单附件信息
type FormFile struct {
	FileName string // file name
	Data     []byte // file bytes content
}

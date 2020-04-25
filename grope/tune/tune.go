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
	"errors"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	ContentTypeJson              = "application/json"
	ContentTypePlain             = "text/plain"
	ContentTypePostForm          = "application/x-www-form-urlencoded"
	ContentTypeMultipartPostForm = "multipart/form-data"
	ContentTypeYaml              = "application/x-yaml"
	ContentTypeMsgPack           = "application/x-msgpack"
	//ContentTypeHtml              = "text/html"
	//ContentTypeXml               = "application/xml"
	//ContentTypeXml2              = "text/xml"
	//ContentTypeProtoBuf          = "application/x-protobuf"
	//ContentTypeMsgPack2          = "application/msgpack"
)

var (
	ErrContentType = errors.New("context type error")
)

// ParseJson 解析请求参数
func ParseJson(r *http.Request, obj interface{}) error {
	contentType := r.Header.Get("Content-Type")
	if strings.Contains(contentType, ContentTypeJson) {
		if err := ValidateStruct(obj); nil != err {
			return err
		}
		return json.NewDecoder(r.Body).Decode(obj)
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
		if reader, err := r.MultipartReader(); nil != err {
			return nil, err
		} else {
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
					if data, err := ioutil.ReadAll(part); nil != err {
						return nil, err
					} else {
						filedMap[part.FormName()] = string(data)
					}
				} else { // This is FileData
					if bytes, err := ioutil.ReadAll(part); nil != err {
						return nil, err
					} else {
						filedMap[part.FormName()] = &FormFile{
							FileName: part.FileName(),
							Data:     bytes,
						}
					}
				}
				func() { _ = part.Close() }()
			}
			return filedMap, nil
		}
	}
	return nil, ErrContentType
}

type FormFile struct {
	FileName string // file name
	Data     []byte // file bytes content
}
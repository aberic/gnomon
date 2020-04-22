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

package grope

import (
	"encoding/json"
	"github.com/aberic/gnomon/grope/tune"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v3"
	"io"
	"net/http"
	"os"
	"strconv"
)

// Context grope 请求处理上下文
type Context struct {
	// writer 原生 net/http 结构
	writer http.ResponseWriter
	// request 原生 net/http 结构
	request *http.Request
	// valueMap 如果需要，这里是与url请求中对应的参数集合，如“/demo/:id/”，则通过 valueMap[id] 获取url中的值
	valueMap map[string]string
	// paramMap 如果需要，这里是请求params
	paramMap map[string]string
	// responded 已经处理过
	responded bool
}

func (c *Context) requestHeader(key string) string {
	return c.request.Header.Get(key)
}

func (c *Context) Request() *http.Request {
	return c.request
}

func (c *Context) HeaderSet(key, value string) {
	if value == "" {
		c.writer.Header().Del(key)
		return
	}
	c.writer.Header().Set(key, value)
}

// GetHeader returns value from request headers.
func (c *Context) HeaderGet(key string) string {
	return c.requestHeader(key)
}

func (c *Context) Status(code int) {
	c.writer.WriteHeader(code)
}

func (c *Context) Values() map[string]string {
	return c.valueMap
}

func (c *Context) Params() map[string]string {
	return c.paramMap
}

func (c *Context) ReceiveJson(model interface{}) (interface{}, error) {
	if err := tune.ParseJson(c.request, model); nil != err {
		return nil, err
	} else {
		return model, nil
	}
}

func (c *Context) ReceiveYaml(model interface{}) (interface{}, error) {
	if err := tune.ParseYaml(c.request, model); nil != err {
		return nil, err
	} else {
		return model, nil
	}
}

func (c *Context) ReceiveMsgPack(model interface{}) (interface{}, error) {
	if err := tune.ParseMsgPack(c.request, model); nil != err {
		return nil, err
	} else {
		return model, nil
	}
}

func (c *Context) ReceiveForm() (map[string]interface{}, error) {
	return tune.ParseForm(c.request)
}

func (c *Context) ReceiveMultipartForm() (map[string]interface{}, error) {
	return tune.ParseMultipartForm(c.request)
}
func (c *Context) response(bytes []byte) error {
	if _, err := c.writer.Write(bytes); nil != err {
		return err
	}
	return nil
}

// statusCode eg:http.StatusOK
func (c *Context) ResponseJson(statusCode int, model interface{}) error {
	if err := tune.ValidateStruct(model); nil != err {
		return err
	}
	c.responded = true
	c.HeaderSet("Content-Type", tune.ContentTypeJson)
	c.Status(statusCode)
	if bytes, err := json.Marshal(model); nil != err {
		return err
	} else {
		return c.response(bytes)
	}
}

// statusCode eg:http.StatusOK
func (c *Context) ResponseYaml(statusCode int, model interface{}) error {
	if err := tune.ValidateStruct(model); nil != err {
		return err
	}
	c.responded = true
	c.HeaderSet("Content-Type", tune.ContentTypeYaml)
	c.Status(statusCode)
	if bytes, err := yaml.Marshal(model); nil != err {
		return err
	} else {
		return c.response(bytes)
	}
}

// statusCode eg:http.StatusOK
func (c *Context) ResponseMsgPack(statusCode int, model interface{}) error {
	if err := tune.ValidateStruct(model); nil != err {
		return err
	}
	c.responded = true
	c.HeaderSet("Content-Type", tune.ContentTypeMsgPack)
	c.Status(statusCode)
	if bytes, err := msgpack.Marshal(model); nil != err {
		return err
	} else {
		return c.response(bytes)
	}
}

// statusCode eg:http.StatusOK
func (c *Context) ResponseText(statusCode int, text string) error {
	c.responded = true
	c.HeaderSet("Content-Type", tune.ContentTypePlain)
	c.Status(statusCode)
	return c.response([]byte(text))
}

// statusCode eg:http.StatusOK
func (c *Context) ResponseFile(statusCode int, filepath string) error {
	var (
		file     *os.File
		fileInfo os.FileInfo
		err      error
	)
	c.responded = true
	if file, err = os.Open(filepath); nil != err {
		return err
	}
	defer func() { _ = file.Close() }()
	fileHeader := make([]byte, 512)
	if _, err = file.Read(fileHeader); nil != err {
		return err
	}
	if fileInfo, err = file.Stat(); nil != err {
		return err
	}
	c.HeaderSet("Content-Disposition", "attachment; filename="+file.Name())
	c.HeaderSet("Content-Type", http.DetectContentType(fileHeader))
	c.HeaderSet("Content-Length", strconv.FormatInt(fileInfo.Size(), 10))
	c.Status(statusCode)

	if _, err := file.Seek(0, 0); nil != err {
		return err
	}
	_, err = io.Copy(c.writer, file)
	return err
}

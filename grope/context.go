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
	"github.com/aberic/gnomon"
	"github.com/aberic/gnomon/grope/tune"
	"github.com/vmihailenco/msgpack"
	"gopkg.in/yaml.v3"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

// Context grope 请求处理上下文
type Context struct {
	// writer 原生 net/http 结构
	writer http.ResponseWriter
	// request 原生 net/http 结构
	request *http.Request
	// sameSite 原生 net/http 结构, SameSite允许服务器定义cookie属性，使得浏览器不可能将此cookie与跨站点请求一起发送。
	sameSite http.SameSite
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

// SetSameSite with cookie
func (c *Context) SetSameSite(sameSite http.SameSite) {
	c.sameSite = sameSite
}

// SetCookie adds a Set-Cookie header to the ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (c *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(c.writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: c.sameSite,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

// Cookie returns the named cookie provided in the request or
// ErrNoCookie if not found. And return the named cookie is unescaped.
// If multiple cookies match the given name, only one cookie will
// be returned.
func (c *Context) Cookie(name string) (string, error) {
	cookie, err := c.request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (c *Context) ClientIP() string {
	return gnomon.IP().Get(c.request)
}

func filterFlags(content string) string {
	for i, char := range content {
		if char == ' ' || char == ';' {
			return content[:i]
		}
	}
	return content
}

// ContentType returns the Content-Type header of the request.
func (c *Context) ContentType() string {
	return filterFlags(c.requestHeader("Content-Type"))
}

// IsWebsocket returns true if the request headers indicate that a websocket
// handshake is being initiated by the client.
func (c *Context) IsWebsocket() bool {
	if strings.Contains(strings.ToLower(c.requestHeader("Connection")), "upgrade") &&
		strings.EqualFold(c.requestHeader("Upgrade"), "websocket") {
		return true
	}
	return false
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

func (c *Context) ReceiveJson(model interface{}) error {
	if err := tune.ParseJson(c.request, model); nil != err {
		return err
	} else {
		return nil
	}
}

func (c *Context) ReceiveYaml(model interface{}) error {
	if err := tune.ParseYaml(c.request, model); nil != err {
		return err
	} else {
		return nil
	}
}

func (c *Context) ReceiveMsgPack(model interface{}) error {
	if err := tune.ParseMsgPack(c.request, model); nil != err {
		return err
	} else {
		return nil
	}
}

func (c *Context) ReceiveForm() (map[string]interface{}, error) {
	return tune.ParseForm(c.request)
}

func (c *Context) ReceiveMultipartForm() (map[string]interface{}, error) {
	return tune.ParseMultipartForm(c.request)
}

// GetRawData return stream data.
func (c *Context) GetRawData() ([]byte, error) {
	return ioutil.ReadAll(c.request.Body)
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

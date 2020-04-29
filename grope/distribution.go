/*
 *  Copyright (c) 2020. aberic - All Rights Reserved.
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
	"crypto/tls"
	"crypto/x509"
	"github.com/aberic/gnomon"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"
)

var (
	clients    = map[string]*http.Client{}
	clientLock sync.Mutex
)

// Fusing 熔断处理
type Fusing func(err error)

// Distribution 请求转发
//
// addr是期望的转发路径，一般可指定为"http://ip:port"、"https://ip:port"、"http://url.com"
func (c *Context) Distribution(url string, fusing Fusing) {
	c.Distributions(url, &Transport{
		Timeout:               30 * time.Second,
		KeepAlive:             30 * time.Second,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   100,
	}, fusing)
}

// DistributionTLS 请求转发
//
// addr是期望的转发路径，一般可指定为"http://ip:port"、"https://ip:port"、"http://url.com"
//
// tLSConfig http tls 请求配置
func (c *Context) DistributionTLS(url string, tLSConfig *TLSConfig, fusing Fusing) {
	c.Distributions(url, &Transport{
		Timeout:               30 * time.Second,
		KeepAlive:             30 * time.Second,
		MaxIdleConns:          100,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
		MaxIdleConnsPerHost:   100,
		TLSConfig:             tLSConfig,
	}, fusing)
}

// Distributions 请求转发
//
// addr是期望的转发路径，一般可指定为"http://ip:port"、"https://ip:port"、"http://url.com"
//
// transport 支持HTTP和HTTPS的传输配置
func (c *Context) Distributions(addr string, transport *Transport, fusing Fusing) {
	var (
		client     *http.Client
		req        *http.Request
		resp       *http.Response
		data       []byte
		patternURL *url.URL
		realURL    string
		err        error
	)
	if patternURL, err = url.Parse(c.request.URL.String()); nil != err {
		goto ERR
	}
	realURL = gnomon.String().StringBuilder(addr, patternURL.String())
	if client, err = getTLSClient(transport); nil != err {
		goto ERR
	}
	if req, err = http.NewRequest(c.request.Method, realURL, c.request.Body); nil != err {
		goto ERR
	}
	// 设置Request头部信息
	for k, v := range c.request.Header {
		for _, vv := range v {
			req.Header.Add(k, vv)
		}
	}
	resp, err = client.Do(req)
	if err != nil {
		goto ERR
	}
	defer func() { _ = resp.Body.Close() }()
	// 设置Response头部信息
	for k, v := range resp.Header {
		for _, vv := range v {
			c.writer.Header().Add(k, vv)
		}
	}
	if data, err = ioutil.ReadAll(resp.Body); nil != err {
		goto ERR
	}
	if _, err = c.writer.Write(data); nil != err {
		goto ERR
	}
ERR:
	fusing(err)
}

func getTLSClient(transport *Transport) (*http.Client, error) {
	var tlsClientKey string
	if nil == transport.TLSConfig {
		tlsClientKey = ""
	} else {
		tlsClientKey = gnomon.CryptoHash().MD516(gnomon.String().StringBuilder(
			transport.TLSConfig.CACrtFilePath,
			transport.TLSConfig.CertFilePath,
			transport.TLSConfig.KeyFilePath))
	}
	if tlsClient, exist := clients[tlsClientKey]; exist {
		return tlsClient, nil
	}
	if tlsClient, exist := clients[tlsClientKey]; exist {
		return tlsClient, nil
	}
	return getTLSClientSync(tlsClientKey, transport)
}

func getTLSClientSync(tlsClientKey string, transport *Transport) (*http.Client, error) {
	defer clientLock.Unlock()
	clientLock.Lock()
	if tlsClient, exist := clients[tlsClientKey]; exist {
		return tlsClient, nil
	}
	var (
		tlsClient = &http.Client{}
		err       error
	)
	if tlsClient.Transport, err = getTlsTransport(transport); nil != err {
		return nil, err
	}
	clients[tlsClientKey] = tlsClient
	return tlsClient, nil
}

func getTlsTransport(transport *Transport) (ht *http.Transport, err error) {
	var (
		pool       = x509.NewCertPool()
		caCrtBytes []byte
		cert       tls.Certificate
	)
	ht = &http.Transport{
		DialContext: (&net.Dialer{
			Timeout:   transport.Timeout,
			KeepAlive: transport.KeepAlive,
		}).DialContext,
		MaxIdleConns:          transport.MaxIdleConns,
		IdleConnTimeout:       transport.IdleConnTimeout,
		TLSHandshakeTimeout:   transport.TLSHandshakeTimeout,
		ExpectContinueTimeout: transport.ExpectContinueTimeout,
		MaxIdleConnsPerHost:   transport.MaxIdleConnsPerHost,
	}
	if nil == transport.TLSConfig {
		return
	}
	if gnomon.String().IsNotEmpty(transport.TLSConfig.CACrtFilePath) {
		// 用于我方验证对方证书合法性
		if caCrtBytes, err = ioutil.ReadFile(transport.TLSConfig.CACrtFilePath); nil != err {
			return
		}
		pool.AppendCertsFromPEM(caCrtBytes)
	} else {
		transport.TLSConfig.InsecureSkipVerify = false
	}
	if gnomon.String().IsNotEmpty(transport.TLSConfig.CertFilePath) && gnomon.String().IsNotEmpty(transport.TLSConfig.KeyFilePath) {
		// 用于对方验证我方证书合法性
		if cert, err = tls.LoadX509KeyPair(transport.TLSConfig.CertFilePath, transport.TLSConfig.KeyFilePath); nil != err {
			return
		}
		ht.TLSClientConfig = &tls.Config{
			RootCAs:            pool,
			InsecureSkipVerify: transport.TLSConfig.InsecureSkipVerify,
			Certificates:       []tls.Certificate{cert},
		}
	} else {
		ht.TLSClientConfig = &tls.Config{RootCAs: pool, InsecureSkipVerify: transport.TLSConfig.InsecureSkipVerify}
	}
	return
}

// Transport 支持HTTP和HTTPS的传输配置
type Transport struct {
	// 等待连接完成的超时时间，默认30s
	Timeout time.Duration
	// 指定网络连接的活动实践间隔，默认30s
	KeepAlive time.Duration
	// 控制最大空闲(保持活动)连接数。0表示没有限制，默认100
	MaxIdleConns int
	// 一个空闲(保持活动)连接在关闭之前保持空闲的最大时间量。0表示没有限制，默认90s
	IdleConnTimeout time.Duration
	// 指定等待TLS握手的最大时间量。0表示没有限制，默认10s
	TLSHandshakeTimeout time.Duration
	// 指定在完全写入请求标头之后等待服务器的第一个响应头的时间量(如果请求具有“Expect: 100-continue”头信息)。
	// 零意味着没有超时，并且会立即发送正文，而无需等待服务器的批准。
	// 此时间不包括发送请求标头的时间。 默认1s
	ExpectContinueTimeout time.Duration
	// 指定每个主机的最大空闲(保持活动)连接，默认100
	MaxIdleConnsPerHost int
	// http tls 请求配置
	TLSConfig *TLSConfig
}

// TLSConfig http tls 请求配置
type TLSConfig struct {
	// 服务端根证书，用于我方验证对方证书合法性
	CACrtFilePath string
	// 服务端签发的子证书，用于对方验证我方证书合法性
	CertFilePath string
	// 客户端私钥，用于对方验证我方证书合法性
	KeyFilePath string
	// 是否验证服务端证书，即双向认证
	InsecureSkipVerify bool
}

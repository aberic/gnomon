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
	"io/ioutil"
	"net/http"
	"sync"
)

// HttpClientCommon http客户端工具
type HttpClientCommon struct{}

func (hc *HttpClientCommon) Get(url string) (resp *http.Response, err error) {
	return hc.request(http.MethodGet, url, nil)
}

func (hc *HttpClientCommon) Post(url string, data []byte) (resp *http.Response, err error) {
	return hc.request(http.MethodPost, url, data)
}

func (hc *HttpClientCommon) Put(url string, data []byte) (resp *http.Response, err error) {
	return hc.request(http.MethodPut, url, data)
}

func (hc *HttpClientCommon) Patch(url string, data []byte) (resp *http.Response, err error) {
	return hc.request(http.MethodPatch, url, data)
}

func (hc *HttpClientCommon) Delete(url string) (resp *http.Response, err error) {
	return hc.request(http.MethodDelete, url, nil)
}

func (hc *HttpClientCommon) Do(req *http.Request) (resp *http.Response, err error) {
	return getClient().Do(req)
}

func (hc *HttpClientCommon) GetTLS(url string, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return hc.requestTLS(http.MethodGet, url, nil, caCrtFilePath, certFilePath, keyFilePath)
}

func (hc *HttpClientCommon) PostTLS(url string, data []byte, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return hc.requestTLS(http.MethodPost, url, data, caCrtFilePath, certFilePath, keyFilePath)
}

func (hc *HttpClientCommon) PutTLS(url string, data []byte, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return hc.requestTLS(http.MethodPut, url, data, caCrtFilePath, certFilePath, keyFilePath)
}

func (hc *HttpClientCommon) PatchTLS(url string, data []byte, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return hc.requestTLS(http.MethodPatch, url, data, caCrtFilePath, certFilePath, keyFilePath)
}

func (hc *HttpClientCommon) DeleteTLS(url string, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return hc.requestTLS(http.MethodDelete, url, nil, caCrtFilePath, certFilePath, keyFilePath)
}

func (hc *HttpClientCommon) DoTLS(req *http.Request, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	return getClient().Do(req)
}

func (hc *HttpClientCommon) request(method, url string, data []byte) (resp *http.Response, err error) {
	var req *http.Request
	if req, err = http.NewRequest(method, url, bytes.NewReader(data)); nil != err {
		return
	}
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		req.Header.Set("content-type", "application/json")
	}
	return getClient().Do(req)
}

func (hc *HttpClientCommon) requestTLS(method, url string, data []byte, caCrtFilePath, certFilePath, keyFilePath string) (resp *http.Response, err error) {
	var req *http.Request
	if req, err = http.NewRequest(method, url, bytes.NewReader(data)); nil != err {
		return
	}
	switch method {
	case http.MethodPost, http.MethodPut, http.MethodPatch:
		req.Header.Set("content-type", "application/json")
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
	client        *http.Client
	onceClient    sync.Once
	tlsClients    map[string]*http.Client
	tlsClientLock sync.Mutex
)

func getClient() *http.Client {
	onceClient.Do(func() {
		client = &http.Client{}
	})
	return client
}

func getTLSClient(tlsClientKey, caCrtFilePath, certFilePath, keyFilePath string) (*http.Client, error) {
	if tlsClient, exist := tlsClients[tlsClientKey]; exist {
		return tlsClient, nil
	}
	defer tlsClientLock.Unlock()
	tlsClientLock.Lock()
	if tlsClient, exist := tlsClients[tlsClientKey]; exist {
		return tlsClient, nil
	}
	var (
		tlsClient = &http.Client{}
		err       error
	)
	if tlsClient.Transport, err = getTlsTransport(caCrtFilePath, certFilePath, keyFilePath); nil != err {
		return nil, err
	}
	tlsClients[tlsClientKey] = tlsClient
	return tlsClient, nil
}

func getTlsTransport(caCrtFilePath, certFilePath, keyFilePath string) (transport *http.Transport, err error) {
	var (
		pool       = x509.NewCertPool()
		caCrtBytes []byte
		cert       tls.Certificate
	)
	if String().IsNotEmpty(caCrtFilePath) {
		// 对方验证我方整数合法性
		if caCrtBytes, err = ioutil.ReadFile(caCrtFilePath); nil != err {
			return
		}
		pool.AppendCertsFromPEM(caCrtBytes)
	}
	if String().IsNotEmpty(certFilePath) && String().IsNotEmpty(keyFilePath) {
		// 我方验证对方整数合法性
		if cert, err = tls.LoadX509KeyPair(certFilePath, keyFilePath); nil != err {
			return
		}
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: true,
				Certificates:       []tls.Certificate{cert},
			},
		}
	} else {
		transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				RootCAs:            pool,
				InsecureSkipVerify: true,
			},
		}
	}
	return
}

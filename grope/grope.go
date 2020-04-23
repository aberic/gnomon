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
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"github.com/aberic/gnomon/grope/log"
	"io/ioutil"
	"net/http"
	"time"
)

// NewHttpServe 新建一个Http服务
func NewHttpServe(filters ...Filter) *GHttpServe {
	return NewGHttpServe(filters...)
}

// ListenAndServe 启动监听
//
// Addr 期望监听的端口号，如“:8080”
func ListenAndServe(Addr string, gs *GHttpServe) {
	err := http.ListenAndServe(Addr, gs) //设置监听的端口
	if err != nil {
		log.Panic("ListenAndServe", log.Err(err))
	}
}

// ListenAndServeTLS 启动监听
//
// Addr 期望监听的端口号，如“:8080”
//
// 必须提供包含证书和与服务器匹配的私钥的文件。如果证书是由证书颁发机构签署的，则certFile应该是服务器证书、任何中间体和CA证书的连接。
func ListenAndServeTLS(gs *GHttpServe, Addr, certFilePath, keyFilePath string, caCertFilePaths ...string) {
	pool := x509.NewCertPool()
	//加载根证书，用于验证对方合法性
	for _, caCertFilePath := range caCertFilePaths {
		//这里读取的是根证书
		if buf, err := ioutil.ReadFile(caCertFilePath); err != nil {
			log.Panic("ListenAndServeTLS ReadFile", log.Err(err))
		} else {
			pool.AppendCertsFromPEM(buf)
		}
	}
	//加载服务端证书，用于对方验证我方合法性
	if cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath); err != nil {
		log.Panic("ListenAndServeTLS LoadX509KeyPair", log.Err(err))
	} else {
		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			ClientAuth:   tls.RequireAndVerifyClientCert,
			ClientCAs:    pool,
		}
		tlsConfig.Time = time.Now
		tlsConfig.Rand = rand.Reader
		if listener, err := tls.Listen("tcp", Addr, tlsConfig); nil != err {
			log.Panic("Serve", log.Err(err))
		} else {
			log.Panic("Serve", log.Err(http.Serve(listener, gs)))
		}
	}
}

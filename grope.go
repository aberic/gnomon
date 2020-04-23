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

//
//import (
//	"crypto/rand"
//	"crypto/tls"
//	"crypto/x509"
//	"github.com/aberic/gnomon/grope"
//	"io/ioutil"
//	"log"
//	"net/http"
//	"time"
//)
//
//// GHTTPCommon Http工具
//type GHTTPCommon struct{}
//
//// NewHttpServe 新建一个Http服务
//func (ghc *GHTTPCommon) NewHttpServe(filters ...grope.Filter) *grope.GHttpServe {
//	return grope.NewGHttpServe(filters...)
//}
//
//// ListenAndServe 启动监听
////
//// Addr 期望监听的端口号，如“:8080”
//func (ghc *GHTTPCommon) ListenAndServe(Addr string, gs *grope.GHttpServe) {
//	err := http.ListenAndServe(Addr, gs) //设置监听的端口
//	if err != nil {
//		log.Panic("ListenAndServe: ", err)
//	}
//}
//
//// ListenAndServeTLS 启动监听
////
//// Addr 期望监听的端口号，如“:8080”
////
//// 必须提供包含证书和与服务器匹配的私钥的文件。如果证书是由证书颁发机构签署的，则certFile应该是服务器证书、任何中间体和CA证书的连接。
//func (ghc *GHTTPCommon) ListenAndServeTLS(gs *grope.GHttpServe, Addr, certFilePath, keyFilePath string, caCertFilePaths ...string) {
//	pool := x509.NewCertPool()
//	//加载根证书，用于验证对方合法性
//	for _, caCertFilePath := range caCertFilePaths {
//		//这里读取的是根证书
//		if buf, err := ioutil.ReadFile(caCertFilePath); err != nil {
//			log.Panic("ListenAndServeTLS ReadFile", err)
//		} else {
//			pool.AppendCertsFromPEM(buf)
//		}
//	}
//	//加载服务端证书，用于对方验证我方合法性
//	if cert, err := tls.LoadX509KeyPair(certFilePath, keyFilePath); err != nil {
//		log.Panic("ListenAndServeTLS LoadX509KeyPair", err)
//	} else {
//		tlsConfig := &tls.Config{
//			Certificates: []tls.Certificate{cert},
//			ClientAuth:   tls.RequireAndVerifyClientCert,
//			ClientCAs:    pool,
//		}
//		tlsConfig.Time = time.Now
//		tlsConfig.Rand = rand.Reader
//		if listener, err := tls.Listen("tcp", Addr, tlsConfig); nil != err {
//			log.Fatalln(err)
//		} else {
//			log.Panic(http.Serve(listener, gs))
//		}
//	}
//}

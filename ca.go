/*
 *
 *  * Copyright (c) 2019. aberic - All Rights Reserved.
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  * http://www.apache.org/licenses/LICENSE-2.0
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 *
 */

package gnomon

import (
	"crypto/elliptic"
	"crypto/x509"
	"io"
)

// CACommon CA工具
type CACommon struct{}

// GenerateRSAPKCS1PrivateKey 生成CA自己的私钥——RSA私钥产生（私钥PKCS1格式）
//
// bits 指定生成位大小
//
// path 指定私钥所在生成目录
//
// fileName 指定私钥的文件名称，如'rootCA.key'
func (ca *CACommon) GenerateRSAPKCS1PrivateKey(bits int, path, fileName string) error {
	return CryptoRSA().GeneratePKCS1PriKey(bits, path, fileName)
}

// GenerateRSAPKCS8PrivateKey 生成CA自己的私钥——RSA私钥产生（私钥PKCS8格式）
//
// bits 指定生成位大小
//
// path 指定私钥所在生成目录
//
// fileName 指定私钥的文件名称，如'rootCA.key'
func (ca *CACommon) GenerateRSAPKCS8PrivateKey(bits int, path, fileName string) error {
	return CryptoRSA().GeneratePKCS8PriKey(bits, path, fileName)
}

// GenerateECCPrivateKey 生成CA自己的私钥——生成ECC私钥
//
// path 指定私钥所在生成目录
//
// fileName 指定生成的密钥名称，如'rootCA.key'
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func (ca *CACommon) GenerateECCPrivateKey(path, fileName string, curve elliptic.Curve) error {
	return CryptoECC().GeneratePemPriKey(path, fileName, curve)
}

// GenerateCertificate 生成证书
func (ca *CACommon) GenerateCertificate(rand io.Reader, template *x509.CertificateRequest, priv interface{}) (csr []byte, err error) {
	return x509.CreateCertificateRequest(rand, template, priv)
}

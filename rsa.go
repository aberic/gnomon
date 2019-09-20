/*
 * Copyright (c) 2019. aberic - All Rights Reserved.
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
 *
 */

package gnomon

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

const (
	publicKeyPemType  = "PUBLIC KEY"
	privateKeyPemType = "RSA PRIVATE KEY"
)

// RSACommon rsa工具
type RSACommon struct{}

// GenerateRsaKey RSA公钥私钥产生
//
// bits 指定生成位大小
//
// path 指定公私钥所在生成目录
func (r *RSACommon) GenerateRsaKey(bits int, path string) (err error) {
	var (
		privateKey *rsa.PrivateKey
		publicKey  *rsa.PublicKey
		fileIO     *os.File
		derPkix    []byte
	)
	// 创建公私钥生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return
		}
	}
	// 生成私钥文件
	if privateKey, err = rsa.GenerateKey(rand.Reader, bits); nil != err {
		return
	}
	// 将私钥转换为ASN.1 DER编码的形式
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	// block表示PEM编码的结构
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	if fileIO, err = os.Create(strings.Join([]string{path, "private.pem"}, "/")); nil != err {
		return
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, block); nil != err {
		return
	}
	// 生成公钥文件
	publicKey = &privateKey.PublicKey
	// 将公钥序列化为der编码的PKIX格式
	if derPkix, err = x509.MarshalPKIXPublicKey(publicKey); nil != err {
		return
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	if fileIO, err = os.Create(strings.Join([]string{path, "public.pem"}, "/")); nil != err {
		return
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, block); nil != err {
		return
	}
	return
}

// RsaPubEncrypt 公钥加密
func (r *RSACommon) RsaPubEncrypt(publicKey, data []byte) ([]byte, error) {
	pub, err := r.parsePublicKey(publicKey)
	if nil != err {
		return nil, err
	}
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// RsaPubEncryptFP 公钥加密
func (r *RSACommon) RsaPubEncryptFP(publicKeyFilePath string, data []byte) ([]byte, error) {
	bs, err := ioutil.ReadFile(publicKeyFilePath)
	if nil != err {
		return nil, err
	}
	return r.RsaPubEncrypt(bs, data)
}

// RsaPriDecrypt 私钥解密
func (r *RSACommon) RsaPriDecrypt(privateKey, data []byte) ([]byte, error) {
	pri, err := r.parsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, pri, data)
}

// RsaPriDecryptFP 私钥解密
func (r *RSACommon) RsaPriDecryptFP(privateKeyFilePath string, data []byte) ([]byte, error) {
	bs, err := ioutil.ReadFile(privateKeyFilePath)
	if nil != err {
		return bs, err
	}
	return r.RsaPriDecrypt(bs, data)
}

// RsaPriSign 签名：采用sha1算法进行签名并输出为hex格式（私钥PKCS8格式）
func (r *RSACommon) RsaPriSign(privateKey, data []byte) ([]byte, error) {
	pri, err := r.parsePrivateKey(privateKey)
	if err != nil {
		return nil, err
	}
	hash := crypto.SHA1
	h := hash.New()
	h.Write(data)
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, pri, hash, hashed)
}

// RsaPriSignFP 签名：采用sha1算法进行签名并输出为hex格式（私钥PKCS8格式）
func (r *RSACommon) RsaPriSignFP(privateKeyPath string, data []byte) ([]byte, error) {
	bs, err := ioutil.ReadFile(privateKeyPath)
	if nil != err {
		return nil, err
	}
	return r.RsaPriSign(bs, data)
}

// RsaPubVerySign 验签：对采用sha1算法进行签名后转base64格式的数据进行验签
func (r *RSACommon) RsaPubVerySign(publicKey, data, signData []byte) error {
	pub, err := r.parsePublicKey(publicKey)
	if nil != err {
		return err
	}
	hash := crypto.SHA1
	h := hash.New()
	h.Write(data)
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(pub, hash, hashed, signData)
}

// RsaPubVerySignFP 验签：对采用sha1算法进行签名后转base64格式的数据进行验签
func (r *RSACommon) RsaPubVerySignFP(privateKeyPath string, data, signData []byte) error {
	bs, err := ioutil.ReadFile(privateKeyPath)
	if nil != err {
		return err
	}
	return r.RsaPubVerySign(bs, data, signData)
}

func (r *RSACommon) parsePrivateKey(key []byte) (*rsa.PrivateKey, error) {
	pemData, err := r.pemParse(key, privateKeyPemType)
	if err != nil {
		return nil, err
	}
	return x509.ParsePKCS1PrivateKey(pemData)
}

func (r *RSACommon) parsePublicKey(key []byte) (*rsa.PublicKey, error) {
	pemData, err := r.pemParse(key, publicKeyPemType)
	if err != nil {
		return nil, err
	}
	keyInterface, err := x509.ParsePKIXPublicKey(pemData)
	if err != nil {
		return nil, err
	}
	pubKey, ok := keyInterface.(*rsa.PublicKey)
	if !ok {
		return nil, errors.New("could not cast parsed key to *rsa.PublicKey")
	}
	return pubKey, nil
}

// pemParse 解密pem格式密钥并验证pem类型
func (r *RSACommon) pemParse(key []byte, pemType string) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("no pem block found")
	}
	if pemType != "" && block.Type != pemType {
		return nil, errors.New(strings.Join([]string{"Key's type is ", block.Type, ", expected ", pemType}, ""))
	}
	return block.Bytes, nil
}

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
	"path/filepath"
	"strings"
)

const (
	pksC1                PKSCType = "PKCS1"
	pksC8                PKSCType = "PKCS8"
	publicRSAKeyPemType           = "PUBLIC KEY"
	privateRSAKeyPemType          = "RSA PRIVATE KEY"
)

// PKSCType 私钥格式，默认提供PKCS1和PKCS8
//
// 通过调用‘CryptoRSA().pksC1()’和‘CryptoRSA().pksC8()’方法赋值
type PKSCType string

// RSACommon RSA工具
type RSACommon struct{}

// GeneratePKCS1Key RSA公钥私钥产生（私钥PKCS1格式）
//
// bits 指定生成位大小
//
// path 指定公私钥所在生成目录
//
// priFileName 指定私钥的文件名称，如'private.pem'
//
// pubFileName 指定公钥的文件名称，如'public.pem'
func (r *RSACommon) GeneratePKCS1Key(bits int, path, priFileName, pubFileName string) error {
	var (
		privateKey *rsa.PrivateKey
		publicKey  *rsa.PublicKey
		fileIOPri  *os.File
		fileIOPub  *os.File
		derPkiX    []byte
		err        error
	)
	// 创建公私钥生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	// 生成私钥文件
	if privateKey, err = rsa.GenerateKey(rand.Reader, bits); nil != err {
		return err
	}
	// 将私钥转换为ASN.1 DER编码的形式
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	// block表示PEM编码的结构
	block := &pem.Block{
		Type:  privateRSAKeyPemType,
		Bytes: derStream,
	}
	defer func() { _ = fileIOPri.Close() }()
	if fileIOPri, err = os.Create(filepath.Join(path, priFileName)); nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIOPri, block); nil != err {
		return err
	}
	// 生成公钥文件
	publicKey = &privateKey.PublicKey
	// 将公钥序列化为der编码的PKIX格式
	if derPkiX, err = x509.MarshalPKIXPublicKey(publicKey); nil != err {
		return err
	}
	block = &pem.Block{
		Type:  publicRSAKeyPemType,
		Bytes: derPkiX,
	}
	defer func() { _ = fileIOPub.Close() }()
	if fileIOPub, err = os.Create(filepath.Join(path, pubFileName)); nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIOPub, block); nil != err {
		return err
	}
	return nil
}

// GeneratePKCS8Key RSA公钥私钥产生（私钥PKCS8格式）
//
// bits 指定生成位大小
//
// path 指定公私钥所在生成目录
func (r *RSACommon) GeneratePKCS8Key(bits int, path, priFileName, pubFileName string) error {
	var (
		privateKey *rsa.PrivateKey
		publicKey  *rsa.PublicKey
		derStream  []byte
		fileIOPri  *os.File
		fileIOPub  *os.File
		derPkiX    []byte
		err        error
	)
	// 创建公私钥生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	// 生成私钥文件
	if privateKey, err = rsa.GenerateKey(rand.Reader, bits); nil != err {
		return err
	}
	// 将私钥转换为ASN.1 DER编码的形式
	if derStream, err = x509.MarshalPKCS8PrivateKey(privateKey); nil != err {
		return err
	}
	// block表示PEM编码的结构
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	defer func() { _ = fileIOPri.Close() }()
	if fileIOPri, err = os.Create(filepath.Join(path, priFileName)); nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIOPri, block); nil != err {
		return err
	}
	// 生成公钥文件
	publicKey = &privateKey.PublicKey
	// 将公钥序列化为der编码的PKIX格式
	if derPkiX, err = x509.MarshalPKIXPublicKey(publicKey); nil != err {
		return err
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkiX,
	}
	defer func() { _ = fileIOPub.Close() }()
	if fileIOPub, err = os.Create(filepath.Join(path, pubFileName)); nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIOPub, block); nil != err {
		return err
	}
	return nil
}

// GeneratePKCS1PriKey RSA私钥产生（私钥PKCS1格式）
//
// bits 指定生成位大小
//
// path 指定私钥所在生成目录
//
// fileName 指定私钥的文件名称，如'private.pem'
func (r *RSACommon) GeneratePKCS1PriKey(bits int, path, fileName string) error {
	var (
		privateKey *rsa.PrivateKey
		fileIO     *os.File
		err        error
	)
	// 创建公私钥生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	// 生成私钥文件
	if privateKey, err = rsa.GenerateKey(rand.Reader, bits); nil != err {
		return err
	}
	// 将私钥转换为ASN.1 DER编码的形式
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	// block表示PEM编码的结构
	block := &pem.Block{
		Type:  privateRSAKeyPemType,
		Bytes: derStream,
	}
	defer func() { _ = fileIO.Close() }()
	if fileIO, err = os.Create(filepath.Join(path, fileName)); nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, block); nil != err {
		return err
	}
	return nil
}

// GeneratePKCS8PriKey RSA私钥产生（私钥PKCS8格式）
//
// bits 指定生成位大小
//
// path 指定私钥所在生成目录
//
// fileName 指定私钥的文件名称，如'private.pem'
func (r *RSACommon) GeneratePKCS8PriKey(bits int, path, fileName string) error {
	var (
		privateKey *rsa.PrivateKey
		derStream  []byte
		fileIO     *os.File
		err        error
	)
	// 创建公私钥生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	// 生成私钥文件
	if privateKey, err = rsa.GenerateKey(rand.Reader, bits); nil != err {
		return err
	}
	// 将私钥转换为ASN.1 DER编码的形式
	if derStream, err = x509.MarshalPKCS8PrivateKey(privateKey); nil != err {
		return err
	}
	// block表示PEM编码的结构
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	defer func() { _ = fileIO.Close() }()
	if fileIO, err = os.Create(filepath.Join(path, fileName)); nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, block); nil != err {
		return err
	}
	return nil
}

// GeneratePubKey RSA公钥产生
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// path 指定公私钥所在生成目录
//
// fileName 指定公钥的文件名称，如'public.pem'
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().pksC1()’和‘CryptoRSA().pksC8()’方法赋值
func (r *RSACommon) GeneratePubKey(privateKey []byte, path, fileName string, pks PKSCType) error {
	var (
		fileIO  *os.File
		derPkiX []byte
		err     error
	)
	pri, err := r.parsePrivateKey(privateKey, pks)
	if err != nil {
		return err
	}
	// 创建公私钥生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	// 生成公钥文件
	publicKey := &pri.PublicKey
	// 将公钥序列化为der编码的PKIX格式
	if derPkiX, err = x509.MarshalPKIXPublicKey(publicKey); nil != err {
		return err
	}
	block := &pem.Block{
		Type:  publicRSAKeyPemType,
		Bytes: derPkiX,
	}
	defer func() { _ = fileIO.Close() }()
	if fileIO, err = os.Create(filepath.Join(path, fileName)); nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, block); nil != err {
		return err
	}
	return nil
}

// GeneratePubKeyFP RSA公钥产生
//
// privateKeyFilePath 私钥地址
//
// path 指定公私钥所在生成目录
//
// fileName 指定公钥的文件名称，如'public.pem'
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().pksC1()’和‘CryptoRSA().pksC8()’方法赋值
func (r *RSACommon) GeneratePubKeyFP(privateKeyFilePath, path, fileName string, pks PKSCType) error {
	bs, err := ioutil.ReadFile(privateKeyFilePath)
	if nil != err {
		return err
	}
	return r.GeneratePubKey(bs, path, fileName, pks)
}

// Encrypt 公钥加密
//
// publicKey 公钥内容，如取出字符串'pubData'，则传入'string(pubData)'即可
//
// data 待加密数据
func (r *RSACommon) Encrypt(publicKey, data []byte) ([]byte, error) {
	pub, err := r.parsePublicKey(publicKey)
	if nil != err {
		return nil, err
	}
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// EncryptFP 公钥加密
//
// publicKeyFilePath 公钥地址
//
// data 待加密数据
func (r *RSACommon) EncryptFP(publicKeyFilePath string, data []byte) ([]byte, error) {
	bs, err := ioutil.ReadFile(publicKeyFilePath)
	if nil != err {
		return nil, err
	}
	return r.Encrypt(bs, data)
}

// Decrypt 私钥解密
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// data 待解密数据
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().pksC1()’和‘CryptoRSA().pksC8()’方法赋值
func (r *RSACommon) Decrypt(privateKey, data []byte, pks PKSCType) ([]byte, error) {
	pri, err := r.parsePrivateKey(privateKey, pks)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, pri, data)
}

// DecryptFP 私钥解密
//
// privateKeyFilePath 私钥地址
//
// data 待解密数据
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().pksC1()’和‘CryptoRSA().pksC8()’方法赋值
func (r *RSACommon) DecryptFP(privateKeyFilePath string, data []byte, pks PKSCType) ([]byte, error) {
	bs, err := ioutil.ReadFile(privateKeyFilePath)
	if nil != err {
		return bs, err
	}
	return r.Decrypt(bs, data, pks)
}

// Sign 签名：采用sha1算法进行签名并输出为hex格式（私钥PKCS8格式）
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().pksC1()’和‘CryptoRSA().pksC8()’方法赋值
func (r *RSACommon) Sign(privateKey, data []byte, hash crypto.Hash, pks PKSCType) ([]byte, error) {
	pri, err := r.parsePrivateKey(privateKey, pks)
	if err != nil {
		return nil, err
	}
	h := hash.New()
	if _, err = h.Write(data); nil != err {
		return nil, err
	}
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, pri, hash, hashed)
}

// SignFP 签名：采用sha1算法进行签名并输出为hex格式（私钥PKCS8格式）
//
// privateKeyPath 私钥文件存储路径
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().pksC1()’和‘CryptoRSA().pksC8()’方法赋值
func (r *RSACommon) SignFP(privateKeyPath string, data []byte, hash crypto.Hash, pks PKSCType) ([]byte, error) {
	bs, err := ioutil.ReadFile(privateKeyPath)
	if nil != err {
		return nil, err
	}
	return r.Sign(bs, data, hash, pks)
}

// Verify 验签：对采用sha1算法进行签名后转base64格式的数据进行验签
//
// publicKey 公钥内容，如取出字符串'pubData'，则传入'string(pubData)'即可
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
func (r *RSACommon) Verify(publicKey, data, signData []byte, hash crypto.Hash) error {
	pub, err := r.parsePublicKey(publicKey)
	if nil != err {
		return err
	}
	h := hash.New()
	if _, err = h.Write(data); nil != err {
		return err
	}
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(pub, hash, hashed, signData)
}

// VerifyFP 验签：对采用sha1算法进行签名后转base64格式的数据进行验签
//
// publicKeyPath 公钥文件存储路径
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
func (r *RSACommon) VerifyFP(publicKeyPath string, data, signData []byte, hash crypto.Hash) error {
	bs, err := ioutil.ReadFile(publicKeyPath)
	if nil != err {
		return err
	}
	return r.Verify(bs, data, signData, hash)
}

// parsePrivateKey 解析私钥
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().pksC1()’和‘CryptoRSA().pksC8()’方法赋值
func (r *RSACommon) parsePrivateKey(key []byte, pks PKSCType) (*rsa.PrivateKey, error) {
	pemData, err := r.pemParse(key, privateRSAKeyPemType)
	if err != nil {
		return nil, err
	}
	switch pks {
	default:
		pri, err := x509.ParsePKCS8PrivateKey(pemData)
		if nil != err {
			return nil, err
		}
		return pri.(*rsa.PrivateKey), nil
	case pksC1:
		return x509.ParsePKCS1PrivateKey(pemData)
	}
}

func (r *RSACommon) parsePublicKey(key []byte) (*rsa.PublicKey, error) {
	pemData, err := r.pemParse(key, publicRSAKeyPemType)
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

func (r *RSACommon) pksC1() PKSCType {
	return pksC1
}

func (r *RSACommon) pksC8() PKSCType {
	return pksC8
}

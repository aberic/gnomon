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
	pksC1        PKSCType = "PKCS1"
	pksC8        PKSCType = "PKCS8"
	signPss      SignMode = "pss"
	signPKCS1v15 SignMode = "pkcs#1 v1.5"
	//publicRSAKeyPemType          = "PUBLIC KEY"
	//privateRSAKeyPemType          = "PRIVATE KEY"
)

// PKSCType 私钥格式，默认提供PKCS1和PKCS8
//
// 通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
type PKSCType string

// SignMode RSA签名模式，默认提供PSS和PKCS1v15
//
// 通过调用‘CryptoRSA().SignPSS()’和‘CryptoRSA().SignPKCS()’方法赋值
type SignMode string

// RSAGenerateKey RSA公钥私钥产生
//
// bits 指定生成位大小
//
// path 指定公私钥所在生成目录
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSAGenerateKey(bits int, path, priFileName, pubFileName, priPemType, pubPemType string, pks PKSCType) error {
	return RSAGenerateKeyWithPass(bits, path, priFileName, pubFileName, "", priPemType, pubPemType, -1, pks)
}

// RSAGenerateKeyWithPass RSA公钥私钥产生
//
// bits 指定生成位大小
//
// path 指定公私钥所在生成目录
//
// passwd 生成密码
//
// alg der编码数据指定算法，如：x509.PEMCipher3DES
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSAGenerateKeyWithPass(bits int, path, priFileName, pubFileName, passwd, priPemType, pubPemType string, alg x509.PEMCipher, pks PKSCType) error {
	privateKey, err := RSAGeneratePriKeyWithPass(bits, path, priFileName, passwd, priPemType, alg, pks)
	if nil != err {
		return err
	}
	return RSAGeneratePubKey(privateKey, path, pubFileName, pubPemType)
}

// RSAGeneratePriKey RSA私钥产生
//
// bits 指定生成位大小
//
// path 指定私钥所在生成目录
//
// fileName 指定私钥的文件名称，如'private.pem'
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSAGeneratePriKey(bits int, path, fileName, pemType string, pks PKSCType) (*rsa.PrivateKey, error) {
	return RSAGeneratePriKeyWithPass(bits, path, fileName, "", pemType, -1, pks)
}

// RSAGeneratePriKeyWithPass RSA私钥产生
//
// bits 指定生成位大小
//
// path 指定私钥所在生成目录
//
// fileName 指定私钥的文件名称，如'private.pem'
//
// passwd 生成密码
//
// alg der编码数据指定算法，如：x509.PEMCipher3DES
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSAGeneratePriKeyWithPass(bits int, path, fileName, passwd, pemType string, alg x509.PEMCipher, pks PKSCType) (*rsa.PrivateKey, error) {
	var (
		privateKey *rsa.PrivateKey
		err        error
	)
	// 创建公私钥生成目录
	if !FilePathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return nil, err
		}
	}
	// 生成私钥文件
	if privateKey, err = rsa.GenerateKey(rand.Reader, bits); nil != err {
		return nil, err
	}
	if err = RSASavePriPemWithPass(privateKey, path, fileName, passwd, pemType, alg, pks); nil != err {
		return nil, err
	}
	return privateKey, nil
}

// RSAGeneratePubKey RSA公钥产生
//
// privateKey 私钥
//
// path 指定公私钥所在生成目录
//
// fileName 指定公钥的文件名称，如'public.pem'
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSAGeneratePubKey(privateKey *rsa.PrivateKey, path, fileName, pemType string) error {
	var (
		fileIO  *os.File
		derPkiX []byte
		err     error
	)
	// 创建公私钥生成目录
	if !FilePathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	// 生成公钥文件
	publicKey := &privateKey.PublicKey
	// 将公钥序列化为der编码的PKIX格式
	if derPkiX, err = x509.MarshalPKIXPublicKey(publicKey); nil != err {
		return err
	}
	block := &pem.Block{
		Type:  pemType,
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

// RSAGeneratePubKeyBytes RSA公钥产生
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// path 指定公私钥所在生成目录
//
// fileName 指定公钥的文件名称，如'public.pem'
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSAGeneratePubKeyBytes(privateKey []byte, path, fileName, priPemType, pubPemType string, pks PKSCType) error {
	return RSAGeneratePubKeyBytesWithPass(privateKey, "", path, fileName, priPemType, pubPemType, pks)
}

// RSAGeneratePubKeyBytesWithPass RSA公钥产生
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// passwd 生成privateKey时输入密码
//
// path 指定公私钥所在生成目录
//
// fileName 指定公钥的文件名称，如'public.pem'
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSAGeneratePubKeyBytesWithPass(privateKey []byte, passwd, priPemType, pubPemType, path, fileName string, pks PKSCType) error {
	pri, err := RSALoadPriWithPass(privateKey, passwd, priPemType, pks)
	if err != nil {
		return err
	}
	return RSAGeneratePubKey(pri, path, fileName, pubPemType)
}

// RSAGeneratePubKeyFP RSA公钥产生
//
// privateKeyFilePath 私钥地址
//
// path 指定公私钥所在生成目录
//
// fileName 指定公钥的文件名称，如'public.pem'
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSAGeneratePubKeyFP(privateKeyFilePath, path, fileName, priPemType, pubPemType string, pks PKSCType) error {
	bs, err := ioutil.ReadFile(privateKeyFilePath)
	if nil != err {
		return err
	}
	return RSAGeneratePubKeyBytes(bs, path, fileName, priPemType, pubPemType, pks)
}

// RSAGeneratePubKeyFPWithPass RSA公钥产生
//
// privateKeyFilePath 私钥地址
//
// passwd 生成privateKey时输入密码
//
// path 指定公私钥所在生成目录
//
// fileName 指定公钥的文件名称，如'public.pem'
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSAGeneratePubKeyFPWithPass(privateKeyFilePath, passwd, path, fileName, priPemType, pubPemType string, pks PKSCType) error {
	bs, err := ioutil.ReadFile(privateKeyFilePath)
	if nil != err {
		return err
	}
	return RSAGeneratePubKeyBytesWithPass(bs, passwd, path, fileName, priPemType, pubPemType, pks)
}

// RSASavePriPem 将私钥保存到给定文件
//
// alg der编码数据指定算法，如：x509.PEMCipher3DES
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSASavePriPem(privateKey *rsa.PrivateKey, path, fileName, pemType string, alg x509.PEMCipher, pks PKSCType) error {
	return RSASavePriPemWithPass(privateKey, path, fileName, "", pemType, alg, pks)
}

// RSASavePriPemWithPass 将私钥保存到给定文件
//
// alg der编码数据指定算法，如：x509.PEMCipher3DES
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSASavePriPemWithPass(privateKey *rsa.PrivateKey, path, fileName, passwd, pemType string, alg x509.PEMCipher, pks PKSCType) error {
	var (
		derStream []byte
		fileIO    *os.File
		block     *pem.Block
		err       error
	)
	// 将私钥转换为ASN.1 DER编码的形式
	switch pks {
	default:
		derStream = x509.MarshalPKCS1PrivateKey(privateKey)
	case pksC8:
		if derStream, err = x509.MarshalPKCS8PrivateKey(privateKey); nil != err {
			return err
		}
	}
	// block表示PEM编码的结构
	if StringIsEmpty(passwd) {
		block = &pem.Block{Type: pemType, Bytes: derStream}
	} else {
		block, err = x509.EncryptPEMBlock(rand.Reader, pemType, derStream, []byte(passwd), alg)
		if nil != err {
			return err
		}
	}
	defer func() { _ = fileIO.Close() }()
	if !FilePathExists(path) {
		if err = os.MkdirAll(path, 0755); nil != err {
			return err
		}
	}
	if fileIO, err = os.Create(filepath.Join(path, fileName)); nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, block); nil != err {
		return err
	}
	return nil
}

// RSAEncrypt 公钥加密
//
// publicKey 公钥内容，如取出字符串'pubData'，则传入'string(pubData)'即可
//
// data 待加密数据
func RSAEncrypt(publicKey, data []byte, pemType string) ([]byte, error) {
	pub, err := RSALoadPub(publicKey, pemType)
	if nil != err {
		return nil, err
	}
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// RSAEncryptFP 公钥加密
//
// publicKeyFilePath 公钥地址
//
// data 待加密数据
func RSAEncryptFP(publicKeyFilePath, pemType string, data []byte) ([]byte, error) {
	pub, err := RSALoadPubFP(publicKeyFilePath, pemType)
	if nil != err {
		return nil, err
	}
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// RSADecrypt 私钥解密
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// data 待解密数据
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSADecrypt(privateKey, data []byte, pemType string, pks PKSCType) ([]byte, error) {
	return RSADecryptWithPass(privateKey, data, "", pemType, pks)
}

// RSADecryptWithPass 私钥解密
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// data 待解密数据
//
// passwd 生成privateKey时输入密码
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSADecryptWithPass(privateKey, data []byte, passwd, pemType string, pks PKSCType) ([]byte, error) {
	pri, err := RSALoadPriWithPass(privateKey, passwd, pemType, pks)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, pri, data)
}

// RSADecryptFP 私钥解密
//
// privateKeyFilePath 私钥地址
//
// data 待解密数据
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSADecryptFP(privateKeyFilePath, pemType string, data []byte, pks PKSCType) ([]byte, error) {
	return RSADecryptFPWithPass(privateKeyFilePath, "", pemType, data, pks)
}

// RSADecryptFPWithPass 私钥解密
//
// privateKeyFilePath 私钥地址
//
// passwd 生成privateKey时输入密码
//
// data 待解密数据
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSADecryptFPWithPass(privateKeyFilePath, passwd, pemType string, data []byte, pks PKSCType) ([]byte, error) {
	pri, err := RSALoadPriFPWithPass(privateKeyFilePath, passwd, pemType, pks)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, pri, data)
}

// RSASign 签名：采用RSA-PKCS#1 v1.5模式
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
//
// mod RSA签名模式，默认提供PSS和PKCS1v15，通过调用‘CryptoRSA().SignPSS()’和‘CryptoRSA().SignPKCS()’方法赋值
func RSASign(privateKey, data []byte, pemType string, hash crypto.Hash, pks PKSCType, mod SignMode) ([]byte, error) {
	return RSASignWithPass(privateKey, data, "", pemType, hash, pks, mod)
}

// RSASignWithPass 签名：采用RSA-PKCS#1 v1.5模式
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// data 待签名数据
//
// passwd 生成privateKey时输入密码
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
//
// mod RSA签名模式，默认提供PSS和PKCS1v15，通过调用‘CryptoRSA().SignPSS()’和‘CryptoRSA().SignPKCS()’方法赋值
func RSASignWithPass(privateKey, data []byte, passwd, pemType string, hash crypto.Hash, pks PKSCType, mod SignMode) ([]byte, error) {
	pri, err := RSALoadPriWithPass(privateKey, passwd, pemType, pks)
	if err != nil {
		return nil, err
	}
	switch mod {
	default:
		return rsaSignPSS(pri, data, hash)
	case signPKCS1v15:
		return rsaSignPKCS1v15(pri, data, hash)
	}
}

// RSASignFP 签名：采用RSA-PKCS#1 v1.5模式
//
// privateKeyPath 私钥文件存储路径
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
//
// mod RSA签名模式，默认提供PSS和PKCS1v15，通过调用‘CryptoRSA().SignPSS()’和‘CryptoRSA().SignPKCS()’方法赋值
func RSASignFP(privateKeyPath, pemType string, data []byte, hash crypto.Hash, pks PKSCType, mod SignMode) ([]byte, error) {
	pri, err := RSALoadPriFP(privateKeyPath, pemType, pks)
	if err != nil {
		return nil, err
	}
	switch mod {
	default:
		return rsaSignPSS(pri, data, hash)
	case signPKCS1v15:
		return rsaSignPKCS1v15(pri, data, hash)
	}
}

// RSASignFPWithPass 签名：采用RSA-PKCS#1 v1.5模式
//
// privateKeyPath 私钥文件存储路径
//
// passwd 生成privateKey时输入密码
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
//
// mod RSA签名模式，默认提供PSS和PKCS1v15，通过调用‘CryptoRSA().SignPSS()’和‘CryptoRSA().SignPKCS()’方法赋值
func RSASignFPWithPass(privateKeyPath, passwd, pemType string, data []byte, hash crypto.Hash, pks PKSCType, mod SignMode) ([]byte, error) {
	pri, err := RSALoadPriFPWithPass(privateKeyPath, passwd, pemType, pks)
	if err != nil {
		return nil, err
	}
	switch mod {
	default:
		return rsaSignPSS(pri, data, hash)
	case signPKCS1v15:
		return rsaSignPKCS1v15(pri, data, hash)
	}
}

// rsaSignPSS 签名：采用PSS模式
//
// privateKey 私钥
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
func rsaSignPSS(privateKey *rsa.PrivateKey, data []byte, hash crypto.Hash) ([]byte, error) {
	h := hash.New()
	if _, err := h.Write(data); nil != err {
		return nil, err
	}
	hashed := h.Sum(nil)
	return rsa.SignPSS(rand.Reader, privateKey, hash, hashed, nil)
}

// rsaSignPKCS1v15 签名：采用RSA-PKCS#1 v1.5模式
//
// privateKey 私钥
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
func rsaSignPKCS1v15(privateKey *rsa.PrivateKey, data []byte, hash crypto.Hash) ([]byte, error) {
	h := hash.New()
	if _, err := h.Write(data); nil != err {
		return nil, err
	}
	hashed := h.Sum(nil)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, hash, hashed)
}

// RSAVerify 验签：采用RSA-PKCS#1 v1.5模式
//
// publicKey 公钥内容，如取出字符串'pubData'，则传入'string(pubData)'即可
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
//
// mod RSA签名模式，默认提供PSS和PKCS1v15，通过调用‘CryptoRSA().SignPSS()’和‘CryptoRSA().SignPKCS()’方法赋值
func RSAVerify(publicKey, data, signData []byte, pemType string, hash crypto.Hash, mod SignMode) error {
	pub, err := RSALoadPub(publicKey, pemType)
	if nil != err {
		return err
	}
	switch mod {
	default:
		return rsaVerifyPSS(pub, data, signData, hash)
	case signPKCS1v15:
		return rsaVerifyPKCS1v15(pub, data, signData, hash)
	}
}

// RSAVerifyFP 验签：采用RSA-PKCS#1 v1.5模式
//
// publicKeyPath 公钥文件存储路径
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
//
// mod RSA签名模式，默认提供PSS和PKCS1v15，通过调用‘CryptoRSA().SignPSS()’和‘CryptoRSA().SignPKCS()’方法赋值
func RSAVerifyFP(publicKeyPath, pemType string, data, signData []byte, hash crypto.Hash, mod SignMode) error {
	pub, err := RSALoadPubFP(publicKeyPath, pemType)
	if nil != err {
		return err
	}
	switch mod {
	default:
		return rsaVerifyPSS(pub, data, signData, hash)
	case signPKCS1v15:
		return rsaVerifyPKCS1v15(pub, data, signData, hash)
	}
}

// rsaVerifyPSS 验签：采用PSS模式
//
// publicKey 公钥
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
func rsaVerifyPSS(publicKey *rsa.PublicKey, data, signData []byte, hash crypto.Hash) error {
	h := hash.New()
	if _, err := h.Write(data); nil != err {
		return err
	}
	hashed := h.Sum(nil)
	opts := &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthAuto, Hash: hash}
	return rsa.VerifyPSS(publicKey, hash, hashed, signData, opts)
}

// rsaVerifyPKCS1v15 验签：采用RSA-PKCS#1 v1.5模式
//
// publicKey 公钥
//
// data 待签名数据
//
// hash 算法，如 crypto.SHA1/crypto.SHA256等
func rsaVerifyPKCS1v15(publicKey *rsa.PublicKey, data, signData []byte, hash crypto.Hash) error {
	h := hash.New()
	if _, err := h.Write(data); nil != err {
		return err
	}
	hashed := h.Sum(nil)
	return rsa.VerifyPKCS1v15(publicKey, hash, hashed, signData)
}

// RSALoadPriFP 加载私钥
//
// privateKeyFilePath 私钥地址
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSALoadPriFP(privateKeyFilePath, pemType string, pks PKSCType) (*rsa.PrivateKey, error) {
	return RSALoadPriFPWithPass(privateKeyFilePath, "", pemType, pks)
}

// RSALoadPriFPWithPass 加载私钥
//
// privateKeyFilePath 私钥地址
//
// passwd 生成privateKey时输入密码
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSALoadPriFPWithPass(privateKeyFilePath, passwd, pemType string, pks PKSCType) (*rsa.PrivateKey, error) {
	bs, err := ioutil.ReadFile(privateKeyFilePath)
	if nil != err {
		return nil, err
	}
	return RSALoadPriWithPass(bs, passwd, pemType, pks)
}

// RSALoadPubFP 加载公钥
//
// publicKeyFilePath 公钥地址
func RSALoadPubFP(publicKeyFilePath, pemType string) (*rsa.PublicKey, error) {
	bs, err := ioutil.ReadFile(publicKeyFilePath)
	if nil != err {
		return nil, err
	}
	return RSALoadPub(bs, pemType)
}

// RSALoadPri 解析私钥
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSALoadPri(privateKey []byte, pemType string, pks PKSCType) (*rsa.PrivateKey, error) {
	return RSALoadPriWithPass(privateKey, "", pemType, pks)
}

// RSALoadPriWithPass 解析私钥
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// passwd 生成privateKey时输入密码
//
// pks 私钥格式，默认提供PKCS1和PKCS8，通过调用‘CryptoRSA().PKSC1()’和‘CryptoRSA().PKSC8()’方法赋值
func RSALoadPriWithPass(privateKey []byte, passwd, pemType string, pks PKSCType) (*rsa.PrivateKey, error) {
	var (
		pemData []byte
		err     error
	)
	if StringIsEmpty(passwd) {
		pemData, err = rsaPemParse(privateKey, pemType)
		if err != nil {
			return nil, err
		}
	} else {
		block, _ := pem.Decode(privateKey)
		pemData, err = x509.DecryptPEMBlock(block, []byte(passwd))
		if err != nil {
			return nil, err
		}
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

// RSALoadPub 加载公钥
//
// publicKey 公钥内容，如取出字符串'pubData'，则传入'string(pubData)'即可
func RSALoadPub(publicKey []byte, pemType string) (*rsa.PublicKey, error) {
	pemData, err := rsaPemParse(publicKey, pemType)
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

// rsaPemParse 解密pem格式密钥并验证pem类型
func rsaPemParse(key []byte, pemType string) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("no pem block found")
	}
	if pemType != "" && block.Type != pemType {
		return nil, errors.New(strings.Join([]string{"Key's type is ", block.Type, ", expected ", pemType}, ""))
	}
	return block.Bytes, nil
}

// RSAPKSC1 私钥格PKCS1
func RSAPKSC1() PKSCType {
	return pksC1
}

// RSAPKSC8 私钥格式PKCS8
func RSAPKSC8() PKSCType {
	return pksC8
}

// RSASignPSS RSA签名模式，PSS
func RSASignPSS() SignMode {
	return signPss
}

// RSASignPKCS RSA签名模式，PKCS1v15
func RSASignPKCS() SignMode {
	return signPKCS1v15
}

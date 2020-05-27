/*
 *  Copyright (c) 2019. aberic - All Rights Reserved.
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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

// ECC，全称椭圆曲线密码学（英语：Elliptic curve cryptography，缩写为 ECC），主要是指相关数学原理
//
// ECIES，在ECC原理的基础上实现的一种公钥加密方法，和RSA类似
//
// ECDSA，在ECC原理上实现的签名方法
//
// ECDH在ECC和DH的基础上实现的密钥交换算法

//const (
//	privateECCKeyPemType = "PRIVATE KEY"
//	publicECCKeyPemType  = "PUBLIC KEY"
//)

// ECCGenerate 生成公私钥对
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func ECCGenerate(curve elliptic.Curve) (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	var (
		privateKey *ecdsa.PrivateKey
		err        error
	)
	privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privateKey, &privateKey.PublicKey, nil
}

// ECCGeneratePemBytes 生成公私钥对
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func ECCGeneratePemBytes(priPemType, pubPemType, passwd string, curve elliptic.Curve) (priBytes, pubBytes []byte, err error) {
	var privateKey *ecdsa.PrivateKey
	if privateKey, err = ecdsa.GenerateKey(curve, rand.Reader); nil == err {
		publicKey := &privateKey.PublicKey
		if priBytes, err = ECCPri2PemBytes(priPemType, passwd, privateKey); nil != err {
			return
		}
		if pubBytes, err = ECCPub2PemBytes(pubPemType, publicKey); nil != err {
			return
		}
	}
	return
}

// ECCGenerateKey 生成公私钥对
//
// path 指定公私钥所在生成目录
//
// priFileName 指定生成的密钥名称
//
// pubFileName 指定生成的密钥名称
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func ECCGenerateKey(path, priFileName, pubFileName string, curve elliptic.Curve) error {
	var (
		privateKey *ecdsa.PrivateKey
		err        error
	)
	// 创建公私钥生成目录
	if !FilePathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return err
	}

	if err = ECCSavePri(filepath.Join(path, priFileName), privateKey); nil != err {
		return err
	}
	if err = ECCSavePub(filepath.Join(path, pubFileName), &privateKey.PublicKey, curve); nil != err {
		return err
	}
	return nil
}

// ECCGeneratePemKey 生成公私钥对
//
// path 指定公私钥所在生成目录
//
// priFileName 指定生成的密钥名称
//
// pubFileName 指定生成的密钥名称
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func ECCGeneratePemKey(path, priFileName, pubFileName, priPemType, pubPemType string, curve elliptic.Curve) error {
	return ECCGeneratePemKeyWithPass(path, priFileName, pubFileName, "", priPemType, pubPemType, curve)
}

// ECCGeneratePemKeyWithPass 生成公私钥对
//
// path 指定公私钥所在生成目录
//
// priFileName 指定生成的密钥名称
//
// pubFileName 指定生成的密钥名称
//
// passwd 生成密码
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func ECCGeneratePemKeyWithPass(path, priFileName, pubFileName, passwd, priPemType, pubPemType string, curve elliptic.Curve) error {
	var (
		privateKey *ecdsa.PrivateKey
		err        error
	)
	// 创建公私钥生成目录
	if !FilePathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return err
	}

	if err = ECCSavePriPemWithPass(privateKey, passwd, path, priFileName, priPemType); nil != err {
		return err
	}
	if err = ECCSavePubPem(filepath.Join(path, pubFileName), pubPemType, &privateKey.PublicKey); nil != err {
		return err
	}
	return nil
}

// ECCGeneratePriKey 生成私钥
//
// path 指定私钥所在生成目录
//
// priFileName 指定生成的密钥名称
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func ECCGeneratePriKey(path, priFileName string, curve elliptic.Curve) error {
	var (
		privateKey *ecdsa.PrivateKey
		err        error
	)
	// 创建公私钥生成目录
	if !FilePathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return err
	}
	if err = ECCSavePri(filepath.Join(path, priFileName), privateKey); nil != err {
		return err
	}
	return nil
}

// ECCGeneratePemPriKey 生成私钥
//
// path 指定私钥所在生成目录
//
// priFileName 指定生成的密钥名称
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func ECCGeneratePemPriKey(path, priFileName, pemType string, curve elliptic.Curve) error {
	return ECCGeneratePemPriKeyWithPass(path, priFileName, "", pemType, curve)
}

// ECCGeneratePemPriKeyWithPass 生成私钥
//
// path 指定私钥所在生成目录
//
// priFileName 指定生成的密钥名称
//
// passwd 生成时输入的密码
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func ECCGeneratePemPriKeyWithPass(path, priFileName, passwd, pemType string, curve elliptic.Curve) error {
	var (
		privateKey *ecdsa.PrivateKey
		err        error
	)
	// 创建公私钥生成目录
	if !FilePathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return err
		}
	}
	privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return err
	}
	if err = ECCSavePriPemWithPass(privateKey, passwd, path, priFileName, pemType); nil != err {
		return err
	}
	return nil
}

// ECCGeneratePubKey 生成公钥
//
// path 指定公钥所在生成目录
func ECCGeneratePubKey(privateKey *ecdsa.PrivateKey, path, pubFileName string, curve elliptic.Curve) error {
	if err := ECCSavePub(filepath.Join(path, pubFileName), &privateKey.PublicKey, curve); nil != err {
		return err
	}
	return nil
}

// ECCGeneratePemPubKey 生成公钥
//
// path 指定公钥所在生成目录
func ECCGeneratePemPubKey(privateKey *ecdsa.PrivateKey, path, pubFileName, pemType string) error {
	if err := ECCSavePubPem(filepath.Join(path, pubFileName), pemType, &privateKey.PublicKey); nil != err {
		return err
	}
	return nil
}

// ECCSavePri 将私钥保存到给定文件，密钥数据保存为hex编码
func ECCSavePri(file string, privateKey *ecdsa.PrivateKey) error {
	k := hex.EncodeToString(ECCPriKey2Bytes(privateKey))
	return ioutil.WriteFile(file, []byte(k), 0600)
	//return crypto.SaveECDSA(file, privateKey)
}

// ECCLoadPri 从文件中加载私钥
//
// file 文件路径
func ECCLoadPri(file string, curve elliptic.Curve) (*ecdsa.PrivateKey, error) {
	bs, err := ioutil.ReadFile(file)
	if nil != err {
		return nil, err
	}
	data, err := hex.DecodeString(string(bs))
	if err != nil {
		return nil, err
	}
	return ECCBytes2PriKey(data, curve), nil
	//return crypto.LoadECDSA(file)
}

// ECCSavePriPem 将私钥保存到给定文件
func ECCSavePriPem(privateKey *ecdsa.PrivateKey, path, fileName, pemType string) error {
	return ECCSavePriPemWithPass(privateKey, "", path, fileName, pemType)
}

// ECCSavePriPemWithPass 将私钥保存到给定文件
func ECCSavePriPemWithPass(privateKey *ecdsa.PrivateKey, passwd, path, fileName, pemType string) error {
	var (
		fileIO *os.File
		block  *pem.Block
	)
	// 将私钥转换为ASN.1 DER编码的形式
	derStream, err := x509.MarshalECPrivateKey(privateKey)
	if nil != err {
		return err
	}
	// block表示PEM编码的结构
	if StringIsEmpty(passwd) {
		block = &pem.Block{Type: pemType, Bytes: derStream}
	} else {
		block, err = x509.EncryptPEMBlock(rand.Reader, pemType, derStream, []byte(passwd), x509.PEMCipher3DES)
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
	if fileIO, err = os.OpenFile(filepath.Join(path, fileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, block); nil != err {
		return err
	}
	return nil
}

// ECCLoadPriPem 解析私钥
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
func ECCLoadPriPem(privateKey []byte, pemType string) (*ecdsa.PrivateKey, error) {
	return ECCLoadPriPemWithPass(privateKey, "", pemType)
}

// ECCLoadPriPemWithPass 解析私钥
//
// privateKey 私钥内容，如取出字符串'priData'，则传入'string(priData)'即可
//
// passwd 生成privateKey时输入密码
func ECCLoadPriPemWithPass(privateKey []byte, passwd, pemType string) (*ecdsa.PrivateKey, error) {
	var (
		pemData []byte
		err     error
	)
	if StringIsEmpty(passwd) {
		pemData, err = eccPemParse(privateKey, pemType)
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
	pri, err := x509.ParseECPrivateKey(pemData)
	if nil != err {
		return nil, err
	}
	return pri, nil
}

// ECCLoadPriPemFP 从文件中加载私钥
//
// file 文件路径
func ECCLoadPriPemFP(file, pemType string) (*ecdsa.PrivateKey, error) {
	return ECCLoadPriPemFPWithPass(file, "", pemType)
}

// ECCLoadPriPemFPWithPass 从文件中加载私钥
//
// file 文件路径
//
// passwd 生成privateKey时输入密码
func ECCLoadPriPemFPWithPass(file, passwd, pemType string) (*ecdsa.PrivateKey, error) {
	var (
		pemData []byte
		err     error
	)
	keyData, err := ioutil.ReadFile(file)
	if nil != err {
		return nil, err
	}
	if StringIsEmpty(passwd) {
		pemData, err = eccPemParse(keyData, pemType)
		if err != nil {
			return nil, err
		}
	} else {
		block, _ := pem.Decode(keyData)
		pemData, err = x509.DecryptPEMBlock(block, []byte(passwd))
		if err != nil {
			return nil, err
		}
	}
	pri, err := x509.ParseECPrivateKey(pemData)
	if nil != err {
		return nil, err
	}
	return pri, nil
}

// ECCSavePubPem 将公钥保存到给定文件
//
// file 文件路径
func ECCSavePubPem(file, pemType string, publicKey *ecdsa.PublicKey) error {
	var fileIO *os.File
	// 将公钥序列化为der编码的PKIX格式
	derPkiX, err := x509.MarshalPKIXPublicKey(publicKey)
	if nil != err {
		return err
	}
	block := &pem.Block{
		Type:  pemType,
		Bytes: derPkiX,
	}
	defer func() { _ = fileIO.Close() }()
	if fileIO, err = os.Create(file); nil != err {
		return err
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, block); nil != err {
		return err
	}
	return nil
}

// ECCPri2PemBytes ECCPri2PemBytes
func ECCPri2PemBytes(priPemType, passwd string, privateKey *ecdsa.PrivateKey) (data []byte, err error) {
	var (
		derStream []byte
		block     *pem.Block
	)
	// 将私钥转换为ASN.1 DER编码的形式
	if derStream, err = x509.MarshalECPrivateKey(privateKey); nil == err {
		// block表示PEM编码的结构
		if StringIsEmpty(passwd) {
			block = &pem.Block{Type: priPemType, Bytes: derStream}
		} else {
			if block, err = x509.EncryptPEMBlock(rand.Reader, priPemType, derStream, []byte(passwd), x509.PEMCipher3DES); nil != err {
				return
			}
		}
		data = pem.EncodeToMemory(block)
	}
	return
}

// ECCPub2PemBytes ECCPub2PemBytes
func ECCPub2PemBytes(pubPemType string, publicKey *ecdsa.PublicKey) (data []byte, err error) {
	var (
		derPkiX []byte
		block   *pem.Block
	)
	// 将公钥序列化为der编码的PKIX格式
	if derPkiX, err = x509.MarshalPKIXPublicKey(publicKey); nil == err {
		// block表示PEM编码的结构
		block = &pem.Block{
			Type:  pubPemType,
			Bytes: derPkiX,
		}
		data = pem.EncodeToMemory(block)
	}
	return
}

// ECCLoadPubPem 从文件中加载公钥
//
// file 文件路径
func ECCLoadPubPem(publicKey []byte, pemType string) (*ecdsa.PublicKey, error) {
	pemData, err := eccPemParse(publicKey, pemType)
	if err != nil {
		return nil, err
	}
	keyInterface, err := x509.ParsePKIXPublicKey(pemData)
	if err != nil {
		return nil, err
	}
	pubKey, ok := keyInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("could not cast parsed key to *rsa.PublicKey")
	}
	return pubKey, nil
}

// ECCLoadPubPemFP 从文件中加载公钥
//
// file 文件路径
func ECCLoadPubPemFP(file, pemType string) (*ecdsa.PublicKey, error) {
	pubData, err := ioutil.ReadFile(file)
	if nil != err {
		return nil, err
	}
	pemData, err := eccPemParse(pubData, pemType)
	if err != nil {
		return nil, err
	}
	keyInterface, err := x509.ParsePKIXPublicKey(pemData)
	if err != nil {
		return nil, err
	}
	pubKey, ok := keyInterface.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("could not cast parsed key to *rsa.PublicKey")
	}
	return pubKey, nil
}

// ECCSavePub 将公钥保存到给定文件，密钥数据保存为hex编码
//
// file 文件路径
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func ECCSavePub(file string, publicKey *ecdsa.PublicKey, curve elliptic.Curve) error {
	k := hex.EncodeToString(ECCPubKey2Bytes(publicKey, curve))
	return ioutil.WriteFile(file, []byte(k), 0600)
}

// ECCLoadPub 从文件中加载公钥
//
// file 文件路径
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func ECCLoadPub(file string, curve elliptic.Curve) (*ecdsa.PublicKey, error) {
	data, err := ioutil.ReadFile(file)
	if nil != err {
		return nil, err
	}
	key, err := hex.DecodeString(string(data))
	if err != nil {
		return nil, err
	}
	return ECCBytes2PubKey(key, curve)
}

// ECCPriKey2Bytes 私钥转[]byte
func ECCPriKey2Bytes(privateKey *ecdsa.PrivateKey) []byte {
	return privateKey.D.Bytes()
}

// ECCBytes2PriKey []byte转私钥
func ECCBytes2PriKey(data []byte, curve elliptic.Curve) *ecdsa.PrivateKey {
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = new(big.Int).SetBytes(data)
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(data)
	return priv
}

// ECCPubKey2Bytes 公钥转[]byte
//
// pub 公钥
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func ECCPubKey2Bytes(publicKey *ecdsa.PublicKey, curve elliptic.Curve) []byte {
	if publicKey == nil || publicKey.X == nil || publicKey.Y == nil {
		return nil
	}
	return elliptic.Marshal(curve, publicKey.X, publicKey.Y)
}

// ECCBytes2PubKey []byte转公钥
func ECCBytes2PubKey(data []byte, curve elliptic.Curve) (*ecdsa.PublicKey, error) {
	x, y := elliptic.Unmarshal(curve, data)
	if x == nil {
		return nil, errors.New("invalid public key")
	}
	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, nil
}

// ECCEncrypt 加密
func ECCEncrypt(data []byte, publicKey *ecies.PublicKey) ([]byte, error) {
	if publicKey.Curve.Params().BitSize != 256 {
		return nil, errors.New("just support P256 and S256")
	}
	return ecies.Encrypt(rand.Reader, publicKey, data, nil, nil)
}

//// EncryptFP 加密
////
//// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
//func ECCEncryptFP(data []byte, publicKeyFilePath string, curve elliptic.Curve) ([]byte, error) {
//	publicKey, err := ECCLoadPub(publicKeyFilePath, curve)
//	if nil != err {
//		return nil, err
//	}
//	return ECCEncrypt(data, ecies.ImportECDSAPublic(publicKey))
//}

// ECCDecrypt 解密
func ECCDecrypt(data []byte, privateKey *ecies.PrivateKey) ([]byte, error) {
	if privateKey.Curve.Params().BitSize != 256 {
		return nil, errors.New("just support P256 and S256")
	}
	return privateKey.Decrypt(data, nil, nil)
}

//// DecryptFP 解密
//func ECCDecryptFP(data []byte, privateKeyFilePath string) ([]byte, error) {
//	privateKey, err := ECCLoadPri(privateKeyFilePath)
//	if nil != err {
//		return nil, err
//	}
//	return ECCDecrypt(data, ecies.ImportECDSA(privateKey))
//}

// ECCSign 签名
func ECCSign(privateKey *ecdsa.PrivateKey, data []byte) (sign []byte, err error) {
	// 根据明文plaintext和私钥，生成两个big.Ing
	r, s, err := ecdsa.Sign(rand.Reader, privateKey, data)
	if err != nil {
		return nil, err
	}
	rs, err := r.MarshalText()
	if err != nil {
		return nil, err
	}
	ss, err := s.MarshalText()
	if err != nil {
		return nil, err
	}
	// 将r，s合并（以“+”分割），作为签名返回
	var b bytes.Buffer
	b.Write(rs)
	b.Write([]byte(`+`))
	b.Write(ss)
	return b.Bytes(), nil
}

// ECCVerify 验签
func ECCVerify(publicKey *ecdsa.PublicKey, data, sign []byte) (bool, error) {
	var rint, sint big.Int
	// 根据sign，解析出r，s
	rs := bytes.Split(sign, []byte("+"))
	if err := rint.UnmarshalText(rs[0]); nil != err {
		return false, err
	}
	if err := sint.UnmarshalText(rs[1]); nil != err {
		return false, err
	}
	// 根据公钥，明文，r，s验证签名
	v := ecdsa.Verify(publicKey, data, &rint, &sint)
	return v, nil
}

// eccPemParse 解密pem格式密钥并验证pem类型
func eccPemParse(key []byte, pemType string) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, errors.New("no pem block found")
	}
	if pemType != "" && block.Type != pemType {
		return nil, errors.New(strings.Join([]string{"Key's type is ", block.Type, ", expected ", pemType}, ""))
	}
	return block.Bytes, nil
}

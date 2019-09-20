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
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"math/big"
	"os"
	"path/filepath"
)

// ECCCommon ECC椭圆加密工具，依赖ETH的包
//
// ECC，全称椭圆曲线密码学（英语：Elliptic curve cryptography，缩写为 ECC），主要是指相关数学原理
//
// ECIES，在ECC原理的基础上实现的一种公钥加密方法，和RSA类似
//
// ECDSA，在ECC原理上实现的签名方法
//
// ECDH在ECC和DH的基础上实现的密钥交换算法
type ECCCommon struct{}

// GenerateKey 生成公私钥对
//
// path 指定公私钥所在生成目录
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func (e *ECCCommon) GenerateKey(path, priFileName, pubFileName string, curve elliptic.Curve) (*ecies.PrivateKey, *ecies.PublicKey, error) {
	var (
		privateKey *ecdsa.PrivateKey
		err        error
	)
	// 创建公私钥生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return nil, nil, err
		}
	}
	privateKey, err = ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, nil, err
	}

	if err = e.SavePri(filepath.Join(path, priFileName), privateKey); nil != err {
		return nil, nil, err
	}
	//if err = e.SavePub(filepath.Join(path, pubFileName), &privateKey.PublicKey, curve); nil != err {
	//	return nil, nil, err
	//}

	pri := ecies.ImportECDSA(privateKey)
	return pri, &pri.PublicKey, nil
}

// SavePri 将私钥保存到给定文件，密钥数据保存为hex编码
func (e *ECCCommon) SavePri(file string, privateKey *ecdsa.PrivateKey) error {
	return crypto.SaveECDSA(file, privateKey)
}

// LoadPri 从文件中加载私钥
//
// file 文件路径
func (e *ECCCommon) LoadPri(file string) (*ecdsa.PrivateKey, error) {
	return crypto.LoadECDSA(file)
}

//// SavePub 将公钥保存到给定文件，密钥数据保存为hex编码
////
//// file 文件路径
////
//// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
//func (e *ECCCommon) SavePub(file string, publicKey *ecdsa.PublicKey, curve elliptic.Curve) error {
//	k := hex.EncodeToString(e.PubKey2Bytes(publicKey, curve))
//	return ioutil.WriteFile(file, []byte(k), 0600)
//}

//// LoadPub 从文件中加载公钥
////
//// file 文件路径
////
//// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
//func (e *ECCCommon) LoadPub(file string, curve elliptic.Curve) (*ecdsa.PublicKey, error) {
//	var lenBuf int
//	switch curve {
//	default:
//		lenBuf = 130
//	case elliptic.P224():
//		lenBuf = 114
//	case elliptic.P256():
//		lenBuf = 130
//	case elliptic.P384():
//		lenBuf = 194
//	case elliptic.P521():
//		lenBuf = 266
//	}
//	buf := make([]byte, lenBuf)
//	fd, err := os.Open(file)
//	if err != nil {
//		return nil, err
//	}
//	defer func() { _ = fd.Close() }()
//	if _, err := io.ReadFull(fd, buf); err != nil {
//		return nil, err
//	}
//
//	key, err := hex.DecodeString(string(buf))
//	if err != nil {
//		return nil, err
//	}
//	return e.Bytes2PubKey(key, curve)
//}

//// PriKey2Bytes 私钥转[]byte
//func (e *ECCCommon) PriKey2Bytes(privateKey *ecdsa.PrivateKey) []byte {
//	return crypto.FromECDSA(privateKey)
//}
//
//// Bytes2PriKey []byte转私钥
//func (e *ECCCommon) Bytes2PriKey(data []byte) (*ecdsa.PrivateKey, error) {
//	return crypto.ToECDSA(data)
//}

//// PubKey2Bytes 公钥转[]byte
////
//// pub 公钥
////
//// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
//func (e *ECCCommon) PubKey2Bytes(publicKey *ecdsa.PublicKey, curve elliptic.Curve) []byte {
//	if publicKey == nil || publicKey.X == nil || publicKey.Y == nil {
//		return nil
//	}
//	return elliptic.Marshal(curve, publicKey.X, publicKey.Y)
//}
//
//// Bytes2PubKey []byte转公钥
//func (e *ECCCommon) Bytes2PubKey(data []byte, curve elliptic.Curve) (*ecdsa.PublicKey, error) {
//	x, y := elliptic.Unmarshal(curve, data)
//	if x == nil {
//		return nil, errors.New("invalid public key")
//	}
//	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, nil
//}

// Encrypt 加密
func (e *ECCCommon) Encrypt(data []byte, publicKey *ecies.PublicKey) ([]byte, error) {
	return ecies.Encrypt(rand.Reader, publicKey, data, nil, nil)
}

//// EncryptFP 加密
////
//// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
//func (e *ECCCommon) EncryptFP(data []byte, publicKeyFilePath string, curve elliptic.Curve) ([]byte, error) {
//	publicKey, err := e.LoadPub(publicKeyFilePath, curve)
//	if nil != err {
//		return nil, err
//	}
//	return e.Encrypt(data, ecies.ImportECDSAPublic(publicKey))
//}

// Decrypt 解密
func (e *ECCCommon) Decrypt(data []byte, privateKey *ecies.PrivateKey) ([]byte, error) {
	return privateKey.Decrypt(data, nil, nil)
}

//// DecryptFP 解密
//func (e *ECCCommon) DecryptFP(data []byte, privateKeyFilePath string) ([]byte, error) {
//	privateKey, err := e.LoadPri(privateKeyFilePath)
//	if nil != err {
//		return nil, err
//	}
//	return e.Decrypt(data, ecies.ImportECDSA(privateKey))
//}

// Sign 签名
func (e *ECCCommon) Sign(privateKey *ecdsa.PrivateKey, data []byte) (sign []byte, err error) {
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

// Verify 验签
func (e *ECCCommon) Verify(publicKey *ecdsa.PublicKey, data, sign []byte) (bool, error) {
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

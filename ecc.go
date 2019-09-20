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
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"io/ioutil"
	"os"
	"path/filepath"
)

// ECCCommon ECC椭圆加密工具，依赖ETH的包
type ECCCommon struct{}

// GenerateKey 生成公私钥对
//
// path 指定公私钥所在生成目录
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P224()/elliptic.P256()/elliptic.P384()/elliptic.P512()
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

	priK := hex.EncodeToString(e.PriKey2Bytes(privateKey))
	if err = ioutil.WriteFile(filepath.Join(path, priFileName), []byte(priK), 0600); nil != err {
		return nil, nil, err
	}
	pubK := hex.EncodeToString(e.PubKey2Bytes(&privateKey.PublicKey, curve))
	if err = ioutil.WriteFile(filepath.Join(path, pubFileName), []byte(pubK), 0600); nil != err {
		return nil, nil, err
	}

	pri := ecies.ImportECDSA(privateKey)
	return pri, &pri.PublicKey, nil
}

// PriKey2Bytes 私钥转[]byte
func (e *ECCCommon) PriKey2Bytes(pri *ecdsa.PrivateKey) []byte {
	if pri == nil {
		return nil
	}
	return math.PaddedBigBytes(pri.D, pri.Params().BitSize/8)
}

// Bytes2PriKey []byte转私钥
func (e *ECCCommon) Bytes2PriKey(data []byte) (*ecdsa.PrivateKey, error) {
	return crypto.ToECDSA(data)
}

// PubKey2Bytes 公钥转[]byte
//
// pub 公钥
//
// curve 曲线生成类型，如 crypto.S256()/elliptic.P224()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func (e *ECCCommon) PubKey2Bytes(pub *ecdsa.PublicKey, curve elliptic.Curve) []byte {
	if pub == nil || pub.X == nil || pub.Y == nil {
		return nil
	}
	return elliptic.Marshal(curve, pub.X, pub.Y)
}

// Bytes2PubKey []byte转公钥
func (e *ECCCommon) Bytes2PubKey(pub []byte, curve elliptic.Curve) *ecdsa.PublicKey {
	if len(pub) == 0 {
		return nil
	}
	x, y := elliptic.Unmarshal(curve, pub)
	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}
}

// ECCEncrypt 加密
func (e *ECCCommon) ECCEncrypt(pt []byte, puk ecies.PublicKey) ([]byte, error) {
	ct, err := ecies.Encrypt(rand.Reader, &puk, pt, nil, nil)
	return ct, err
}

// ECCDecrypt 解密
func (e *ECCCommon) ECCDecrypt(ct []byte, prk ecies.PrivateKey) ([]byte, error) {
	pt, err := prk.Decrypt(ct, nil, nil)
	return pt, err
}

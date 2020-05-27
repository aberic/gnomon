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
	"encoding/pem"
	"github.com/tjfoc/gmsm/sm2"
)

// SM2Generate SM2公钥私钥产生
//
// path 指定公私钥所在生成目录
func SM2Generate() (pri *sm2.PrivateKey, pub *sm2.PublicKey, err error) {
	// 生成密钥对
	if pri, err = sm2.GenerateKey(); err == nil {
		pub = &pri.PublicKey
	}
	return
}

// SM2GeneratePemBytes SM2公钥私钥产生
func SM2GeneratePemBytes(priPemType, pubPemType, passwd string) (priBytes, pubBytes []byte, err error) {
	var (
		pri *sm2.PrivateKey
		pub *sm2.PublicKey
	)
	// 生成密钥对
	if pri, pub, err = SM2Generate(); nil == err {
		if priBytes, err = SM2Pri2Bytes(priPemType, passwd, pri); nil != err {
			return
		}
		if pubBytes, err = SM2Pub2Bytes(pubPemType, pub); nil != err {
			return
		}
	}
	return
}

func SM2Pri2Bytes(priPemType, passwd string, pri *sm2.PrivateKey) (data []byte, err error) {
	var (
		der, pw []byte
		block   *pem.Block
	)
	if StringIsEmpty(passwd) {
		pw = nil
	} else {
		pw = []byte(passwd)
	}
	if der, err = sm2.MarshalSm2PrivateKey(pri, pw); err == nil {
		if StringIsEmpty(priPemType) {
			block = &pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: der,
			}
		} else {
			block = &pem.Block{
				Type:  priPemType,
				Bytes: der,
			}
		}
		data = pem.EncodeToMemory(block)
	}
	return
}

func SM2Pub2Bytes(pubPemType string, pub *sm2.PublicKey) (data []byte, err error) {
	var (
		der   []byte
		block *pem.Block
	)
	if der, err = sm2.MarshalSm2PublicKey(pub); err == nil {
		if StringIsEmpty(pubPemType) {
			block = &pem.Block{
				Type:  "PUBLIC KEY",
				Bytes: der,
			}
		} else {
			block = &pem.Block{
				Type:  pubPemType,
				Bytes: der,
			}
		}
		data = pem.EncodeToMemory(block)
	}
	return
}

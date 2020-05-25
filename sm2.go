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
	"os"
	"path/filepath"
)

// SM2Generate SM2公钥私钥产生
//
// path 指定公私钥所在生成目录
func SM2Generate(path, priFileName, pubFileName string) error {
	priv, err := sm2.GenerateKey() // 生成密钥对
	if err != nil {
		return err
	}
	pub := &priv.PublicKey

	priFilePath := filepath.Join(path, priFileName)
	pubFilePath := filepath.Join(path, pubFileName)
	if err = writePrivateKeytoPem(path, priFilePath, "", priv, nil); nil != err { // 生成密钥文件
		return err
	}
	return writePublicKeytoPem(path, pubFilePath, "", pub, nil) // 生成公钥文件
}

// SM2GenerateCustom SM2公钥私钥产生
//
// path 指定公私钥所在生成目录
func SM2GenerateCustom(path, priFileName, pubFileName, priPemType, pubPemType string) error {
	priv, err := sm2.GenerateKey() // 生成密钥对
	if err != nil {
		return err
	}
	pub := &priv.PublicKey

	priFilePath := filepath.Join(path, priFileName)
	pubFilePath := filepath.Join(path, pubFileName)
	if err = writePrivateKeytoPem(path, priFilePath, priPemType, priv, nil); nil != err { // 生成密钥文件
		return err
	}
	return writePublicKeytoPem(path, pubFilePath, pubPemType, pub, nil) // 生成公钥文件
}

func writePrivateKeytoPem(path, filePath, priPemType string, key *sm2.PrivateKey, pwd []byte) error {
	var block *pem.Block

	der, err := sm2.MarshalSm2PrivateKey(key, pwd)
	if err != nil {
		return err
	}
	if StringIsEmpty(priPemType) {
		if pwd != nil {
			block = &pem.Block{
				Type:  "ENCRYPTED PRIVATE KEY",
				Bytes: der,
			}
		} else {
			block = &pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: der,
			}
		}
	} else {
		block = &pem.Block{
			Type:  priPemType,
			Bytes: der,
		}
	}
	if !FilePathExists(path) {
		if err = os.MkdirAll(path, 0755); nil != err {
			return err
		}
	}
	var fileIO *os.File
	if fileIO, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); nil != err {
		return err
	}
	defer func() { _ = fileIO.Close() }()
	return pem.Encode(fileIO, block)
}

func writePublicKeytoPem(path, filePath, pubPemType string, key *sm2.PublicKey, _ []byte) error {
	der, err := sm2.MarshalSm2PublicKey(key)
	if err != nil {
		return err
	}
	var block *pem.Block
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
	if !FilePathExists(path) {
		if err = os.MkdirAll(path, 0755); nil != err {
			return err
		}
	}
	var fileIO *os.File
	if fileIO, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); nil != err {
		return err
	}
	defer func() { _ = fileIO.Close() }()
	return pem.Encode(fileIO, block)
}

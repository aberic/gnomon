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
	"encoding/hex"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/ecies"
	"gotest.tools/assert"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var (
	contentECC = "this is a test"

	patheccs256 = "./tmp/example/ecc/s256"
	patheccp224 = "./tmp/example/ecc/p224"
	patheccp256 = "./tmp/example/ecc/p256"
	patheccp384 = "./tmp/example/ecc/p384"
	patheccp521 = "./tmp/example/ecc/p521"

	patheccpemp224 = "./tmp/example/ecc/pemp224"
	patheccpemp256 = "./tmp/example/ecc/pemp256"
	patheccpemp384 = "./tmp/example/ecc/pemp384"
	patheccpemp521 = "./tmp/example/ecc/pemp521"

	privateECCName = "private.key"
	publicECCName  = "public.key"

	priKeyS256 *ecdsa.PrivateKey
	priKeyP224 *ecdsa.PrivateKey
	priKeyP256 *ecdsa.PrivateKey
	priKeyP384 *ecdsa.PrivateKey
	priKeyP521 *ecdsa.PrivateKey

	pubKeyS256 *ecdsa.PublicKey
	pubKeyP224 *ecdsa.PublicKey
	pubKeyP256 *ecdsa.PublicKey
	pubKeyP384 *ecdsa.PublicKey
	pubKeyP521 *ecdsa.PublicKey

	pri *ecies.PrivateKey
	pub *ecies.PublicKey

	dataECC       []byte
	dataECCEncode []byte
	signECCResult []byte
	valid         bool
	errECC        error
)

// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func TestECCCommon_Generate(t *testing.T) {
	if _, _, errECC = ECCGenerate(crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
}

func TestECCGeneratePemBytes(t *testing.T) {
	priBytes, pubBytes, err := ECCGeneratePemBytes("SM2 PRIVATE KEY", "SM2 PUBLIC KEY", "", elliptic.P256())
	assert.NilError(t, err)
	t.Log(string(priBytes))
	t.Log(string(pubBytes))
}

// curve 曲线生成类型，如 crypto.S256()/elliptic.P256()/elliptic.P384()/elliptic.P512()
func TestECCCommon_GenerateKey(t *testing.T) {
	if errECC = ECCGenerateKey(patheccs256, privateECCName, publicECCName, crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
	if errECC = ECCGenerateKey(patheccp224, privateECCName, publicECCName, elliptic.P224()); nil != errECC {
		t.Skip(errECC)
	}
	if errECC = ECCGenerateKey(patheccp256, privateECCName, publicECCName, elliptic.P256()); nil != errECC {
		t.Skip(errECC)
	}
	if errECC = ECCGenerateKey(patheccp384, privateECCName, publicECCName, elliptic.P384()); nil != errECC {
		t.Skip(errECC)
	}

	if errECC = ECCGeneratePriKey(patheccp521, privateECCName, elliptic.P521()); nil != errECC {
		t.Skip(errECC)
	}
	if priKeyP521, errECC = ECCLoadPri(filepath.Join(patheccp521, privateECCName), elliptic.P521()); nil != errECC {
		t.Skip(errECC)
	}
	if errECC = ECCGeneratePubKey(priKeyP521, patheccp521, publicECCName, elliptic.P521()); nil != errECC {
		t.Skip(errECC)
	}

	if errECC = ECCGeneratePemKey(patheccpemp224, privateECCName, publicECCName, "PRIVATE KEY", "PUBLIC KEY", elliptic.P224()); nil != errECC {
		t.Skip(errECC)
	}

	if errECC = ECCGeneratePemPriKey(patheccpemp224, privateECCName, "PRIVATE KEY", elliptic.P224()); nil != errECC {
		t.Skip(errECC)
	}
	dataECC, errECC = ioutil.ReadFile(filepath.Join(patheccpemp224, privateECCName))
	if nil != errECC {
		t.Skip(errECC)
	}
	if priKeyP224, errECC = ECCLoadPriPem(dataECC, "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if errECC = ECCGeneratePemPubKey(priKeyP224, patheccpemp224, publicECCName, "PUBLIC KEY"); nil != errECC {
		t.Skip(errECC)
	}

	if errECC = ECCGeneratePemPriKeyWithPass(patheccpemp256, privateECCName, "123456", "PRIVATE KEY", elliptic.P256()); nil != errECC {
		t.Skip(errECC)
	}
	dataECC, errECC = ioutil.ReadFile(filepath.Join(patheccpemp256, privateECCName))
	if nil != errECC {
		t.Skip(errECC)
	}
	if priKeyP256, errECC = ECCLoadPriPemWithPass(dataECC, "123456", "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if errECC = ECCGeneratePemPubKey(priKeyP256, patheccpemp256, publicECCName, "PUBLIC KEY"); nil != errECC {
		t.Skip(errECC)
	}

	if errECC = ECCGeneratePemKey(patheccpemp256, privateECCName, publicECCName, "PRIVATE KEY", "PUBLIC KEY", elliptic.P256()); nil != errECC {
		t.Skip(errECC)
	}

	if errECC = ECCGeneratePemPriKeyWithPass(patheccpemp384, privateECCName, "123456", "PRIVATE KEY", elliptic.P384()); nil != errECC {
		t.Skip(errECC)
	}
	if priKeyP384, errECC = ECCLoadPriPemFPWithPass(filepath.Join(patheccpemp384, privateECCName), "123456", "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if errECC = ECCGeneratePemPubKey(priKeyP384, patheccpemp384, publicECCName, "PUBLIC KEY"); nil != errECC {
		t.Skip(errECC)
	}

	if errECC = ECCGeneratePemPriKey(patheccpemp521, privateECCName, "PRIVATE KEY", elliptic.P521()); nil != errECC {
		t.Skip(errECC)
	}
	if priKeyP521, errECC = ECCLoadPriPemFP(filepath.Join(patheccpemp521, privateECCName), "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if errECC = ECCGeneratePemPubKey(priKeyP521, patheccpemp521, publicECCName, "PUBLIC KEY"); nil != errECC {
		t.Skip(errECC)
	}

	t.Log(ECCSavePriPem(priKeyP521, patheccpemp521, "test.pem", "PRIVATE KEY"))
	t.Log(ECCSavePriPem(priKeyP521, "/etc/ecc", "test.pem", "PRIVATE KEY"))
}

func TestLog(t *testing.T) {
	if priKeyP521, errECC = ECCLoadPriPemFP(filepath.Join(patheccpemp256, privateECCName), "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	t.Log(ECCSavePriPem(priKeyP521, patheccpemp521, "test.pem", "PRIVATE KEY"))
}

func TestECCCommon_GenerateKey_FailPathPermissionDenied(t *testing.T) {
	if errECC = ECCGenerateKey("/etc/test", privateECCName, publicECCName, crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
}

func TestECCCommon_GeneratePriKey_FailPathPermissionDenied(t *testing.T) {
	if errECC = ECCGeneratePriKey("/etc/test", privateECCName, crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
}

func TestECCCommon_GeneratePemKey_FailPathPermissionDenied(t *testing.T) {
	if errECC = ECCGeneratePemKey("/etc/test", privateECCName, publicECCName, "PRIVATE KEY", "PUBLIC KEY", crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
}

func TestECCCommon_GeneratePemPriKey_FailPathPermissionDenied(t *testing.T) {
	if errECC = ECCGeneratePemPriKey("/etc/test", "PRIVATE KEY", privateECCName, crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
}

func TestECCCommon_GenerateKey_FailCreate(t *testing.T) {
	if errECC = ECCGenerateKey("/etc", privateECCName, publicECCName, crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
}

func TestECCCommon_GeneratePriKey_FailCreate(t *testing.T) {
	if errECC = ECCGeneratePriKey("/etc", privateECCName, crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
}

func TestECCCommon_GeneratePemKey_FailCreate(t *testing.T) {
	if errECC = ECCGeneratePemKey("/etc", privateECCName, publicECCName, "PRIVATE KEY", "PUBLIC KEY", crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
}

func TestECCCommon_GeneratePemPriKey_FailCreate(t *testing.T) {
	if errECC = ECCGeneratePemPriKey("/etc", privateECCName, "PRIVATE KEY", crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
}

func TestECCCommon_GeneratePubKey_Fail(t *testing.T) {
	if errECC = ECCGeneratePriKey(patheccp521, privateECCName, elliptic.P521()); nil != errECC {
		t.Skip(errECC)
	}
	if priKeyP521, errECC = ECCLoadPri(filepath.Join(patheccp521, privateECCName), elliptic.P521()); nil != errECC {
		t.Skip(errECC)
	}
	t.Log(ECCGeneratePubKey(priKeyP521, "/etc/ecc", publicECCName, elliptic.P521()))
}

func TestECCCommon_GeneratePemPubKey_Fail(t *testing.T) {
	if errECC = ECCGeneratePemPriKey(patheccpemp521, privateECCName, "PRIVATE KEY", elliptic.P521()); nil != errECC {
		t.Skip(errECC)
	}
	if priKeyP521, errECC = ECCLoadPriPemFP(filepath.Join(patheccpemp521, privateECCName), "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	t.Log(ECCGeneratePemPubKey(priKeyP521, "/etc/ecc", publicECCName, "PUBLIC KEY"))
	t.Log(ECCSavePubPem("/etc/ecc", "PRIVATE KEY", &priKeyP521.PublicKey))
}

func TestECCCommon_LoadPri_Fail(t *testing.T) {
	if priKeyS256, errECC = ECCLoadPri(filepath.Join(patheccs256, "test.ha"), crypto.S256()); nil != errECC {
		t.Log(errECC)
	}
}

func TestECCCommon_LoadPub_Fail1(t *testing.T) {
	if _, errECC = ECCLoadPub("/etc/ecc", elliptic.P224()); nil != errECC {
		t.Log(errECC)
	}
}

func TestECCCommon_LoadPubPem_Fail1(t *testing.T) {
	if _, errECC = ECCLoadPubPemFP("/etc/ecc", "PUBLIC KEY"); nil != errECC {
		t.Log(errECC)
	}
}

func TestECCCommon_LoadPubPem_Fail2(t *testing.T) {
	if _, errECC = ECCLoadPubPemFP(filepath.Join(patheccs256, publicECCName), "PUBLIC KEY"); nil != errECC {
		t.Log(errECC)
	}
}

func TestECCCommon_PubKey2Bytes_Fail(t *testing.T) {
	t.Log(ECCPubKey2Bytes(&ecdsa.PublicKey{}, elliptic.P224()))
}

func TestECCCommon_Bytes2PubKey_Fail(t *testing.T) {
	_, errECC = ECCBytes2PubKey([]byte{}, elliptic.P224())
	t.Log(errECC)
}

func TestECCCommon_Verify_Fail(t *testing.T) {
	if priKeyP384, errECC = ECCLoadPriPemFPWithPass(filepath.Join(patheccpemp384, privateECCName), "123456", "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	_, errECC = ECCVerify(&priKeyP384.PublicKey, []byte("yes"), []byte{})
	t.Log(errECC)
}

func TestCryptoECC_Crypt(t *testing.T) {
	t.Log("加密前：", contentECC)
	t.Log("=================================")

	if priKeyS256, errECC = ECCLoadPri(filepath.Join(patheccs256, privateECCName), crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
	pri = ecies.ImportECDSA(priKeyS256)
	if dataECCEncode, errECC = ECCEncrypt([]byte(contentECC), &pri.PublicKey); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("加密后S256：", hex.EncodeToString(dataECCEncode))
	if dataECC, errECC = ECCDecrypt(dataECCEncode, pri); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("解密后S256：", string(dataECC))
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")

	if priKeyP256, errECC = ECCLoadPri(filepath.Join(patheccp256, privateECCName), elliptic.P256()); nil != errECC {
		t.Skip(errECC)
	}
	pri = ecies.ImportECDSA(priKeyP256)
	if dataECCEncode, errECC = ECCEncrypt([]byte(contentECC), &pri.PublicKey); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("加密后P256：", hex.EncodeToString(dataECCEncode))
	if dataECC, errECC = ECCDecrypt(dataECCEncode, pri); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("解密后P256：", string(dataECC))
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")

	if priKeyP384, errECC = ECCLoadPri(filepath.Join(patheccp384, privateECCName), elliptic.P384()); nil != errECC {
		t.Skip(errECC)
	}
	pri = ecies.ImportECDSA(priKeyP384)
	if dataECCEncode, errECC = ECCEncrypt([]byte(contentECC), &pri.PublicKey); nil != errECC {
		t.Log(errECC)
	}
	t.Log("加密后P384：", hex.EncodeToString(dataECCEncode))
	if dataECC, errECC = ECCDecrypt(dataECCEncode, pri); nil != errECC {
		t.Log(errECC)
	}
	t.Log("解密后P384：", string(dataECC))
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")

	if priKeyP521, errECC = ECCLoadPri(filepath.Join(patheccp521, privateECCName), elliptic.P521()); nil != errECC {
		t.Skip(errECC)
	}
	pri = ecies.ImportECDSA(priKeyP521)
	if dataECCEncode, errECC = ECCEncrypt([]byte(contentECC), &pri.PublicKey); nil != errECC {
		t.Log(errECC)
	}
	t.Log("加密后P521：", hex.EncodeToString(dataECCEncode))
	if dataECC, errECC = ECCDecrypt(dataECCEncode, pri); nil != errECC {
		t.Log(errECC)
	}
	t.Log("解密后P521：", string(dataECC))
	t.Log("=================================")
}

func TestCryptoECCPem_Crypt(t *testing.T) {
	t.Log("加密前：", contentECC)
	t.Log("=================================")

	if priKeyP256, errECC = ECCLoadPriPemFP(filepath.Join(patheccpemp256, privateECCName), "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	pri = ecies.ImportECDSA(priKeyP256)
	if dataECCEncode, errECC = ECCEncrypt([]byte(contentECC), &pri.PublicKey); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("加密后PemP256：", hex.EncodeToString(dataECCEncode))
	if dataECC, errECC = ECCDecrypt(dataECCEncode, pri); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("解密后PemP256：", string(dataECC))
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")
}

func TestCryptoECC_CryptFile(t *testing.T) {
	t.Log("加密前：", contentECC)
	t.Log("=================================")

	if pubKeyS256, errECC = ECCLoadPub(filepath.Join(patheccs256, publicECCName), crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
	pub = ecies.ImportECDSAPublic(pubKeyS256)
	if dataECCEncode, errECC = ECCEncrypt([]byte(contentECC), pub); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("加密后S256：", hex.EncodeToString(dataECCEncode))
	if priKeyS256, errECC = ECCLoadPri(filepath.Join(patheccs256, privateECCName), crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
	pri = ecies.ImportECDSA(priKeyS256)
	if dataECC, errECC = ECCDecrypt(dataECCEncode, pri); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("解密后S256：", string(dataECC))
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")

	if pubKeyP256, errECC = ECCLoadPub(filepath.Join(patheccp256, publicECCName), elliptic.P256()); nil != errECC {
		t.Skip(errECC)
	}
	pub = ecies.ImportECDSAPublic(pubKeyP256)
	if dataECCEncode, errECC = ECCEncrypt([]byte(contentECC), pub); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("加密后P256：", hex.EncodeToString(dataECCEncode))
	if priKeyP256, errECC = ECCLoadPri(filepath.Join(patheccp256, privateECCName), elliptic.P256()); nil != errECC {
		t.Skip(errECC)
	}
	pri = ecies.ImportECDSA(priKeyP256)
	if dataECC, errECC = ECCDecrypt(dataECCEncode, pri); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("解密后P256：", string(dataECC))
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")
}

func TestCryptoECCPem_CryptFile(t *testing.T) {
	t.Log("加密前：", contentECC)
	t.Log("=================================")

	if pubKeyP256, errECC = ECCLoadPubPemFP(filepath.Join(patheccpemp256, publicECCName), "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	pub = ecies.ImportECDSAPublic(pubKeyP256)
	if dataECCEncode, errECC = ECCEncrypt([]byte(contentECC), pub); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("加密后PemP256：", hex.EncodeToString(dataECCEncode))
	if priKeyP256, errECC = ECCLoadPriPemFP(filepath.Join(patheccpemp256, privateECCName), "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	pri = ecies.ImportECDSA(priKeyP256)
	if dataECC, errECC = ECCDecrypt(dataECCEncode, pri); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("解密后PemP256：", string(dataECC))
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")
}

func TestCryptoECC_Sign(t *testing.T) {
	t.Log("签名内容：", contentECC)
	t.Log("=================================")

	if priKeyS256, errECC = ECCLoadPri(filepath.Join(patheccs256, privateECCName), crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyS256, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码S256", string(signECCResult))
	if valid, errECC = ECCVerify(&priKeyS256.PublicKey, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过S256")
		} else {
			t.Log("验签错误S256")
		}
	}
	t.Log("=================================")

	if priKeyP256, errECC = ECCLoadPri(filepath.Join(patheccp256, privateECCName), elliptic.P256()); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP256, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P256", string(signECCResult))
	if valid, errECC = ECCVerify(&priKeyP256.PublicKey, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P256")
		} else {
			t.Log("验签错误P256")
		}
	}
	t.Log("=================================")

	if priKeyP384, errECC = ECCLoadPri(filepath.Join(patheccp384, privateECCName), elliptic.P384()); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP384, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P384", string(signECCResult))
	if valid, errECC = ECCVerify(&priKeyP384.PublicKey, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P384")
		} else {
			t.Log("验签错误P384")
		}
	}
	t.Log("=================================")

	if priKeyP521, errECC = ECCLoadPri(filepath.Join(patheccp521, privateECCName), elliptic.P521()); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP521, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P521", string(signECCResult))
	if valid, errECC = ECCVerify(&priKeyP521.PublicKey, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P521")
		} else {
			t.Log("验签错误P521")
		}
	}
	t.Log("=================================")
}

func TestCryptoECCPem_Sign(t *testing.T) {
	t.Log("签名内容：", contentECC)
	t.Log("=================================")

	if priKeyP224, errECC = ECCLoadPriPemFP(filepath.Join(patheccpemp256, privateECCName), "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP224, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P224", string(signECCResult))
	if valid, errECC = ECCVerify(&priKeyP224.PublicKey, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P224")
		} else {
			t.Log("验签错误P224")
		}
	}
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")

	if priKeyP256, errECC = ECCLoadPriPemFP(filepath.Join(patheccpemp256, privateECCName), "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP256, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P256", string(signECCResult))
	if valid, errECC = ECCVerify(&priKeyP256.PublicKey, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P256")
		} else {
			t.Log("验签错误P256")
		}
	}
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")

	if priKeyP384, errECC = ECCLoadPriPemFPWithPass(filepath.Join(patheccpemp384, privateECCName), "123456", "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP384, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P384", string(signECCResult))
	if valid, errECC = ECCVerify(&priKeyP384.PublicKey, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P384")
		} else {
			t.Log("验签错误P384")
		}
	}
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")

	if priKeyP521, errECC = ECCLoadPriPemFP(filepath.Join(patheccpemp521, privateECCName), "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP521, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P521", string(signECCResult))
	if valid, errECC = ECCVerify(&priKeyP521.PublicKey, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P521")
		} else {
			t.Log("验签错误P521")
		}
	}
	t.Log("=================================")
}

func TestCryptoECC_Sign_File(t *testing.T) {
	t.Log("签名内容：", contentECC)
	t.Log("=================================")

	if priKeyS256, errECC = ECCLoadPri(filepath.Join(patheccs256, privateECCName), crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyS256, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码S256", string(signECCResult))
	if pubKeyS256, errECC = ECCLoadPub(filepath.Join(patheccs256, publicECCName), crypto.S256()); nil != errECC {
		t.Skip(errECC)
	}
	if valid, errECC = ECCVerify(pubKeyS256, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过S256")
		} else {
			t.Log("验签错误S256")
		}
	}
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")

	if priKeyP224, errECC = ECCLoadPri(filepath.Join(patheccp224, privateECCName), elliptic.P224()); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP224, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P224", string(signECCResult))
	if pubKeyP224, errECC = ECCLoadPub(filepath.Join(patheccp224, publicECCName), elliptic.P224()); nil != errECC {
		t.Skip(errECC)
	}
	if valid, errECC = ECCVerify(pubKeyP224, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P224")
		} else {
			t.Log("验签错误P224")
		}
	}
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")

	if priKeyP256, errECC = ECCLoadPri(filepath.Join(patheccp256, privateECCName), elliptic.P256()); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP256, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P256", string(signECCResult))
	if pubKeyP256, errECC = ECCLoadPub(filepath.Join(patheccp256, publicECCName), elliptic.P256()); nil != errECC {
		t.Skip(errECC)
	}
	if valid, errECC = ECCVerify(pubKeyP256, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P256")
		} else {
			t.Log("验签错误P256")
		}
	}
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")

	if priKeyP384, errECC = ECCLoadPri(filepath.Join(patheccp384, privateECCName), elliptic.P384()); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP384, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P384", string(signECCResult))
	if pubKeyP384, errECC = ECCLoadPub(filepath.Join(patheccp384, publicECCName), elliptic.P384()); nil != errECC {
		t.Skip(errECC)
	}
	if valid, errECC = ECCVerify(pubKeyP384, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P384")
		} else {
			t.Log("验签错误P384")
		}
	}
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")

	if priKeyP521, errECC = ECCLoadPri(filepath.Join(patheccp521, privateECCName), elliptic.P521()); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP521, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P521", string(signECCResult))
	if pubKeyP521, errECC = ECCLoadPub(filepath.Join(patheccp521, publicECCName), elliptic.P521()); nil != errECC {
		t.Skip(errECC)
	}
	if valid, errECC = ECCVerify(pubKeyP521, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P521")
		} else {
			t.Log("验签错误P521")
		}
	}
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")
}

func TestCryptoECCPem_Sign_File(t *testing.T) {
	t.Log("签名内容：", contentECC)
	t.Log("=================================")

	if priKeyP224, errECC = ECCLoadPriPemFP(filepath.Join(patheccpemp224, privateECCName), "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP224, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P224", string(signECCResult))
	if pubKeyP224, errECC = ECCLoadPubPemFP(filepath.Join(patheccpemp224, publicECCName), "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if valid, errECC = ECCVerify(pubKeyP224, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P224")
		} else {
			t.Log("验签错误P224")
		}
	}
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")

	if priKeyP256, errECC = ECCLoadPriPemFP(filepath.Join(patheccpemp256, privateECCName), "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP256, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P256", string(signECCResult))
	if pubKeyP256, errECC = ECCLoadPubPemFP(filepath.Join(patheccpemp256, publicECCName), "PUBLIC KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if valid, errECC = ECCVerify(pubKeyP256, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P256")
		} else {
			t.Log("验签错误P256")
		}
	}
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")

	if priKeyP384, errECC = ECCLoadPriPemFPWithPass(filepath.Join(patheccpemp384, privateECCName), "123456", "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP384, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P384", string(signECCResult))
	if pubKeyP384, errECC = ECCLoadPubPemFP(filepath.Join(patheccpemp384, publicECCName), "PUBLIC KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if valid, errECC = ECCVerify(pubKeyP384, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P384")
		} else {
			t.Log("验签错误P384")
		}
	}
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")

	if priKeyP521, errECC = ECCLoadPriPemFP(filepath.Join(patheccpemp521, privateECCName), "PRIVATE KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if signECCResult, errECC = ECCSign(priKeyP521, []byte(contentECC)); nil != errECC {
		t.Skip(errECC)
	}
	t.Log("签名码P521", string(signECCResult))
	if pubKeyP521, errECC = ECCLoadPubPemFP(filepath.Join(patheccpemp521, publicECCName), "PUBLIC KEY"); nil != errECC {
		t.Skip(errECC)
	}
	if valid, errECC = ECCVerify(pubKeyP521, []byte(contentECC), signECCResult); nil != errECC {
		t.Skip(errECC)
	} else {
		if valid {
			t.Log("验签通过P521")
		} else {
			t.Log("验签错误P521")
		}
	}
	t.Log("=================================")
	t.Log("=================================")
	t.Log("=================================")
}

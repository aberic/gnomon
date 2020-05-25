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
	"crypto"
	"crypto/rsa"
	"crypto/x509"
	"encoding/hex"
	"io/ioutil"
	"path/filepath"
	"testing"
)

var (
	contentRSA       = "this is a test"
	pathrsapksc1256  = "./tmp/example/rsa/pksc1/256"
	pathrsapksc1512  = "./tmp/example/rsa/pksc1/512"
	pathrsapksc11024 = "./tmp/example/rsa/pksc1/1024"
	pathrsapksc12048 = "./tmp/example/rsa/pksc1/2048"
	pathrsapksc8256  = "./tmp/example/rsa/pksc8/256"
	pathrsapksc8512  = "./tmp/example/rsa/pksc8/512"
	pathrsapksc81024 = "./tmp/example/rsa/pksc8/1024"
	pathrsapksc82048 = "./tmp/example/rsa/pksc8/2048"
	privateRSAName   = "private.pem"
	publicRSAName    = "public.pem"

	priRSAKey *rsa.PrivateKey
	//pubRSAKey *rsa.PublicKey

	priRSAKeyData []byte
	pubRSAKeyData []byte
	dataRSA       []byte
	dataRSAEncode []byte
	signRSAResult []byte
	errRSA        error
)

func TestRSACommon_GenerateRsaKey(t *testing.T) {
	t.Log(RSAGenerateKey(256, pathrsapksc1256, privateRSAName, publicRSAName, "PRIVATE KEY", "PUBLIC KEY", RSAPKSC1()))
	t.Log(RSAGenerateKeyWithPass(512, pathrsapksc1512, privateRSAName, publicRSAName, "123456", "PRIVATE KEY", "PUBLIC KEY", x509.PEMCipher3DES, RSAPKSC1()))
	t.Log(RSAGenerateKey(1024, pathrsapksc11024, privateRSAName, publicRSAName, "PRIVATE KEY", "PRIVATE KEY", RSAPKSC1()))
	t.Log(RSAGenerateKeyWithPass(2048, pathrsapksc12048, privateRSAName, publicRSAName, "123456", "PRIVATE KEY", "PUBLIC KEY", x509.PEMCipher3DES, RSAPKSC1()))

	t.Log(RSAGenerateKeyWithPass(256, pathrsapksc8256, privateRSAName, publicRSAName, "123456", "PUBLIC KEY", "PUBLIC KEY", x509.PEMCipher3DES, RSAPKSC8()))
	t.Log(RSAGenerateKey(512, pathrsapksc8512, privateRSAName, publicRSAName, "PRIVATE KEY", "PUBLIC KEY", RSAPKSC8()))
	t.Log(RSAGenerateKeyWithPass(1024, pathrsapksc81024, privateRSAName, publicRSAName, "123456", "PRIVATE KEY", "PUBLIC KEY", x509.PEMCipher3DES, RSAPKSC8()))
	t.Log(RSAGenerateKey(2048, pathrsapksc82048, privateRSAName, publicRSAName, "PRIVATE KEY", "PUBLIC KEY", RSAPKSC8()))
}

func TestRSACommon_GenerateRsaCustomPriKey(t *testing.T) {
	t.Log(RSAGeneratePriKeyWithPass(256, pathrsapksc1256, "private1.pem", "123456", "PRIVATE KEY", x509.PEMCipher3DES, RSAPKSC1()))
	t.Log(RSAGeneratePriKey(512, pathrsapksc1512, "private1.pem", "PUBLIC KEY", RSAPKSC1()))
	t.Log(RSAGeneratePriKeyWithPass(1024, pathrsapksc11024, "private1.pem", "123456", "PRIVATE KEY", -1, RSAPKSC1()))
	t.Log(RSAGeneratePriKeyWithPass(1024, pathrsapksc11024, "private1.pem", "123456", "PRIVATE KEY", x509.PEMCipher3DES, RSAPKSC1()))
	t.Log(RSAGeneratePriKey(2048, pathrsapksc12048, "private1.pem", "PRIVATE KEY", RSAPKSC1()))
	t.Log()

	t.Log(RSAGeneratePriKey(256, pathrsapksc8256, "private1.pem", "PRIVATE KEY", RSAPKSC8()))
	t.Log(RSAGeneratePriKeyWithPass(512, pathrsapksc8512, "private1.pem", "123456", "PRIVATE KEY", x509.PEMCipher3DES, RSAPKSC8()))
	t.Log(RSAGeneratePriKey(1024, pathrsapksc81024, "private1.pem", "PRIVATE KEY", RSAPKSC8()))
	t.Log(RSAGeneratePriKeyWithPass(2048, pathrsapksc82048, "private1.pem", "123456", "PRIVATE KEY", -1, RSAPKSC8()))
	t.Log(RSAGeneratePriKeyWithPass(2048, pathrsapksc82048, "private1.pem", "123456", "PRIVATE KEY", x509.PEMCipher3DES, RSAPKSC8()))
	t.Log(RSAGeneratePubKey(nil, "/etc/pub", "public1.pem", "PUBLIC KEY"))
}

func TestRSACommon_GenerateRsaCustomPubKey(t *testing.T) {
	t.Log(RSAGeneratePubKeyFPWithPass(pathrsapksc1256+"/"+"private1.pem", "123456", pathrsapksc1256, "public1.pem", "PRIVATE KEY", "PUBLIC KEY", RSAPKSC1()))
	t.Log(RSAGeneratePubKeyFP(pathrsapksc1512+"/"+"private1.pem", pathrsapksc1512, "public1.pem", "PRIVATE KEY", "PUBLIC KEY", RSAPKSC1()))
	t.Log(RSAGeneratePubKeyFPWithPass(pathrsapksc11024+"/"+"private1.pem", "123456", pathrsapksc11024, "public1.pem", "PRIVATE KEY", "PUBLIC KEY", RSAPKSC1()))
	t.Log(RSAGeneratePubKeyFP(pathrsapksc12048+"/"+"private1.pem", pathrsapksc12048, "public1.pem", "PRIVATE KEY", "PUBLIC KEY", RSAPKSC1()))
	t.Log()

	t.Log(RSAGeneratePubKeyFP(pathrsapksc8256+"/"+"private1.pem", pathrsapksc8256, "public1.pem", "PRIVATE KEY", "PUBLIC KEY", RSAPKSC8()))
	t.Log(RSAGeneratePubKeyFPWithPass(pathrsapksc8512+"/"+"private1.pem", "123456", pathrsapksc8512, "public1.pem", "PRIVATE KEY", "PUBLIC KEY", RSAPKSC8()))
	t.Log(RSAGeneratePubKeyFP(pathrsapksc81024+"/"+"private1.pem", pathrsapksc81024, "public1.pem", "PRIVATE KEY", "PUBLIC KEY", RSAPKSC8()))
	t.Log(RSAGeneratePubKeyFPWithPass(pathrsapksc82048+"/"+"private1.pem", "123456", pathrsapksc82048, "public1.pem", "PRIVATE KEY", "PUBLIC KEY", RSAPKSC8()))
	t.Log()

	t.Log(RSAGeneratePubKeyFP(pathrsapksc12048+"/"+"private100.pem", pathrsapksc12048, "public1.pem", "PRIVATE KEY", "PUBLIC KEY", RSAPKSC1()))
	t.Log(RSAGeneratePubKeyFPWithPass(pathrsapksc82048+"/"+"private100.pem", "123456", pathrsapksc82048, "public1.pem", "PRIVATE KEY", "PUBLIC KEY", RSAPKSC8()))
	t.Log(RSAGeneratePubKeyFP(pathrsapksc82048+"/"+"private1.pem", pathrsapksc82048, "public2.pem", "PRIVATE KEY", "PUBLIC KEY", RSAPKSC1()))
}

func TestRSACommon_GenerateRsaKey_FailPathExists(t *testing.T) {
	t.Log(RSAGenerateKey(256, "/etc/test", privateRSAName, publicRSAName, "PRIVATE KEY", "PUBLIC KEY", RSAPKSC1()))
}

func TestRSACommon_GenerateRsaKey_FailGenerate(t *testing.T) {
	t.Log(RSAGenerateKey(-1, pathrsapksc1256, privateRSAName, publicRSAName, "PRIVATE KEY", "PUBLIC KEY", RSAPKSC1()))
}

func TestRSACommon_GenerateRsaKey_FailCreate(t *testing.T) {
	t.Log(RSAGenerateKey(256, "/etc", privateRSAName, publicRSAName, "PRIVATE KEY", "PUBLIC KEY", RSAPKSC1()))
}

func TestRSACommon_GenerateRsaPKSC8Key_FailPathExists(t *testing.T) {
	t.Log(RSAGenerateKey(256, "/etc/test", privateRSAName, publicRSAName, "PRIVATE KEY", "PUBLIC KEY", RSAPKSC8()))
}

func TestRSACommon_GenerateRsaPKSC8Key_FailGenerate(t *testing.T) {
	t.Log(RSAGenerateKey(-1, pathrsapksc1256, privateRSAName, publicRSAName, "PRIVATE KEY", "PUBLIC KEY", RSAPKSC8()))
}

func TestRSACommon_GenerateRsaPKSC8Key_FailCreate(t *testing.T) {
	t.Log(RSAGenerateKey(256, "/etc", privateRSAName, publicRSAName, "PRIVATE KEY", "PUBLIC KEY", RSAPKSC8()))
}

func TestRSACommon_RsaEncryptDecrypt_Fail(t *testing.T) {
	dataRSA, errRSA = RSAEncrypt([]byte{}, []byte(contentRSA), "PUBLIC KEY")
	if nil != errRSA {
		t.Log(errRSA)
	}

	dataRSA, errRSA = RSAEncryptFP(pathrsapksc1256+"/a/"+publicRSAName, "PUBLIC KEY", []byte(contentRSA))
	if nil != errRSA {
		t.Log(errRSA)
	}

	dataRSA, errRSA = RSADecrypt([]byte{}, []byte(contentRSA), "PUBLIC KEY", RSAPKSC8())
	if nil != errRSA {
		t.Log(errRSA)
	}
}

func TestRSACommon_RsaEncryptDecrypt(t *testing.T) {
	t.Log("加密前：", contentRSA)
	t.Log("=================================")

	pubRSAKeyData, errRSA = ioutil.ReadFile(filepath.Join(pathrsapksc1256, publicRSAName))
	if nil != errRSA {
		t.Error(errRSA)
	}
	dataRSA, errRSA = RSAEncrypt(pubRSAKeyData, []byte(contentRSA), "PUBLIC KEY")
	if nil != errRSA {
		t.Error(errRSA)
	}
	dataRSA, errRSA = RSAEncryptFP(pathrsapksc1256+"/"+publicRSAName, "PUBLIC KEY", []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后256：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = RSADecryptFP(pathrsapksc1256+"/"+privateRSAName, "PRIVATE KEY", dataRSA, RSAPKSC1())
	t.Log("解密后256：", string(dataRSAEncode))
	priRSAKeyData, errRSA = ioutil.ReadFile(filepath.Join(pathrsapksc1256, privateRSAName))
	if nil != errRSA {
		t.Error(errRSA)
	}
	dataRSAEncode, errRSA = RSADecrypt(priRSAKeyData, dataRSA, "PRIVATE KEY", RSAPKSC1())
	if nil != errRSA {
		t.Log(errRSA)
	}
	t.Log("解密后256：", string(dataRSAEncode))
	t.Log("=================================")

	dataRSA, errRSA = RSAEncryptFP(pathrsapksc1512+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后512：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = RSADecryptFP(pathrsapksc1512+"/"+privateRSAName, "PRIVATE KEY", dataRSA, RSAPKSC1())
	t.Log("解密后512：", string(dataRSAEncode))
	t.Log("=================================")

	dataRSA, errRSA = RSAEncryptFP(pathrsapksc11024+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后1024：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = RSADecryptFP(pathrsapksc11024+"/"+privateRSAName, "PRIVATE KEY", dataRSA, RSAPKSC1())
	t.Log("解密后1024：", string(dataRSAEncode))
	t.Log("=================================")

	dataRSA, errRSA = RSAEncryptFP(pathrsapksc12048+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后2048：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = RSADecryptFP(pathrsapksc12048+"/"+privateRSAName, "PRIVATE KEY", dataRSA, RSAPKSC1())
	t.Log("解密后2048：", string(dataRSAEncode))
}

func TestRSACommon_RsaPKSC8EncryptDecrypt(t *testing.T) {
	t.Log("加密前：", contentRSA)
	t.Log("=================================")
	dataRSA, errRSA = RSAEncryptFP(pathrsapksc8256+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后256：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = RSADecryptFP(pathrsapksc8256+"/"+privateRSAName, "PRIVATE KEY", dataRSA, RSAPKSC8())
	t.Log("解密后256：", string(dataRSAEncode))
	t.Log("=================================")

	dataRSA, errRSA = RSAEncryptFP(pathrsapksc8512+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后512：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = RSADecryptFP(pathrsapksc8512+"/"+privateRSAName, "PRIVATE KEY", dataRSA, RSAPKSC8())
	t.Log("解密后512：", string(dataRSAEncode))
	t.Log("=================================")

	dataRSA, errRSA = RSAEncryptFP(pathrsapksc81024+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后1024：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = RSADecryptFP(pathrsapksc81024+"/"+privateRSAName, "PRIVATE KEY", dataRSA, RSAPKSC8())
	t.Log("解密后1024：", string(dataRSAEncode))
	t.Log("=================================")

	dataRSA, errRSA = RSAEncryptFP(pathrsapksc82048+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA))
	if nil != errRSA {
		t.Skip(errRSA)
	}
	t.Log("加密后2048：", hex.EncodeToString(dataRSA))
	dataRSAEncode, errRSA = RSADecryptFP(pathrsapksc82048+"/"+privateRSAName, "PRIVATE KEY", dataRSA, RSAPKSC8())
	t.Log("解密后2048：", string(dataRSAEncode))
}

func TestRSACommon_RsaSign(t *testing.T) {
	t.Log("签名：", contentRSA)
	t.Log("=================================")
	priRSAKeyData, errRSA = ioutil.ReadFile(filepath.Join(pathrsapksc11024, privateRSAName))
	if nil != errRSA {
		t.Error(errRSA)
	}
	if signRSAResult, errRSA = RSASign(priRSAKeyData, []byte(contentRSA), "PRIVATE KEY", crypto.SHA256, RSAPKSC1(), RSASignPKCS()); nil != errRSA {
		t.Skip("签名错误1024：", errRSA)
	} else {
		t.Log("验签1024：", signRSAResult)
		pubRSAKeyData, errRSA = ioutil.ReadFile(filepath.Join(pathrsapksc11024, publicRSAName))
		if nil != errRSA {
			t.Error(errRSA)
		}
		if errRSA = RSAVerify(pubRSAKeyData, []byte(contentRSA), signRSAResult, "PRIVATE KEY", crypto.SHA256, RSASignPKCS()); nil != errRSA {
			t.Skip("验签错误1024：", errRSA)
		} else {
			t.Log("验签通过1024")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = RSASignFPWithPass(pathrsapksc1512+"/"+privateRSAName, "123456", "PRIVATE KEY", []byte(contentRSA), crypto.SHA256, RSAPKSC1(), RSASignPKCS()); nil != errRSA {
		t.Skip("签名错误512：", errRSA)
	} else {
		t.Log("验签512：", signRSAResult)
		if errRSA = RSAVerifyFP(pathrsapksc1512+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA), signRSAResult, crypto.SHA256, RSASignPKCS()); nil != errRSA {
			t.Skip("验签错误512：", errRSA)
		} else {
			t.Log("验签通过512")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = RSASignFP(pathrsapksc11024+"/"+privateRSAName, "PRIVATE KEY", []byte(contentRSA), crypto.SHA512, RSAPKSC1(), RSASignPKCS()); nil != errRSA {
		t.Skip("签名错误1024：", errRSA)
	} else {
		t.Log("验签1024：", signRSAResult)
		if errRSA = RSAVerifyFP(pathrsapksc11024+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA), signRSAResult, crypto.SHA512, RSASignPKCS()); nil != errRSA {
			t.Skip("验签错误1024：", errRSA)
		} else {
			t.Log("验签通过1024")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = RSASignFP(pathrsapksc12048+"/"+privateRSAName, "PRIVATE KEY", []byte(contentRSA), crypto.SHA384, RSAPKSC1(), RSASignPKCS()); nil != errRSA {
		t.Skip("签名错误2048：", errRSA)
	} else {
		t.Log("验签2048：", signRSAResult)
		if errRSA = RSAVerifyFP(pathrsapksc12048+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA), signRSAResult, crypto.SHA384, RSASignPKCS()); nil != errRSA {
			t.Skip("验签错误2048：", errRSA)
		} else {
			t.Log("验签通过2048")
		}
	}
}

func TestRSACommon_RsaSignPSS(t *testing.T) {
	t.Log("签名：", contentRSA)
	t.Log("=================================")
	if signRSAResult, errRSA = RSASignFP(pathrsapksc1512+"/"+privateRSAName, "PRIVATE KEY", []byte(contentRSA), crypto.SHA256, RSAPKSC1(), RSASignPSS()); nil != errRSA {
		t.Skip("签名错误512：", errRSA)
	} else {
		t.Log("验签512：", signRSAResult)
		if errRSA = RSAVerifyFP(pathrsapksc1512+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA), signRSAResult, crypto.SHA256, RSASignPSS()); nil != errRSA {
			t.Skip("验签错误512：", errRSA)
		} else {
			t.Log("验签通过512")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = RSASignFP(pathrsapksc11024+"/"+privateRSAName, "PRIVATE KEY", []byte(contentRSA), crypto.SHA512, RSAPKSC1(), RSASignPSS()); nil != errRSA {
		t.Skip("签名错误1024：", errRSA)
	} else {
		t.Log("验签1024：", signRSAResult)
		if errRSA = RSAVerifyFP(pathrsapksc11024+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA), signRSAResult, crypto.SHA512, RSASignPSS()); nil != errRSA {
			t.Skip("验签错误1024：", errRSA)
		} else {
			t.Log("验签通过1024")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = RSASignFP(pathrsapksc12048+"/"+privateRSAName, "PRIVATE KEY", []byte(contentRSA), crypto.SHA384, RSAPKSC1(), RSASignPSS()); nil != errRSA {
		t.Skip("签名错误2048：", errRSA)
	} else {
		t.Log("验签2048：", signRSAResult)
		if errRSA = RSAVerifyFP(pathrsapksc12048+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA), signRSAResult, crypto.SHA384, RSASignPSS()); nil != errRSA {
			t.Skip("验签错误2048：", errRSA)
		} else {
			t.Log("验签通过2048")
		}
	}
}

func TestRSACommon_RsaPKSC8Sign(t *testing.T) {
	t.Log("签名：", contentRSA)
	t.Log("=================================")
	if signRSAResult, errRSA = RSASignFP(pathrsapksc8512+"/"+privateRSAName, "PRIVATE KEY", []byte(contentRSA), crypto.SHA256, RSAPKSC8(), RSASignPKCS()); nil != errRSA {
		t.Skip("签名错误512：", errRSA)
	} else {
		t.Log("验签512：", signRSAResult)
		if errRSA = RSAVerifyFP(pathrsapksc8512+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA), signRSAResult, crypto.SHA256, RSASignPKCS()); nil != errRSA {
			t.Skip("验签错误512：", errRSA)
		} else {
			t.Log("验签通过512")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = RSASignFP(pathrsapksc81024+"/"+privateRSAName, "PRIVATE KEY", []byte(contentRSA), crypto.SHA512, RSAPKSC8(), RSASignPKCS()); nil != errRSA {
		t.Skip("签名错误1024：", errRSA)
	} else {
		t.Log("验签1024：", signRSAResult)
		if errRSA = RSAVerifyFP(pathrsapksc81024+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA), signRSAResult, crypto.SHA512, RSASignPKCS()); nil != errRSA {
			t.Skip("验签错误1024：", errRSA)
		} else {
			t.Log("验签通过1024")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = RSASignFP(pathrsapksc82048+"/"+privateRSAName, "PRIVATE KEY", []byte(contentRSA), crypto.SHA384, RSAPKSC8(), RSASignPKCS()); nil != errRSA {
		t.Skip("签名错误2048：", errRSA)
	} else {
		t.Log("验签2048：", signRSAResult)
		if errRSA = RSAVerifyFP(pathrsapksc82048+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA), signRSAResult, crypto.SHA384, RSASignPKCS()); nil != errRSA {
			t.Skip("验签错误2048：", errRSA)
		} else {
			t.Log("验签通过2048")
		}
	}
}

func TestRSACommon_RsaPKSC8SignPSS(t *testing.T) {
	t.Log("签名：", contentRSA)
	t.Log("=================================")
	if signRSAResult, errRSA = RSASignFP(pathrsapksc8512+"/"+privateRSAName, "PRIVATE KEY", []byte(contentRSA), crypto.SHA256, RSAPKSC8(), RSASignPSS()); nil != errRSA {
		t.Skip("签名错误512：", errRSA)
	} else {
		t.Log("验签512：", signRSAResult)
		if errRSA = RSAVerifyFP(pathrsapksc8512+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA), signRSAResult, crypto.SHA256, RSASignPSS()); nil != errRSA {
			t.Skip("验签错误512：", errRSA)
		} else {
			t.Log("验签通过512")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = RSASignFP(pathrsapksc81024+"/"+privateRSAName, "PRIVATE KEY", []byte(contentRSA), crypto.SHA512, RSAPKSC8(), RSASignPSS()); nil != errRSA {
		t.Skip("签名错误1024：", errRSA)
	} else {
		t.Log("验签1024：", signRSAResult)
		if errRSA = RSAVerifyFP(pathrsapksc81024+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA), signRSAResult, crypto.SHA512, RSASignPSS()); nil != errRSA {
			t.Skip("验签错误1024：", errRSA)
		} else {
			t.Log("验签通过1024")
		}
	}
	t.Log("=================================")
	if signRSAResult, errRSA = RSASignFP(pathrsapksc82048+"/"+privateRSAName, "PRIVATE KEY", []byte(contentRSA), crypto.SHA384, RSAPKSC8(), RSASignPSS()); nil != errRSA {
		t.Skip("签名错误2048：", errRSA)
	} else {
		t.Log("验签2048：", signRSAResult)
		if errRSA = RSAVerifyFP(pathrsapksc82048+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA), signRSAResult, crypto.SHA384, RSASignPSS()); nil != errRSA {
			t.Skip("验签错误2048：", errRSA)
		} else {
			t.Log("验签通过2048")
		}
	}
}

func TestRSACommon_RsaSign_Fail(t *testing.T) {
	t.Log("签名：", contentRSA)
	t.Log("=================================")
	if signRSAResult, errRSA = RSASignFP(pathrsapksc1256+"/"+privateRSAName, "PRIVATE KEY", []byte(contentRSA), crypto.SHA384, RSAPKSC1(), RSASignPKCS()); nil != errRSA {
		t.Skip("签名错误256：", errRSA)
	} else {
		t.Log("验签256：", signRSAResult)
		if errRSA = RSAVerifyFP(pathrsapksc1256+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA), signRSAResult, crypto.SHA384, RSASignPKCS()); nil != errRSA {
			t.Skip("验签错误256：", errRSA)
		} else {
			t.Log("验签通过256")
		}
	}
}

func TestRSACommon_RsaPKSC8Sign_Fail(t *testing.T) {
	t.Log("签名：", contentRSA)
	t.Log("=================================")
	if signRSAResult, errRSA = RSASignFP(pathrsapksc8256+"/"+privateRSAName, "PRIVATE KEY", []byte(contentRSA), crypto.SHA384, RSAPKSC8(), RSASignPKCS()); nil != errRSA {
		t.Skip("签名错误256：", errRSA)
	} else {
		t.Log("验签256：", signRSAResult)
		if errRSA = RSAVerifyFP(pathrsapksc8256+"/"+publicRSAName, "PRIVATE KEY", []byte(contentRSA), signRSAResult, crypto.SHA384, RSASignPKCS()); nil != errRSA {
			t.Skip("验签错误256：", errRSA)
		} else {
			t.Log("验签通过256")
		}
	}
}

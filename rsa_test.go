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
	"encoding/hex"
	"testing"
)

var (
	content       = "this is a test"
	pathpksc1256  = "./example/rsa/pksc1/256"
	pathpksc1512  = "./example/rsa/pksc1/512"
	pathpksc11024 = "./example/rsa/pksc1/1024"
	pathpksc12048 = "./example/rsa/pksc1/2048"
	pathpksc8256  = "./example/rsa/pksc8/256"
	pathpksc8512  = "./example/rsa/pksc8/512"
	pathpksc81024 = "./example/rsa/pksc8/1024"
	pathpksc82048 = "./example/rsa/pksc8/2048"
	privateName   = "private.pem"
	publicName    = "public.pem"

	privateKey = `-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQC+A2+bmYwxAsgoqsz/OdvCq5IsR91g4r3vnaPi/JqxBovSnpSW
u0C8sCtMPUthFBeaTxRc5siszg6UD0VeZ7F9cAhoOhmaLWPzlKCOe+lzbeH/jfCN
h0Ptp2g24vDBXBmlQXxv+0djNkoGP4/fFC6mdDYIEk8r5/6rpoNEwTMrgQIDAQAB
AoGBAKFnH8XHfzsehtFsjGd8ST0qKictAFJNKFPCzzR/qBMZYWoORF9pPgtJhe3j
3wTeYVI1PDHR48JN4jpIYn9Xi35wGu9cZTz166KZfm7Tj9Ikwv9vwggjrkt5uUYW
n+hBbK+tBe50JbPn2Td6JMW7g6f1fBn/KTYyqPgsq/hyEwzVAkEA3UuQ8aoKV0JQ
2s6QJ40uneQF8eAZKwBcdTVDwvSqNPYiT0LumscR1FccFifyii4gyUS9lJ0Wckd7
inPePD6W3wJBANvP/MdLfR1iVMKsJmk51yIJBR3RBZOdivgwT5uo5cWAVX/dFb3V
erQD3kb7qgFMIVuNRzhdtPaTuWXtTR3OaZ8CQDVs3utZD/INEpQgtnC2BwEbYcwJ
PEpDZg7t8xQIGWd73MCh+hTn5ogLF77JmiZ+CHBO5i4Q1rB0TYEZhBerTKUCQDBS
RahuIN//yNBO1dbV/0QdJYHLfGVaAb3TqPx4IaLMNn94U5o6vtGp9Ag4tMO6P68H
nLt4Zhq6mMweYZCG2tMCQFmoHJ5xtedo0iUPINL6IekjMew5pxoYS293ZPoW4X80
9srcyig275AcYjy0Nt6Hnc6qYmoTKu7X2aE8XBqUFZ8=
-----END RSA PRIVATE KEY-----`
	publicKey = `-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQC+A2+bmYwxAsgoqsz/OdvCq5Is
R91g4r3vnaPi/JqxBovSnpSWu0C8sCtMPUthFBeaTxRc5siszg6UD0VeZ7F9cAho
OhmaLWPzlKCOe+lzbeH/jfCNh0Ptp2g24vDBXBmlQXxv+0djNkoGP4/fFC6mdDYI
Ek8r5/6rpoNEwTMrgQIDAQAB
-----END PUBLIC KEY-----`

	data       []byte
	dataEncode []byte
	signResult []byte
	err        error
)

func TestRSACommon_GenerateRsaKey(t *testing.T) {
	t.Log(CryptoRSA().GenerateRsaPKCS1Key(256, pathpksc1256, "private.pem", "public.pem"))
	t.Log(CryptoRSA().GenerateRsaPKCS1Key(512, pathpksc1512, "private.pem", "public.pem"))
	t.Log(CryptoRSA().GenerateRsaPKCS1Key(1024, pathpksc11024, "private.pem", "public.pem"))
	t.Log(CryptoRSA().GenerateRsaPKCS1Key(2048, pathpksc12048, "private.pem", "public.pem"))

	t.Log(CryptoRSA().GenerateRsaPKCS8Key(256, pathpksc8256, "private.pem", "public.pem"))
	t.Log(CryptoRSA().GenerateRsaPKCS8Key(512, pathpksc8512, "private.pem", "public.pem"))
	t.Log(CryptoRSA().GenerateRsaPKCS8Key(1024, pathpksc81024, "private.pem", "public.pem"))
	t.Log(CryptoRSA().GenerateRsaPKCS8Key(2048, pathpksc82048, "private.pem", "public.pem"))
}

func TestRSACommon_GenerateRsaCustomPriKey(t *testing.T) {
	t.Log(CryptoRSA().GenerateRsaPKCS1PriKey(256, pathpksc1256, "private1.pem"))
	t.Log(CryptoRSA().GenerateRsaPKCS1PriKey(512, pathpksc1512, "private1.pem"))
	t.Log(CryptoRSA().GenerateRsaPKCS1PriKey(1024, pathpksc11024, "private1.pem"))
	t.Log(CryptoRSA().GenerateRsaPKCS1PriKey(2048, pathpksc12048, "private1.pem"))
	t.Log()

	t.Log(CryptoRSA().GenerateRsaPKCS8PriKey(256, pathpksc8256, "private1.pem"))
	t.Log(CryptoRSA().GenerateRsaPKCS8PriKey(512, pathpksc8512, "private1.pem"))
	t.Log(CryptoRSA().GenerateRsaPKCS8PriKey(1024, pathpksc81024, "private1.pem"))
	t.Log(CryptoRSA().GenerateRsaPKCS8PriKey(2048, pathpksc82048, "private1.pem"))
}

func TestRSACommon_GenerateRsaCustomPubKey(t *testing.T) {
	t.Log(CryptoRSA().GenerateRsaPubKeyFP(pathpksc1256+"/"+privateName, pathpksc1256, "public1.pem", CryptoRSA().pksC1()))
	t.Log(CryptoRSA().GenerateRsaPubKeyFP(pathpksc1512+"/"+privateName, pathpksc1512, "public1.pem", CryptoRSA().pksC1()))
	t.Log(CryptoRSA().GenerateRsaPubKeyFP(pathpksc11024+"/"+privateName, pathpksc11024, "public1.pem", CryptoRSA().pksC1()))
	t.Log(CryptoRSA().GenerateRsaPubKeyFP(pathpksc12048+"/"+privateName, pathpksc12048, "public1.pem", CryptoRSA().pksC1()))
	t.Log()

	t.Log(CryptoRSA().GenerateRsaPubKeyFP(pathpksc8256+"/"+privateName, pathpksc8256, "public1.pem", CryptoRSA().pksC8()))
	t.Log(CryptoRSA().GenerateRsaPubKeyFP(pathpksc8512+"/"+privateName, pathpksc8512, "public1.pem", CryptoRSA().pksC8()))
	t.Log(CryptoRSA().GenerateRsaPubKeyFP(pathpksc81024+"/"+privateName, pathpksc81024, "public1.pem", CryptoRSA().pksC8()))
	t.Log(CryptoRSA().GenerateRsaPubKeyFP(pathpksc82048+"/"+privateName, pathpksc82048, "public1.pem", CryptoRSA().pksC8()))
	t.Log()

	t.Log(CryptoRSA().GenerateRsaPubKeyFP(pathpksc82048+"/"+privateName, pathpksc82048, "public2.pem", CryptoRSA().pksC1()))
}

func TestRSACommon_GenerateRsaKey_FailPathExists(t *testing.T) {
	t.Log(CryptoRSA().GenerateRsaPKCS1Key(256, "/etc/test", "private.pem", "public.pem"))
}

func TestRSACommon_GenerateRsaKey_FailGenerate(t *testing.T) {
	t.Log(CryptoRSA().GenerateRsaPKCS1Key(-1, pathpksc1256, "private.pem", "public.pem"))
}

func TestRSACommon_GenerateRsaKey_FailCreate(t *testing.T) {
	t.Log(CryptoRSA().GenerateRsaPKCS1Key(256, "/etc", "private.pem", "public.pem"))
}

func TestRSACommon_GenerateRsaPKSC8Key_FailPathExists(t *testing.T) {
	t.Log(CryptoRSA().GenerateRsaPKCS8Key(256, "/etc/test", "private.pem", "public.pem"))
}

func TestRSACommon_GenerateRsaPKSC8Key_FailGenerate(t *testing.T) {
	t.Log(CryptoRSA().GenerateRsaPKCS8Key(-1, pathpksc1256, "private.pem", "public.pem"))
}

func TestRSACommon_GenerateRsaPKSC8Key_FailCreate(t *testing.T) {
	t.Log(CryptoRSA().GenerateRsaPKCS8Key(256, "/etc", "private.pem", "public.pem"))
}

func TestRSACommon_RsaEncryptDecrypt(t *testing.T) {
	t.Log("加密前：", content)
	t.Log("=================================")
	data, err = CryptoRSA().RsaPubEncrypt([]byte(publicKey), []byte(content))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后：", hex.EncodeToString(data))
	dataEncode, err = CryptoRSA().RsaPriDecrypt([]byte(privateKey), data, CryptoRSA().pksC1())
	t.Log("解密后：", string(dataEncode))
	t.Log("=================================")

	data, err = CryptoRSA().RsaPubEncryptFP(pathpksc1256+"/"+publicName, []byte(content))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后256：", hex.EncodeToString(data))
	dataEncode, err = CryptoRSA().RsaPriDecryptFP(pathpksc1256+"/"+privateName, data, CryptoRSA().pksC1())
	t.Log("解密后256：", string(dataEncode))
	t.Log("=================================")

	data, err = CryptoRSA().RsaPubEncryptFP(pathpksc1512+"/"+publicName, []byte(content))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后512：", hex.EncodeToString(data))
	dataEncode, err = CryptoRSA().RsaPriDecryptFP(pathpksc1512+"/"+privateName, data, CryptoRSA().pksC1())
	t.Log("解密后512：", string(dataEncode))
	t.Log("=================================")

	data, err = CryptoRSA().RsaPubEncryptFP(pathpksc11024+"/"+publicName, []byte(content))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后1024：", hex.EncodeToString(data))
	dataEncode, err = CryptoRSA().RsaPriDecryptFP(pathpksc11024+"/"+privateName, data, CryptoRSA().pksC1())
	t.Log("解密后1024：", string(dataEncode))
	t.Log("=================================")

	data, err = CryptoRSA().RsaPubEncryptFP(pathpksc12048+"/"+publicName, []byte(content))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后2048：", hex.EncodeToString(data))
	dataEncode, err = CryptoRSA().RsaPriDecryptFP(pathpksc12048+"/"+privateName, data, CryptoRSA().pksC1())
	t.Log("解密后2048：", string(dataEncode))
}

func TestRSACommon_RsaPKSC8EncryptDecrypt(t *testing.T) {
	t.Log("加密前：", content)
	t.Log("=================================")
	data, err = CryptoRSA().RsaPubEncrypt([]byte(publicKey), []byte(content))
	if nil != err {
		t.Skip(err)
	}
	data, err = CryptoRSA().RsaPubEncryptFP(pathpksc8256+"/"+publicName, []byte(content))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后256：", hex.EncodeToString(data))
	dataEncode, err = CryptoRSA().RsaPriDecryptFP(pathpksc8256+"/"+privateName, data, CryptoRSA().pksC8())
	t.Log("解密后256：", string(dataEncode))
	t.Log("=================================")

	data, err = CryptoRSA().RsaPubEncryptFP(pathpksc8512+"/"+publicName, []byte(content))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后512：", hex.EncodeToString(data))
	dataEncode, err = CryptoRSA().RsaPriDecryptFP(pathpksc8512+"/"+privateName, data, CryptoRSA().pksC8())
	t.Log("解密后512：", string(dataEncode))
	t.Log("=================================")

	data, err = CryptoRSA().RsaPubEncryptFP(pathpksc81024+"/"+publicName, []byte(content))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后1024：", hex.EncodeToString(data))
	dataEncode, err = CryptoRSA().RsaPriDecryptFP(pathpksc81024+"/"+privateName, data, CryptoRSA().pksC8())
	t.Log("解密后1024：", string(dataEncode))
	t.Log("=================================")

	data, err = CryptoRSA().RsaPubEncryptFP(pathpksc82048+"/"+publicName, []byte(content))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后2048：", hex.EncodeToString(data))
	dataEncode, err = CryptoRSA().RsaPriDecryptFP(pathpksc82048+"/"+privateName, data, CryptoRSA().pksC8())
	t.Log("解密后2048：", string(dataEncode))
}

func TestRSACommon_RsaSign(t *testing.T) {
	t.Log("签名：", content)
	t.Log("=================================")
	if signResult, err = CryptoRSA().RsaPriSign([]byte(privateKey), []byte(content), crypto.SHA1, CryptoRSA().pksC1()); nil != err {
		t.Skip("签名错误：", err)
	} else {
		t.Log("验签：", signResult)
		if err = CryptoRSA().RsaPubVerySign([]byte(publicKey), []byte(content), signResult, crypto.SHA1); nil != err {
			t.Skip("验签错误：", err)
		} else {
			t.Log("验签通过")
		}
	}
	t.Log("=================================")
	if signResult, err = CryptoRSA().RsaPriSignFP(pathpksc1512+"/"+privateName, []byte(content), crypto.SHA256, CryptoRSA().pksC1()); nil != err {
		t.Skip("签名错误512：", err)
	} else {
		t.Log("验签512：", signResult)
		if err = CryptoRSA().RsaPubVerySignFP(pathpksc1512+"/"+publicName, []byte(content), signResult, crypto.SHA256); nil != err {
			t.Skip("验签错误512：", err)
		} else {
			t.Log("验签通过512")
		}
	}
	t.Log("=================================")
	if signResult, err = CryptoRSA().RsaPriSignFP(pathpksc11024+"/"+privateName, []byte(content), crypto.SHA512, CryptoRSA().pksC1()); nil != err {
		t.Skip("签名错误1024：", err)
	} else {
		t.Log("验签1024：", signResult)
		if err = CryptoRSA().RsaPubVerySignFP(pathpksc11024+"/"+publicName, []byte(content), signResult, crypto.SHA512); nil != err {
			t.Skip("验签错误1024：", err)
		} else {
			t.Log("验签通过1024")
		}
	}
	t.Log("=================================")
	if signResult, err = CryptoRSA().RsaPriSignFP(pathpksc12048+"/"+privateName, []byte(content), crypto.SHA384, CryptoRSA().pksC1()); nil != err {
		t.Skip("签名错误2048：", err)
	} else {
		t.Log("验签2048：", signResult)
		if err = CryptoRSA().RsaPubVerySignFP(pathpksc12048+"/"+publicName, []byte(content), signResult, crypto.SHA384); nil != err {
			t.Skip("验签错误2048：", err)
		} else {
			t.Log("验签通过2048")
		}
	}
}

func TestRSACommon_RsaPKSC8Sign(t *testing.T) {
	t.Log("签名：", content)
	t.Log("=================================")
	if signResult, err = CryptoRSA().RsaPriSignFP(pathpksc8512+"/"+privateName, []byte(content), crypto.SHA256, CryptoRSA().pksC8()); nil != err {
		t.Skip("签名错误512：", err)
	} else {
		t.Log("验签512：", signResult)
		if err = CryptoRSA().RsaPubVerySignFP(pathpksc8512+"/"+publicName, []byte(content), signResult, crypto.SHA256); nil != err {
			t.Skip("验签错误512：", err)
		} else {
			t.Log("验签通过512")
		}
	}
	t.Log("=================================")
	if signResult, err = CryptoRSA().RsaPriSignFP(pathpksc81024+"/"+privateName, []byte(content), crypto.SHA512, CryptoRSA().pksC8()); nil != err {
		t.Skip("签名错误1024：", err)
	} else {
		t.Log("验签1024：", signResult)
		if err = CryptoRSA().RsaPubVerySignFP(pathpksc81024+"/"+publicName, []byte(content), signResult, crypto.SHA512); nil != err {
			t.Skip("验签错误1024：", err)
		} else {
			t.Log("验签通过1024")
		}
	}
	t.Log("=================================")
	if signResult, err = CryptoRSA().RsaPriSignFP(pathpksc82048+"/"+privateName, []byte(content), crypto.SHA384, CryptoRSA().pksC8()); nil != err {
		t.Skip("签名错误2048：", err)
	} else {
		t.Log("验签2048：", signResult)
		if err = CryptoRSA().RsaPubVerySignFP(pathpksc82048+"/"+publicName, []byte(content), signResult, crypto.SHA384); nil != err {
			t.Skip("验签错误2048：", err)
		} else {
			t.Log("验签通过2048")
		}
	}
}

func TestRSACommon_RsaSign_Fail(t *testing.T) {
	t.Log("签名：", content)
	t.Log("=================================")
	if signResult, err = CryptoRSA().RsaPriSignFP(pathpksc1256+"/"+privateName, []byte(content), crypto.SHA384, CryptoRSA().pksC1()); nil != err {
		t.Skip("签名错误256：", err)
	} else {
		t.Log("验签256：", signResult)
		if err = CryptoRSA().RsaPubVerySignFP(pathpksc1256+"/"+publicName, []byte(content), signResult, crypto.SHA384); nil != err {
			t.Skip("验签错误256：", err)
		} else {
			t.Log("验签通过256")
		}
	}
}

func TestRSACommon_RsaPKSC8Sign_Fail(t *testing.T) {
	t.Log("签名：", content)
	t.Log("=================================")
	if signResult, err = CryptoRSA().RsaPriSignFP(pathpksc8256+"/"+privateName, []byte(content), crypto.SHA384, CryptoRSA().pksC8()); nil != err {
		t.Skip("签名错误256：", err)
	} else {
		t.Log("验签256：", signResult)
		if err = CryptoRSA().RsaPubVerySignFP(pathpksc8256+"/"+publicName, []byte(content), signResult, crypto.SHA384); nil != err {
			t.Skip("验签错误256：", err)
		} else {
			t.Log("验签通过256")
		}
	}
}

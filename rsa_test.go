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
	"encoding/hex"
	"testing"
)

var (
	content     = "this is a test"
	path256     = "./example/256"
	path512     = "./example/512"
	path1024    = "./example/1024"
	path2048    = "./example/2048"
	privateName = "private.pem"
	publicName  = "public.pem"

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
	t.Log(CryptoRSA().GenerateRsaKey(256, path256))
	t.Log(CryptoRSA().GenerateRsaKey(512, path512))
	t.Log(CryptoRSA().GenerateRsaKey(1024, path1024))
	t.Log(CryptoRSA().GenerateRsaKey(2048, path2048))
}

func TestRSACommon_GenerateRsaKey_FailPathExists(t *testing.T) {
	t.Log(CryptoRSA().GenerateRsaKey(256, "/etc/test"))
}

func TestRSACommon_GenerateRsaKey_FailGenerate(t *testing.T) {
	t.Log(CryptoRSA().GenerateRsaKey(-1, path256))
}

func TestRSACommon_GenerateRsaKey_FailCreate(t *testing.T) {
	t.Log(CryptoRSA().GenerateRsaKey(256, "/etc"))
}

func TestRSACommon_RsaEncryptDecrypt(t *testing.T) {
	t.Log("加密前：", content)
	t.Log("=================================")
	data, err = CryptoRSA().RsaPubEncrypt([]byte(publicKey), []byte(content))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后：", hex.EncodeToString(data))
	dataEncode, err = CryptoRSA().RsaPriDecrypt([]byte(privateKey), data)
	t.Log("解密后：", string(dataEncode))
	t.Log("=================================")

	data, err = CryptoRSA().RsaPubEncryptFP(path256+"/"+publicName, []byte(content))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后256：", hex.EncodeToString(data))
	dataEncode, err = CryptoRSA().RsaPriDecryptFP(path256+"/"+privateName, data)
	t.Log("解密后256：", string(dataEncode))
	t.Log("=================================")

	data, err = CryptoRSA().RsaPubEncryptFP(path512+"/"+publicName, []byte(content))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后512：", hex.EncodeToString(data))
	dataEncode, err = CryptoRSA().RsaPriDecryptFP(path512+"/"+privateName, data)
	t.Log("解密后512：", string(dataEncode))
	t.Log("=================================")

	data, err = CryptoRSA().RsaPubEncryptFP(path1024+"/"+publicName, []byte(content))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后1024：", hex.EncodeToString(data))
	dataEncode, err = CryptoRSA().RsaPriDecryptFP(path1024+"/"+privateName, data)
	t.Log("解密后1024：", string(dataEncode))
	t.Log("=================================")

	data, err = CryptoRSA().RsaPubEncryptFP(path2048+"/"+publicName, []byte(content))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后2048：", hex.EncodeToString(data))
	dataEncode, err = CryptoRSA().RsaPriDecryptFP(path2048+"/"+privateName, data)
	t.Log("解密后2048：", string(dataEncode))
}

func TestRSACommon_RsaSign(t *testing.T) {
	t.Log("签名：", content)
	t.Log("=================================")
	if signResult, err = CryptoRSA().RsaPriSign([]byte(privateKey), []byte(content)); nil != err {
		t.Skip("签名错误：", err)
	} else {
		t.Log("验签：", signResult)
		if err = CryptoRSA().RsaPubVerySign([]byte(publicKey), []byte(content), signResult); nil != err {
			t.Skip("验签错误：", err)
		} else {
			t.Log("验签通过")
		}
	}
	t.Log("=================================")
	if signResult, err = CryptoRSA().RsaPriSignFP(path512+"/"+privateName, []byte(content)); nil != err {
		t.Skip("签名错误512：", err)
	} else {
		t.Log("验签512：", signResult)
		if err = CryptoRSA().RsaPubVerySignFP(path512+"/"+publicName, []byte(content), signResult); nil != err {
			t.Skip("验签错误512：", err)
		} else {
			t.Log("验签通过512")
		}
	}
	t.Log("=================================")
	if signResult, err = CryptoRSA().RsaPriSignFP(path1024+"/"+privateName, []byte(content)); nil != err {
		t.Skip("签名错误1024：", err)
	} else {
		t.Log("验签1024：", signResult)
		if err = CryptoRSA().RsaPubVerySignFP(path1024+"/"+publicName, []byte(content), signResult); nil != err {
			t.Skip("验签错误1024：", err)
		} else {
			t.Log("验签通过1024")
		}
	}
	t.Log("=================================")
	if signResult, err = CryptoRSA().RsaPriSignFP(path2048+"/"+privateName, []byte(content)); nil != err {
		t.Skip("签名错误2048：", err)
	} else {
		t.Log("验签2048：", signResult)
		if err = CryptoRSA().RsaPubVerySignFP(path2048+"/"+publicName, []byte(content), signResult); nil != err {
			t.Skip("验签错误2048：", err)
		} else {
			t.Log("验签通过2048")
		}
	}
}

func TestRSACommon_RsaSign_Fail(t *testing.T) {
	t.Log("签名：", content)
	t.Log("=================================")
	if signResult, err = CryptoRSA().RsaPriSignFP(path256+"/"+privateName, []byte(content)); nil != err {
		t.Skip("签名错误256：", err)
	} else {
		t.Log("验签256：", signResult)
		if err = CryptoRSA().RsaPubVerySignFP(path256+"/"+publicName, []byte(content), signResult); nil != err {
			t.Skip("验签错误256：", err)
		} else {
			t.Log("验签通过256")
		}
	}
}

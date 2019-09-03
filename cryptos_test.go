package common

import (
	"strings"
	"testing"
)

var (
	data        = "this is a test"
	path256     = "./example/256"
	path512     = "./example/512"
	path1024    = "./example/1024"
	path2048    = "./example/2048"
	privateName = "private.pem"
	publicName  = "public.pem"

	bs           []byte
	err          error
	strAfter256  string
	strAfter512  string
	strAfter1024 string
	strAfter2048 string
)

func TestHash(t *testing.T) {
	t.Log("------------- mad5 -------------")
	t.Log(CryptoHash().MD5("haha"))
	t.Log(CryptoHash().MD5("haha"))
	t.Log()
	t.Log("------------- mad516 -------------")
	t.Log(CryptoHash().MD516("haha"))
	t.Log(CryptoHash().MD516("haha"))
	t.Log()
	t.Log("------------- sha1 -------------")
	t.Log(CryptoHash().Sha1("haha"))
	t.Log(CryptoHash().Sha1("haha"))
	t.Log()
	t.Log("------------- sha224 -------------")
	t.Log(CryptoHash().Sha224("haha"))
	t.Log(CryptoHash().Sha224("haha"))
	t.Log()
	t.Log("------------- sha256 -------------")
	t.Log(CryptoHash().Sha256("haha"))
	t.Log(CryptoHash().Sha256("haha"))
	t.Log()
	t.Log("------------- sha384 -------------")
	t.Log(CryptoHash().Sha384("haha"))
	t.Log(CryptoHash().Sha384("haha"))
	t.Log()
	t.Log("------------- sha512 -------------")
	t.Log(CryptoHash().Sha512("haha"))
	t.Log(CryptoHash().Sha512("haha"))
	t.Log()
}

func TestRsa(t *testing.T) {
	t.Log(CryptoRSA().GenerateRsaKey(256, path256))
	t.Log(CryptoRSA().GenerateRsaKey(512, path512))
	t.Log(CryptoRSA().GenerateRsaKey(1024, path1024))
	t.Log(CryptoRSA().GenerateRsaKey(2048, path2048))
}

func TestRsaCrypt(t *testing.T) {
	t.Log("加密前：", data)
	strAfter256, err = CryptoRSA().EncryptRsaPub1(keyPath(path256, publicName), []byte(data))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后：", strAfter256)
	bs, err = CryptoRSA().DecryptRsaPri1(keyPath(path256, privateName), strAfter256)
	t.Log("解密后：", string(bs))
	t.Log("=================================")

	t.Log("加密前：", data)
	strAfter512, err = CryptoRSA().EncryptRsaPub1(keyPath(path512, publicName), []byte(data))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后：", strAfter512)
	bs, err = CryptoRSA().DecryptRsaPri1(keyPath(path512, privateName), strAfter512)
	t.Log("解密后：", string(bs))
	t.Log("=================================")

	t.Log("加密前：", data)
	strAfter1024, err = CryptoRSA().EncryptRsaPub1(keyPath(path1024, publicName), []byte(data))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后：", strAfter1024)
	bs, err = CryptoRSA().DecryptRsaPri1(keyPath(path1024, privateName), strAfter1024)
	t.Log("解密后：", string(bs))
	t.Log("=================================")

	t.Log("加密前：", data)
	strAfter2048, err = CryptoRSA().EncryptRsaPub1(keyPath(path2048, publicName), []byte(data))
	if nil != err {
		t.Skip(err)
	}
	t.Log("加密后：", strAfter2048)
	bs, err = CryptoRSA().DecryptRsaPri1(keyPath(path2048, privateName), strAfter2048)
	t.Log("解密后：", string(bs))
	t.Log("=================================")
}

func TestRsaSign(t *testing.T) {
	var signResult string

	t.Log("签名：", data)
	if signResult, err = CryptoRSA().RsaSign1(keyPath(path256, privateName), data); nil != err {
		t.Skip(err)
	} else {
		t.Log("验签：", signResult)
		if err = CryptoRSA().RsaVerySign1(keyPath(path256, publicName), data, signResult); nil != err {
			t.Skip(err)
		} else {
			t.Log("验签通过")
		}
	}
}

func keyPath(root, key string) string {
	return strings.Join([]string{root, key}, "/")
}

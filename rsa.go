package gnomon

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

type rsaCommon struct{}

// GenerateRsaKey RSA公钥私钥产生
//
// bits 指定生成位大小
//
// path 指定公私钥所在生成目录
func (r *rsaCommon) GenerateRsaKey(bits int, path string) (err error) {
	var (
		privateKey *rsa.PrivateKey
		publicKey  *rsa.PublicKey
		fileIO     *os.File
		derPkix    []byte
	)
	// 创建公私钥生成目录
	if !File().PathExists(path) {
		if err = os.MkdirAll(path, os.ModePerm); nil != err {
			return
		}
	}
	// 生成私钥文件
	if privateKey, err = rsa.GenerateKey(rand.Reader, bits); nil != err {
		return
	}
	// 将私钥转换为ASN.1 DER编码的形式
	derStream := x509.MarshalPKCS1PrivateKey(privateKey)
	// block表示PEM编码的结构
	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: derStream,
	}
	if fileIO, err = os.Create(strings.Join([]string{path, "private.pem"}, "/")); nil != err {
		return
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, block); nil != err {
		return
	}
	// 生成公钥文件
	publicKey = &privateKey.PublicKey
	// 将公钥序列化为der编码的PKIX格式
	if derPkix, err = x509.MarshalPKIXPublicKey(publicKey); nil != err {
		return
	}
	block = &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: derPkix,
	}
	if fileIO, err = os.Create(strings.Join([]string{path, "public.pem"}, "/")); nil != err {
		return
	}
	// 将block的PEM编码写入fileIO
	if err = pem.Encode(fileIO, block); nil != err {
		return
	}
	return
}

// EncryptRsaPub1 公钥加密
func (r *rsaCommon) EncryptRsaPub1(publicKeyPath string, data []byte) (string, error) {
	if bs, err := r.EncryptRsaPub3(publicKeyPath, data); nil != err {
		return "", err
	} else {
		return hex.EncodeToString(bs), nil
	}
}

// EncryptRsaPub2 公钥加密
func (r *rsaCommon) EncryptRsaPub2(publicKey, data []byte) (string, error) {
	if bs, err := r.EncryptRsaPub4(publicKey, []byte(data)); nil != err {
		return "", err
	} else {
		return hex.EncodeToString(bs), nil
	}
}

// EncryptRsaPub3 公钥加密
func (r *rsaCommon) EncryptRsaPub3(publicKeyPath string, data []byte) (bs []byte, err error) {
	if bs, err = ioutil.ReadFile(publicKeyPath); nil != err {
		return
	} else {
		return r.EncryptRsaPub4(bs, data)
	}
}

// EncryptRsaPub4 公钥加密
func (r *rsaCommon) EncryptRsaPub4(publicKey, data []byte) ([]byte, error) {
	// 解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, data)
}

// DecryptRsaPri1 私钥解密
func (r *rsaCommon) DecryptRsaPri1(privateKeyPath, data string) ([]byte, error) {
	if bs, err := hex.DecodeString(data); nil != err {
		return bs, err
	} else {
		return r.DecryptRsaPri3(privateKeyPath, bs)
	}
}

// DecryptRsaPri2 私钥解密
func (r *rsaCommon) DecryptRsaPri2(privateKey []byte, data string) (bs []byte, err error) {
	if bs, err := hex.DecodeString(data); nil != err {
		return bs, err
	} else {
		return r.DecryptRsaPri4(privateKey, bs)
	}
}

// DecryptRsaPri3 私钥解密
func (r *rsaCommon) DecryptRsaPri3(privateKeyPath string, data []byte) ([]byte, error) {
	if bs, err := ioutil.ReadFile(privateKeyPath); nil != err {
		return bs, err
	} else {
		return r.DecryptRsaPri4(bs, data)
	}
}

// DecryptRsaPri4 私钥解密
func (r *rsaCommon) DecryptRsaPri4(privateKey, data []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error")
	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, data)
}

// RsaSign1 签名：采用sha1算法进行签名并输出为hex格式（私钥PKCS8格式）
func (r *rsaCommon) RsaSign1(privateKeyPath, data string) (string, error) {
	return r.RsaSign3(privateKeyPath, []byte(data))
}

// RsaSign2 签名：采用sha1算法进行签名并输出为hex格式（私钥PKCS8格式）
func (r *rsaCommon) RsaSign2(privateKey []byte, data string) (string, error) {
	return r.RsaSign4(privateKey, []byte(data))
}

// RsaSign3 签名：采用sha1算法进行签名并输出为hex格式（私钥PKCS8格式）
func (r *rsaCommon) RsaSign3(privateKeyPath string, data []byte) (string, error) {
	if bs, err := ioutil.ReadFile(privateKeyPath); nil != err {
		return "", err
	} else {
		return r.RsaSign4(bs, data)
	}
}

// RsaSign4 签名：采用sha1算法进行签名并输出为hex格式（私钥PKCS8格式）
func (r *rsaCommon) RsaSign4(privateKey, data []byte) (string, error) {
	priv, err := x509.ParsePKCS8PrivateKey(privateKey)
	if err != nil {
		Log().Error("RsaSign4 ParsePKCS8PrivateKey Error", Log().Err(err))
		return "", err
	}
	h := sha1.New()
	h.Write([]byte(data))
	hash := h.Sum(nil)
	signature, err := rsa.SignPKCS1v15(rand.Reader, priv.(*rsa.PrivateKey), crypto.SHA1, hash[:])
	if err != nil {
		Log().Error("RsaSign4 Error from signing", Log().Err(err))
		return "", err
	}
	out := hex.EncodeToString(signature)
	return out, nil
}

// RsaVerySign1 验签：对采用sha1算法进行签名后转base64格式的数据进行验签
func (r *rsaCommon) RsaVerySign1(publicKeyPath, data, signData string) error {
	return r.RsaVerySign3(publicKeyPath, []byte(data), signData)
}

// RsaVerySign2 验签：对采用sha1算法进行签名后转base64格式的数据进行验签
func (r *rsaCommon) RsaVerySign2(publicKey []byte, data, signData string) error {
	return r.RsaVerySign4(publicKey, []byte(data), signData)
}

// RsaVerySign3 验签：对采用sha1算法进行签名后转base64格式的数据进行验签
func (r *rsaCommon) RsaVerySign3(publicKeyPath string, data []byte, signData string) error {
	if bs, err := ioutil.ReadFile(publicKeyPath); nil != err {
		return err
	} else {
		return r.RsaVerySign4(bs, data, signData)
	}
}

// RsaVerySign4 验签：对采用sha1算法进行签名后转base64格式的数据进行验签
func (r *rsaCommon) RsaVerySign4(publicKey, data []byte, signData string) error {
	sign, err := base64.StdEncoding.DecodeString(signData)
	if err != nil {
		return err
	}
	pub, err := x509.ParsePKIXPublicKey(publicKey)
	if err != nil {
		return err
	}
	hash := sha1.New()
	hash.Write(data)
	return rsa.VerifyPKCS1v15(pub.(*rsa.PublicKey), crypto.SHA1, hash.Sum(nil), sign)
}

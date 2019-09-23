/*
 *
 *  * Copyright (c) 2019. aberic - All Rights Reserved.
 *  *
 *  * Licensed under the Apache License, Version 2.0 (the "License");
 *  * you may not use this file except in compliance with the License.
 *  * You may obtain a copy of the License at
 *  * http://www.apache.org/licenses/LICENSE-2.0
 *  * Unless required by applicable law or agreed to in writing, software
 *  * distributed under the License is distributed on an "AS IS" BASIS,
 *  * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  * See the License for the specific language governing permissions and
 *  * limitations under the License.
 *
 *
 */

package gnomon

import (
	"crypto/elliptic"
	"crypto/x509"
	"crypto/x509/pkix"
	"io/ioutil"
	"net"
	"path/filepath"
	"testing"
	"time"
)

var (
	pathcarsapksc1512  = "./tmp/example/ca/pksc1/512"
	pathcarsapksc11024 = "./tmp/example/ca/pksc1/1024"
	//pathcarsapksc12048 = "./tmp/example/ca/pksc1/2048"

	//pathcarsapksc8512 = "./tmp/example/ca/pksc8/512"
	pathcarsapksc81024 = "./tmp/example/ca/pksc8/1024"
	pathcarsapksc82048 = "./tmp/example/ca/pksc8/2048"
	//
	pathcaeccpemp224 = "./tmp/example/ca/pemp224"
	pathcaeccpemp256 = "./tmp/example/ca/pemp256"
	pathcaeccpemp384 = "./tmp/example/ca/pemp384"
	pathcaeccpemp521 = "./tmp/example/ca/pemp521"

	priData []byte

	caPriKeyFileName             = "rootCA.key" // ca 私钥
	caCertificateRequestFileName = "rootCA.csr" // 证书签名请求文件
	caCertificateFileName        = "rootCA.crt"

	errCA error
)

var CAMockSubject = pkix.Name{
	Country:            []string{"CN"},
	Organization:       []string{"Gnomon"},
	OrganizationalUnit: []string{"GnomonRD"},
	Locality:           []string{"Beijing"},
	Province:           []string{"Beijing"},
	CommonName:         "aberic.cn",
}

func TestCACommon_GenerateRSAPKCS1PrivateKey(t *testing.T) {
	if _, errCA = CryptoRSA().GeneratePKCS1PriKey(512, pathcarsapksc1512, caPriKeyFileName); nil != errCA {
		t.Error(errCA)
	}
	priData, errCA = ioutil.ReadFile(filepath.Join(pathcarsapksc1512, caPriKeyFileName))
	if nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateRSACertificateRequest(&CertRequest{
		PrivateKeyData:             priData,
		CertificateRequestFilePath: filepath.Join(pathcarsapksc1512, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA256WithRSAPSS,
		Subject:                    CAMockSubject,
	}, CryptoRSA().PKSC1()); nil != errCA {
		t.Error(errCA)
	}

	if _, errCA = CryptoRSA().GeneratePKCS1PriKeyWithPass(1024, pathcarsapksc11024, caPriKeyFileName, "123456"); nil != errCA {
		t.Error(errCA)
	}
	priData, errCA = ioutil.ReadFile(filepath.Join(pathcarsapksc11024, caPriKeyFileName))
	if nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateRSACertificateRequestWithPass(&CertRequest{
		PrivateKeyData:             priData,
		CertificateRequestFilePath: filepath.Join(pathcarsapksc11024, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA384WithRSAPSS,
		Subject:                    CAMockSubject,
	}, "123456", CryptoRSA().PKSC1()); nil != errCA {
		t.Error(errCA)
	}
}

func TestCACommon_GenerateRSAPKCS1PrivateKeyFP(t *testing.T) {
	if _, errCA = CryptoRSA().GeneratePKCS1PriKey(512, pathcarsapksc1512, caPriKeyFileName); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateRSACertificateRequestFP(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcarsapksc1512, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcarsapksc1512, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA256WithRSAPSS,
		Subject:                    CAMockSubject,
	}, CryptoRSA().PKSC1()); nil != errCA {
		t.Error(errCA)
	}

	if _, errCA = CryptoRSA().GeneratePKCS1PriKeyWithPass(1024, pathcarsapksc11024, caPriKeyFileName, "123456"); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateRSACertificateRequestFPWithPass(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcarsapksc11024, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcarsapksc11024, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA384WithRSAPSS,
		Subject:                    CAMockSubject,
	}, "123456", CryptoRSA().PKSC1()); nil != errCA {
		t.Error(errCA)
	}
}

func TestCACommon_GenerateRSAPKCS8PrivateKeyFP(t *testing.T) {
	if priRSAKey, errCA = CryptoRSA().GeneratePKCS8PriKey(1024, pathcarsapksc81024, caPriKeyFileName); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateRSACertificateRequestFP(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcarsapksc81024, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcarsapksc81024, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA384WithRSAPSS,
		Subject:                    CAMockSubject,
		EmailAddresses:             []string{"test@test.com"},
		IPAddresses:                []net.IP{net.ParseIP("192.168.1.59")},
	}, CryptoRSA().PKSC8()); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateCertificate(&Cert{
		CertificateFilePath:   filepath.Join(pathcarsapksc81024, caCertificateFileName),
		Subject:               CAMockSubject,
		PrivateKey:            priRSAKey,
		PublicKey:             priRSAKey.Public(),
		NotAfterDays:          time.Now(),
		NotBeforeDays:         time.Now().Add(5000 * 24 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, //证书用途(客户端认证，数据加密)
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
		SignatureAlgorithm:    x509.SHA384WithRSAPSS,
	}); nil != errCA {
		t.Error(errCA)
	}

	if _, errCA = CryptoRSA().GeneratePKCS8PriKeyWithPass(2048, pathcarsapksc82048, caPriKeyFileName, "123456"); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateRSACertificateRequestFPWithPass(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcarsapksc82048, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcarsapksc82048, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA512WithRSAPSS,
		Subject:                    CAMockSubject,
	}, "123456", CryptoRSA().PKSC8()); nil != errCA {
		t.Error(errCA)
	}
}

func TestCACommon_GenerateECCPrivateKey(t *testing.T) {
	if errCA = CryptoECC().GeneratePemPriKey(pathcaeccpemp224, caPriKeyFileName, elliptic.P224()); nil != errCA {
		t.Error(errCA)
	}
	priData, errECC = ioutil.ReadFile(filepath.Join(pathcaeccpemp224, caPriKeyFileName))
	if nil != errECC {
		t.Error(errECC)
	}
	if _, errCA = CA().GenerateECCCertificateRequest(&CertRequest{
		PrivateKeyData:             priData,
		CertificateRequestFilePath: filepath.Join(pathcaeccpemp224, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.ECDSAWithSHA256,
		Subject:                    CAMockSubject,
		EmailAddresses:             []string{"test@test.com"},
		IPAddresses:                []net.IP{net.ParseIP("192.168.1.59")},
	}); nil != errCA {
		t.Error(errCA)
	}
	if priKeyP224, errCA = CryptoECC().LoadPriPem(priData); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateCertificate(&Cert{
		CertificateFilePath:   filepath.Join(pathcaeccpemp224, caCertificateFileName),
		Subject:               CAMockSubject,
		PrivateKey:            priKeyP224,
		PublicKey:             priKeyP224.Public(),
		NotAfterDays:          time.Now(),
		NotBeforeDays:         time.Now().Add(5000 * 24 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
	}); nil != errCA {
		t.Error(errCA)
	}

	if errCA = CryptoECC().GeneratePemPriKey(pathcaeccpemp256, caPriKeyFileName, elliptic.P256()); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateECCCertificateRequestFP(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcaeccpemp256, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcaeccpemp256, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.ECDSAWithSHA256,
		Subject:                    CAMockSubject,
		EmailAddresses:             []string{"test@test.com"},
		IPAddresses:                []net.IP{net.ParseIP("192.168.1.59")},
	}); nil != errCA {
		t.Error(errCA)
	}
	if priKeyP256, errCA = CryptoECC().LoadPriPemFP(filepath.Join(pathcaeccpemp256, caPriKeyFileName)); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateCertificate(&Cert{
		CertificateFilePath:   filepath.Join(pathcaeccpemp256, caCertificateFileName),
		Subject:               CAMockSubject,
		PrivateKey:            priKeyP256,
		PublicKey:             priKeyP256.Public(),
		NotAfterDays:          time.Now(),
		NotBeforeDays:         time.Now().Add(5000 * 24 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
	}); nil != errCA {
		t.Error(errCA)
	}

	if errCA = CryptoECC().GeneratePemPriKeyWithPass(pathcaeccpemp384, caPriKeyFileName, "123456", elliptic.P384()); nil != errCA {
		t.Error(errCA)
	}
	priData, errECC = ioutil.ReadFile(filepath.Join(pathcaeccpemp384, caPriKeyFileName))
	if nil != errECC {
		t.Error(errECC)
	}
	if _, errCA = CA().GenerateECCCertificateRequestWithPass(&CertRequest{
		PrivateKeyData:             priData,
		CertificateRequestFilePath: filepath.Join(pathcaeccpemp384, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.ECDSAWithSHA384,
		Subject:                    CAMockSubject,
	}, "123456"); nil != errCA {
		t.Error(errCA)
	}
	if priKeyP384, errCA = CryptoECC().LoadPriPemFPWithPass(filepath.Join(pathcaeccpemp384, caPriKeyFileName), "123456"); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateCertificate(&Cert{
		CertificateFilePath:   filepath.Join(pathcaeccpemp384, caCertificateFileName),
		Subject:               CAMockSubject,
		PrivateKey:            priKeyP384,
		PublicKey:             priKeyP384.Public(),
		NotAfterDays:          time.Now(),
		NotBeforeDays:         time.Now().Add(5000 * 24 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
		SignatureAlgorithm:    x509.ECDSAWithSHA384,
	}); nil != errCA {
		t.Error(errCA)
	}

	if errCA = CryptoECC().GeneratePemPriKeyWithPass(pathcaeccpemp521, caPriKeyFileName, "123456", elliptic.P521()); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateECCCertificateRequestFPWithPass(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcaeccpemp521, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcaeccpemp521, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.ECDSAWithSHA512,
		Subject:                    CAMockSubject,
	}, "123456"); nil != errCA {
		t.Error(errCA)
	}
	if priKeyP521, errCA = CryptoECC().LoadPriPemFPWithPass(filepath.Join(pathcaeccpemp521, caPriKeyFileName), "123456"); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateCertificate(&Cert{
		CertificateFilePath:   filepath.Join(pathcaeccpemp521, caCertificateFileName),
		Subject:               CAMockSubject,
		PrivateKey:            priKeyP521,
		PublicKey:             priKeyP521.Public(),
		NotAfterDays:          time.Now(),
		NotBeforeDays:         time.Now().Add(5000 * 24 * time.Hour),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
		SignatureAlgorithm:    x509.ECDSAWithSHA512,
	}); nil != errCA {
		t.Error(errCA)
	}
}

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
	"path/filepath"
	"testing"
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
	//caCertificateFileName        = "rootCA.crt"

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
	if _, errCA = CryptoRSA().GeneratePKCS8PriKey(1024, pathcarsapksc81024, caPriKeyFileName); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateRSACertificateRequestFP(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcarsapksc81024, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcarsapksc81024, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA384WithRSAPSS,
		Subject:                    CAMockSubject,
	}, CryptoRSA().PKSC8()); nil != errCA {
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
	if _, errCA = CA().GenerateECCCertificateRequestFP(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcaeccpemp224, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcaeccpemp224, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.ECDSAWithSHA256,
		Subject:                    CAMockSubject,
	}); nil != errCA {
		t.Error(errCA)
	}
	//if _, errCA = CA().GenerateCertificate(&Cert{
	//	CertificateFilePath:filepath.Join(pathcaeccpemp224, caCertificateFileName),
	//	Subject:CAMockSubject,
	//});nil!=errCA {
	//	t.Error(errCA)
	//}

	if errCA = CryptoECC().GeneratePemPriKey(pathcaeccpemp256, caPriKeyFileName, elliptic.P256()); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateECCCertificateRequestFP(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcaeccpemp256, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcaeccpemp256, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.ECDSAWithSHA256,
		Subject:                    CAMockSubject,
	}); nil != errCA {
		t.Error(errCA)
	}

	if errCA = CryptoECC().GeneratePemPriKeyWithPass(pathcaeccpemp384, caPriKeyFileName, "123456", elliptic.P384()); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CA().GenerateECCCertificateRequestFPWithPass(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcaeccpemp384, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcaeccpemp384, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.ECDSAWithSHA384,
		Subject:                    CAMockSubject,
	}, "123456"); nil != errCA {
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
}

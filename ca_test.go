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
	//pathcarsapksc1256  = "./tmp/example/ca/pksc1/256"
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

	pathcarsapksc1fabric2048 = "./tmp/example/ca/fabric/pksc1/2048"
	pathcaeccpempfabric384   = "./tmp/example/ca/fabric/pemp384"

	parentCert *x509.Certificate
	priData    []byte
	certData   []byte

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
	//CommonName:         "aberic.cn",
	CommonName: "localhost",
}

func TestCACommon_GenerateRSAPKCS1PrivateKey(t *testing.T) {
	if _, errCA = RSAGeneratePriKey(512, pathcarsapksc1512, caPriKeyFileName, "PRIVATE KEY", RSAPKSC1()); nil != errCA {
		t.Error(errCA)
	}
	priData, errCA = ioutil.ReadFile(filepath.Join(pathcarsapksc1512, caPriKeyFileName))
	if nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateRSACertificateRequest(&CertRequest{
		PrivateKeyData:             priData,
		CertificateRequestFilePath: filepath.Join(pathcarsapksc1512, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA256WithRSAPSS,
		Subject:                    CAMockSubject,
	}, "PRIVATE KEY", RSAPKSC1()); nil != errCA {
		t.Error(errCA)
	}

	if _, errCA = RSAGeneratePriKeyWithPass(1024, pathcarsapksc11024, caPriKeyFileName, "123456", "PRIVATE KEY", x509.PEMCipher3DES, RSAPKSC1()); nil != errCA {
		t.Error(errCA)
	}
	priData, errCA = ioutil.ReadFile(filepath.Join(pathcarsapksc11024, caPriKeyFileName))
	if nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateRSACertificateRequestWithPass(&CertRequest{
		PrivateKeyData:             priData,
		CertificateRequestFilePath: filepath.Join(pathcarsapksc11024, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA384WithRSAPSS,
		Subject:                    CAMockSubject,
	}, "123456", "PRIVATE KEY", RSAPKSC1()); nil != errCA {
		t.Error(errCA)
	}
}

func TestCACommon_GenerateRSAPKCS1PrivateKeyFP(t *testing.T) {
	if _, errCA = RSAGeneratePriKey(512, pathcarsapksc1512, caPriKeyFileName, "PRIVATE KEY", RSAPKSC1()); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateRSACertificateRequestFP(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcarsapksc1512, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcarsapksc1512, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA256WithRSAPSS,
		Subject:                    CAMockSubject,
	}, "PRIVATE KEY", RSAPKSC1()); nil != errCA {
		t.Error(errCA)
	}

	if _, errCA = RSAGeneratePriKeyWithPass(1024, pathcarsapksc11024, caPriKeyFileName, "123456", "PRIVATE KEY", x509.PEMCipher3DES, RSAPKSC1()); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateRSACertificateRequestFPWithPass(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcarsapksc11024, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcarsapksc11024, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA384WithRSAPSS,
		Subject:                    CAMockSubject,
	}, "123456", "PRIVATE KEY", RSAPKSC1()); nil != errCA {
		t.Error(errCA)
	}
}

func TestCACommon_GenerateRSAPKCS1PrivateKeyFPFabricCA(t *testing.T) {
	if priRSAKey, errCA = RSAGeneratePriKey(2048, pathcarsapksc1fabric2048, caPriKeyFileName, "PRIVATE KEY", RSAPKSC1()); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateRSACertificateRequestFP(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcarsapksc1fabric2048, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcarsapksc1fabric2048, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA256WithRSA,
		Subject:                    CAMockSubject,
		EmailAddresses:             []string{"test@test.com"},
		IPAddresses:                []net.IP{net.ParseIP("localhost"), net.ParseIP("127.0.0.1")},
	}, "PRIVATE KEY", RSAPKSC1()); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateCertificateSelf(&CertSelf{
		CertificateFilePath:   filepath.Join(pathcarsapksc1fabric2048, caCertificateFileName),
		Subject:               CAMockSubject,
		ParentPrivateKey:      priRSAKey,
		PublicKey:             priRSAKey.Public(),
		NotAfterDays:          time.Now().Add(5000 * 24 * time.Hour),
		NotBeforeDays:         time.Now(),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, //证书用途(客户端认证，数据加密)
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
		SignatureAlgorithm:    x509.SHA256WithRSA,
	}); nil != errCA {
		t.Error(errCA)
	}

	if errCA = ECCGeneratePemPriKey(pathcaeccpempfabric384, caPriKeyFileName, "PRIVATE KEY", elliptic.P384()); nil != errCA {
		t.Error(errCA)
	}
	priData, errECC = ioutil.ReadFile(filepath.Join(pathcaeccpempfabric384, caPriKeyFileName))
	if nil != errECC {
		t.Error(errECC)
	}
	if _, errCA = CAGenerateECCCertificateRequest(&CertRequest{
		PrivateKeyData:             priData,
		CertificateRequestFilePath: filepath.Join(pathcaeccpempfabric384, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.ECDSAWithSHA256,
		Subject:                    CAMockSubject,
	}, "PRIVATE KEY"); nil != errCA {
		t.Error(errCA)
	}
	if priKeyP384, errCA = ECCLoadPriPemFP(filepath.Join(pathcaeccpempfabric384, caPriKeyFileName), "PRIVATE KEY"); nil != errCA {
		t.Error(errCA)
	}
	if certData, errCA = CAGenerateCertificateSelf(&CertSelf{
		CertificateFilePath:   filepath.Join(pathcaeccpempfabric384, caCertificateFileName),
		Subject:               CAMockSubject,
		ParentPrivateKey:      priKeyP384,
		PublicKey:             priKeyP384.Public(),
		NotAfterDays:          time.Now().Add(5000 * 24 * time.Hour),
		NotBeforeDays:         time.Now(),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
	}); nil != errCA {
		t.Error(errCA)
	}

}

func TestCACommon_GenerateRSAPKCS8PrivateKeyFP(t *testing.T) {
	if priRSAKey, errCA = RSAGeneratePriKey(1024, pathcarsapksc81024, caPriKeyFileName, "PRIVATE KEY", RSAPKSC8()); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateRSACertificateRequestFP(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcarsapksc81024, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcarsapksc81024, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA384WithRSAPSS,
		Subject:                    CAMockSubject,
		EmailAddresses:             []string{"test@test.com"},
		IPAddresses:                []net.IP{net.ParseIP("192.168.1.59")},
	}, "PRIVATE KEY", RSAPKSC8()); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateCertificateSelf(&CertSelf{
		CertificateFilePath:   filepath.Join(pathcarsapksc81024, caCertificateFileName),
		Subject:               CAMockSubject,
		ParentPrivateKey:      priRSAKey,
		PublicKey:             priRSAKey.Public(),
		NotAfterDays:          time.Now().Add(5000 * 24 * time.Hour),
		NotBeforeDays:         time.Now(),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth}, //证书用途(客户端认证，数据加密)
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
		SignatureAlgorithm:    x509.SHA384WithRSAPSS,
	}); nil != errCA {
		t.Error(errCA)
	}

	if _, errCA = RSAGeneratePriKeyWithPass(2048, pathcarsapksc82048, caPriKeyFileName, "123456", "PRIVATE KEY", x509.PEMCipher3DES, RSAPKSC8()); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateRSACertificateRequestFPWithPass(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcarsapksc82048, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcarsapksc82048, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA512WithRSAPSS,
		Subject:                    CAMockSubject,
	}, "123456", "PRIVATE KEY", RSAPKSC8()); nil != errCA {
		t.Error(errCA)
	}
}

func TestCACommon_GenerateECCPrivateKey(t *testing.T) {
	if errCA = ECCGeneratePemPriKey(pathcaeccpemp224, caPriKeyFileName, "PRIVATE KEY", elliptic.P224()); nil != errCA {
		t.Error(errCA)
	}
	priData, errECC = ioutil.ReadFile(filepath.Join(pathcaeccpemp224, caPriKeyFileName))
	if nil != errECC {
		t.Error(errECC)
	}
	if _, errCA = CAGenerateECCCertificateRequest(&CertRequest{
		PrivateKeyData:             priData,
		CertificateRequestFilePath: filepath.Join(pathcaeccpemp224, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.ECDSAWithSHA256,
		Subject:                    CAMockSubject,
		EmailAddresses:             []string{"test@test.com"},
		IPAddresses:                []net.IP{net.ParseIP("192.168.1.59")},
	}, "PRIVATE KEY"); nil != errCA {
		t.Error(errCA)
	}
	if priKeyP224, errCA = ECCLoadPriPem(priData, "PRIVATE KEY"); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateCertificateSelf(&CertSelf{
		CertificateFilePath:   filepath.Join(pathcaeccpemp224, caCertificateFileName),
		Subject:               CAMockSubject,
		ParentPrivateKey:      priKeyP224,
		PublicKey:             priKeyP224.Public(),
		NotAfterDays:          time.Now().Add(5000 * 24 * time.Hour),
		NotBeforeDays:         time.Now(),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
	}); nil != errCA {
		t.Error(errCA)
	}

	if errCA = ECCGeneratePemPriKey(pathcaeccpemp256, caPriKeyFileName, "PRIVATE KEY", elliptic.P256()); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateECCCertificateRequestFP(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcaeccpemp256, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcaeccpemp256, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.ECDSAWithSHA256,
		Subject:                    CAMockSubject,
		EmailAddresses:             []string{"test@test.com"},
		IPAddresses:                []net.IP{net.ParseIP("192.168.1.59")},
	}, "PRIVATE KEY"); nil != errCA {
		t.Error(errCA)
	}
	if priKeyP256, errCA = ECCLoadPriPemFP(filepath.Join(pathcaeccpemp256, caPriKeyFileName), "PRIVATE KEY"); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateCertificateSelf(&CertSelf{
		CertificateFilePath:   filepath.Join(pathcaeccpemp256, caCertificateFileName),
		Subject:               CAMockSubject,
		ParentPrivateKey:      priKeyP256,
		PublicKey:             priKeyP256.Public(),
		NotAfterDays:          time.Now().Add(5000 * 24 * time.Hour),
		NotBeforeDays:         time.Now(),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
		SignatureAlgorithm:    x509.ECDSAWithSHA256,
	}); nil != errCA {
		t.Error(errCA)
	}

	if errCA = ECCGeneratePemPriKeyWithPass(pathcaeccpemp384, caPriKeyFileName, "123456", "PRIVATE KEY", elliptic.P384()); nil != errCA {
		t.Error(errCA)
	}
	priData, errECC = ioutil.ReadFile(filepath.Join(pathcaeccpemp384, caPriKeyFileName))
	if nil != errECC {
		t.Error(errECC)
	}
	if _, errCA = CAGenerateECCCertificateRequestWithPass(&CertRequest{
		PrivateKeyData:             priData,
		CertificateRequestFilePath: filepath.Join(pathcaeccpemp384, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.ECDSAWithSHA384,
		Subject:                    CAMockSubject,
	}, "123456", "PRIVATE KEY"); nil != errCA {
		t.Error(errCA)
	}
	if priKeyP384, errCA = ECCLoadPriPemFPWithPass(filepath.Join(pathcaeccpemp384, caPriKeyFileName), "123456", "PRIVATE KEY"); nil != errCA {
		t.Error(errCA)
	}
	if certData, errCA = CAGenerateCertificateSelf(&CertSelf{
		CertificateFilePath:   filepath.Join(pathcaeccpemp384, caCertificateFileName),
		Subject:               CAMockSubject,
		ParentPrivateKey:      priKeyP384,
		PublicKey:             priKeyP384.Public(),
		NotAfterDays:          time.Now().Add(5000 * 24 * time.Hour),
		NotBeforeDays:         time.Now(),
		BasicConstraintsValid: true,
		IsCA:                  true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
		SignatureAlgorithm:    x509.ECDSAWithSHA384,
	}); nil != errCA {
		t.Error(errCA)
	}

	if errCA = ECCGeneratePemPriKeyWithPass(pathcaeccpemp521, caPriKeyFileName, "123456", "PRIVATE KEY", elliptic.P521()); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateECCCertificateRequestFPWithPass(&CertRequestFP{
		PrivateKeyFilePath:         filepath.Join(pathcaeccpemp521, caPriKeyFileName),
		CertificateRequestFilePath: filepath.Join(pathcaeccpemp521, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.ECDSAWithSHA512,
		Subject:                    CAMockSubject,
	}, "123456", "PRIVATE KEY"); nil != errCA {
		t.Error(errCA)
	}
	if priKeyP521, errCA = ECCLoadPriPemFPWithPass(filepath.Join(pathcaeccpemp521, caPriKeyFileName), "123456", "PRIVATE KEY"); nil != errCA {
		t.Error(errCA)
	}
	if parentCert, errCA = x509.ParseCertificate(certData); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateCertificate(&Cert{
		ParentCert: parentCert,
		CertSelf: CertSelf{
			CertificateFilePath:   filepath.Join(pathcaeccpemp521, caCertificateFileName),
			Subject:               CAMockSubject,
			ParentPrivateKey:      priKeyP384,
			PublicKey:             priKeyP521.Public(),
			NotAfterDays:          time.Now().Add(5000 * 24 * time.Hour),
			NotBeforeDays:         time.Now(),
			BasicConstraintsValid: true,
			IsCA:                  true,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
			SignatureAlgorithm:    x509.ECDSAWithSHA384,
		},
	}); nil != errCA {
		t.Error(errCA)
	}

	if parentCert, errCA = CALoadCrtFromFP(filepath.Join(pathcaeccpemp384, caCertificateFileName)); nil != errCA {
		t.Error(errCA)
	}
	if _, errCA = CAGenerateCertificate(&Cert{
		ParentCert: parentCert,
		CertSelf: CertSelf{
			CertificateFilePath:   filepath.Join(pathcaeccpemp521, caCertificateFileName),
			Subject:               CAMockSubject,
			ParentPrivateKey:      priKeyP384,
			PublicKey:             priKeyP521.Public(),
			NotAfterDays:          time.Now().Add(5000 * 24 * time.Hour),
			NotBeforeDays:         time.Now(),
			BasicConstraintsValid: true,
			IsCA:                  true,
			ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
			KeyUsage:              x509.KeyUsageCertSign | x509.KeyUsageCRLSign | x509.KeyUsageDataEncipherment,
			SignatureAlgorithm:    x509.ECDSAWithSHA384,
		},
	}); nil != errCA {
		t.Error(errCA)
	}
}

func TestCACommon_GenerateRSACertificateRequest_Fail(t *testing.T) {
	_, errCA = CAGenerateRSACertificateRequest(&CertRequest{
		PrivateKeyData:             priData,
		CertificateRequestFilePath: filepath.Join(pathcarsapksc1512, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA256WithRSAPSS,
		Subject:                    CAMockSubject,
	}, "PRIVATE KEY", RSAPKSC1())
	t.Log(errCA)
}

func TestCACommon_GenerateRSACertificateRequestFP_Fail(t *testing.T) {
	_, errCA = CAGenerateRSACertificateRequestFP(&CertRequestFP{
		PrivateKeyFilePath:         "",
		CertificateRequestFilePath: filepath.Join(pathcarsapksc1512, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA256WithRSAPSS,
		Subject:                    CAMockSubject,
	}, "PRIVATE KEY", RSAPKSC1())
	t.Log(errCA)
}

func TestCACommon_GenerateRSACertificateRequestWithPass_Fail(t *testing.T) {
	priData, errECC = ioutil.ReadFile(filepath.Join(pathcaeccpemp384, caPriKeyFileName))
	if nil != errECC {
		t.Error(errECC)
	}
	_, errCA = CAGenerateRSACertificateRequestWithPass(&CertRequest{
		PrivateKeyData:             priData,
		CertificateRequestFilePath: filepath.Join(pathcarsapksc1512, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA256WithRSAPSS,
		Subject:                    CAMockSubject,
	}, "PRIVATE KEY", "123", RSAPKSC1())
	t.Log(errCA)
}

func TestCACommon_GenerateRSACertificateRequestFPWithPass_Fail(t *testing.T) {
	_, errCA = CAGenerateRSACertificateRequestFPWithPass(&CertRequestFP{
		PrivateKeyFilePath:         "",
		CertificateRequestFilePath: filepath.Join(pathcarsapksc1512, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA256WithRSAPSS,
		Subject:                    CAMockSubject,
	}, "PRIVATE KEY", "123", RSAPKSC1())
	t.Log(errCA)
}

func TestCACommon_GenerateECCCertificateRequest_Fail(t *testing.T) {
	_, errCA = CAGenerateECCCertificateRequest(&CertRequest{
		PrivateKeyData:             priData,
		CertificateRequestFilePath: filepath.Join(pathcarsapksc1512, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA256WithRSAPSS,
		Subject:                    CAMockSubject,
	}, "PRIVATE KEY")
	t.Log(errCA)
}

func TestCACommon_GenerateECCCertificateRequestFP_Fail(t *testing.T) {
	_, errCA = CAGenerateECCCertificateRequestFP(&CertRequestFP{
		PrivateKeyFilePath:         "",
		CertificateRequestFilePath: filepath.Join(pathcarsapksc1512, caCertificateRequestFileName),
		SignatureAlgorithm:         x509.SHA256WithRSAPSS,
		Subject:                    CAMockSubject,
	}, "PRIVATE KEY")
	t.Log(errCA)
}

func TestCACommon_GenerateCertificateRequest_Fail(t *testing.T) {
	_, errCA = CAGenerateCertificateRequest(&CertRequestModel{
		CertificateRequestFilePath: filepath.Join(pathcarsapksc1512, caCertificateRequestFileName),
	})
	t.Log(errCA)
}

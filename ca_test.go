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
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
)

var (
	pathcarsapksc1256 = "./tmp/example/ca/pksc1/256"
	//pathcarsapksc1512  = "./tmp/example/ca/pksc1/512"
	//pathcarsapksc11024 = "./tmp/example/ca/pksc1/1024"
	//pathcarsapksc12048 = "./tmp/example/ca/pksc1/2048"

	//pathcarsapksc8256  = "./tmp/example/ca/pksc8/256"
	//pathcarsapksc8512  = "./tmp/example/ca/pksc8/512"
	//pathcarsapksc81024 = "./tmp/example/ca/pksc8/1024"
	//pathcarsapksc82048 = "./tmp/example/ca/pksc8/2048"
	//
	//pathcaeccpemp224 = "./tmp/example/ca/pemp224"
	//pathcaeccpemp256 = "./tmp/example/ca/pemp256"
	//pathcaeccpemp384 = "./tmp/example/ca/pemp384"
	//pathcaeccpemp521 = "./tmp/example/ca/pemp521"

	caPriKeyFileName      = "rootCA.key" // ca 私钥
	caCertificateFileName = "rootCA.crt" // 根证书文件

	fileIO *os.File

	priRSAKey *rsa.PrivateKey
	//priECCKey *ecdsa.PrivateKey

	csrData []byte

	errCA error
)

func TestCACommon_GenerateRSAPKCS1PrivateKey(t *testing.T) {
	if errCA = CA().GenerateRSAPKCS1PrivateKey(256, pathcarsapksc1256, caPriKeyFileName); nil != errCA {
		t.Error(errCA)
	}
	priRSAKey, errCA = CryptoRSA().LoadPriFP(filepath.Join(pathcarsapksc1256, caPriKeyFileName), CryptoRSA().pksC1())
	if errCA != nil {
		t.Error(errCA)
	}
	if csrData, errCA = CA().GenerateCertificate(rand.Reader,
		&x509.CertificateRequest{
			SignatureAlgorithm: x509.SHA256WithRSA,
			Subject: pkix.Name{
				Country:            []string{"CN"},
				Organization:       []string{"Gnomon"},
				OrganizationalUnit: []string{"GnomonRD"},
				Locality:           []string{"Beijing"},
				Province:           []string{"Beijing"},
				CommonName:         "aberic.cn",
			},
		},
		priRSAKey); nil != errCA {
		t.Skip(errCA)
	}
	if fileIO, errCA = os.OpenFile(filepath.Join(pathcarsapksc1256, caCertificateFileName), os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600); nil != errCA {
		t.Error(errCA)
	}
	// 将block的PEM编码写入fileIO
	if errCA = pem.Encode(fileIO, &pem.Block{Type: "CERTIFICATE REQUEST", Bytes: csrData}); nil != errCA {
		t.Error(errCA)
	}
}

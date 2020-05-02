/*
 * Copyright (c) 2019. aberic - All Rights Reserved.
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
 *
 */

package gnomon

import (
	"crypto"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
)

// HashMD5Bytes MD5Bytes
func HashMD5Bytes(bytes []byte) string {
	hash := md5.New()
	_, _ = hash.Write(bytes)
	return hex.EncodeToString(hash.Sum(nil))
}

// HashMD5 MD5
func HashMD5(text string) string {
	return HashMD5Bytes([]byte(text))
}

// HashMD516Bytes MD516Bytes
func HashMD516Bytes(bytes []byte) string {
	md516 := string([]rune(HashMD5Bytes(bytes))[8:24])
	return md516
}

// HashMD516 MD516
func HashMD516(text string) string {
	md516 := string([]rune(HashMD5(text))[8:24])
	return md516
}

// HashSha1Bytes Sha1Bytes
func HashSha1Bytes(bytes []byte) string {
	hash := sha1.New()
	_, _ = hash.Write(bytes)
	return hex.EncodeToString(hash.Sum(nil))
}

// HashSha1 Sha1
func HashSha1(text string) string {
	return HashSha1Bytes([]byte(text))
}

// HashSha224Bytes Sha224Bytes
func HashSha224Bytes(bytes []byte) string {
	hash := crypto.SHA224.New()
	_, _ = hash.Write(bytes)
	return hex.EncodeToString(hash.Sum(nil))
}

// HashSha224 Sha224
func HashSha224(text string) string {
	return HashSha224Bytes([]byte(text))
}

// HashSha256Bytes Sha256Bytes
func HashSha256Bytes(bytes []byte) string {
	hash := sha256.New()
	_, _ = hash.Write(bytes)
	return hex.EncodeToString(hash.Sum(nil))
}

// HashSha256 Sha256
func HashSha256(text string) string {
	return HashSha256Bytes([]byte(text))
}

// HashSha384Bytes Sha384Bytes
func HashSha384Bytes(bytes []byte) string {
	hash := sha512.New384()
	_, _ = hash.Write(bytes)
	return hex.EncodeToString(hash.Sum(nil))
}

// HashSha384 Sha384
func HashSha384(text string) string {
	return HashSha384Bytes([]byte(text))
}

// HashSha512Bytes Sha512Bytes
func HashSha512Bytes(bytes []byte) string {
	hash := sha512.New()
	_, _ = hash.Write(bytes)
	return hex.EncodeToString(hash.Sum(nil))
}

// HashSha512 Sha512
func HashSha512(text string) string {
	return HashSha512Bytes([]byte(text))
}

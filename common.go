package common

import "sync"

var (
	bc          *byteCommon
	cc          *commandCommon
	ec          *envCommon
	fc *fileCommon
	ic *ipCommon
	jc *jwtCommon
	sc *stringCommon
	hc *hashCommon
	rc *rsaCommon
	onceByte    sync.Once
	onceCommand sync.Once
	onceEnv     sync.Once
	onceFile sync.Once
	onceIp sync.Once
	onceJwt sync.Once
	onceString sync.Once
	onceHash sync.Once
	onceRSA sync.Once
)

func Byte() *byteCommon {
	onceByte.Do(func() {
		bc = &byteCommon{}
	})
	return bc
}

func Command() *commandCommon {
	onceCommand.Do(func() {
		cc = &commandCommon{}
	})
	return cc
}

func Env() *envCommon {
	onceEnv.Do(func() {
		ec = &envCommon{}
	})
	return ec
}

func File() *fileCommon {
	onceFile.Do(func() {
		fc = &fileCommon{}
	})
	return fc
}

func IP() *ipCommon {
	onceIp.Do(func() {
		ic = &ipCommon{}
	})
	return ic
}

func Jwt() *jwtCommon {
	onceJwt.Do(func() {
		jc = &jwtCommon{}
	})
	return jc
}

func String() *stringCommon {
	onceString.Do(func() {
		sc = &stringCommon{}
	})
	return sc
}

func CryptoHash() *hashCommon {
	onceHash.Do(func() {
		hc = &hashCommon{}
	})
	return hc
}

func CryptoRSA() *rsaCommon {
	onceRSA.Do(func() {
		rc = &rsaCommon{}
	})
	return rc
}

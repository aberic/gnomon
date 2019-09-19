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
	"sync"
)

var (
	bc          *byteCommon
	cc          *commandCommon
	ec          *envCommon
	fc          *fileCommon
	ic          *ipCommon
	jc          *jwtCommon
	sc          *stringCommon
	hc          *hashCommon
	rc          *rsaCommon
	lc          *logCommon
	scc         *scaleCommon
	tc          *timeCommon
	onceByte    sync.Once
	onceCommand sync.Once
	onceEnv     sync.Once
	onceFile    sync.Once
	onceIp      sync.Once
	onceJwt     sync.Once
	onceString  sync.Once
	onceHash    sync.Once
	onceRSA     sync.Once
	onceLog     sync.Once
	onceScale   sync.Once
	onceTime    sync.Once
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

func Log() *logCommon {
	onceLog.Do(func() {
		lc = &logCommon{level: debugLevel, production: false}
	})
	return lc
}

func Scale() *scaleCommon {
	onceScale.Do(func() {
		scc = &scaleCommon{}
	})
	return scc
}

func Time() *timeCommon {
	onceTime.Do(func() {
		tc = &timeCommon{}
	})
	return tc
}

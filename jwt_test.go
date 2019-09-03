/*
 * Copyright (c) 2019. ENNOO - All Rights Reserved.
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

package common

import (
	"fmt"
	"testing"
	"time"
)

func TestBuild(t *testing.T) {
	key := []byte("Hello World！This is secret!")
	tokenString1, err1 := Jwt().Build(SigningMethodHS256, key, "1", "rivet", "userMD5", time.Now().Unix(), time.Now().Unix(), time.Now().Unix()+1000)
	tokenString2, err2 := Jwt().Build(SigningMethodHS384, key, "1", "rivet", "userMD5", time.Now().Unix(), time.Now().Unix(), time.Now().Unix()+1000)
	tokenString3, err3 := Jwt().Build(SigningMethodHS512, key, "1", "rivet", "userMD5", time.Now().Unix(), time.Now().Unix(), time.Now().Unix()+1000)
	if nil == err1 && nil == err2 && nil == err3 {
		fmt.Println("tokenString1 = ", tokenString1)
		fmt.Println("tokenString2 = ", tokenString2)
		fmt.Println("tokenString3 = ", tokenString3)
		time.Sleep(1 * time.Second)
		bo1 := Jwt().Check(key, tokenString1)
		bo2 := Jwt().Check(key, tokenString2)
		bo3 := Jwt().Check(key, tokenString3)
		fmt.Println("bo1 = ", bo1)
		fmt.Println("bo2 = ", bo2)
		fmt.Println("bo3 = ", bo3)
	}
}

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

package gnomon

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestPathExists(t *testing.T) {
	path := "/etc/profile"
	exist := File().PathExists(path)
	fmt.Println(path, "exist =", exist)

	path = "/etc/hello"
	exist = File().PathExists(path)
	fmt.Println(path, "exist =", exist)

	path = "/ha/oo"
	err = os.MkdirAll(path, os.ModePerm)
	if nil == err {
		exist = File().PathExists(path)
		fmt.Println(path, "exist =", exist)
		err = os.Remove(path)
		if nil != err {
			fmt.Println(err.Error())
		}
	} else {
		fmt.Println(err.Error())
	}
}

func TestReadFileFirstLine(t *testing.T) {
	txt, err := File().ReadFileFirstLine("../../../a.txt")
	if nil != err {
		fmt.Println(err.Error())
	} else {
		fmt.Println("txt =", txt)
	}

	profile, err := File().ReadFileFirstLine("/etc/profile")
	if nil != err {
		fmt.Println(err.Error())
	} else {
		fmt.Println("profile =", profile)
	}

	hello, err := File().ReadFileFirstLine("/etc/hello")
	if nil != err {
		fmt.Println(err.Error())
	} else {
		fmt.Println("hello =", hello)
	}
}

func TestReadFileByLine(t *testing.T) {
	hosts, err := File().ReadFileByLine("/etc/hostname")
	if nil != err {
		fmt.Println(err.Error())
	} else {
		fmt.Println("hosts =", hosts)
	}

	profile, err := File().ReadFileByLine("/etc/profile")
	if nil != err {
		fmt.Println(err.Error())
	} else {
		fmt.Println("profile =", profile)
	}

	hello, err := File().ReadFileByLine("/etc/hello")
	if nil != err {
		fmt.Println(err.Error())
	} else {
		fmt.Println("hello =", hello)
	}
}

func TestCreateAndWrite(t *testing.T) {
	if err := File().CreateAndWrite("/etc/yes/go/test.txt", []byte("haha"), false); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestGetAllFile(t *testing.T) {
	var s []string
	if arr, err := File().LoopAllFileFromDir("./log", s); nil != err {
		t.Skip(err)
	} else {
		t.Log(arr)
	}
}

func TestFileCommon_LoopDirFromDir(t *testing.T) {
	if arr, err := File().LoopDirFromDir("./log"); nil != err {
		t.Skip(err)
	} else {
		t.Log(arr)
	}
	t.Log(time.Now().Local().Format("20060102"))
}

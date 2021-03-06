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
 */

package gnomon

import (
	"gotest.tools/assert"
	"os"
	"testing"
)

func TestFileCommon_PathExists(t *testing.T) {
	path := "/etc"
	exist := FilePathExists(path)
	t.Log(path, "exist =", exist)

	path = "/etc/hello"
	exist = FilePathExists(path)
	t.Log(path, "exist =", exist)
}

func TestFileCommon_ReadFirstLine(t *testing.T) {
	profile, err := FileReadFirstLine("/etc/profile")
	if nil != err {
		t.Skip(err)
	} else {
		t.Log("profile =", profile)
	}
}

func TestFileCommon_ReadFirstLine_Fail(t *testing.T) {
	_, err := FileReadFirstLine("/etc/hello")
	t.Skip(err)
}

func TestFileCommon_ReadPointLine(t *testing.T) {
	profile, err := FileReadPointLine("/etc/profile", 1)
	if nil != err {
		t.Skip(err)
	} else {
		t.Log("profile =", profile)
	}
}

func TestFileCommon_ReadPointLine_KeyPoint(t *testing.T) {
	_, _ = FileAppend("./tmp/log/yes/go/point.txt", []byte("haha"), false)
	profile, err := FileReadPointLine("./tmp/log/yes/go/point.txt", 1)
	if nil != err {
		t.Skip(err)
	} else {
		t.Log("profile =", profile)
	}
}

func TestFileCommon_ReadPointLine_Fail_IndexOut(t *testing.T) {
	_, err := FileReadPointLine("/etc/profile", 300)
	t.Skip(err)
}

func TestFileCommon_ReadPointLine_Fail_NotExist(t *testing.T) {
	_, err := FileReadPointLine("/etc/hello", 1)
	t.Skip(err)
}

func TestFileCommon_ReadLines(t *testing.T) {
	profile, err := FileReadLines("/etc/profile")
	if nil != err {
		t.Skip(err)
	} else {
		t.Log("profile =", profile)
	}
}

func TestFileCommon_ReadLines_Fail(t *testing.T) {
	_, err := FileReadLines("/etc/hello")
	t.Skip(err)
}

func TestFileCommon_ParentPath(t *testing.T) {
	t.Log(FileParentPath("/etc/yes/go/test.txt"))
}

func TestFileCommon_Append(t *testing.T) {
	if _, err := FileAppend("./tmp/log/yes/go/test.txt", []byte("haha"), false); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Append_Force(t *testing.T) {
	if _, err := FileAppend("./tmp/log/yes/go/test.txt", []byte("haha"), true); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Append_UnForce(t *testing.T) {
	if _, err := FileAppend("./tmp/log/yes/go/test.txt", []byte("haha"), false); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Append_Fail_PermissionFileForce(t *testing.T) {
	_, err := FileAppend("/etc/www.json", []byte("haha"), true)
	t.Skip(err)
}

func TestFileCommon_Append_Fail_PermissionFileUnForce(t *testing.T) {
	_, err := FileAppend("/etc/www.json", []byte("haha"), false)
	t.Skip(err)
}

func TestFileCommon_Modify(t *testing.T) {
	if _, err := FileModify("./tmp/log/yes/go/test.txt", 1, []byte("haha"), false); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Modify_Force(t *testing.T) {
	if _, err := FileModify("./tmp/log/yes/go/test.txt", 1, []byte("haha"), true); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Modify_UnForce(t *testing.T) {
	if _, err := FileModify("./tmp/log/yes/go/test.txt", 1, []byte("haha"), false); nil != err {
		t.Skip(err)
	} else {
		t.Log("success")
	}
}

func TestFileCommon_Modify_Fail_PermissionFileForce(t *testing.T) {
	_, err := FileModify("/etc/www.json", 1, []byte("haha"), true)
	t.Skip(err)
}

func TestFileCommon_Modify_Fail_PermissionFileUnForce(t *testing.T) {
	_, err := FileModify("/etc/www.json", 1, []byte("haha"), false)
	t.Skip(err)
}

func TestFileCommon_LoopDirs(t *testing.T) {
	if arr, err := FileLoopDirs("./tmp/log"); nil != err {
		t.Skip(err)
	} else {
		t.Log(arr)
	}
}

func TestFileCommon_LoopDirs_Fail(t *testing.T) {
	_, err := FileLoopDirs("./tmp/logger")
	t.Skip(err)
}

func TestFileCommon_LoopFiles(t *testing.T) {
	if arr, err := FileLoopFiles("./tmp/log"); nil != err {
		t.Skip(err)
	} else {
		t.Log(arr)
	}
}

func TestFileCommon_LoopFiles_Fail(t *testing.T) {
	_, err := FileLoopFiles("./tmp/logger")
	t.Skip(err)
}

func TestFileCommon_LoopOneDirs(t *testing.T) {
	array, err := FileLoopOneDirs("./tmp")
	t.Skip(array, err)
}

func TestFileCommon_Copy(t *testing.T) {
	if _, err := FileAppend("./tmp/copy/1.txt", []byte("hello"), true); nil != err {
		t.Skip(err)
	}
	if _, err := FileCopy("./tmp/copy/1.txt", "./tmp/copy/2.txt"); nil != err {
		t.Skip(err)
	}
}

func TestFileCompress(t *testing.T) {
	f, err := os.Open("./example/ca")
	assert.NilError(t, err)
	err = FileCompressZip([]*os.File{f}, "./example/cas.zip")
	assert.NilError(t, err)
}

func TestFileCompressTar(t *testing.T) {
	f, err := os.Open("./example/ca")
	assert.NilError(t, err)
	err = FileCompressTar([]*os.File{f}, "./example/cas.tar")
	assert.NilError(t, err)
	//err = FileDeCompressTar("./example/cas.tar", "./example/castar")
	//assert.NilError(t, err)
}

func TestFileDeCompressZip(t *testing.T) {
	err := FileDeCompressZip("./example/ca_grope_sql.zip", "./example/ca_grope_sql_de")
	assert.NilError(t, err)
	err = FileDeCompressZip("./example/sql.zip", "./example/sql_de")
	assert.NilError(t, err)
}

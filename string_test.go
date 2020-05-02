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
	"strings"
	"testing"
	"time"
)

func TestIsEmpty(t *testing.T) {
	t.Log("haha empty =", StringIsEmpty("haha"))
	t.Log("'' empty =", StringIsEmpty(""))
}

func TestIsNotEmpty(t *testing.T) {
	t.Log("haha empty =", StringIsNotEmpty("haha"))
	t.Log("'' empty =", StringIsNotEmpty(""))
}

func TestConvert(t *testing.T) {
	t.Log("uu_xx_aa =", StringConvert("uu_xx_aa"))
}

func TestRandSeq(t *testing.T) {
	t.Log("13 =", StringRandSeq(13))
	t.Log("23 =", StringRandSeq(23))
	t.Log("33 =", StringRandSeq(33))
}

func TestRandSeq16(t *testing.T) {
	t.Log("RandSeq16 =", StringRandSeq16())
}

func TestTrim(t *testing.T) {
	s := "kjsdhfj ajsd\nksjhdka sjkh"
	t.Log(s, "=", StringTrim(s))
}

func TestString_ToString(t *testing.T) {
	t.Log(ToString(&BTest{Name: "test", Age: 18, Male: true}))
	t.Log(ToString(nil))
}

func TestString_SingleValue(t *testing.T) {
	t.Log(StringSingleValue("ksjdf/////lksjdf/////lkjlksdf/////lkjl/lkjasldj kjnkj ///", "/"))
}

func TestString_SingleSpace(t *testing.T) {
	t.Log(StringSingleSpace("ksjdf     lksjdf  lkjlksdf        lkjl   lkjasldj kjnkj     "))
}

func TestString_PrefixSupplementZero(t *testing.T) {
	t.Log("ui64", "92873890910928019", StringPrefixSupplementZero("92873890910928019", 10))
	t.Log("ui64", "92873890910928019", StringPrefixSupplementZero("92873890910928019", 20))
	t.Log("ui64", "92873890910928019", StringPrefixSupplementZero("92873890910928019", 30))
	t.Log("ui64", "92873890910928019", StringPrefixSupplementZero("92873890910928019", 40))
	t.Log("ui64", "92873890910928019", StringPrefixSupplementZero("92873890910928019", 50))
	t.Log("ui64", "92873890910928019", StringPrefixSupplementZero("92873890910928019", 60))
}

func TestString_SubString(t *testing.T) {
	t.Log(SubString("110xxxxxxx", 0, 3))
}

func TestString_StringBuilder(t *testing.T) {
	t.Log(StringBuild("a", "b"))
}

func TestSlice(t *testing.T) {
	var url = "/test/a/b/c/d"
	t.Log(strings.Split(url, "/")[1:])
}

func TestUrl(t *testing.T) {
	var url1 = "http://127.0.0.1:8080/test/demo/1/g?name=hello&pass=work"
	var url2 = "http://127.0.0.1:8080/test/demo/1/g"
	t.Log(strings.Split(url1, "?"))
	t.Log(strings.Split(url2, "?"))
}

func TestString_String2Timestamp(t *testing.T) {
	i64, err := String2Timestamp("2019/09/17 10:16:56", "2006/01/02 15:04:05", time.Local)
	t.Log("i64", i64, err)
	i64, err = String2Timestamp("2019/09/17 10:16:56", "2006/01/02 15:04:05", time.UTC)
	t.Log("i64", i64, err)
}

func TestString_String2Timestamp_Fail(t *testing.T) {
	_, err := String2Timestamp("hello world", "2006/01/02 15:04:05", time.Local)
	t.Skip(err)
}

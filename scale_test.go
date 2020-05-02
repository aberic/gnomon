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

import "testing"

func TestScaleHex(t *testing.T) {
	var (
		i    int
		ui   uint
		i8   int8
		ui8  uint8
		i16  int16
		ui16 uint16
		i32  int32
		ui32 uint32
		i64  int64
		ui64 uint64
	)
	i = 2999999
	ui = 2888888
	i8 = 127
	ui8 = 255
	i16 = 32767
	ui16 = 65535
	i32 = 99922299
	ui32 = 88811188
	i64 = 827639847923879
	ui64 = 92873890910928019

	iStr := ScaleIntToHexString(i)
	t.Log("i", "2999999", iStr)
	t.Log("i", "2999999", ScaleHexStringToUint64(iStr))
	t.Log("")

	iStr = ScaleUintToHexString(uint(i))
	t.Log("i", "2999999", iStr)
	t.Log("i", "2999999", ScaleHexStringToUint64(iStr))
	t.Log("")

	uiStr := ScaleUintToHexString(uint(ui))
	t.Log("ui", "2888888", uiStr)
	t.Log("i", "2888888", ScaleHexStringToUint64(uiStr))
	t.Log("")

	i8Str := ScaleUintToHexString(uint(i8))
	t.Log("i8", "127", i8Str)
	t.Log("i8", "127", ScaleHexStringToUint64(i8Str))
	t.Log("")

	ui8Str := ScaleUintToHexString(uint(ui8))
	t.Log("ui8", "255", ui8Str)
	t.Log("ui8", "255", ScaleHexStringToUint64(ui8Str))
	t.Log("")

	i16Str := ScaleUintToHexString(uint(i16))
	t.Log("i16", "32767", i16Str)
	t.Log("i16", "32767", ScaleHexStringToUint64(i16Str))
	t.Log("")

	ui16Str := ScaleUintToHexString(uint(ui16))
	t.Log("ui16", "65535", ui16Str)
	t.Log("ui16", "65535", ScaleHexStringToUint64(ui16Str))
	t.Log("")

	i32Str := ScaleInt32ToHexString(i32)
	t.Log("i32", "99922299", i32Str)
	t.Log("i32", "99922299", ScaleHexStringToUint64(i32Str))
	t.Log("")

	ui32Str := ScaleUint32ToHexString(ui32)
	t.Log("ui32", "88811188", ui32Str)
	t.Log("ui32", "88811188", ScaleHexStringToUint64(ui32Str))
	t.Log("")

	i64Str := ScaleInt64ToHexString(i64)
	t.Log("i64", "827639847923879", i64Str)
	t.Log("i64", "827639847923879", ScaleHexStringToInt64(i64Str))
	t.Log("")

	ui64Str := ScaleUint64ToHexString(ui64)
	t.Log("ui64", "92873890910928019", ui64Str)
	t.Log("ui64", "92873890910928019", ScaleHexStringToUint64(ui64Str))
}

func TestScaleDuo(t *testing.T) {
	var (
		i    int
		ui   uint
		i8   int8
		ui8  uint8
		i16  int16
		ui16 uint16
		i32  int32
		ui32 uint32
		i64  int64
		ui64 uint64
	)
	i = 2999999
	ui = 2888888
	i8 = 127
	ui8 = 255
	i16 = 32767
	ui16 = 65535
	i32 = 99922299
	ui32 = 88811188
	i64 = 827639847923879
	ui64 = 92873890910928019

	iStr := ScaleIntToDuoString(i)
	t.Log("i", "2999999", iStr)
	t.Log("i", "2999999", ScaleDuoStringToUint64(iStr))
	t.Log("")

	iStr = ScaleUintToDuoString(uint(i))
	t.Log("i", "2999999", iStr)
	t.Log("i", "2999999", ScaleDuoStringToUint64(iStr))
	t.Log("")

	uiStr := ScaleUintToDuoString(uint(ui))
	t.Log("ui", "2888888", uiStr)
	t.Log("i", "2888888", ScaleDuoStringToUint64(uiStr))
	t.Log("")

	i8Str := ScaleUintToDuoString(uint(i8))
	t.Log("i8", "127", i8Str)
	t.Log("i8", "127", ScaleDuoStringToUint64(i8Str))
	t.Log("")

	ui8Str := ScaleUintToDuoString(uint(ui8))
	t.Log("ui8", "255", ui8Str)
	t.Log("ui8", "255", ScaleDuoStringToUint64(ui8Str))
	t.Log("")

	i16Str := ScaleUintToDuoString(uint(i16))
	t.Log("i16", "32767", i16Str)
	t.Log("i16", "32767", ScaleDuoStringToUint64(i16Str))
	t.Log("")

	ui16Str := ScaleUintToDuoString(uint(ui16))
	t.Log("ui16", "65535", ui16Str)
	t.Log("ui16", "65535", ScaleDuoStringToUint64(ui16Str))
	t.Log("")

	i32Str := ScaleInt32ToDuoString(i32)
	t.Log("i32", "99922299", i32Str)
	t.Log("i32", "99922299", ScaleDuoStringToUint64(i32Str))
	t.Log("")

	ui32Str := ScaleUint32ToDuoString(ui32)
	t.Log("ui32", "88811188", ui32Str)
	t.Log("ui32", "88811188", ScaleDuoStringToUint64(ui32Str))
	t.Log("")

	i64Str := ScaleInt64ToDuoString(i64)
	t.Log("i64", "827639847923879", i64Str)
	t.Log("i64", "827639847923879", ScaleDuoStringToInt64(i64Str))
	t.Log("")

	ui64Str := ScaleUint64ToDuoString(ui64)
	t.Log("ui64", "92873890910928019", ui64Str)
	t.Log("ui64", "92873890910928019", ScaleDuoStringToUint64(ui64Str))
}

func TestScaleDDuo(t *testing.T) {
	var (
		i    int
		ui   uint
		i8   int8
		ui8  uint8
		i16  int16
		ui16 uint16
		i32  int32
		ui32 uint32
		i64  int64
		ui64 uint64
	)
	i = 2999999
	ui = 2888888
	i8 = 127
	ui8 = 255
	i16 = 32767
	ui16 = 65535
	i32 = 99922299
	ui32 = 88811188
	i64 = 9223372036854770018
	ui64 = 92873890910928019

	iStr := ScaleIntToDDuoString(i)
	t.Log("i", "2999999", iStr)
	t.Log("i", "2999999", ScaleDDuoStringToUint64(iStr))
	t.Log("")

	iStr = ScaleUintToDDuoString(uint(i))
	t.Log("i", "2999999", iStr)
	t.Log("i", "2999999", ScaleDDuoStringToUint64(iStr))
	t.Log("")

	uiStr := ScaleUintToDDuoString(uint(ui))
	t.Log("ui", "2888888", uiStr)
	t.Log("i", "2888888", ScaleDDuoStringToUint64(uiStr))
	t.Log("")

	i8Str := ScaleUintToDDuoString(uint(i8))
	t.Log("i8", "127", i8Str)
	t.Log("i8", "127", ScaleDDuoStringToUint64(i8Str))
	t.Log("")

	ui8Str := ScaleUintToDDuoString(uint(ui8))
	t.Log("ui8", "255", ui8Str)
	t.Log("ui8", "255", ScaleDDuoStringToUint64(ui8Str))
	t.Log("")

	i16Str := ScaleUintToDDuoString(uint(i16))
	t.Log("i16", "32767", i16Str)
	t.Log("i16", "32767", ScaleDDuoStringToUint64(i16Str))
	t.Log("")

	ui16Str := ScaleUintToDDuoString(uint(ui16))
	t.Log("ui16", "65535", ui16Str)
	t.Log("ui16", "65535", ScaleDDuoStringToUint64(ui16Str))
	t.Log("")

	i32Str := ScaleInt32ToDDuoString(i32)
	t.Log("i32", "99922299", i32Str)
	t.Log("i32", "99922299", ScaleDDuoStringToUint64(i32Str))
	t.Log("")

	ui32Str := ScaleUint32ToDDuoString(ui32)
	t.Log("ui32", "88811188", ui32Str)
	t.Log("ui32", "88811188", ScaleDDuoStringToUint64(ui32Str))
	t.Log("")

	i64Str := ScaleInt64ToDDuoString(i64)
	t.Log("i64", "827639847923879", i64Str)
	t.Log("i64", "827639847923879", ScaleDDuoStringToInt64(i64Str))
	t.Log("")

	ui64Str := ScaleUint64ToDDuoString(ui64)
	t.Log("ui64", "92873890910928019", ui64Str)
	t.Log("ui64", "92873890910928019", ScaleDDuoStringToUint64(ui64Str))

	ui64 = 18446744073709551615
	ui64Str = ScaleUint64ToDDuoString(ui64)
	t.Log("ui64", "18446744073709551615", ui64Str)
	t.Log("ui64", "18446744073709551615", ScaleDDuoStringToUint64(ui64Str))

	for i := 0; i < 100000; i++ {
		ui64Str = ScaleUint64ToDDuoString(ui64)
		ScaleDDuoStringToUint64(ui64Str)
		//t.Log("ui64", "18446744073709551615", ui64Str)
		//t.Log("ui64", "18446744073709551615", ScaleDDuoStringToUint64(ui64Str))
	}
}

func TestScaleLen(t *testing.T) {
	var (
		i    int
		ui   uint
		i32  int32
		ui32 uint32
		i64  int64
		ui64 uint64
	)
	i = 2999999
	ui = 2888888
	i32 = 99922299
	ui32 = 88811188
	i64 = 827639847923879
	ui64 = 92873890910928019
	t.Log(ScaleUint64Len(ui64))
	t.Log(ScaleInt64Len(i64))
	t.Log(ScaleUint32Len(ui32))
	t.Log(ScaleInt32Len(i32))
	t.Log(ScaleUintLen(ui))
	t.Log(ScaleIntLen(i))
}

func TestScaleFullState(t *testing.T) {
	var (
		ui8  uint8
		ui32 uint32
	)
	ui8 = 25
	ui32 = 88811188
	t.Log(ScaleUint8toFullState(ui8))
	t.Log(ScaleUint32toFullState(ui32))
}

func TestScaleFloat64(t *testing.T) {
	var (
		i64 int64
		f64 float64
	)
	i64 = 87372
	f64 = 92837.87263876498
	t.Log(ScaleInt64toFloat64(i64, 1))
	t.Log(ScaleInt64toFloat64(i64, 2))
	t.Log(ScaleInt64toFloat64(i64, 3))
	t.Log(ScaleInt64toFloat64(i64, 4))
	t.Log(ScaleInt64toFloat64(i64, 5))
	t.Log()
	t.Log(ScaleFloat64toInt64(f64, 1))
	t.Log(ScaleFloat64toInt64(f64, 2))
	t.Log(ScaleFloat64toInt64(f64, 3))
	t.Log(ScaleFloat64toInt64(f64, 4))
	t.Log(ScaleFloat64toInt64(f64, 5))
}

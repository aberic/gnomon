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
 *
 */

package gnomon

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
)

type byteCommon struct{}

// GetBytes 获取接口字节数组
func (b *byteCommon) GetBytes(key interface{}) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	err := enc.Encode(key)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// IntToBytes 整形转换成字节
func (b *byteCommon) IntToBytes(n int) ([]byte, error) {
	x := int32(n)

	bytesBuffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(bytesBuffer, binary.BigEndian, x); nil != err {
		return nil, err
	}
	return bytesBuffer.Bytes(), nil
}

// BytesToInt 字节转换成整形
func (b *byteCommon) BytesToInt(byte []byte) (int, error) {
	bytesBuffer := bytes.NewBuffer(byte)

	var x int32
	if err := binary.Read(bytesBuffer, binary.BigEndian, &x); nil != err {
		return 0, err
	}

	return int(x), nil
}

// Int64ToBytes 整形64转换成字节
func (b *byteCommon) Int64ToBytes(i int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(i))
	return buf
}

// BytesToInt64 字节转换成整形64
func (b *byteCommon) BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

// Append 字节数组追加
func (b *byteCommon) Append(bs2 []byte, bsf []byte) []byte {
	for _, b := range bsf {
		bs2 = append(bs2, b)
	}
	return bs2
}

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

// IntToBytes int转换成字节
func (b *byteCommon) IntToBytes(n int) ([]byte, error) {
	x := int32(n)
	bytesBuffer := bytes.NewBuffer([]byte{})
	if err := binary.Write(bytesBuffer, binary.BigEndian, x); nil != err {
		return nil, err
	}
	return bytesBuffer.Bytes(), nil
}

// BytesToInt 字节转换成int
func (b *byteCommon) BytesToInt(byte []byte) (int, error) {
	bytesBuffer := bytes.NewBuffer(byte)
	var x int32
	if err := binary.Read(bytesBuffer, binary.BigEndian, &x); nil != err {
		return 0, err
	}
	return int(x), nil
}

// Uint16ToBytes uint16转换成字节
func (b *byteCommon) Uint16ToBytes(i uint16) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint16(buf, i)
	return buf
}

// BytesToUint16 字节转换成uint16
func (b *byteCommon) BytesToUint16(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf)
}

// Uint32ToBytes uint32转换成字节
func (b *byteCommon) Uint32ToBytes(i uint32) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint32(buf, i)
	return buf
}

// BytesToUint32 字节转换成uint32
func (b *byteCommon) BytesToUint32(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}

// Uint64ToBytes uint64转换成字节
func (b *byteCommon) Uint64ToBytes(i uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, i)
	return buf
}

// BytesToUint64 字节转换成uint64
func (b *byteCommon) BytesToUint64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}

// Append 字节数组追加
func (b *byteCommon) Append(bs2 []byte, bsf []byte) []byte {
	for _, b := range bsf {
		bs2 = append(bs2, b)
	}
	return bs2
}

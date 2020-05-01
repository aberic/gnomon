/*
 *  Copyright (c) 2020. aberic - All Rights Reserved.
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

package balance

import (
	"errors"
	"fmt"
	"hash/crc32"
	"math/rand"
	"sync"
)

type hash struct {
	interSlice []interface{}
	lock       sync.RWMutex
}

func newHash() *hash {
	return &hash{
		interSlice: []interface{}{},
	}
}

// Add 新增负载对象
func (h *hash) Add(obj interface{}) {
	defer h.lock.Unlock()
	h.lock.Lock()
	h.interSlice = append(h.interSlice, obj)
}

// Weight 设置负载对象权重
func (h *hash) Weight(obj interface{}, weight int) {
	h.Remove(obj)
	for i := 0; i < weight; i++ {
		h.Add(obj)
	}
}

// Remove 移除负载对象
func (h *hash) Remove(obj interface{}) {
	defer h.lock.Unlock()
	h.lock.Lock()
	for index, i := range h.interSlice {
		if i == obj {
			h.interSlice = append(h.interSlice[:index], h.interSlice[index+1:]...)
		}
	}
}

// Class 获取负载均衡分类
func (h *hash) Class() Class {
	return Hash
}

// Acquire 执行负载均衡算法得到期望对象
func (h *hash) Acquire() (interface{}, error) {
	lens := len(h.interSlice)
	if lens == 0 {
		return nil, errors.New("no instance")
	}
	defKey := fmt.Sprintf("%d", rand.Int())
	hashVal := crc32.Checksum([]byte(defKey), crc32.MakeTable(crc32.IEEE))
	index := int(hashVal) % lens
	return h.interSlice[index], nil
}

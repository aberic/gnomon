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
	"math/rand"
	"sync"
)

type random struct {
	interSlice []interface{}
	lock       sync.RWMutex
}

func newRandom() *random {
	return &random{
		interSlice: []interface{}{},
	}
}

// Add 新增负载对象
func (r *random) Add(obj interface{}) {
	defer r.lock.Unlock()
	r.lock.Lock()
	r.interSlice = append(r.interSlice, obj)
}

// Weight 设置负载对象权重
func (r *random) Weight(obj interface{}, weight int) {
	r.Remove(obj)
	for i := 0; i < weight; i++ {
		r.Add(obj)
	}
}

// Remove 移除负载对象
func (r *random) Remove(obj interface{}) {
	defer r.lock.Unlock()
	r.lock.Lock()
	for index, i := range r.interSlice {
		if i == obj {
			r.interSlice = append(r.interSlice[:index], r.interSlice[index+1:]...)
		}
	}
}

// Class 获取负载均衡分类
func (r *random) Class() Class {
	return Random
}

// Acquire 执行负载均衡算法得到期望对象
func (r *random) Acquire() (interface{}, error) {
	lens := len(r.interSlice)
	if lens == 0 {
		return nil, errors.New("no instance")
	}
	return r.interSlice[rand.Intn(lens)], nil
}

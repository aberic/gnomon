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
	"sync"
)

type round struct {
	interSlice []interface{}
	lock       sync.RWMutex
	ch         chan int
}

func newRound() *round {
	r := &round{
		interSlice: []interface{}{},
		ch:         make(chan int, 65535),
	}
	go r.generaCount()
	return r
}

// generaCount 自增生成一个0~65535的数，到达65535则重0开始计数
func (r *round) generaCount() {
	i := 0
	for {
		r.ch <- i // 等待索要数据
		if i == 65535 {
			i = -1
		}
		i++
	}
}

// Add 新增负载对象
func (r *round) Add(obj interface{}) {
	defer r.lock.Unlock()
	r.lock.Lock()
	r.interSlice = append(r.interSlice, obj)
}

// Remove 移除负载对象
func (r *round) Remove(obj interface{}) {
	defer r.lock.Unlock()
	r.lock.Lock()
	for index, i := range r.interSlice {
		if i == obj {
			r.interSlice = append(r.interSlice[:index], r.interSlice[index+1:]...)
			break
		}
	}
}

// Class 获取负载均衡分类
func (r *round) Class() Class {
	return Round
}

// Run 执行负载均衡算法得到期望对象
func (r *round) Run() (interface{}, error) {
	lens := len(r.interSlice)
	if lens == 0 {
		return nil, errors.New("no instance")
	}
	position := <-r.ch
	if position >= lens {
		position = position % lens
	}
	return r.interSlice[position], nil
}

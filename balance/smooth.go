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
	"sort"
	"sync"
)

type smooth struct {
	totalWeight int
	params      []*param
	lock        sync.RWMutex
}

func newSmooth() *smooth {
	r := &smooth{
		params: []*param{},
	}
	return r
}

// Add 新增负载对象
func (s *smooth) Add(obj interface{}) {
	defer s.lock.Unlock()
	s.lock.Lock()
	for _, p := range s.params {
		if p.obj == obj {
			differ := 1 - p.staticWeight
			p.staticWeight = 1
			p.dynamicWeight = differ + p.dynamicWeight
			s.totalWeight += differ
			return
		}
	}
	s.params = append(s.params, &param{obj: obj, staticWeight: 1, dynamicWeight: 1})
	sort.Slice(s.params, func(i, j int) bool {
		return s.params[i].dynamicWeight > s.params[j].dynamicWeight
	})
	s.totalWeight++
}

// Weight 设置负载对象权重
func (s *smooth) Weight(obj interface{}, weight int) {
	for _, p := range s.params {
		if p.obj == obj {
			differ := weight - p.staticWeight
			p.staticWeight = weight
			p.dynamicWeight = differ + p.dynamicWeight
			s.totalWeight += differ
			break
		}
	}
	sort.Slice(s.params, func(i, j int) bool {
		return s.params[i].dynamicWeight > s.params[j].dynamicWeight
	})
}

// Remove 移除负载对象
func (s *smooth) Remove(obj interface{}) {
	defer s.lock.Unlock()
	s.lock.Lock()
	for index, p := range s.params {
		if p.obj == obj {
			s.params = append(s.params[:index], s.params[index+1:]...)
			s.totalWeight -= p.staticWeight
		}
	}
}

// Class 获取负载均衡分类
func (s *smooth) Class() Class {
	return Smooth
}

// Acquire 执行负载均衡算法得到期望对象
func (s *smooth) Acquire() (interface{}, error) {
	param := s.params[0]
	param.dynamicWeight = param.dynamicWeight - s.totalWeight + param.staticWeight
	sort.Slice(s.params, func(i, j int) bool {
		return s.params[i].dynamicWeight > s.params[j].dynamicWeight
	})
	return param.obj, nil
}

// param 平滑加权聚合对象
type param struct {
	obj           interface{}
	staticWeight  int
	dynamicWeight int
}

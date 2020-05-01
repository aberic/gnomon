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

const (
	// Round 轮询负载
	Round Class = iota
	// Random 随机负载
	Random
	// Hash hash负载
	Hash
)

// Class 负载均衡分类
type Class int

// Balancer 负载均衡器
type Balancer interface {
	// Add 新增负载对象
	Add(interface{})
	// Weight 设置负载对象权重
	Weight(interface{}, int)
	// Remove 移除负载对象
	Remove(interface{})
	// Class 获取负载均衡分类
	Class() Class
	// Acquire 执行负载均衡算法得到期望对象
	Acquire() (interface{}, error)
}

// NewBalance 新建负载均衡器
func NewBalance(c Class) Balancer {
	switch c {
	default:
		return newRound()
	case Round:
		return newRound()
	case Random:
		return newRandom()
	case Hash:
		return newHash()
	}
}

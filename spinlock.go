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

package gnomon

import (
	"runtime"
	"sync"
	"sync/atomic"
)

// NewSpinLock 新建自旋锁
func NewSpinLock() sync.Locker {
	var lock spinLock
	return &lock
}

// spinLock 是一个自旋锁实现
//
// spinLock 不可在首次使用后复制
type spinLock struct {
	lock uintptr
}

// Lock 锁定
//
// 如果锁已经在使用，则调用goroutine，直到锁可用为止
func (sl *spinLock) Lock() {
	for !atomic.CompareAndSwapUintptr(&sl.lock, 0, 1) {
		runtime.Gosched()
	}
}

// Unlock 解锁
func (sl *spinLock) Unlock() {
	atomic.StoreUintptr(&sl.lock, 0)
}

func (sl *spinLock) TryLock() bool {
	return atomic.CompareAndSwapUintptr(&sl.lock, 0, 1)
}

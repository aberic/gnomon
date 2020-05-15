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
	"fmt"
	"sync"
	"testing"
	"time"
)

func testSpinLock(threads, n int, l sync.Locker) time.Duration {
	var wg sync.WaitGroup
	wg.Add(threads)

	var count1 int
	var count2 int

	start := time.Now()
	for i := 0; i < threads; i++ {
		go func() {
			for i := 0; i < n; i++ {
				l.Lock()
				count1++
				count2 += 2
				l.Unlock()
			}
			wg.Done()
		}()
	}
	wg.Wait()
	dur := time.Since(start)
	if count1 != threads*n {
		panic("mismatch")
	}
	if count2 != threads*n*2 {
		panic("mismatch")
	}
	return dur
}

func TestNewSpinLock(t *testing.T) {
	fmt.Printf("[1] spinlock %4.0fms\n", testSpinLock(1, 1000000, NewSpinLock()).Seconds()*1000)
	fmt.Printf("[1] mutex    %4.0fms\n", testSpinLock(1, 1000000, &sync.Mutex{}).Seconds()*1000)
	fmt.Printf("[4] spinlock %4.0fms\n", testSpinLock(4, 1000000, NewSpinLock()).Seconds()*1000)
	fmt.Printf("[4] mutex    %4.0fms\n", testSpinLock(4, 1000000, &sync.Mutex{}).Seconds()*1000)
	fmt.Printf("[8] spinlock %4.0fms\n", testSpinLock(8, 1000000, NewSpinLock()).Seconds()*1000)
	fmt.Printf("[8] mutex    %4.0fms\n", testSpinLock(8, 1000000, &sync.Mutex{}).Seconds()*1000)
}

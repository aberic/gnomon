/*
 * Copyright (c) 2020. Aberic - All Rights Reserved.
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

package grope

import (
	"testing"
)

func TestLimit(t *testing.T) {
	l := &Limit{
		LimitMillisecond:         1000,
		LimitCount:               5,
		LimitIntervalMillisecond: 100,
	}
	l.init()
	go l.limit()
	loop(l, t)
}

func loop(limit *Limit, t *testing.T) {
	i := 0
	for i <= 20 {
		t.Log("被堵住了 c len = ", len(limit.limitChan))
		limit.limitChan <- struct{}{}
		t.Log("被放行了 ", limit.times)
		i++
		//time.Sleep(100 * time.Millisecond)
	}
}

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
	"time"
)

// Limit 限流策略
type Limit struct {
	LimitMillisecond         int64         // 请求限定的时间段（毫秒）
	LimitCount               int           // 请求限定的时间段内允许的请求次数
	LimitIntervalMillisecond int64         // 请求允许的最小间隔时间（毫秒），0表示不限
	limitChan                chan struct{} // 限流通道
	times                    []int64       // 请求时间数组
}

// new 新建限流策略
func (l *Limit) init() {
	l.limitChan = make(chan struct{}, l.LimitCount)
	l.times = []int64{}
	for i := 0; i < l.LimitCount; i++ {
		time.Sleep(10 * time.Nanosecond)
		l.add(time.Now().UnixNano() / 1e6)
	}
}

func (l *Limit) limit() {
	for {
		timeNow := time.Now().UnixNano() / 1e6
		// 如果当前时间与时间数组第一时间差大于限定时间段，并且当前时间与时间数组最后时间差大于最小请求间隔，则放行新的请求
		if timeNow-l.times[0] > l.LimitMillisecond && timeNow-l.times[len(l.times)-1] > l.LimitIntervalMillisecond {
			<-l.limitChan // 取出一个元素，放行
			l.resetTimes(time.Now().UnixNano() / 1e6)
		} else {
			time.Sleep(10 * time.Nanosecond)
		}
	}
}

// add 新增一个元素
func (l *Limit) add(time int64) {
	if len(l.times) < l.LimitCount {
		l.times = append(l.times, time)
	}
}

// remove 移除第一个元素
//func (l *Limit) remove() {
//	l.times = l.times[1:]
//}

func (l *Limit) resetTimes(time int64) {
	l.times = l.times[1:]
	l.times = append(l.times, time)
}

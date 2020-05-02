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

package log

import (
	"os"
	"sync"
	"time"
)

// filed 日志文件操作对象
type filed struct {
	fileIndex string       // fileIndex 日志文件相同日期编号，根据文件新建规则确定
	file      *os.File     // 日志文件对象
	tasks     chan string  // 任务队列，默认1000个缓存
	lock      sync.RWMutex // lock 每次做io开销的安全锁
}

// running 循环执行文件写入，默认60秒超时
func (f *filed) running() {
	to := time.NewTimer(60 * time.Second)
	for {
		select {
		case task := <-f.tasks:
			f.lock.RLock()
			to.Reset(time.Second)
			if _, err := f.file.WriteString(task); nil != err {
				panic(err)
			}
			f.lock.RUnlock()
		case <-to.C:
			_ = f.file.Close()
			f.file = nil
			return
		}
	}
}

// FieldInter field 接口
type FieldInter interface {
	GetKey() string
	GetValue() interface{}
}

// field 日志输出子集对象
type field struct {
	key   string
	value interface{}
}

// GetKey 获取日志打印 pairs key
func (f *field) GetKey() string {
	return f.key
}

// GetValue 获取日志打印 pairs value
func (f *field) GetValue() interface{} {
	return f.value
}

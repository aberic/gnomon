/*
 * Copyright (c) 2019. Aberic - All Rights Reserved.
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
	"github.com/panjf2000/ants"
	"sync"
)

var (
	pc       *poolCommon
	oncePool sync.Once
)

func pool() *poolCommon {
	oncePool.Do(func() {
		pc = &poolCommon{}
		pc.init(1000)
	})
	return pc
}

type poolCommon struct {
	pool *ants.Pool
	once sync.Once
}

func (p *poolCommon) init(size int) {
	p.once.Do(func() {
		p.pool, _ = ants.NewPool(size)
	})
}

// tune 动态变更协程池数量
func (p *poolCommon) tune(poolSize int) {
	p.pool.Tune(poolSize)
}

func (p *poolCommon) submit(task func()) error {
	return p.pool.Submit(func() {
		task()
	})
}

func (p *poolCommon) submitObj(i interface{}, task func(i interface{})) error {
	return p.pool.Submit(func() {
		task(i)
	})
}

func (p *poolCommon) submitField(
	task func(timeString, fileString, stackString, levelName, msg string, level Level, fields ...*field),
	timeString, fileString, stackString, levelName, msg string, level Level, fields ...*field) error {
	return p.pool.Submit(func() {
		task(timeString, fileString, stackString, levelName, msg, level, fields...)
	})
}

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
	"errors"
	"io"
	"sync"
)

var errPoolClosed = errors.New("pond closed")

// Conn 连接单体接口
type Conn interface {
	io.Closer // 实现io.Closer接口的对象都可以使用该连接池
}

// factory 创建连接的方法
type factory func() (Conn, error)

// NewPond 新建一个支持所有实现 io.Closer 接口的连接池
//
// minOpen 池中最少资源数
//
// maxOpen 池中最大资源数
//
// factory
func NewPond(minOpen, maxOpen int, factory factory) *Pond {
	if maxOpen <= 0 {
		maxOpen = 5
	}
	if minOpen > maxOpen {
		maxOpen = minOpen + 1
	}
	p := &Pond{
		maxOpen: maxOpen,
		minOpen: minOpen,
		factory: factory,
		conn:    make(chan Conn, maxOpen),
	}

	for i := 0; i < minOpen; i++ {
		connect, err := factory()
		if err != nil {
			continue
		}
		p.nowOpen++
		p.conn <- connect
	}
	return p
}

// Pond 连接池对象
type Pond struct {
	sync.Mutex
	conn    chan Conn
	maxOpen int     // 池中最大资源数
	nowOpen int     // 当前池中资源数
	minOpen int     // 池中最少资源数
	closed  bool    // 池是否已关闭
	factory factory // 创建连接的方法
}

func (p *Pond) getOrCreate() (Conn, error) {
	//select {
	//case connect := <-p.Pond:
	//	return connect, nil
	//default:
	//}
	defer p.Unlock()
	p.Lock()
	if p.nowOpen >= p.maxOpen {
		return <-p.conn, nil
	}
	// 新建连接
	connect, err := p.factory()
	if err != nil {
		return nil, err
	}
	p.nowOpen++
	return connect, nil
}

// Acquire 获取资源
func (p *Pond) Acquire() (Conn, error) {
	if p.closed {
		return nil, errPoolClosed
	}
	//for {
	connect, err := p.getOrCreate()
	if err != nil {
		return nil, err
	}
	//// 如果设置了超时且当前连接的活跃时间+超时时间早于现在，则当前连接已过期
	//if p.maxLifetime > 0 && connect.lastActiveTime().Add(p.maxLifetime).Before(time.Now()) {
	//	_ = p.close(connect)
	//	continue
	//}
	return connect, nil
	//}
}

// Release 释放单个资源到连接池
func (p *Pond) Release(conn Conn) error {
	if p.closed {
		return errPoolClosed
	}
	if len(p.conn) < p.minOpen {
		p.conn <- conn
	} else {
		p.Close(conn)
	}
	return nil
}

// Close 关闭单个资源
func (p *Pond) Close(conn Conn) {
	p.Lock()
	_ = conn.Close()
	p.nowOpen--
	p.Unlock()
}

// Shutdown 关闭连接池，释放所有资源
func (p *Pond) Shutdown() error {
	if p.closed {
		return errPoolClosed
	}
	p.Lock()
	close(p.conn)
	for connect := range p.conn {
		_ = connect.Close()
		p.nowOpen--
	}
	p.closed = true
	p.Unlock()
	return nil
}

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
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/peer"
	"net"
	"strconv"
	"strings"
	"sync"
)

var (
	connections = map[string]*grpc.ClientConn{}
	muConn      sync.Mutex
	reqs        map[string]*Pond
	mu          sync.Mutex
)

func init() {
	reqs = map[string]*Pond{}
}

// Business 真实业务逻辑
type Business func(conn *grpc.ClientConn) (interface{}, error)

// GRPCRequest RPC 通过rpc进行通信 protoc --go_out=plugins=grpc:. grpc/proto/*.proto
func GRPCRequest(url string, business Business) (interface{}, error) {
	var (
		conn *grpc.ClientConn
		err  error
	)
	// 创建一个grpc连接器
	if conn, err = grpc.Dial(url, grpc.WithInsecure()); nil != err {
		return nil, err
	}
	// 请求完毕后关闭连接
	defer func() { _ = conn.Close() }()
	return business(conn)
}

// GRPCRequestSingleConn RPC 通过rpc进行通信 protoc --go_out=plugins=grpc:. grpc/proto/*.proto
func GRPCRequestSingleConn(url string, business Business) (interface{}, error) {
	return business(getGRPCConn(url))
}

// GRPCRequestPools 通过rpc进行通信 protoc --go_out=plugins=grpc:. grpc/proto/*.proto
func GRPCRequestPools(url string, business Business) (interface{}, error) {
	var (
		c    Conn
		conn *grpc.ClientConn
		pond = reqs[url]
		err  error
	)
	if nil == pond {
		mu.Lock()
		pond = NewPond(1, 10, func() (conn Conn, e error) {
			return grpc.Dial(url, grpc.WithInsecure())
		})
		reqs[url] = pond
		mu.Unlock()
	}
	for {
		// 创建一个grpc连接器
		if c, err = pond.Acquire(); nil != err {
			return nil, err
		}
		conn = c.(*grpc.ClientConn)
		if conn.GetState() != connectivity.Shutdown && conn.GetState() != connectivity.TransientFailure {
			break
		}
		pond.Close(c)
	}
	// 请求完毕后释放连接
	defer func() { _ = pond.Release(c) }()
	return business(conn)
}

// GRPCRequestPool 通过rpc进行通信 protoc --go_out=plugins=grpc:. grpc/proto/*.proto
func GRPCRequestPool(pool *Pond, business Business) (interface{}, error) {
	var (
		c    Conn
		conn *grpc.ClientConn
		err  error
	)
	for {
		// 创建一个grpc连接器
		if c, err = pool.Acquire(); nil != err {
			return nil, err
		}
		conn = c.(*grpc.ClientConn)
		if conn.GetState() != connectivity.Shutdown && conn.GetState() != connectivity.TransientFailure {
			break
		}
		pool.Close(c)
	}
	// 请求完毕后释放连接
	defer func() { _ = pool.Release(c) }()
	return business(conn)
}

// GRPCGetClientIP 取出gRPC客户端的ip地址和端口号
//
// string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
func GRPCGetClientIP(ctx context.Context) (address string, port int, err error) {
	var (
		pr *peer.Peer
		ok bool
	)
	if pr, ok = peer.FromContext(ctx); !ok {
		err = fmt.Errorf("[getGRPCClientIP] invoke FromContext() failed")
		return
	}
	if pr.Addr == net.Addr(nil) {
		err = fmt.Errorf("[getGRPCClientIP] peer.Addr is nil")
		return
	}
	addSlice := strings.Split(pr.Addr.String(), ":")
	address = addSlice[0]
	if port, err = strconv.Atoi(addSlice[1]); nil != err {
		err = fmt.Errorf("[getGRPCClientIP] peer.Addr.port parse int error, err: %s", err.Error())
		return
	}
	return
}

func getGRPCConn(url string) *grpc.ClientConn {
	if conn, ok := connections[url]; ok && conn.GetState() != connectivity.Shutdown && conn.GetState() != connectivity.TransientFailure {
		return conn
	}
	defer muConn.Unlock()
	muConn.Lock()
	// 创建一个grpc连接器
	if conn, err := grpc.Dial(url, grpc.WithInsecure()); nil != err {
		panic(err)
	} else {
		connections[url] = conn
		return conn
	}
}

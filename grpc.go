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
)

// GRPCCommon gRPC工具
type GRPCCommon struct{}

type Business func(conn *grpc.ClientConn) (interface{}, error)

// RPC 通过rpc进行通信 protoc --go_out=plugins=grpc:. grpc/proto/*.proto
func (g *GRPCCommon) Request(url string, business Business) (interface{}, error) {
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

// RPCPool 通过rpc进行通信 protoc --go_out=plugins=grpc:. grpc/proto/*.proto
func (g *GRPCCommon) RequestPool(pool *Pond, business Business) (interface{}, error) {
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

// GetClientIP 取出gRPC客户端的ip地址和端口号
//
// string form of address (for example, "192.0.2.1:25", "[2001:db8::1]:80")
func (g *GRPCCommon) GetClientIP(ctx context.Context) (address string, port int, err error) {
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

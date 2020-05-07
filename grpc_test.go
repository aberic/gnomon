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

package gnomon

import (
	"google.golang.org/grpc"
	"testing"
)

func TestGRPCRequest(t *testing.T) {
	_, _ = GRPCRequest("http://127.0.0.1", func(conn *grpc.ClientConn) (i interface{}, err error) {
		return nil, nil
	})
}

func TestGRPCRequestSingleConn(t *testing.T) {
	_, _ = GRPCRequestSingleConn("http://127.0.0.1", func(conn *grpc.ClientConn) (i interface{}, err error) {
		return nil, nil
	})
}

func TestGRPCRequestPool(t *testing.T) {
	_, _ = GRPCRequestPool(NewPond(1, 5, func() (conn Conn, e error) {
		return grpc.Dial("http://127.0.0.1", grpc.WithInsecure())
	}), func(conn *grpc.ClientConn) (i interface{}, err error) {
		return nil, nil
	})
}

func TestGRPCRequestPools(t *testing.T) {
	_, _ = GRPCRequestPools("http://127.0.0.1", func(conn *grpc.ClientConn) (i interface{}, err error) {
		return nil, nil
	})
}

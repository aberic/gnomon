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

package rpc

import (
	"google.golang.org/grpc"
	"testing"
	"time"
)

func TestNewPond(t *testing.T) {
	t.Log(NewPond(10, 100, 5*time.Second, func() (conn Conn, e error) {
		return grpc.Dial("http://wwww.gnomon.com", grpc.WithInsecure())
	}))
}

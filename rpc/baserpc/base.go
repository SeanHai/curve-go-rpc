/*
*  Copyright (c) 2023 NetEase Inc.
*
*  Licensed under the Apache License, Version 2.0 (the "License");
*  you may not use this file except in compliance with the License.
*  You may obtain a copy of the License at
*
*      http://www.apache.org/licenses/LICENSE-2.0
*
*  Unless required by applicable law or agreed to in writing, software
*  distributed under the License is distributed on an "AS IS" BASIS,
*  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
*  See the License for the specific language governing permissions and
*  limitations under the License.
 */

/*
* Project: Curve-Go-RPC
* Created Date: 2023-03-03
* Author: wanghai (SeanHai)
 */

package baserpc

import (
	"context"
	"fmt"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	pool "github.com/processout/grpc-go-pool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	POOL_INIT_CONN_NUM = 0
	POOL_CONN_CAPACITY = 10
	IDLE_TIMEOUT_SEC   = 10
)

type BaseRpc struct {
	Timeout    time.Duration
	RetryTimes uint32
	ConnPools  map[string]*pool.Pool
}

type RpcContext struct {
	addrs []string // endpoint: 127.0.0.1:6666
	name  string
}

type RpcResult struct {
	Key    interface{}
	Err    error
	Result interface{}
}

type Rpc interface {
	NewRpcClient(cc grpc.ClientConnInterface)
	Stub_Func(ctx context.Context, opt ...grpc.CallOption) (interface{}, error)
}

func NewRpcContext(addrs []string, funcName string) *RpcContext {
	return &RpcContext{
		addrs: addrs,
		name:  funcName,
	}
}

func NewBaseRpc(addrs []string, timeoutMs int, retryTimes uint32) (*BaseRpc) {
	client := &BaseRpc{
		Timeout:    time.Duration(timeoutMs * int(time.Millisecond)),
		RetryTimes: retryTimes,
		ConnPools:  make(map[string]*pool.Pool),
	}
	for _, addr := range addrs {
		p, err := pool.New(func() (*grpc.ClientConn, error) { return client.createConn(addr) },
			POOL_INIT_CONN_NUM, POOL_CONN_CAPACITY, time.Duration(IDLE_TIMEOUT_SEC*int(time.Second)))
		if err == nil {
			client.ConnPools[addr] = p
		}
	}
	return client
}

func (cli *BaseRpc) createConn(addr string) (*grpc.ClientConn, error) {
	ctx, cancel := context.WithTimeout(context.Background(), cli.Timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (cli *BaseRpc) SendRpc(opt *RpcContext, rpcFunc Rpc) *RpcResult {
	size := len(opt.addrs)
	if size == 0 {
		return &RpcResult{
			Key:    "",
			Err:    fmt.Errorf("empty addr"),
			Result: nil,
		}
	}
	results := make(chan RpcResult, size)
	for _, addr := range opt.addrs {
		go func(address string) {
			var conn *grpc.ClientConn
			var err error
			ctx, cancel := context.WithTimeout(context.Background(), cli.Timeout)
			defer cancel()
			if _, ok := cli.ConnPools[address]; !ok {
				conn, err = cli.createConn(address)
				if err != nil {
					results <- RpcResult{
						Key:    address,
						Err:    err,
						Result: nil,
					}
					return
				}
				defer conn.Close()
			} else {
				cconn, err := cli.ConnPools[address].Get(ctx)
				if err != nil {
					results <- RpcResult{
						Key:    address,
						Err:    err,
						Result: nil,
					}
					return
				}
				conn = cconn.ClientConn
				defer cconn.Close()
			}
			rpcFunc.NewRpcClient(conn)
			res, err := rpcFunc.Stub_Func(ctx, grpc_retry.WithMax(uint(cli.RetryTimes)),
				grpc_retry.WithCodes(codes.Unknown, codes.Unavailable, codes.DeadlineExceeded))
			results <- RpcResult{
				Key:    address,
				Err:    err,
				Result: res,
			}
		}(addr)
	}
	count := 0
	var rpcErr string
	for res := range results {
		if res.Err == nil {
			return &res
		}
		count++
		rpcErr = fmt.Sprintf("%s;%s:%s", rpcErr, res.Key, res.Err.Error())
		if count >= size {
			break
		}
	}
	return &RpcResult{
		Key:    "",
		Err:    fmt.Errorf(rpcErr),
		Result: nil,
	}
}

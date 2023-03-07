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
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
)

type BaseRpc struct {
	Timeout    time.Duration
	RetryTimes uint32
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

func (cli *BaseRpc) getOrCreateConn(addr string, ctx context.Context) (*grpc.ClientConn, error) {
	// TODO: conn pool maybe needed
	ctx, cancel := context.WithTimeout(context.Background(), cli.Timeout)
	defer cancel()
	conn, err := grpc.DialContext(ctx, addr, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return nil, err
	}
	return conn, nil
}

func (cli *BaseRpc) SendRpc(ctx *RpcContext, rpcFunc Rpc) *RpcResult {
	size := len(ctx.addrs)
	if size == 0 {
		return &RpcResult{
			Key:    "",
			Err:    fmt.Errorf("empty addr"),
			Result: nil,
		}
	}
	results := make(chan RpcResult, size)
	for _, addr := range ctx.addrs {
		go func(address string) {
			ctx, cancel := context.WithTimeout(context.Background(), cli.Timeout)
			defer cancel()
			conn, err := cli.getOrCreateConn(address, ctx)
			if err != nil {
				results <- RpcResult{
					Key:    address,
					Err:    err,
					Result: nil,
				}
			} else {
				rpcFunc.NewRpcClient(conn)
				res, err := rpcFunc.Stub_Func(ctx, grpc_retry.WithMax(uint(cli.RetryTimes)),
					grpc_retry.WithCodes(codes.Unknown, codes.Unavailable, codes.DeadlineExceeded))
				results <- RpcResult{
					Key:    address,
					Err:    err,
					Result: res,
				}
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

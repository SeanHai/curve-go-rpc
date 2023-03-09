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

package curvebs

import (
	"context"
	"fmt"
	"net"
	"os"
	"testing"

	"github.com/SeanHai/curve-go-rpc/curvebs_proto/proto/topology"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":12900"
)

var (
	gs *grpc.Server

	status_success     int32  = 0
	physical_pool_id   uint32 = 1
	physical_pool_name string = "physical_pool"

	clientOption MdsClientOption = MdsClientOption{
		TimeoutMs:  500,
		RetryTimes: 3,
		Addrs:      []string{"127.0.0.1:12900"},
	}
)

type server struct {
	topology.UnimplementedTopologyServiceServer
}

func (s *server) CreateLogicalPool(ctx context.Context, req *topology.CreateLogicalPoolRequest) (
	*topology.CreateLogicalPoolResponse, error) {
	return nil, nil
}

func (s *server) CreatePhysicalPool(ctx context.Context, req *topology.PhysicalPoolRequest) (
	*topology.PhysicalPoolResponse, error) {
	return nil, nil
}

func (s *server) CreateZone(ctx context.Context, req *topology.ZoneRequest) (
	*topology.ZoneResponse, error) {
	return nil, nil
}

func (s *server) DeleteChunkServer(ctx context.Context, req *topology.DeleteChunkServerRequest) (
	*topology.DeleteChunkServerResponse, error) {
	return nil, nil
}

func (s *server) DeleteLogicalPool(ctx context.Context, req *topology.DeleteLogicalPoolRequest) (
	*topology.DeleteLogicalPoolResponse, error) {
	return nil, nil
}

func (s *server) DeletePhysicalPool(ctx context.Context, req *topology.PhysicalPoolRequest) (
	*topology.PhysicalPoolResponse, error) {
	return nil, nil
}

func (s *server) DeleteServer(ctx context.Context, req *topology.DeleteServerRequest) (
	*topology.DeleteServerResponse, error) {
	return nil, nil
}

func (s *server) DeleteZone(ctx context.Context, req *topology.ZoneRequest) (
	*topology.ZoneResponse, error) {
	return nil, nil
}

func (s *server) GetChunkServer(ctx context.Context, req *topology.GetChunkServerInfoRequest) (
	*topology.GetChunkServerInfoResponse, error) {
	return nil, nil
}

func (s *server) GetChunkServerInCluster(ctx context.Context, req *topology.GetChunkServerInClusterRequest) (
	*topology.GetChunkServerInClusterResponse, error) {
	return nil, nil
}

func (s *server) GetChunkServerListInCopySets(ctx context.Context, req *topology.GetChunkServerListInCopySetsRequest) (
	*topology.GetChunkServerListInCopySetsResponse, error) {
	return nil, nil
}

func (s *server) GetClusterInfo(ctx context.Context, req *topology.GetClusterInfoRequest) (
	*topology.GetClusterInfoResponse, error) {
	return nil, nil
}

func (s *server) GetCopySetsInChunkServer(ctx context.Context, req *topology.GetCopySetsInChunkServerRequest) (
	*topology.GetCopySetsInChunkServerResponse, error) {
	return nil, nil
}

func (s *server) GetCopySetsInCluster(ctx context.Context, req *topology.GetCopySetsInClusterRequest) (
	*topology.GetCopySetsInClusterResponse, error) {
	return nil, nil
}

func (s *server) GetCopyset(ctx context.Context, req *topology.GetCopysetRequest) (
	*topology.GetCopysetResponse, error) {
	return nil, nil
}

func (s *server) GetLogicalPool(ctx context.Context, req *topology.GetLogicalPoolRequest) (
	*topology.GetLogicalPoolResponse, error) {
	return nil, nil
}

func (s *server) GetPhysicalPool(ctx context.Context, req *topology.PhysicalPoolRequest) (
	*topology.PhysicalPoolResponse, error) {
	return nil, nil
}

func (s *server) GetServer(ctx context.Context, req *topology.GetServerRequest) (
	*topology.GetServerResponse, error) {
	return nil, nil
}

func (s *server) GetZone(ctx context.Context, req *topology.ZoneRequest) (
	*topology.ZoneResponse, error) {
	return nil, nil
}

func (s *server) ListChunkServer(ctx context.Context, req *topology.ListChunkServerRequest) (
	*topology.ListChunkServerResponse, error) {
	return nil, nil
}

func (s *server) ListLogicalPool(ctx context.Context, req *topology.ListLogicalPoolRequest) (
	*topology.ListLogicalPoolResponse, error) {
	return nil, nil
}

func (s *server) ListPoolZone(ctx context.Context, req *topology.ListPoolZoneRequest) (
	*topology.ListPoolZoneResponse, error) {
	return nil, nil
}

func (s *server) ListUnAvailCopySets(ctx context.Context, req *topology.ListUnAvailCopySetsRequest) (
	*topology.ListUnAvailCopySetsResponse, error) {
	return nil, nil
}

func (s *server) ListZoneServer(ctx context.Context, req *topology.ListZoneServerRequest) (
	*topology.ListZoneServerResponse, error) {
	return nil, nil
}

func (s *server) RegistChunkServer(ctx context.Context, req *topology.ChunkServerRegistRequest) (
	*topology.ChunkServerRegistResponse, error) {
	return nil, nil
}

func (s *server) RegistServer(ctx context.Context, req *topology.ServerRegistRequest) (
	*topology.ServerRegistResponse, error) {
	return nil, nil
}

func (s *server) SetChunkServer(ctx context.Context, req *topology.SetChunkServerStatusRequest) (
	*topology.SetChunkServerStatusResponse, error) {
	return nil, nil
}

func (s *server) SetCopysetsAvailFlag(ctx context.Context, req *topology.SetCopysetsAvailFlagRequest) (
	*topology.SetCopysetsAvailFlagResponse, error) {
	return nil, nil
}

func (s *server) SetLogicalPool(ctx context.Context, req *topology.SetLogicalPoolRequest) (
	*topology.SetLogicalPoolResponse, error) {
	return nil, nil
}

func (s *server) SetLogicalPoolScanState(ctx context.Context, req *topology.SetLogicalPoolScanStateRequest) (
	*topology.SetLogicalPoolScanStateResponse, error) {
	return nil, nil
}

func (s *server) ListPhysicalPool(ctx context.Context, req *topology.ListPhysicalPoolRequest) (
	*topology.ListPhysicalPoolResponse, error) {
	info := topology.PhysicalPoolInfo{
		PhysicalPoolID:   &physical_pool_id,
		PhysicalPoolName: &physical_pool_name,
	}
	return &topology.ListPhysicalPoolResponse{
		StatusCode:        &status_success,
		PhysicalPoolInfos: []*topology.PhysicalPoolInfo{&info},
	}, nil
}

func init() {
	go func() {
		lis, err := net.Listen("tcp", port)
		if err != nil {
			fmt.Printf("failed to listen: %v", err)
		}
		gs = grpc.NewServer()
		topology.RegisterTopologyServiceServer(gs, &server{})
		// Register reflection service on gRPC server.
		reflection.Register(gs)
		if err := gs.Serve(lis); err != nil {
			fmt.Printf("failed to serve: %v", err)
		}
	}()
}

func teardown() {
	gs.Stop()
}

func TestListPhysicalPool(t *testing.T) {
	mdsClient := NewMdsClient(clientOption)
	pools, err := mdsClient.ListPhysicalPool()
	if err != nil {
		t.Errorf("TestListPhysicalPool rpc failed, error = %v", err)
	}
	if len(pools) != 1 {
		t.Errorf("TestListPhysicalPool response failed, expected size = 1, but actual len = %d", len(pools))
	}
	if pools[0].Id != physical_pool_id || pools[0].Name != physical_pool_name {
		t.Errorf("TestListPhysicalPool response failed, expected id = %d, name = %s; actual id = %d, name = %s",
			physical_pool_id, physical_pool_name, pools[0].Id, pools[0].Name)
	}
}

func TestMain(m *testing.M) {
	code := m.Run()
	teardown()
	os.Exit(code)
}

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
	"fmt"
	"time"

	"github.com/SeanHai/curve-go-rpc/proto/nameserver2"
	"github.com/SeanHai/curve-go-rpc/rpc/baserpc"
	"github.com/SeanHai/curve-go-rpc/rpc/common"
)

const (
	// file type
	INODE_DIRECTORY         = "INODE_DIRECTORY"
	INODE_PAGEFILE          = "INODE_PAGEFILE"
	INODE_APPENDFILE        = "INODE_APPENDFILE"
	INODE_APPENDECFILE      = "INODE_APPENDECFILE"
	INODE_SNAPSHOT_PAGEFILE = "INODE_SNAPSHOT_PAGEFILE"

	// file status
	FILE_CREATED             = "kFileCreated"
	FILE_DELETING            = "kFileDeleting"
	FILE_CLONING             = "kFileCloning"
	FILE_CLONEMETA_INSTALLED = "kFileCloneMetaInstalled"
	FILE_CLONED              = "kFileCloned"
	FILE_BEIING_CLONED       = "kFileBeingCloned"

	// throttle type
	IOPS_TOTAL = "IOPS_TOTAL"
	IOPS_READ  = "IOPS_READ"
	IOPS_WRITE = "IOPS_WRITE"
	BPS_TOTAL  = "BPS_TOTAL"
	BPS_READ   = "BPS_READ"
	BPS_WRITE  = "BPS_WRITE"

	// apis
	GET_FILE_ALLOC_SIZE_FUNC = "GetAllocatedSize"
	LIST_DIR_FUNC            = "ListDir"
	GET_FILE_INFO            = "GetFileInfo"
	GET_FILE_SIZE            = "GetFileSize"
)

type ThrottleParams struct {
	Type        string `json:"type"`
	Limit       uint64 `json:"limit"`
	Burst       uint64 `json:"burst"`
	BurstLength uint64 `json:"burstLength"`
}

type FileInfo struct {
	Id                   uint64           `json:"id"`
	FileName             string           `json:"fileName"`
	ParentId             uint64           `json:"parentId"`
	FileType             string           `json:"fileType"`
	Owner                string           `json:"owner"`
	ChunkSize            uint32           `json:"chunkSize"`
	SegmentSize          uint32           `json:"segmentSize"`
	Length               uint64           `json:"length"`
	AllocateSize         uint64           `json:"alloc"`
	Ctime                string           `json:"ctime"`
	SeqNum               uint64           `json:"seqNum"`
	FileStatus           string           `json:"fileStatus"`
	OriginalFullPathName string           `json:"originalFullPathName"`
	CloneSource          string           `json:"cloneSource"`
	CloneLength          uint64           `json:"cloneLength"`
	StripeUnit           uint64           `json:"stripeUnit"`
	StripeCount          uint64           `json:"stripeCount"`
	ThrottleParams       []ThrottleParams `json:"throttleParams"`
	Epoch                uint64           `json:"epoch"`
}

func (cli *MdsClient) GetFileAllocatedSize(filename string) (uint64, map[uint32]uint64, error) {
	Rpc := &GetFileAllocatedSize{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, GET_FILE_ALLOC_SIZE_FUNC)
	Rpc.Request = &nameserver2.GetAllocatedSizeRequest{
		FileName: &filename,
	}
	ret := cli.baseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return 0, nil, ret.Err
	}

	response := ret.Result.(*nameserver2.GetAllocatedSizeResponse)
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return 0, nil, fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	infos := make(map[uint32]uint64)
	for k, v := range response.GetAllocSizeMap() {
		infos[k] = v / common.GiB
	}
	return response.GetAllocatedSize() / common.GiB, infos, nil
}

func getFileType(t nameserver2.FileType) string {
	switch t {
	case nameserver2.FileType_INODE_DIRECTORY:
		return INODE_DIRECTORY
	case nameserver2.FileType_INODE_PAGEFILE:
		return INODE_PAGEFILE
	case nameserver2.FileType_INODE_APPENDFILE:
		return INODE_APPENDFILE
	case nameserver2.FileType_INODE_APPENDECFILE:
		return INODE_APPENDECFILE
	case nameserver2.FileType_INODE_SNAPSHOT_PAGEFILE:
		return INODE_SNAPSHOT_PAGEFILE
	default:
		return INVALID
	}
}

func getFileStatus(s nameserver2.FileStatus) string {
	switch s {
	case nameserver2.FileStatus_kFileCreated:
		return FILE_CREATED
	case nameserver2.FileStatus_kFileDeleting:
		return FILE_DELETING
	case nameserver2.FileStatus_kFileCloning:
		return FILE_CLONING
	case nameserver2.FileStatus_kFileCloneMetaInstalled:
		return FILE_CLONEMETA_INSTALLED
	case nameserver2.FileStatus_kFileCloned:
		return FILE_CLONED
	case nameserver2.FileStatus_kFileBeingCloned:
		return FILE_BEIING_CLONED
	default:
		return INVALID
	}
}

func getThrottleType(t nameserver2.ThrottleType) string {
	switch t {
	case nameserver2.ThrottleType_IOPS_TOTAL:
		return IOPS_TOTAL
	case nameserver2.ThrottleType_IOPS_READ:
		return IOPS_READ
	case nameserver2.ThrottleType_IOPS_WRITE:
		return IOPS_WRITE
	case nameserver2.ThrottleType_BPS_TOTAL:
		return BPS_TOTAL
	case nameserver2.ThrottleType_BPS_READ:
		return BPS_READ
	case nameserver2.ThrottleType_BPS_WRITE:
		return BPS_WRITE
	default:
		return INVALID
	}
}

func (cli *MdsClient) ListDir(filename, owner, sig string, date uint64) ([]FileInfo, error) {
	Rpc := &ListDir{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, LIST_DIR_FUNC)
	Rpc.Request = &nameserver2.ListDirRequest{
		FileName: &filename,
		Owner:    &owner,
		Date:     &date,
	}
	if sig != "" {
		Rpc.Request.Signature = &sig
	}
	ret := cli.baseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return nil, ret.Err
	}

	response := ret.Result.(*nameserver2.ListDirResponse)
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return nil, fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	infos := []FileInfo{}
	for _, v := range response.GetFileInfo() {
		var info FileInfo
		info.Id = v.GetId()
		info.FileName = v.GetFileName()
		info.ParentId = v.GetParentId()
		info.FileType = getFileType(v.GetFileType())
		info.Owner = v.GetOwner()
		info.ChunkSize = v.GetChunkSize()
		info.SegmentSize = v.GetSegmentSize()
		info.Length = v.GetLength() / common.GiB
		info.Ctime = time.Unix(int64(v.GetCtime()/1000000), 0).Format(common.TIME_FORMAT)
		info.SeqNum = v.GetSeqNum()
		info.FileStatus = getFileStatus(v.GetFileStatus())
		info.OriginalFullPathName = v.GetOriginalFullPathName()
		info.CloneSource = v.GetCloneSource()
		info.CloneLength = v.GetCloneLength()
		info.StripeUnit = v.GetStripeUnit()
		info.StripeCount = v.GetStripeCount()
		info.ThrottleParams = []ThrottleParams{}
		for _, p := range v.GetThrottleParams().GetThrottleParams() {
			var param ThrottleParams
			param.Type = getThrottleType(p.GetType())
			param.Limit = p.GetLimit()
			param.Burst = p.GetBurst()
			param.BurstLength = p.GetBurstLength()
			info.ThrottleParams = append(info.ThrottleParams, param)
		}
		info.Epoch = v.GetEpoch()
		infos = append(infos, info)
	}
	return infos, nil
}

func (cli *MdsClient) GetFileInfo(filename, owner, sig string, date uint64) (FileInfo, error) {
	info := FileInfo{}
	Rpc := &GetFileInfo{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, GET_FILE_INFO)
	Rpc.Request = &nameserver2.GetFileInfoRequest{
		FileName: &filename,
		Owner:    &owner,
		Date:     &date,
	}
	if sig != "" {
		Rpc.Request.Signature = &sig
	}
	ret := cli.baseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return info, ret.Err
	}

	response := ret.Result.(*nameserver2.GetFileInfoResponse)
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return info, fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	v := response.GetFileInfo()
	info.Id = v.GetId()
	info.FileName = v.GetFileName()
	info.ParentId = v.GetParentId()
	info.FileType = getFileType(v.GetFileType())
	info.Owner = v.GetOwner()
	info.ChunkSize = v.GetChunkSize()
	info.SegmentSize = v.GetSegmentSize()
	info.Length = v.GetLength() / common.GiB
	info.Ctime = time.Unix(int64(v.GetCtime()/1000000), 0).Format(common.TIME_FORMAT)
	info.SeqNum = v.GetSeqNum()
	info.FileStatus = getFileStatus(v.GetFileStatus())
	info.OriginalFullPathName = v.GetOriginalFullPathName()
	info.CloneSource = v.GetCloneSource()
	info.CloneLength = v.GetCloneLength()
	info.StripeUnit = v.GetStripeUnit()
	info.StripeCount = v.GetStripeCount()
	info.ThrottleParams = []ThrottleParams{}
	for _, p := range v.GetThrottleParams().GetThrottleParams() {
		var param ThrottleParams
		param.Type = getThrottleType(p.GetType())
		param.Limit = p.GetLimit()
		param.Burst = p.GetBurst()
		param.BurstLength = p.GetBurstLength()
		info.ThrottleParams = append(info.ThrottleParams, param)
	}
	info.Epoch = v.GetEpoch()
	return info, nil
}

func (cli *MdsClient) GetFileSize(fileName string) (uint64, error) {
	var size uint64
	Rpc := &GetFileSize{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, GET_FILE_SIZE)
	Rpc.Request = &nameserver2.GetFileSizeRequest{
		FileName: &fileName,
	}
	ret := cli.baseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return size, ret.Err
	}

	response := ret.Result.(*nameserver2.GetFileSizeResponse)
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return size, fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	size = response.GetFileSize() / common.GiB
	return size, nil
}

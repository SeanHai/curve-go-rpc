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

	"github.com/SeanHai/curve-go-rpc/curvebs_proto/proto/nameserver2"
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
	GET_FILE_ALLOC_SIZE_FUNC    = "GetAllocatedSize"
	LIST_DIR_FUNC               = "ListDir"
	GET_FILE_INFO               = "GetFileInfo"
	GET_FILE_SIZE               = "GetFileSize"
	DELETE_FILE                 = "DeleteFile"
	CREATE_FILE                 = "CreateFile"
	EXTEND_FILE                 = "ExtendFile"
	RECOVER_FILE                = "RecoverFile"
	UPDATE_FILE_THROTTLE_PARAMS = "UpdateFileThrottleParams"
	FIND_FILE_MOUNTPOINT        = "FindFileMountPoint"
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
	MountPoints          []string         `json:"mountPoints"`
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

func getFileType(t string) nameserver2.FileType {
	switch t {
	case INODE_DIRECTORY:
		return nameserver2.FileType_INODE_DIRECTORY
	case INODE_PAGEFILE:
		return nameserver2.FileType_INODE_PAGEFILE
	case INODE_APPENDFILE:
		return nameserver2.FileType_INODE_APPENDFILE
	case INODE_APPENDECFILE:
		return nameserver2.FileType_INODE_APPENDECFILE
	case INODE_SNAPSHOT_PAGEFILE:
		return nameserver2.FileType_INODE_SNAPSHOT_PAGEFILE
	default:
		return -1
	}
}

func getFileTypeStr(t nameserver2.FileType) string {
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

func getThrottleTypeStr(t nameserver2.ThrottleType) string {
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

func getThrottleType(t string) nameserver2.ThrottleType {
	switch t {
	case IOPS_TOTAL:
		return nameserver2.ThrottleType_IOPS_TOTAL
	case IOPS_READ:
		return nameserver2.ThrottleType_IOPS_READ
	case IOPS_WRITE:
		return nameserver2.ThrottleType_IOPS_WRITE
	case BPS_TOTAL:
		return nameserver2.ThrottleType_BPS_TOTAL
	case BPS_READ:
		return nameserver2.ThrottleType_BPS_READ
	case BPS_WRITE:
		return nameserver2.ThrottleType_BPS_WRITE
	default:
		return 0
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
		info.FileType = getFileTypeStr(v.GetFileType())
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
			param.Type = getThrottleTypeStr(p.GetType())
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
	info.FileType = getFileTypeStr(v.GetFileType())
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
		param.Type = getThrottleTypeStr(p.GetType())
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

func (cli *MdsClient) DeleteFile(filename, owner, sig string, fileId, date uint64, forceDelete bool) error {
	Rpc := &DeleteFile{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, DELETE_FILE)
	Rpc.Request = &nameserver2.DeleteFileRequest{
		FileName:    &filename,
		Owner:       &owner,
		Date:        &date,
		ForceDelete: &forceDelete,
	}
	if sig != "" {
		Rpc.Request.Signature = &sig
	}
	if fileId != 0 {
		Rpc.Request.FileId = &fileId
	}

	ret := cli.baseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return ret.Err
	}
	response := ret.Result.(*nameserver2.DeleteFileResponse)
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	return nil
}

func (cli *MdsClient) RecoverFile(filename, owner, sig string, fileId, date uint64) error {
	Rpc := &RecoverFile{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, RECOVER_FILE)
	Rpc.Request = &nameserver2.RecoverFileRequest{
		FileName: &filename,
		Owner:    &owner,
		Date:     &date,
	}
	if sig != "" {
		Rpc.Request.Signature = &sig
	}
	if fileId != 0 {
		Rpc.Request.FileId = &fileId
	}

	ret := cli.baseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return ret.Err
	}
	response := ret.Result.(*nameserver2.RecoverFileResponse)
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	return nil
}

func (cli *MdsClient) CreateFile(filename, ftype, owner, sig string, length, date, stripeUnit, stripeCount uint64) error {
	Rpc := &CreateFile{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, CREATE_FILE)
	fileType := getFileType(ftype)
	Rpc.Request = &nameserver2.CreateFileRequest{
		FileName: &filename,
		FileType: &fileType,
		Owner:    &owner,
		Date:     &date,
	}
	if ftype != INODE_DIRECTORY {
		Rpc.Request.FileLength = &length
	}
	if sig != "" {
		Rpc.Request.Signature = &sig
	}
	if stripeCount != 0 && stripeUnit != 0 {
		Rpc.Request.StripeCount = &stripeCount
		Rpc.Request.StripeUnit = &stripeUnit
	}

	ret := cli.baseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return ret.Err
	}
	response := ret.Result.(*nameserver2.CreateFileResponse)
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	return nil
}

func (cli *MdsClient) ExtendFile(filename, owner, sig string, newSize, date uint64) error {
	Rpc := &ExtendFile{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, EXTEND_FILE)
	Rpc.Request = &nameserver2.ExtendFileRequest{
		FileName: &filename,
		NewSize:  &newSize,
		Owner:    &owner,
		Date:     &date,
	}
	if sig != "" {
		Rpc.Request.Signature = &sig
	}

	ret := cli.baseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return ret.Err
	}
	response := ret.Result.(*nameserver2.ExtendFileResponse)
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	return nil
}

func (cli *MdsClient) UpdateFileThrottleParams(filename, owner, sig string, date uint64, params ThrottleParams) error {
	Rpc := &UpdateFileThrottleParams{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, UPDATE_FILE_THROTTLE_PARAMS)
	burstType := getThrottleType(params.Type)
	Rpc.Request = &nameserver2.UpdateFileThrottleParamsRequest{
		FileName: &filename,
		Owner:    &owner,
		Date:     &date,
		ThrottleParams: &nameserver2.ThrottleParams{
			Type:        &burstType,
			Limit:       &params.Limit,
			Burst:       &params.Burst,
			BurstLength: &params.BurstLength,
		},
	}
	if sig != "" {
		Rpc.Request.Signature = &sig
	}

	ret := cli.baseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return ret.Err
	}
	response := ret.Result.(*nameserver2.UpdateFileThrottleParamsResponse)
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	return nil
}

func (cli *MdsClient) FindFileMountPoint(filename string) ([]string, error) {
	info := []string{}
	Rpc := &FindFileMountPoint{}
	Rpc.ctx = baserpc.NewRpcContext(cli.addrs, FIND_FILE_MOUNTPOINT)
	Rpc.Request = &nameserver2.FindFileMountPointRequest{
		FileName: &filename,
	}

	ret := cli.baseClient.SendRpc(Rpc.ctx, Rpc)
	if ret.Err != nil {
		return nil, ret.Err
	}
	response := ret.Result.(*nameserver2.FindFileMountPointResponse)
	statusCode := response.GetStatusCode()
	if statusCode != nameserver2.StatusCode_kOK {
		return info, fmt.Errorf(nameserver2.StatusCode_name[int32(statusCode)])
	}
	for _, v := range response.GetClientInfo() {
		info = append(info, fmt.Sprintf("%s:%d", v.GetIp(), v.GetPort()))
	}
	return info, nil
}

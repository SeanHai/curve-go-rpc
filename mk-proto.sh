bs_proto_path=external/_curve
fs_proto_path=external/_curve/curvefs

out_bs_proto_path="./curvebs_proto"
out_fs_proto_path="./curvefs_proto"

mkdir -p $out_bs_proto_path
mkdir -p $out_fs_proto_path

# curvebs proto
## common.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    $bs_proto_path/proto/common.proto

## chunk.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    $bs_proto_path/proto/chunk.proto

## chunkserver.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    $bs_proto_path/proto/chunkserver.proto

## cli.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    $bs_proto_path/proto/cli.proto

## cli2.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    --go_opt=Mproto/common.proto=github.com/SeanHai/curve-go-rpc/curvebs_proto/proto/common \
    $bs_proto_path/proto/cli2.proto

## configuration.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    $bs_proto_path/proto/configuration.proto

## copyset.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    --go_opt=Mproto/common.proto=github.com/SeanHai/curve-go-rpc/curvebs_proto/proto/common \
    $bs_proto_path/proto/copyset.proto

## curve_storage.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    $bs_proto_path/proto/curve_storage.proto

## scan.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    $bs_proto_path/proto/scan.proto

## heartbeat.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    --go_opt=Mproto/common.proto=github.com/SeanHai/curve-go-rpc/curvebs_proto/proto/common \
    --go_opt=Mproto/scan.proto=github.com/SeanHai/curve-go-rpc/curvebs_proto/proto/scan \
    $bs_proto_path/proto/heartbeat.proto

## integrity.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    $bs_proto_path/proto/integrity.proto

## namespace2.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    --go_opt=Mproto/common.proto=github.com/SeanHai/curve-go-rpc/curvebs_proto/proto/common \
   $bs_proto_path/proto/nameserver2.proto

## schedule.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    $bs_proto_path/proto/schedule.proto

## snapshotcloneserver.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    $bs_proto_path/proto/snapshotcloneserver.proto

## topology.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    --go_opt=Mproto/common.proto=github.com/SeanHai/curve-go-rpc/curvebs_proto/proto/common \
    $bs_proto_path/proto/topology.proto

## statuscode.proto
protoc --go_out=$out_bs_proto_path --proto_path=$bs_proto_path/tools-v2 \
    $bs_proto_path/tools-v2/internal/proto/curvebs/topology/statuscode.proto

protoc --go-grpc_out=$out_bs_proto_path --proto_path=$bs_proto_path \
    $bs_proto_path/proto/*.proto

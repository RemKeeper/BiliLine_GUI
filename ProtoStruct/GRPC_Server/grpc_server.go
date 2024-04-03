package GRPC_Server

import (
	GlobalType "BiliLine_Windows/Global"
	grpcServer "BiliLine_Windows/ProtoStruct/GenProtoStructGo"
	"context"
)

type server struct {
	grpcServer.DanmuServerServer
}

var msgCount = 0

func (s *server) SendDanmu(ctx context.Context, in *grpcServer.CmdDanmuData) (*grpcServer.Response, error) {
	GlobalType.Broad.Broadcast(GrpcDmTypeToSdkDmType(in))
	msgCount++
	return &grpcServer.Response{
		DataCount: uint64(msgCount),
		Status:    0,
	}, nil
}

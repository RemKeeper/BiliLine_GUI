package GRPC_Server

import (
	grpcServer "BiliLine_Windows/ProtoStruct/GenProtoStructGo"
	"github.com/vtb-link/bianka/proto"
)

// GrpcDmTypeToSdkDmType grpc弹幕类型转换为sdk弹幕类型
func GrpcDmTypeToSdkDmType(DanmuData *grpcServer.CmdDanmuData) *proto.CmdDanmuData {
	return &proto.CmdDanmuData{
		RoomID:                 int(DanmuData.RoomId),
		OpenID:                 DanmuData.OpenId,
		Uid:                    int(DanmuData.Uid),
		Uname:                  DanmuData.Uname,
		Msg:                    DanmuData.Msg,
		MsgID:                  DanmuData.MsgId,
		FansMedalLevel:         int(DanmuData.FansMedalLevel),
		FansMedalName:          DanmuData.FansMedalName,
		FansMedalWearingStatus: DanmuData.FansMedalWearingStatus,
		GuardLevel:             int(DanmuData.GuardLevel),
		Timestamp:              int(DanmuData.Timestamp),
		UFace:                  DanmuData.Uface,
		EmojiImgUrl:            DanmuData.EmojiImgUrl,
		DmType:                 int(DanmuData.DmType),
	}
}

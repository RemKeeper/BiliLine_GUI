package main

import (
	"fmt"
	"github.com/vtb-link/bianka/live"
	"github.com/vtb-link/bianka/proto"
	"log"
)

func messageHandle(msg *proto.Message) error {
	line.GuardIndex = make(map[int]int)
	line.GiftIndex = make(map[int]int)
	line.CommonIndex = make(map[int]int)

	lineTemp, err := GetLine()
	if err == nil && !lineTemp.IsEmpty() {
		line = lineTemp
	}

	// 单条消息raw 如果需要自己解析可以使用
	//log.Println(string(msg.Payload()))

	// sdk提供了自动解析消息的方法，可以快速解析为对应的cmd和data
	// 具体的cmd 可以参考 live/cmd.go
	cmd, data, err := live.AutomaticParsingMessageCommand(msg.Payload())
	if err != nil {
		return err
	}
	// 你可以使用cmd进行switch
	switch cmd {
	case live.CmdLiveOpenPlatformDanmu:
		//log.Println(cmd, data.(*live.CmdLiveOpenPlatformDanmuData))
		if globalConfiguration.IsOnlyGift {
			break
		}
		DanmuData := data.(*live.CmdLiveOpenPlatformDanmuData)

		if DanmuData.Msg == "取消排队" {
			DeleteLine(DanmuData.Uid)
		}

		if DanmuData.Msg != globalConfiguration.LineKey {
			break
		}
		uid := DanmuData.Uid

		if line.GuardIndex[uid] != 0 || line.GiftIndex[uid] != 0 || line.CommonIndex[uid] != 0 {
			break
		}
		switch {
		case DanmuData.GuardLevel <= 3 && DanmuData.GuardLevel != 0:
			fmt.Println(DanmuData)
			lineTemp := Line{
				Uid:        DanmuData.Uid,
				UserName:   DanmuData.Uname,
				Avatar:     DanmuData.UFace,
				PrintColor: globalConfiguration.GuardPrintColor,
			}
			line.GuardLine = append(line.GuardLine, lineTemp)
			line.GuardIndex[DanmuData.Uid] = len(line.GuardLine)
			SendLineToWs(lineTemp, GiftLine{}, GuardLineType)
			SetLine(line)

		case len(line.CommonLine) <= globalConfiguration.MaxLineCount:
			lineTemp := Line{
				Uid:        DanmuData.Uid,
				UserName:   DanmuData.Uname,
				Avatar:     DanmuData.UFace,
				PrintColor: globalConfiguration.CommonPrintColor,
			}
			line.CommonLine = append(line.CommonLine, lineTemp)
			line.CommonIndex[DanmuData.Uid] = len(line.CommonLine)
			SendLineToWs(lineTemp, GiftLine{}, CommonLineType)
			SetLine(line)
		}
	case live.CmdLiveOpenPlatformSendGift:
		//log.Println(cmd, data.(*live.CmdLiveOpenPlatformSendGiftData))
		GiftData := data.(*live.CmdLiveOpenPlatformSendGiftData)
		if !globalConfiguration.AutoJoinGiftLine {
			break
		}
		if line.GiftIndex[GiftData.Uid] != 0 {
			break
		}
		lineTemp := GiftLine{
			Uid:        GiftData.Uid,
			UserName:   GiftData.Uname,
			Avatar:     GiftData.Uface,
			PrintColor: globalConfiguration.GiftPrintColor,
			GiftPrice:  float64(GiftData.Price),
		}
		if (float64(GiftData.GiftNum*GiftData.Price))/1000 >= globalConfiguration.GiftLinePrice {
			line.GiftLine = append(line.GiftLine, lineTemp)
			line.GiftIndex[GiftData.Uid] = len(line.GiftLine)
			SendLineToWs(Line{}, lineTemp, GiftLineType)
			SetLine(line)
		}
	case live.CmdLiveOpenPlatformGuard:
		log.Println(cmd, data.(*live.CmdLiveOpenPlatformGuardData))

	}

	return nil
}

func RoomConnect(IdCode string) {
	//<-CloseConn
	sdkConfig := live.NewConfig(AccessKey, AccessSecret, AppID)

	// 创建sdk实例
	sdk := live.NewClient(sdkConfig)

	// app start
	startResp, err := sdk.AppStart(IdCode)
	if err != nil {
		panic(err)
	}
	RoomId <- startResp.AnchorInfo.RoomID

	// app end
	defer func() {
		//tk.Stop()
		sdk.AppEnd(startResp.GameInfo.GameID)
	}()

	dispatcherHandle := map[uint32]live.DispatcherHandle{
		proto.OperationMessage: messageHandle,
	}
	onCloseCallback := func(startResp *live.AppStartResponse) {
		log.Println("WebsocketClient onClose", startResp)
	}
	wsClient, err := sdk.StartWebsocket(startResp, dispatcherHandle, onCloseCallback)
	if err != nil {
		panic(err)
	}

	defer wsClient.Close()

	<-CloseConn
	log.Println("接收到退出信号")
	return

}

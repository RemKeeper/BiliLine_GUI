package main

import (
	"regexp"

	"golang.org/x/exp/slog"

	"github.com/vtb-link/bianka/basic"

	"github.com/vtb-link/bianka/live"
	"github.com/vtb-link/bianka/proto"
)

func messageHandle(ws *basic.WsClient, msg *proto.Message) error {
	line.GuardIndex = make(map[string]int)
	line.GiftIndex = make(map[string]int)
	line.CommonIndex = make(map[string]int)

	lineTemp, err := GetLine()
	if err == nil && !lineTemp.IsEmpty() {
		line = lineTemp
	}
	cmd, data, err := proto.AutomaticParsingMessageCommand(msg.Payload())
	if err != nil {
		return err
	}
	// 你可以使用cmd进行switch
	switch cmd {
	case proto.CmdLiveOpenPlatformDanmu:
		if globalConfiguration.IsOnlyGift {
			break
		}

		DanmuData := data.(*proto.CmdDanmuData)
		slog.Info(DanmuData.Uname, DanmuData.Msg)
		DanmuDataChan <- DanmuData

	case proto.CmdLiveOpenPlatformSendGift:
		GiftData := data.(*proto.CmdSendGiftData)

		if !globalConfiguration.AutoJoinGiftLine {
			break
		}
		if line.GiftIndex[GiftData.OpenID] != 0 {
			break
		}

		if !GiftData.Paid {
			break
		}

		lineTemp := GiftLine{
			OpenID: GiftData.OpenID,
			// OpenID:     strconv.Itoa(GiftData.Uid),
			UserName:   GiftData.Uname,
			Avatar:     GiftData.Uface,
			PrintColor: globalConfiguration.GiftPrintColor,
			GiftPrice:  float64(GiftData.Price * GiftData.GiftNum / 1000),
		}
		if (float64(GiftData.GiftNum*GiftData.Price))/1000 >= globalConfiguration.GiftLinePrice {
			line.GiftLine = append(line.GiftLine, lineTemp)
			line.GiftIndex[GiftData.OpenID] = len(line.GiftLine)
			SendLineToWs(Line{}, lineTemp, GiftLineType)
			SetLine(line)
		}
	case live.CmdLiveOpenPlatformGuard:
		slog.Info(cmd, data.(*proto.CmdGuardData))
	}

	return nil
}

func RoomConnect(IdCode string) (AppClient *live.Client, GameId string, WsClient *basic.WsClient, HeartbeatCloseChan chan bool) {
	//	初始化应用连接信息配置，自编译请申明以下3个值
	LinkConfig := live.NewConfig(AccessKey, AccessSecret, AppID)
	//	创建Api连接实例
	client := live.NewClient(LinkConfig)
	//	开始身份码认证流程

	AppStart, err := client.AppStart(IdCode)
	RoomId = AppStart.AnchorInfo.RoomID
	if err != nil {
		slog.Error("应用流程开启失败", err)
		return nil, "", nil, nil
	}
	// 开启心跳
	HeartbeatCloseChan = make(chan bool, 1)
	NewHeartbeat(client, AppStart.GameInfo.GameID, HeartbeatCloseChan)

	dispatcherHandleMap := basic.DispatcherHandleMap{
		proto.OperationMessage: messageHandle,
	}
	onCloseCallback := func(wcs *basic.WsClient, startResp basic.StartResp, closeType int) {
		slog.Info("WebsocketClient onClose", startResp)
		// 注意检查关闭类型, 避免无限重连
		if closeType == live.CloseReceivedShutdownMessage || closeType == live.CloseAuthFailed {
			slog.Info("WebsocketClient exit")
			return
		}
		err := wcs.Reconnection(startResp)
		if err != nil {
			slog.Error("Reconnection fail", err)
		}
	}
	// 一键开启websocket
	wsClient, err := basic.StartWebsocket(AppStart, dispatcherHandleMap, onCloseCallback, logger)
	if err != nil {
		panic(err)
	}
	return client, AppStart.GameInfo.GameID, wsClient, HeartbeatCloseChan
}

var KeyWordMatchMap = make(map[string]bool)

func KeyWordMatchInit(keyWord string) {
	reg := regexp.MustCompile(`[^.,!！；：’“‘”?？;:，。、-]+`)
	matches := reg.FindAllString(keyWord, -1)
	for _, match := range matches {
		KeyWordMatchMap[match] = true
	}
}

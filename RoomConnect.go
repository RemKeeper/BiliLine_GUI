package main

import (
	"errors"
	"fyne.io/fyne/v2/dialog"
	"regexp"

	"golang.org/x/exp/slog"

	"github.com/vtb-link/bianka/basic"

	"github.com/vtb-link/bianka/live"
	"github.com/vtb-link/bianka/proto"
)

func messageHandle(ws *basic.WsClient, msg *proto.Message) error {
	cmd, data, err := proto.AutomaticParsingMessageCommand(msg.Payload())
	if err != nil {
		return err
	}
	// 你可以使用cmd进行switch
	switch cmd {
	case proto.CmdLiveOpenPlatformDanmu:

		DanmuData := data.(*proto.CmdDanmuData)
		slog.Info(DanmuData.Uname, DanmuData.Msg)
		ResponseQueCtrl(DanmuData)

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

		//FindAndModifyDiscountGift(GiftData)

		lineTemp := GiftLine{
			OpenID: GiftData.OpenID,
			// OpenID:     strconv.Itoa(DiscountGiftData.Uid),
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

var (
	AccessSecret       = ""
	AppID        int64 = 123456
	AccessKey          = ""
)

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
		switch closeType {
		case live.CloseAuthFailed:
			dialog.ShowError(errors.New("请求鉴权失败，请检查身份码"), MainWindows)
		case live.CloseReadingConnError:
			dialog.ShowError(errors.New("读取链接错误,请重启客户端"), MainWindows)
		case live.CloseReceivedShutdownMessage:
			dialog.ShowError(errors.New("收到来自于B站的关闭消息，连接被关闭"), MainWindows)
		case live.CloseTypeUnknown:
			dialog.ShowError(errors.New("未知原因导致连接关闭"), MainWindows)
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

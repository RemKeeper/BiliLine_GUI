package main

import (
	"BiliLine_Windows/Global"
	"BiliLine_Windows/key"
	"fmt"
	"log"
	"os"
	"regexp"
	"time"

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

		Broadcast.Broadcast(DanmuData)

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

		lineTemp := GlobalType.GiftLine{
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
			SendLineToWs(GlobalType.Line{}, lineTemp, GlobalType.GiftLineType)
			SetLine(line)
		}
	case live.CmdLiveOpenPlatformGuard:
		log.Println(cmd, data.(*proto.CmdGuardData))
	}

	return nil
}

func RoomConnect(IdCode string) {
	sdkConfig := live.NewConfig(key.AccessKey, key.AccessSecret, key.AppID)
	// 创建sdk实例
	sdk := live.NewClient(sdkConfig)
	// app start
	startResp, err := sdk.AppStart(IdCode)
	if err != nil {
		panic(err)
	}

	fmt.Println("当前连接用户为", startResp.AnchorInfo)

	// 启用项目心跳 20s一次
	// see https://open-live.bilibili.com/document/eba8e2e1-847d-e908-2e5c-7a1ec7d9266f
	tk := time.NewTicker(time.Second * 20)
	go func(GameID string) {
		for {
			select {
			case <-tk.C:
				// 心跳
				if err := sdk.AppHeartbeat(GameID); err != nil {
					log.Println("Heartbeat fail", err)
				}
			}
		}
	}(startResp.GameInfo.GameID)
	RoomId = startResp.AnchorInfo.RoomID
	// app end
	defer func() {
		tk.Stop()
		sdk.AppEnd(startResp.GameInfo.GameID)
	}()

	dispatcherHandleMap := basic.DispatcherHandleMap{
		proto.OperationMessage: messageHandle,
	}

	// 关闭回调事件
	// 此事件会在websocket连接关闭后触发
	// 时序如下：
	// 0. send close message // 主动发送关闭消息
	// 1. close eventLoop // 不再处理任何消息
	// 2. close websocket // 关闭websocket连接
	// 3. onCloseCallback // 触发关闭回调事件
	// 增加了closeType 参数, 用于区分关闭类型
	onCloseCallback := func(wcs *basic.WsClient, startResp basic.StartResp, closeType int) {
		// 注册关闭回调
		log.Println("WebsocketClient onClose", startResp)

		// 注意检查关闭类型, 避免无限重连
		if closeType == live.CloseActively || closeType == live.CloseReceivedShutdownMessage || closeType == live.CloseAuthFailed {
			log.Println("WebsocketClient exit")
			os.Exit(0)
		}

		err := wcs.Reconnection(startResp)
		if err != nil {
			log.Println("Reconnection fail", err)
		}
	}

	wsClient, err := basic.StartWebsocket(startResp, dispatcherHandleMap, onCloseCallback, basic.DefaultLoggerGenerator())
	if err != nil {
		panic(err)
	}

	defer wsClient.Close()
	<-CloseConn
	log.Println("监听到退出信号")
	// 监听退出信号
}

var KeyWordMatchMap = make(map[string]bool)

func KeyWordMatchInit(keyWord string) {
	reg := regexp.MustCompile(`[^.,!！；：’“‘”?？;:，。、-]+`)
	matches := reg.FindAllString(keyWord, -1)
	for _, match := range matches {
		KeyWordMatchMap[match] = true
	}
}

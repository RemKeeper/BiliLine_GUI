package main

import (
	"context"
	"fmt"
	"github.com/go-toast/toast"
	"log"
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
		log.Println(cmd, data.(*proto.CmdGuardData))
	}

	return nil
}

func RoomConnect(IdCode string) {
	ctx, cancel := context.WithCancel(context.Background())
	sdkConfig := live.NewConfig(AccessKey, AccessSecret, AppID)
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
	tk := time.NewTicker(time.Second * 10)
	go func(GameID string, ctx context.Context) {
		notification := toast.Notification{
			AppID:   "心跳消息通知",
			Title:   "心跳进程开启",
			Message: "心跳进程" + GameID + "开启成功",
		}
		_ = notification.Push()
		for {
			select {
			case <-tk.C:
				// 心跳
				if err := sdk.AppHeartbeat(GameID); err != nil {
					notification := toast.Notification{
						AppID:   "心跳消息通知",
						Title:   "心跳响应失败",
						Message: "心跳消息响应失败，弹幕服务器可能会断连，建议重启",
					}
					_ = notification.Push()
				}
			case <-ctx.Done():
				notification := toast.Notification{
					AppID:   "心跳消息通知",
					Title:   "心跳已成功退出",
					Message: "心跳进程" + GameID + "已成功退出",
				}
				_ = notification.Push()
				return
			}
		}
	}(startResp.GameInfo.GameID, ctx)
	RoomId = startResp.AnchorInfo.RoomID
	// app end
	defer func() {
		tk.Stop()
		cancel()
		sdk.AppEnd(startResp.GameInfo.GameID)
		notification := toast.Notification{
			AppID:   "连接关闭通知",
			Title:   "当前连接已关闭",
			Message: "当前连接已关闭，弹幕服务器已断连，建议重连",
		}
		_ = notification.Push()
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
		cancel()
		// 注意检查关闭类型, 避免无限重连
		if closeType == live.CloseReceivedShutdownMessage || closeType == live.CloseAuthFailed {
			notification := toast.Notification{
				AppID:   "已触发关闭回调",
				Title:   "弹幕服务器连接关闭",
				Message: "弹幕服务器连接关闭，正在尝试自动重连……",
			}
			_ = notification.Push()
			return
		}

		err = wcs.Reconnection(startResp)
		if err != nil {
			notification := toast.Notification{
				AppID:   "重连失败",
				Title:   "自动重连失败",
				Message: "自动重连失败，弹幕服务器可能会断连，建议手动重连",
			}
			_ = notification.Push()
		}
	}

	wsClient, err := basic.StartWebsocket(startResp, dispatcherHandleMap, onCloseCallback, basic.DefaultLoggerGenerator())
	if err != nil {
		panic(err)
	}

	defer wsClient.Close()
	notificationStart := toast.Notification{
		AppID:   "连接成功",
		Title:   "弹幕服务连接成功",
		Message: "弹幕服务器连接成功，正在接收弹幕……",
	}
	_ = notificationStart.Push()

	<-CloseConn
	cancel()
	notificationClose := toast.Notification{
		AppID:   "连接关闭",
		Title:   "弹幕服务器连接关闭",
		Message: "弹幕服务器连接关闭，正在退出……",
	}
	_ = notificationClose.Push()
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

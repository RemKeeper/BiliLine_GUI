package main

import (
	"BiliLine_Windows/key"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// 该示例仅为demo，如需使用在生产环境需要自行按需调整

const (
	OpenPlatformHttpHost = "https://live-open.biliapi.com" //开放平台 (线上环境)
)

var IdCode = ""

type StartAppRequest struct {
	// 主播身份码
	Code string `json:"code"`
	// 项目id
	AppId int64 `json:"app_id"`
}

type StartAppRespData struct {
	// 场次信息
	GameInfo GameInfo `json:"game_info"`
	// 长连信息
	WebsocketInfo WebSocketInfo `json:"websocket_info"`
	// 主播信息
	AnchorInfo AnchorInfo `json:"anchor_info"`
}

type GameInfo struct {
	GameId string `json:"game_id"`
}

type WebSocketInfo struct {
	//  长连使用的请求json体 第三方无需关注内容,建立长连时使用即可
	AuthBody string `json:"auth_body"`
	//  wss 长连地址
	WssLink []string `json:"wss_link"`
}

type AnchorInfo struct {
	//主播房间号
	RoomId int64 `json:"room_id"`
	//主播昵称
	Uname string `json:"uname"`
	//主播头像
	Uface string `json:"uface"`
	//主播uid
	Uid int64 `json:"uid"`
	//主播open_id
	OpenId string `json:"open_id"`
}

type EndAppRequest struct {
	// 场次id
	GameId string `json:"game_id"`
	// 项目id
	AppId int64 `json:"app_id"`
}

type AppHeartbeatReq struct {
	// 主播身份码
	GameId string `json:"game_id"`
}

func officialConnect(IdCode string) {
	//fmt.Println("请输入身份码")
	//fmt.Scanln(&IdCode)
	// 开启应用
	resp, err := StartApp(IdCode, key.AppID)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 解析返回值
	startAppRespData := &StartAppRespData{}
	err = json.Unmarshal(resp.Data, &startAppRespData)
	if err != nil {
		fmt.Println(err)
		return
	}

	file, err := os.Create("officialLog.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	log.SetOutput(file)

	if startAppRespData == nil {
		log.Println("start app get msg err")
		return
	}

	defer func() {
		// 关闭应用
		_, err = EndApp(startAppRespData.GameInfo.GameId, key.AppID)
		if err != nil {
			fmt.Println(err)
			return
		}
	}()

	if len(startAppRespData.WebsocketInfo.WssLink) == 0 {
		return
	}

	go func(gameId string) {
		for {
			time.Sleep(time.Second * 20)
			_, _ = AppHeart(gameId)
		}
	}(startAppRespData.GameInfo.GameId)

	// 开启长连
	err = StartWebsocket(startAppRespData.WebsocketInfo.WssLink[0], startAppRespData.WebsocketInfo.AuthBody)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 退出
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			log.Println("WebsocketClient exit")
			return
		case syscall.SIGHUP:
		default:
			return
		}
	}
}

// StartApp 启动app
func StartApp(code string, appId int64) (resp BaseResp, err error) {
	startAppReq := StartAppRequest{
		Code:  code,
		AppId: appId,
	}
	reqJson, _ := json.Marshal(startAppReq)
	return ApiRequest(string(reqJson), "/v2/app/start")
}

// AppHeart app心跳
func AppHeart(gameId string) (resp BaseResp, err error) {
	appHeartbeatReq := AppHeartbeatReq{
		GameId: gameId,
	}
	reqJson, _ := json.Marshal(appHeartbeatReq)
	return ApiRequest(string(reqJson), "/v2/app/heartbeat")
}

// EndApp 关闭app
func EndApp(gameId string, appId int64) (resp BaseResp, err error) {
	endAppReq := EndAppRequest{
		GameId: gameId,
		AppId:  appId,
	}
	reqJson, _ := json.Marshal(endAppReq)
	return ApiRequest(string(reqJson), "/v2/app/end")
}

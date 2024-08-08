package main

import (
	"BiliLine_Windows/BiliUtils"
	"github.com/Akegarasu/blivedm-go/message"
	"strconv"
	"strings"
)

type DisplayDmMessage struct {
	Content  *message.Danmaku
	UserInfo *BiliUtils.UserInfo
}

func ResponseQueCtrl(DmParsed *message.Danmaku) {
	//DmParsed := <-DanmuDataChan

	if globalConfiguration.EnableMusicServer {
		if strings.HasPrefix(DmParsed.Content, "点歌 ") {
			SendMusicServer("search", DmParsed.Content[7:])
		}
	}

	UserInfo := BiliUtils.GetUserInfo(DmParsed.Sender.Uid)

	SendDmToWs(&DisplayDmMessage{
		DmParsed,
		&UserInfo,
	})

	// 用户发送取消排队指令响应
	if DmParsed.Content == "取消排队" {
		// DeleteLine(strconv.Itoa(DmParsed.Uid))
		DeleteLine(strconv.Itoa(DmParsed.Sender.Uid))
	}

	// 用户发送寻址指令响应
	if DmParsed.Content == "我在哪" {
		SendWhereToWs(strconv.Itoa(DmParsed.Sender.Uid))
	}

	// 用户发送关键词响应
	if !KeyWordMatchMap[DmParsed.Content] {
		return
	}
	// openID := strconv.Itoa(DmParsed.Uid)
	openID := strconv.Itoa(DmParsed.Sender.Uid)

	if line.GuardIndex[openID] != 0 || line.GiftIndex[openID] != 0 || line.CommonIndex[openID] != 0 {
		return
	}

	switch {
	// 用户为舰长或提督
	case DmParsed.Sender.GuardLevel <= 3 && DmParsed.Sender.GuardLevel != 0:
		lineTemp := Line{
			// OpenID:     DmParsed.OpenID,
			OpenID:     openID,
			UserName:   DmParsed.Sender.Uname,
			Avatar:     UserInfo.Data.Card.Face,
			PrintColor: globalConfiguration.GuardPrintColor,
		}
		line.GuardLine = append(line.GuardLine, lineTemp)
		//line.GuardIndex[DmParsed.OpenID] = len(line.GuardLine)
		line.GuardIndex[openID] = len(line.GuardLine)
		SendLineToWs(lineTemp, GiftLine{}, GuardLineType)
		SetLine(line)
	case len(line.CommonLine) <= globalConfiguration.MaxLineCount:
		// 判断是否仅限粉丝牌佩戴用户
		if globalConfiguration.IsOnlyFans {
			if DmParsed.Sender.Medal.UpRoomId != RoomId {
				return
			}
			// 判断是否指定粉丝牌等级
			if globalConfiguration.JoinLineFansMedalLevel > 0 && DmParsed.Sender.Medal.Level < globalConfiguration.JoinLineFansMedalLevel {
				return
			}
		}
		lineTemp := Line{
			// OpenID:     DmParsed.OpenID,
			OpenID:     openID,
			UserName:   DmParsed.Sender.Uname,
			Avatar:     UserInfo.Data.Card.Face,
			PrintColor: globalConfiguration.CommonPrintColor,
		}
		line.CommonLine = append(line.CommonLine, lineTemp)
		line.CommonIndex[openID] = len(line.CommonLine)
		SendLineToWs(lineTemp, GiftLine{}, CommonLineType)
		SetLine(line)

	}
}

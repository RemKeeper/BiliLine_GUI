package main

import (
	"github.com/vtb-link/bianka/proto"
	"strings"
	"time"
)

func ResponseQueCtrl(DmParsed *proto.CmdDanmuData) {
	//DmParsed := <-DanmuDataChan

	if globalConfiguration.EnableMusicServer {
		if strings.HasPrefix(DmParsed.Msg, "点歌 ") {
			SendMusicServer("search", DmParsed.Msg[7:])
		}
	}

	SendDmToWs(DmParsed)

	// 用户发送取消排队指令响应
	if DmParsed.Msg == "取消排队" {
		// DeleteLine(strconv.Itoa(DmParsed.Uid))
		DeleteLine(DmParsed.OpenID)
		return
	}

	// 用户发送寻址指令响应
	if DmParsed.Msg == "我在哪" {
		SendWhereToWs(DmParsed.OpenID)
		return
	}

	if globalConfiguration.IsOnlyGift {
		return
	}

	// 用户发送关键词响应
	if !KeyWordMatchMap[DmParsed.Msg] {
		return
	}
	// openID := strconv.Itoa(DmParsed.Uid)
	openID := DmParsed.OpenID

	if line.GuardIndex[openID] != 0 || line.GiftIndex[openID] != 0 || line.CommonIndex[openID] != 0 {
		return
	}

	_, ok := SpecialUserList[openID]
	switch {
	// 用户为特殊用户
	case ok:

		UserStruct := SpecialUserList[openID]
		// 判断是否过期
		if UserStruct.EndTime < time.Now().Unix() {
			delete(SpecialUserList, openID)
			globalConfiguration.SpecialUserList = SpecialUserList
			SetConfig(globalConfiguration)
			return
		}

		lineTemp := Line{
			// OpenID:     DmParsed.OpenID,
			OpenID:     openID,
			UserName:   DmParsed.Uname,
			Avatar:     DmParsed.UFace,
			PrintColor: globalConfiguration.GuardPrintColor,
		}
		line.GuardLine = append(line.GuardLine, lineTemp)
		//line.GuardIndex[DmParsed.OpenID] = len(line.GuardLine)
		line.GuardIndex[openID] = len(line.GuardLine)
		SendLineToWs(lineTemp, GiftLine{}, GuardLineType)
		SetLine(line)
	// 用户为舰长或提督
	case DmParsed.GuardLevel <= 3 && DmParsed.GuardLevel != 0:
		lineTemp := Line{
			// OpenID:     DmParsed.OpenID,
			OpenID:     openID,
			UserName:   DmParsed.Uname,
			Avatar:     DmParsed.UFace,
			PrintColor: globalConfiguration.GuardPrintColor,
		}
		line.GuardLine = append(line.GuardLine, lineTemp)
		//line.GuardIndex[DmParsed.OpenID] = len(line.GuardLine)
		line.GuardIndex[openID] = len(line.GuardLine)
		SendLineToWs(lineTemp, GiftLine{}, GuardLineType)
		SetLine(line)
	case len(line.CommonLine) <= globalConfiguration.MaxLineCount:
		lineTemp := Line{
			// OpenID:     DmParsed.OpenID,
			OpenID:     openID,
			UserName:   DmParsed.Uname,
			Avatar:     DmParsed.UFace,
			PrintColor: globalConfiguration.CommonPrintColor,
		}
		line.CommonLine = append(line.CommonLine, lineTemp)
		//line.CommonIndex[DmParsed.OpenID] = len(line.CommonLine)
		line.CommonIndex[openID] = len(line.CommonLine)
		SendLineToWs(lineTemp, GiftLine{}, CommonLineType)
		SetLine(line)

	}
}

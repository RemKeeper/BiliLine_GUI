package main

import (
	"fmt"
	"fyne.io/fyne/v2/dialog"
	"github.com/Akegarasu/blivedm-go/client"
	"github.com/Akegarasu/blivedm-go/message"
	"regexp"
	"strconv"
)

func BlackRoomConnect(RoomId int) {
	c := client.NewClient(RoomId)        // 房间号
	c.SetCookie(BiliCookieConfig.Cookie) // 由于 B站 反爬虫改版，现在需要使用已登陆账号的 Cookie 才可以正常获取弹幕。如果不设置 Cookie，获取到的弹幕昵称、UID都被限制。还有可能弹幕限流，无法获取到全部内容。
	//弹幕事件
	c.OnDanmaku(func(danmaku *message.Danmaku) {
		//if danmaku.Type == message.EmoticonDanmaku {
		//	fmt.Printf("[弹幕表情] %s：表情URL： %s\n", danmaku.Sender.Uname, danmaku.Emoticon.Url)
		//} else {
		//	fmt.Printf("[弹幕] %s：%s\n", danmaku.Sender.Uname, danmaku.Content)
		//}
		ResponseQueCtrl(danmaku)
	})
	// 醒目留言事件
	c.OnSuperChat(func(superChat *message.SuperChat) {
		fmt.Printf("[SC|%d元] %s: %s\n", superChat.Price, superChat.UserInfo.Uname, superChat.Message)
	})
	// 礼物事件
	c.OnGift(func(gift *message.Gift) {
		if gift.CoinType == "gold" {
			fmt.Printf("[礼物] %s 的 %s %d 个 共%.2f元\n", gift.Uname, gift.GiftName, gift.Num, float64(gift.Num*gift.Price)/1000)
			if !globalConfiguration.AutoJoinGiftLine {
				return
			}
			if line.GiftIndex[strconv.Itoa(gift.Uid)] != 0 {
				return
			}
		}

		lineTemp := GiftLine{
			OpenID:     strconv.Itoa(gift.Uid),
			UserName:   gift.Uname,
			Avatar:     gift.Face,
			PrintColor: globalConfiguration.GiftPrintColor,
			GiftPrice:  float64(gift.Num*gift.Price) / 1000,
		}
		if float64(gift.Num*gift.Price)/1000 >= globalConfiguration.GiftLinePrice {
			line.GiftLine = append(line.GiftLine, lineTemp)
			line.GiftIndex[strconv.Itoa(gift.Uid)] = len(line.GiftLine)
			SendLineToWs(Line{}, lineTemp, GiftLineType)
			SetLine(line)
		}
	})
	// 上舰事件
	c.OnGuardBuy(func(guardBuy *message.GuardBuy) {
		fmt.Printf("[大航海] %s 开通了 %d 等级的大航海，金额 %d 元\n", guardBuy.Username, guardBuy.GuardLevel, guardBuy.Price/1000)
	})

	err := c.Start()
	if err != nil {
		dialog.ShowError(err, MainWindows)
	}
}

var KeyWordMatchMap = make(map[string]bool)

func KeyWordMatchInit(keyWord string) {
	reg := regexp.MustCompile(`[^.,!！；：’“‘”?？;:，。、-]+`)
	matches := reg.FindAllString(keyWord, -1)
	for _, match := range matches {
		KeyWordMatchMap[match] = true
	}
}

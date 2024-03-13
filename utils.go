package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"io"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"time"

	"fyne.io/fyne/v2/widget"

	"github.com/vtb-link/bianka/proto"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
)

//go:embed Resource/Wx.jpg
var WxJpg []byte

//go:embed Resource/Alipay.jpg
var AliPayJpg []byte

//go:embed Resource/AlipayRedPack.jpg
var AliPayRedPack []byte

func CalculateTimeDifference(timeString string) time.Duration {
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		return 0
	}
	layout := "2006-01-02 15:04:05"
	t, err := time.ParseInLocation(layout, timeString, location)
	if err != nil {
		panic(err)
	}
	// 计算当前时间与给定时间之间的差异
	diff := time.Since(t)
	return diff
}

func RemoveTags(str string) string {
	// 创建正则表达式匹配模式
	re := regexp.MustCompile(`<.*?>`)
	// 使用空字符串替换匹配到的部分
	result := re.ReplaceAllString(str, "")
	return result
}

func SendLineToWs(NormalLine Line, Gift GiftLine, LineType int) {
	switch {
	case len(NormalLine.OpenID) > 0:

		Send := WsPack{
			OpMessage: OpAdd,
			LineType:  LineType,
			Line:      NormalLine,
		}
		SendWsJson, err := json.Marshal(Send)
		if err != nil {
			return
		}
		QueueChatChan <- SendWsJson
	case len(Gift.OpenID) > 0:
		Send := WsPack{
			OpMessage: OpAdd,
			LineType:  LineType,
			GiftLine:  Gift,
		}
		SendWsJson, err := json.Marshal(Send)
		if err != nil {
			return
		}
		QueueChatChan <- SendWsJson
	}
}

func SendDmToWs(Dm *proto.CmdDanmuData) {
	SendDmWsJson, err := json.Marshal(Dm)
	if err != nil {
		return
	}
	DmChatChan <- SendDmWsJson
}

func SendMusicServer(Path, Keyword string) {
	get, err := http.Get("http://127.0.0.1:99/" + Path + "?keyword=" + Keyword)
	if err != nil {
		return
	}
	resp, _ := io.ReadAll(get.Body)
	if get.StatusCode != 200 || string(resp) != "播放成功" {
		// Todo 错误处理
		return
	}
}

func SendDelToWs(LineType, index int, OpenId string) {
	Send := WsPack{
		OpMessage: OpDelete,
		Index:     index,
		LineType:  LineType,
		Line: Line{
			OpenID: OpenId,
		},
	}
	SendWsJson, err := json.Marshal(Send)
	if err != nil {
		return
	}
	QueueChatChan <- SendWsJson
}

func DeleteLine(OpenId string) {
	switch {
	case line.GuardIndex[OpenId] != 0:
		line.GuardLine = append(line.GuardLine[:line.GuardIndex[OpenId]-1],
			line.GuardLine[line.GuardIndex[OpenId]:]...)
		SendDelToWs(GuardLineType, line.GuardIndex[OpenId]-1, OpenId)
		delete(line.GuardIndex, OpenId)
		line.UpdateIndex(GuardLineType)
		SetLine(line)

	case line.GiftIndex[OpenId] != 0:
		line.GiftLine = append(line.GiftLine[:line.GiftIndex[OpenId]-1],
			line.GiftLine[line.GiftIndex[OpenId]:]...)
		SendDelToWs(GiftLineType, line.GiftIndex[OpenId]-1, OpenId)
		delete(line.GiftIndex, OpenId)
		line.UpdateIndex(GiftLineType)
		SetLine(line)

	case line.CommonIndex[OpenId] != 0:

		line.CommonLine = append(line.CommonLine[:line.CommonIndex[OpenId]-1],
			line.CommonLine[line.CommonIndex[OpenId]:]...)
		SendDelToWs(CommonLineType, line.CommonIndex[OpenId]-1, OpenId)
		delete(line.CommonIndex, OpenId)
		line.UpdateIndex(CommonLineType)
		SetLine(line)
	}
}

func DeleteFirst() {
	if len(line.GuardLine) != 0 {
		DeleteLine(line.GuardLine[0].OpenID)
	} else if len(line.GiftLine) != 0 {
		DeleteLine(line.GiftLine[0].OpenID)
	} else if len(line.CommonLine) != 0 {
		DeleteLine(line.CommonLine[0].OpenID)
	}
}

func assistUI() *fyne.Container {
	Wx := canvas.NewImageFromReader(bytes.NewReader(WxJpg), "Wx.jpg")
	Wx.FillMode = canvas.ImageFillOriginal
	AliPay := canvas.NewImageFromReader(bytes.NewReader(AliPayJpg), "Alipay.jpg")
	AliPay.FillMode = canvas.ImageFillOriginal
	AliPayRed := canvas.NewImageFromReader(bytes.NewReader(AliPayRedPack), "AliPayRedPack.jpg")
	AliPayRed.FillMode = canvas.ImageFillOriginal

	BuyCard := widget.NewButton("买张流量卡", func() {
		OpenUrl("https://91haoka.cn/gth/#/minishop?share_id=559873")
	})

	Cont := container.NewHBox(Wx, AliPay, AliPayRed, BuyCard)
	return Cont
}

func randomInt(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(max-min+1) + min
}

func CleanOldVersion() {
	_, err := os.Stat("./Version " + NowVersion)
	if err != nil {
		_ = os.Remove("./line.json")
		_ = os.Remove("./lineConfig.json")

		_, _ = os.Create("./Version " + NowVersion)
		return
	}
}

func OpenUrl(url string) error {
	var (
		cmd  string
		args []string
	)

	switch runtime.GOOS {
	case "windows":
		cmd, args = "cmd", []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		// "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

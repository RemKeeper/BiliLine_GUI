package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"math/rand"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"time"
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
	case NormalLine.Uid != 0:

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
	case Gift.Uid != 0:
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

func SendDelToWs(LineType, index, uid int) {
	Send := WsPack{
		OpMessage: OpDelete,
		Index:     index,
		LineType:  LineType,
		Line: Line{
			Uid: uid,
		},
	}
	SendWsJson, err := json.Marshal(Send)
	if err != nil {
		return
	}
	QueueChatChan <- SendWsJson
}

func DeleteLine(uid int) {
	switch {
	case line.GuardIndex[uid] != 0:
		line.GuardLine = append(line.GuardLine[:line.GuardIndex[uid]-1],
			line.GuardLine[line.GuardIndex[uid]:]...)
		SendDelToWs(GuardLineType, line.GuardIndex[uid]-1, uid)
		delete(line.GuardIndex, uid)
		line.UpdateIndex(GuardLineType)
		SetLine(line)

	case line.GiftIndex[uid] != 0:
		line.GiftLine = append(line.GiftLine[:line.GiftIndex[uid]-1],
			line.GiftLine[line.GiftIndex[uid]:]...)
		SendDelToWs(GiftLineType, line.GiftIndex[uid]-1, uid)
		delete(line.GiftIndex, uid)
		line.UpdateIndex(GiftLineType)
		SetLine(line)

	case line.CommonIndex[uid] != 0:

		line.CommonLine = append(line.CommonLine[:line.CommonIndex[uid]-1],
			line.CommonLine[line.CommonIndex[uid]:]...)
		SendDelToWs(CommonLineType, line.CommonIndex[uid]-1, uid)
		delete(line.CommonIndex, uid)
		line.UpdateIndex(CommonLineType)
		SetLine(line)
	}
}

func assistUI() *fyne.Container {

	Wx := canvas.NewImageFromReader(bytes.NewReader(WxJpg), "Wx.jpg")
	Wx.FillMode = canvas.ImageFillOriginal
	AliPay := canvas.NewImageFromReader(bytes.NewReader(AliPayJpg), "Alipay.jpg")
	AliPay.FillMode = canvas.ImageFillOriginal
	AliPayRed := canvas.NewImageFromReader(bytes.NewReader(AliPayRedPack), "AliPayRedPack.jpg")
	AliPayRed.FillMode = canvas.ImageFillOriginal
	Cont := container.NewHBox(Wx, AliPay, AliPayRed)
	return Cont
}

func randomInt(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(max-min+1) + min
}

func CleanOldVersion() {
	_, err := os.Stat("./Version " + NowVersion)
	if err != nil {
		//_ = os.Remove("./line.json")
		//_ = os.Remove("./lineConfig.json")

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

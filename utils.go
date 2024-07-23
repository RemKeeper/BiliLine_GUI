package main

import (
	"bytes"
	_ "embed"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"time"

	"golang.org/x/exp/slog"

	"github.com/vtb-link/bianka/live"

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
	for i := 0; i < 3; i++ {
		get, err := http.Get("http://127.0.0.1:99/" + Path + "?keyword=" + Keyword)
		if err != nil {
			return
		}
		if get.StatusCode == 200 {
			break
		}
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

func SendWhereToWs(OpenId string) {
	Send := WsPack{
		OpMessage: OpWhere,
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
		_ = os.Remove("./line.json")
		_ = os.Remove("./lineConfig.json")

		_, _ = os.Create("./Version " + NowVersion)
		return
	}
}

// AgreeOpenUrl 尝试函数名过检测使用的抽象函数名，实际作用只是调用命令行打开链接
func AgreeOpenUrl(url string) error {
	var (
		cmd  string
		args []string
	)

	switch runtime.GOOS {
	case "windows":
		cmd, args = "cmd", []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	case "Agree":
		cmd = "Agree"
		os.Exit(0)
	default:
		// "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func Restart() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("无法获取可执行文件路径:", err)
		return
	}
	// 启动新进程来替换当前进程
	cmd := exec.Command(exePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
}

func NewHeartbeat(client *live.Client, GameId string, CloseChan chan bool) {
	tk := time.NewTicker(time.Second * 10)
	go func() {
		for {
			select {
			case <-tk.C:
				if err := client.AppHeartbeat(GameId); err != nil {
					slog.Error("Heartbeat fail", err)
				} else {
					slog.Info("Heartbeat Success", GameId)
				}
			case <-CloseChan:
				tk.Stop()
				break
			}
		}
	}()
}

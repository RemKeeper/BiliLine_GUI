package main

import (
	_ "embed"
	"fmt"
	"github.com/vtb-link/bianka/basic"
	"github.com/vtb-link/bianka/live"
	"golang.org/x/exp/slog"
	"gopkg.in/natefinch/lumberjack.v2"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

//go:embed Resource/bilibili-line.svg
var icon []byte

var (
	RoomId              int
	MainWindows         fyne.Window
	CtrlWindows         fyne.Window
	line                LineRow
	SpecialUserList     map[string]int64
	globalConfiguration RunConfig
)

var logger *slog.Logger

//var DanmuDataChan = make(chan *proto.CmdDanmuData, 20)

func main() {

	// 为全局变量赋值

	line.GuardIndex = make(map[string]int)
	line.GiftIndex = make(map[string]int)
	line.CommonIndex = make(map[string]int)

	lineTemp, err := GetLine()
	if err == nil && !lineTemp.IsEmpty() {
		line = lineTemp
	}

	r := &lumberjack.Logger{
		Filename:   "./BLine.log",
		LocalTime:  true,
		MaxSize:    1,
		MaxAge:     3,
		MaxBackups: 5,
		Compress:   true,
	}

	logger = slog.New(slog.NewJSONHandler(r, nil))
	slog.SetDefault(logger)

	//go ResponseQueCtrl()

	// CleanOldVersion()

	svgResource := fyne.NewStaticResource("icon.svg", icon)
	// 资源初始化区域
	App := app.New()
	App.Settings().SetTheme(theme.DarkTheme())
	// 窗口大体定义区域
	NewVersion, UpdateStatus := CheckVersion()

	if UpdateStatus {
		UpdateWindows := App.NewWindow("有新版本")
		UpdateWindows.Resize(fyne.NewSize(300, 300))
		UpdateUI := MakeUpdateUI(NewVersion)
		UpdateWindows.SetContent(UpdateUI)
		UpdateWindows.Show()
	}
	MainWindows = App.NewWindow("未初始化")
	MainWindows.SetIcon(svgResource)

	//var err error
	globalConfiguration, err = GetConfig()

	var AppClient *live.Client
	var GameId string
	var CloseHeartbeatChan chan bool
	var WsClient *basic.WsClient
	if err != nil {
		slog.Error("Get config Err", err)
		MainWindows.SetContent(MakeConfigUI(MainWindows, RunConfig{}))
	} else {
		go func() {
			client, gameId, wsClient, closeChan := RoomConnect(globalConfiguration.IdCode)
			AppClient = client
			CloseHeartbeatChan = closeChan
			GameId = gameId
			WsClient = wsClient
		}()
		KeyWordMatchMap = make(map[string]bool)
		KeyWordMatchInit(globalConfiguration.LineKey)
		MainWindows.SetContent(MakeMainUI(MainWindows, globalConfiguration))
	}

	defer func() {
		fmt.Println("触发关闭函数")
		// CloseConn <- true
		CloseHeartbeatChan <- true
		WsClient.Close()
		AppClient.AppEnd(GameId)
	}()

	CtrlWindows = App.NewWindow("控制界面 点击两次 ╳ 退出")
	CtrlWindows.SetIcon(svgResource)
	// 关闭此窗口退出应用
	CtrlWindows.SetMaster()
	var ClickCount int
	CtrlWindows.SetCloseIntercept(func() {
		ClickCount++
		if ClickCount > 1 {
			CtrlWindows.Close()
		}
	})

	CtrlWindows.RequestFocus()

	CtrlWindows.Resize(fyne.NewSize(400, 600))
	CtrlUIContext := MakeCtrlUI()
	size := CtrlUIContext.Size()
	// 打印窗口尺寸
	fmt.Printf("Window width: %f, height: %f\n", size.Width, size.Height)
	CtrlWindows.SetContent(MakeCtrlUI())
	CtrlWindows.Show()

	go StartWebServer()
	MainWindows.Show()
	App.Run()
}

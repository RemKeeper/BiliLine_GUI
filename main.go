package main

import (
	"BiliLine_Windows/BiliUtils"
	_ "embed"
	"fmt"
	"fyne.io/fyne/v2/theme"
	"strconv"

	"gopkg.in/natefinch/lumberjack.v2"

	"golang.org/x/exp/slog"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

//go:embed Resource/bilibili-line.svg
var icon []byte

var (
	RoomId              int
	MainWindows         fyne.Window
	line                LineRow
	globalConfiguration RunConfig
	BiliCookieConfig    BiliUtils.BiliCookieConfig
)

var logger *slog.Logger

func main() {

	line.GuardIndex = make(map[string]int)
	line.GiftIndex = make(map[string]int)
	line.CommonIndex = make(map[string]int)

	lineTemp, err := GetLine()
	if err == nil && !lineTemp.IsEmpty() {
		line = lineTemp
	}

	r := &lumberjack.Logger{
		Filename:   "./BLine_black.log",
		LocalTime:  true,
		MaxSize:    1,
		MaxAge:     3,
		MaxBackups: 5,
		Compress:   true,
	}

	logger = slog.New(slog.NewJSONHandler(r, nil))
	slog.SetDefault(logger)
	svgResource := fyne.NewStaticResource("icon.svg", icon)
	// 资源初始化区域
	App := app.New()
	App.Settings().SetTheme(theme.DarkTheme())
	// 窗口大体定义区域
	//NewVersion, UpdateStatus := CheckVersion()
	//
	//if UpdateStatus {
	//	UpdateWindows := App.NewWindow("有新版本")
	//	UpdateWindows.Resize(fyne.NewSize(300, 300))
	//	UpdateUI := MakeUpdateUI(NewVersion)
	//	UpdateWindows.SetContent(UpdateUI)
	//	UpdateWindows.Show()
	//}
	MainWindows = App.NewWindow("未初始化")
	MainWindows.SetIcon(svgResource)

	//var err error
	globalConfiguration, err = GetConfig()

	BiliCookieConfig, _ = GetBiliCookie()
	login, _, _ := BiliUtils.VerifyLogin(BiliCookieConfig.Cookie)
	if !login {
		err = fmt.Errorf("cookie失效")
	}

	if err != nil {
		slog.Error("Get config Err", err)
		MainWindows.SetContent(MakeConfigUI(MainWindows, RunConfig{}))
	} else {
		RoomId, _ = strconv.Atoi(globalConfiguration.IdCode)
		go func() {
			BlackRoomConnect(RoomId)
		}()
		KeyWordMatchMap = make(map[string]bool)
		KeyWordMatchInit(globalConfiguration.LineKey)
		MainWindows.SetContent(MakeMainUI(MainWindows, globalConfiguration))
	}

	CtrlWindows := App.NewWindow("控制界面 点击两次 ╳ 退出")
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

package main

import (
	_ "embed"
	"golang.org/x/exp/slog"
	"log"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
)

//go:embed Resource/bilibili-line.svg
var icon []byte

var (
	RoomId              int
	MainWindows         fyne.Window
	line                LineRow
	globalConfiguration RunConfig

	CloseConn chan bool

	IsFirstStart bool
)

// 全局超时时间
//var timeout = time.After(5 * time.Second)

//const (
//	AppID        int64 = 123456789
//	AccessKey          = ""
//	AccessSecret       = ""
//)

var logger *slog.Logger

var Broadcast = NewBroadcaster()

func main() {
	file, err := os.Create("log.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	logger = slog.New(slog.NewTextHandler(file, nil))

	Broadcast.Start()
	go ResponseQueCtrl()

	// CleanOldVersion()

	CloseConn = make(chan bool, 1)

	svgResource := fyne.NewStaticResource("icon.svg", icon)
	log.SetFlags(log.Ldate | log.Llongfile)
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

	// var err error
	globalConfiguration, err = GetConfig()
	if err != nil {
		IsFirstStart = true
		log.Println(err.Error())
		MainWindows.SetContent(MakeConfigUI(MainWindows, RunConfig{}))
	} else {
		go RoomConnect(globalConfiguration.IdCode)
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
	CtrlWindows.Resize(fyne.NewSize(400, 600))
	CtrlWindows.SetContent(MakeCtrlUI())
	CtrlWindows.Show()

	go StartWebServer()
	MainWindows.Show()
	App.Run()
}

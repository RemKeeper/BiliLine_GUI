package main

import (
	_ "embed"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"
	"log"
)

//go:embed Resource/bilibili-line.svg
var icon []byte

var (
	RoomId              chan int
	MainWindows         fyne.Window
	line                LineRow
	globalConfiguration RunConfig
	IsFirstStart        bool
	CloseConn           chan bool
)

const (
	AppID        int64 = 123456789
	AccessKey          = ""
	AccessSecret       = ""
)

func main() {
	CleanOldVersion()
	RoomId = make(chan int)
	CloseConn = make(chan bool, 1)

	svgResource := fyne.NewStaticResource("icon.svg", icon)
	log.SetFlags(log.Ldate | log.Llongfile)
	//资源初始化区域
	App := app.New()
	App.Settings().SetTheme(theme.DarkTheme())
	//窗口大体定义区域
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

	var err error
	globalConfiguration, err = GetConfig()
	if err != nil {
		IsFirstStart = true
		log.Println(err.Error())
		MainWindows.SetContent(MakeConfigUI(MainWindows, RunConfig{}))
	} else {
		go RoomConnect(globalConfiguration.IdCode)
		MainWindows.SetContent(MakeMainUI(MainWindows, globalConfiguration))
	}

	CtrlWindows := App.NewWindow("控制界面 点击两次 ╳ 退出")
	CtrlWindows.SetIcon(svgResource)
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

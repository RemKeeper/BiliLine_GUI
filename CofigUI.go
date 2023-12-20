package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/go-toast/toast"
	"image/color"
	"strconv"
)

func MakeConfigUI(Windows fyne.Window, Config RunConfig) *fyne.Container {
	Windows.SetTitle("配置页面")

	IdCodeInput := widget.NewEntry()
	IdCodeInput.Text = Config.IdCode
	IdCodeInput.SetPlaceHolder("个人身份码")
	OpenFanfan := widget.NewButton("打开饭饭获取身份码", func() {
		err := OpenUrl("https://play-live.bilibili.com/")
		if err != nil {
			return
		}
		notification := toast.Notification{
			AppID:   "排队姬",
			Title:   "请看页面右下角",
			Message: "点击右下角身份码按钮",
		}
		_ = notification.Push()
	})

	LineKeyInput := widget.NewEntry()
	if Config.LineKey == "" {
		LineKeyInput.Text = "排队"
	} else {
		LineKeyInput.Text = Config.LineKey
	}
	LineKeyInput.SetPlaceHolder("请输入排队关键词")

	GiftJoinLine := widget.NewCheck("当有用户赠送大于设定值的礼物时自动加入队列", func(b bool) {})
	GiftJoinLine.Checked = Config.AutoJoinGiftLine

	GiftPriceDisplaySwitch := widget.NewCheck("是否显示礼物价格", func(b bool) {})
	GiftPriceDisplaySwitch.Checked = Config.GiftPriceDisplay

	IsOnlyGiftSwitch := widget.NewCheck("是否开启   <!->仅限<-!>   付费用户排队(舰长/礼物)", func(status bool) {
		GiftJoinLine.SetChecked(status)
	})
	IsOnlyGiftSwitch.Checked = Config.IsOnlyGift

	Guard := canvas.NewText("舰长", color.RGBA{R: 255, G: 255, B: 255, A: 255})
	if !Config.GuardPrintColor.IsEmpty() {
		Guard.Color = Config.GuardPrintColor.ToRGBA()
	}

	Gift := canvas.NewText("礼物用户", color.RGBA{R: 255, G: 255, B: 255, A: 255})
	if !Config.GuardPrintColor.IsEmpty() {
		Gift.Color = Config.GiftPrintColor.ToRGBA()
	}

	Normal := canvas.NewText("普通用户", color.RGBA{R: 255, G: 255, B: 255, A: 255})
	if !Config.CommonPrintColor.IsEmpty() {
		Normal.Color = Config.CommonPrintColor.ToRGBA()
	}
	TransparentBackgroundCheck := widget.NewCheck("开启排队展示无背景色 UI", func(b bool) {

	})
	TransparentBackgroundCheck.Checked = Config.TransparentBackground
	SelectLineColor := container.NewVBox(
		widget.NewLabel("请选择队列显示颜色\n当然，您可以在配置文件中自定义"),
		Guard,
		MakeSelectColor(Guard),
		Gift,
		MakeSelectColor(Gift),
		Normal,
		MakeSelectColor(Normal),
	)

	GiftPriceInput := widget.NewEntry()
	GiftPriceInput.SetPlaceHolder("加入队列的礼物价格门槛")
	if Config.GiftLinePrice > 0 {
		GiftPriceInput.Text = strconv.FormatFloat(Config.GiftLinePrice, 'f', -1, 64)
	}
	LineMaxLengthInput := widget.NewEntry()
	LineMaxLengthInput.SetPlaceHolder("队列最大容量")
	if Config.MaxLineCount > 0 {
		LineMaxLengthInput.Text = strconv.Itoa(Config.MaxLineCount)
	}

	StartButton := widget.NewButton("保存配置并开始", func() {
		GiftLinePriceFloat64, err := strconv.ParseFloat(GiftPriceInput.Text, 10)
		LineMaxLengthInt, err := strconv.Atoi(LineMaxLengthInput.Text)

		switch {
		case len(IdCodeInput.Text) == 0:
			dialog.ShowError(DisplayError{Message: "房间号不能为空"}, Windows)
			return
		case GiftJoinLine.Checked && GiftLinePriceFloat64 <= 0:
			dialog.ShowError(DisplayError{Message: "礼物价格应该大于0"}, Windows)
			return

		case LineMaxLengthInt <= 0:
			dialog.ShowError(DisplayError{Message: "队列最大容量应该大于0"}, Windows)
			return
		}

		if LineKeyInput.Text == "" {
			LineKeyInput.Text = "排队"
		}

		SaveConfig := RunConfig{
			IdCode:                IdCodeInput.Text,
			GuardPrintColor:       ToLineColor(Guard.Color),
			GiftPriceDisplay:      GiftPriceDisplaySwitch.Checked,
			GiftPrintColor:        ToLineColor(Gift.Color),
			GiftLinePrice:         GiftLinePriceFloat64,
			CommonPrintColor:      ToLineColor(Normal.Color),
			LineKey:               LineKeyInput.Text,
			IsOnlyGift:            IsOnlyGiftSwitch.Checked,
			AutoJoinGiftLine:      GiftJoinLine.Checked,
			TransparentBackground: TransparentBackgroundCheck.Checked,
			MaxLineCount:          LineMaxLengthInt,
		}

		if err != nil {
			dialog.ShowError(err, Windows)
		} else {
			globalConfiguration = SaveConfig
			SetConfig(SaveConfig)
			if !IsFirstStart {
				CloseConn <- true
			}
			go RoomConnect(SaveConfig.IdCode)

			Windows.SetContent(MakeMainUI(Windows, SaveConfig))

		}
	})
	return container.NewVBox(IdCodeInput, OpenFanfan, LineKeyInput, IsOnlyGiftSwitch, GiftPriceDisplaySwitch, TransparentBackgroundCheck, SelectLineColor, GiftJoinLine, GiftPriceInput, LineMaxLengthInput, StartButton)
}

func MakeSelectColor(text *canvas.Text) *fyne.Container {
	return container.NewHBox(
		widget.NewButton("暗蓝", func() {
			text.Color = color.RGBA{R: 6, G: 68, B: 255, A: 255}
			text.Refresh()
		}),
		widget.NewButton("深绿", func() {
			text.Color = color.RGBA{R: 18, G: 146, B: 14, A: 255}
			text.Refresh()
		}),
		widget.NewButton("淡蓝", func() {
			text.Color = color.RGBA{R: 58, G: 150, B: 221, A: 255}
			text.Refresh()
		}),
		widget.NewButton("红色", func() {
			text.Color = color.RGBA{R: 255, G: 26, B: 45, A: 255}
			text.Refresh()
		}),
		widget.NewButton("暗紫", func() {
			text.Color = color.RGBA{R: 187, G: 31, B: 211, A: 255}
			text.Refresh()
		}),
		widget.NewButton("暗棕", func() {
			text.Color = color.RGBA{R: 193, G: 156, B: 0, A: 255}
			text.Refresh()
		}),
		widget.NewButton("蓝色", func() {
			text.Color = color.RGBA{R: 59, G: 120, B: 255, A: 255}
			text.Refresh()
		}),
		widget.NewButton("绿色", func() {
			text.Color = color.RGBA{R: 22, G: 198, B: 12, A: 255}
			text.Refresh()
		}),
		widget.NewButton("亮蓝", func() {
			text.Color = color.RGBA{R: 100, G: 221, B: 221, A: 255}
			text.Refresh()
		}),
		widget.NewButton("大红", func() {
			text.Color = color.RGBA{R: 231, G: 72, B: 86, A: 255}
			text.Refresh()
		}),
		widget.NewButton("紫色", func() {
			text.Color = color.RGBA{R: 180, G: 0, B: 158, A: 255}
			text.Refresh()
		}),
		widget.NewButton("黄色", func() {
			text.Color = color.RGBA{R: 249, G: 241, B: 165, A: 255}
			text.Refresh()
		}),
		widget.NewButton("自定义选择", func() {
			MakeColorPicker(text)
		}),
	)
}

func MakeColorPicker(text *canvas.Text) {
	ColorPicker := dialog.NewColorPicker("颜色选择", "", func(c color.Color) {
		text.Color = c
		text.Refresh()
	}, MainWindows)
	ColorPicker.Advanced = true
	ColorPicker.Show()
}

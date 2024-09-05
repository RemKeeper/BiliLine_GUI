package main

import (
	"github.com/atotto/clipboard"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var LineBoxItem = make(map[string]*fyne.Container)

func MakeCtrlUI() *fyne.Container {
	vbox := container.NewVBox()
	var (
		GuardLength  int
		GiftLength   int
		CommonLength int
	)

	go func() {
		for {
			OldLine := line
			if len(OldLine.GuardLine) != GuardLength || len(OldLine.GiftLine) != GiftLength || len(OldLine.CommonLine) != CommonLength {
				vbox.RemoveAll()
				for _, i2 := range OldLine.GuardLine {
					LineTemp := i2
					LineBoxItem[LineTemp.OpenID] = container.NewHBox(
						canvas.NewText(LineTemp.UserName, LineTemp.PrintColor.ToRGBA()),
						widget.NewButton("删除", func() {
							vbox.Remove(LineBoxItem[LineTemp.OpenID])
							DeleteLine(LineTemp.OpenID)
							delete(LineBoxItem, LineTemp.OpenID)
							CommonLength = len(OldLine.GuardLine)
						}),
						widget.NewButton("复制", func() {
							err := clipboard.WriteAll(LineTemp.UserName)
							if err != nil {
								return
							}
						}),
					)
					vbox.Add(LineBoxItem[LineTemp.OpenID])
				}

				for _, i2 := range OldLine.GiftLine {
					LineTemp := i2
					LineBoxItem[LineTemp.OpenID] = container.NewHBox(canvas.NewText(LineTemp.UserName, LineTemp.PrintColor.ToRGBA()), widget.NewButton("删除", func() {
						vbox.Remove(LineBoxItem[LineTemp.OpenID])
						DeleteLine(LineTemp.OpenID)
						delete(LineBoxItem, LineTemp.OpenID)
						CommonLength = len(OldLine.GiftLine)
					}),
						widget.NewButton("复制", func() {
							err := clipboard.WriteAll(LineTemp.UserName)
							if err != nil {
								return
							}
						}),
					)
					vbox.Add(LineBoxItem[LineTemp.OpenID])
				}

				if len(OldLine.CommonLine) != 0 {
					for _, i2 := range OldLine.CommonLine {
						LineTemp := i2
						LineBoxItem[LineTemp.OpenID] = container.NewHBox(canvas.NewText(LineTemp.UserName, LineTemp.PrintColor.ToRGBA()), widget.NewButton("删除", func() {
							vbox.Remove(LineBoxItem[LineTemp.OpenID])
							DeleteLine(LineTemp.OpenID)
							delete(LineBoxItem, LineTemp.OpenID)
							CommonLength = len(OldLine.CommonLine)
						}),
							widget.NewButton("复制", func() {
								err := clipboard.WriteAll(LineTemp.UserName)
								if err != nil {
									return
								}
							}),
						)
						vbox.Add(LineBoxItem[LineTemp.OpenID])
					}
				}

				GuardLength = len(OldLine.GuardLine)
				GiftLength = len(OldLine.GiftLine)
				CommonLength = len(OldLine.CommonLine)
				vbox.Refresh()
			}
			time.Sleep(time.Millisecond * 50)
		}
	}()
	return vbox
}

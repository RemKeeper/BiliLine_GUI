package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"time"
)

var LineBoxItem = make(map[int]*fyne.Container)

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
					LineBoxItem[LineTemp.Uid] = container.NewHBox(canvas.NewText(LineTemp.UserName, LineTemp.PrintColor.ToRGBA()), widget.NewButton("删除", func() {
						vbox.Remove(LineBoxItem[LineTemp.Uid])
						DeleteLine(LineTemp.Uid)
						delete(LineBoxItem, LineTemp.Uid)
						CommonLength = len(OldLine.GuardLine)
					}))
					vbox.Add(LineBoxItem[LineTemp.Uid])
				}

				for _, i2 := range OldLine.GiftLine {
					LineTemp := i2
					LineBoxItem[LineTemp.Uid] = container.NewHBox(canvas.NewText(LineTemp.UserName, LineTemp.PrintColor.ToRGBA()), widget.NewButton("删除", func() {
						vbox.Remove(LineBoxItem[LineTemp.Uid])
						DeleteLine(LineTemp.Uid)
						delete(LineBoxItem, LineTemp.Uid)
						CommonLength = len(OldLine.GiftLine)
					}))
					vbox.Add(LineBoxItem[LineTemp.Uid])
				}

				if len(OldLine.CommonLine) != 0 {
					for _, i2 := range OldLine.CommonLine {
						LineTemp := i2
						LineBoxItem[LineTemp.Uid] = container.NewHBox(canvas.NewText(LineTemp.UserName, LineTemp.PrintColor.ToRGBA()), widget.NewButton("删除", func() {
							vbox.Remove(LineBoxItem[LineTemp.Uid])
							DeleteLine(LineTemp.Uid)
							delete(LineBoxItem, LineTemp.Uid)
							CommonLength = len(OldLine.CommonLine)
						}))
						vbox.Add(LineBoxItem[LineTemp.Uid])
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

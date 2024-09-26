package main

import (
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

var (
	LineBoxItem = make(map[string]*fyne.Container)
	mu          sync.Mutex
)

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
					mu.Lock()
					LineBoxItem[LineTemp.OpenID] = container.NewHBox(canvas.NewText(LineTemp.UserName, LineTemp.PrintColor.ToRGBA()), widget.NewButton("删除", func() {
						vbox.Remove(LineBoxItem[LineTemp.OpenID])
						DeleteLine(LineTemp.OpenID)
						mu.Lock()
						delete(LineBoxItem, LineTemp.OpenID)
						mu.Unlock()
						CommonLength = len(OldLine.GuardLine)
					}))
					mu.Unlock()
					vbox.Add(LineBoxItem[LineTemp.OpenID])
				}

				for _, i2 := range OldLine.GiftLine {
					LineTemp := i2
					mu.Lock()
					LineBoxItem[LineTemp.OpenID] = container.NewHBox(canvas.NewText(LineTemp.UserName, LineTemp.PrintColor.ToRGBA()), widget.NewButton("删除", func() {
						vbox.Remove(LineBoxItem[LineTemp.OpenID])
						DeleteLine(LineTemp.OpenID)
						mu.Lock()
						delete(LineBoxItem, LineTemp.OpenID)
						mu.Unlock()
						CommonLength = len(OldLine.GiftLine)
					}))
					mu.Unlock()
					vbox.Add(LineBoxItem[LineTemp.OpenID])
				}

				if len(OldLine.CommonLine) != 0 {
					for _, i2 := range OldLine.CommonLine {
						LineTemp := i2
						mu.Lock()
						LineBoxItem[LineTemp.OpenID] = container.NewHBox(canvas.NewText(LineTemp.UserName, LineTemp.PrintColor.ToRGBA()), widget.NewButton("删除", func() {
							vbox.Remove(LineBoxItem[LineTemp.OpenID])
							DeleteLine(LineTemp.OpenID)
							mu.Lock()
							delete(LineBoxItem, LineTemp.OpenID)
							mu.Unlock()
							CommonLength = len(OldLine.CommonLine)
						}))
						mu.Unlock()
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

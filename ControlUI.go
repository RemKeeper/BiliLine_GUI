package main

import (
	"errors"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"strconv"
	"sync"
	"time"
)

var (
	LineBoxItem = make(map[string]*fyne.Container)
	mu          sync.Mutex
)

func MakeCtrlUI() *fyne.Container {

	SpecialUserList = make(map[string]SpecialUserStruct)

	if globalConfiguration.SpecialUserList != nil {
		SpecialUserList = globalConfiguration.SpecialUserList
	}

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

					var MarkBth *widget.Button

					MarkBth = widget.NewButton("标记为特殊用户", func() {})
					//判断用户是否位于特殊列表中
					_, ok := SpecialUserList[LineTemp.OpenID]
					if ok {
						MarkBth.Disable()
					}
					var selectedYear, selectedMonth, selectedDay string
					MarkBth.OnTapped = func() {
						dialog.ShowCustomConfirm("选择截止日期", "确定", "取消", NewDatePicker(&selectedYear, &selectedMonth, &selectedDay), func(b bool) {
							timestamp, err := ConvertToTimestamp(selectedYear, selectedMonth, selectedDay)
							if err != nil {
								dialog.ShowError(errors.New("时间选择错误"), CtrlWindows)
								return
							}
							SpecialUserList[LineTemp.OpenID] = SpecialUserStruct{
								EndTime:  timestamp,
								UserName: LineTemp.UserName,
							}
							globalConfiguration.SpecialUserList = SpecialUserList
							SetConfig(globalConfiguration)
							MarkBth.Disable()

						}, CtrlWindows)
					}

					mu.Lock()
					LineBoxItem[LineTemp.OpenID] = container.NewHBox(canvas.NewText(LineTemp.UserName, LineTemp.PrintColor.ToRGBA()),
						widget.NewButton("删除", func() {
							vbox.Remove(LineBoxItem[LineTemp.OpenID])
							DeleteLine(LineTemp.OpenID)
							mu.Lock()
							delete(LineBoxItem, LineTemp.OpenID)
							mu.Unlock()
							CommonLength = len(OldLine.GuardLine)

						}), MarkBth)
					mu.Unlock()
					vbox.Add(LineBoxItem[LineTemp.OpenID])
				}

				for _, i2 := range OldLine.GiftLine {
					LineTemp := i2

					var MarkBth *widget.Button
					MarkBth = widget.NewButton("标记为特殊用户", func() {})
					//判断用户是否位于特殊列表中
					_, ok := SpecialUserList[LineTemp.OpenID]
					if ok {
						MarkBth.Disable()
					}
					var selectedYear, selectedMonth, selectedDay string
					MarkBth.OnTapped = func() {
						dialog.ShowCustomConfirm("选择截止日期", "确定", "取消", NewDatePicker(&selectedYear, &selectedMonth, &selectedDay), func(b bool) {
							timestamp, err := ConvertToTimestamp(selectedYear, selectedMonth, selectedDay)
							if err != nil {
								dialog.ShowError(errors.New("时间选择错误"), CtrlWindows)
								return
							}
							SpecialUserList[LineTemp.OpenID] = SpecialUserStruct{
								EndTime:  timestamp,
								UserName: LineTemp.UserName,
							}
							globalConfiguration.SpecialUserList = SpecialUserList
							SetConfig(globalConfiguration)
							MarkBth.Disable()

						}, CtrlWindows)
					}

					mu.Lock()
					LineBoxItem[LineTemp.OpenID] = container.NewHBox(canvas.NewText(LineTemp.UserName, LineTemp.PrintColor.ToRGBA()),
						widget.NewButton("删除", func() {
							vbox.Remove(LineBoxItem[LineTemp.OpenID])
							DeleteLine(LineTemp.OpenID)
							mu.Lock()
							delete(LineBoxItem, LineTemp.OpenID)
							mu.Unlock()
							CommonLength = len(OldLine.GiftLine)

						}), MarkBth)
					mu.Unlock()
					vbox.Add(LineBoxItem[LineTemp.OpenID])
				}

				if len(OldLine.CommonLine) != 0 {
					for _, i2 := range OldLine.CommonLine {
						LineTemp := i2

						var MarkBth *widget.Button
						MarkBth = widget.NewButton("标记为特殊用户", func() {})
						//判断用户是否位于特殊列表中
						_, ok := SpecialUserList[LineTemp.OpenID]
						if ok {
							MarkBth.Disable()
						}
						var selectedYear, selectedMonth, selectedDay string
						MarkBth.OnTapped = func() {
							dialog.ShowCustomConfirm("选择截止日期", "确定", "取消", NewDatePicker(&selectedYear, &selectedMonth, &selectedDay), func(b bool) {
								timestamp, err := ConvertToTimestamp(selectedYear, selectedMonth, selectedDay)
								if err != nil {
									dialog.ShowError(errors.New("时间选择错误"), CtrlWindows)
									return
								}
								SpecialUserList[LineTemp.OpenID] = SpecialUserStruct{
									EndTime:  timestamp,
									UserName: LineTemp.UserName,
								}
								globalConfiguration.SpecialUserList = SpecialUserList
								SetConfig(globalConfiguration)
								MarkBth.Disable()
							}, CtrlWindows)
						}

						mu.Lock()
						LineBoxItem[LineTemp.OpenID] = container.NewHBox(canvas.NewText(LineTemp.UserName, LineTemp.PrintColor.ToRGBA()),

							widget.NewButton("删除", func() {
								vbox.Remove(LineBoxItem[LineTemp.OpenID])
								DeleteLine(LineTemp.OpenID)
								mu.Lock()
								delete(LineBoxItem, LineTemp.OpenID)
								mu.Unlock()
								CommonLength = len(OldLine.CommonLine)
							}), MarkBth,
						)
						mu.Unlock()
						vbox.Add(LineBoxItem[LineTemp.OpenID])
					}
				}

				GuardLength = len(OldLine.GuardLine)
				GiftLength = len(OldLine.GiftLine)
				CommonLength = len(OldLine.CommonLine)
				vbox.Refresh()
			}
			time.Sleep(1 * time.Second)
		}
	}()
	return vbox
}

func NewDatePicker(selectedYear, selectedMonth, selectedDay *string) *fyne.Container {
	currentYear := time.Now().Year()
	years := make([]string, 0)
	for i := currentYear + 20; i >= currentYear; i-- {
		years = append(years, strconv.Itoa(i))
	}

	months := make([]string, 12)
	for i := 1; i <= 12; i++ {
		months[i-1] = strconv.Itoa(i)
	}

	days := make([]string, 31)
	for i := 1; i <= 31; i++ {
		days[i-1] = strconv.Itoa(i)
	}

	yearSelect := widget.NewSelect(years, func(value string) {
		*selectedYear = value
	})
	monthSelect := widget.NewSelect(months, func(value string) {
		*selectedMonth = value
	})
	daySelect := widget.NewSelect(days, func(value string) {
		*selectedDay = value
	})

	return container.NewHBox(yearSelect, monthSelect, daySelect)
}

func ConvertToTimestamp(year, month, day string) (int64, error) {
	y, err := strconv.Atoi(year)
	if err != nil {
		return 0, err
	}
	m, err := strconv.Atoi(month)
	if err != nil {
		return 0, err
	}
	d, err := strconv.Atoi(day)
	if err != nil {
		return 0, err
	}

	date := time.Date(y, time.Month(m), d, 0, 0, 0, 0, time.Local)
	return date.Unix(), nil
}

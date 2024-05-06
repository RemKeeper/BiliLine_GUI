package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"golang.org/x/image/colornames"
)

func MakeUpdateUI(NewVersion VersionSct) *fyne.Container {

	versionLabel := widget.NewLabel(NewVersion.Version)
	versionCountLabel := widget.NewLabel(fmt.Sprintf("版本号:%d", NewVersion.VersionCount))
	updateDateLabel := widget.NewLabel(NewVersion.UpdateDate)

	changelogTitle := widget.NewLabel("变更记录:")
	var changelogItems []fyne.CanvasObject
	for _, msg := range NewVersion.Changelog {
		changelogItems = append(changelogItems, widget.NewLabel(msg))
	}

	content := container.NewVBox(
		canvas.NewText("检测到新版本", colornames.Aqua),
		versionLabel,
		versionCountLabel,
		updateDateLabel,
		changelogTitle,
		container.NewVBox(changelogItems...),
		widget.NewButton("下载更新", func() {
			err := AgreeOpenUrl(NewVersion.UpdateUrl)
			if err != nil {
				return
			}
		}),
	)
	return content
}

package gui

import (
	"Go-Cleaner/cleaner"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func RunUI(rules []cleaner.Rule) {
	myApp := app.New()
	myWindow := myApp.NewWindow("🧹 Go-Cleaner by kukuqi666")

	filesList := widget.NewMultiLineEntry()
	filesList.SetPlaceHolder("Matched files will appear here")

	scanButton := widget.NewButton("🔍 Scan", func() {
		filesList.SetText("")
		for _, rule := range rules {
			files := cleaner.MatchFiles(rule)
			for _, f := range files {
				filesList.SetText(filesList.Text + f + "\n")
			}
		}
	})

	cleanButton := widget.NewButton("🗑️ Clean", func() {
		dialog.ShowConfirm("Confirm", "Delete listed files?", func(ok bool) {
			if ok {
				for _, rule := range rules {
					files := cleaner.MatchFiles(rule)
					for _, f := range files {
						cleaner.DeleteFile(f)
					}
				}
				dialog.ShowInformation("Done", "Cleanup finished.", myWindow)
			}
		}, myWindow)
	})

	aboutButton := widget.NewButton("ℹ️ About", func() {
		dialog.ShowInformation("About", "Go-Cleaner v1.0.0\nAuthor: kukuqi666\nEmail: kukuqi666@gmail.com", myWindow)
	})

	content := container.NewVBox(scanButton, cleanButton, aboutButton, filesList)

	myWindow.SetContent(content)
	myWindow.Resize(fyne.NewSize(700, 600))
	myWindow.ShowAndRun()
}

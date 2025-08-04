package gui

import (
	"fmt"
	"time"

	"Go-Cleaner/cleaner"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)


var scanButton *widget.Button
var cleanButton *widget.Button

func RunUI(rules []cleaner.Rule) {
	myApp := app.New()
	myWindow := myApp.NewWindow("🧹 Go-Cleaner by kukuqi666")

	// 状态标签
	statusLabel := widget.NewLabel("Ready to scan")
	
	// 文件计数标签
	fileCountLabel := widget.NewLabel("Files found: 0")
	
	// 进度条
	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	// 文件列表
	filesList := widget.NewMultiLineEntry()
	filesList.SetPlaceHolder("Click 'Scan' to find files matching cleanup rules")
	filesList.Wrapping = fyne.TextWrapWord

	// 滚动容器包装文件列表
	filesScroll := container.NewScroll(filesList)
	filesScroll.SetMinSize(fyne.NewSize(650, 300))

	var matchedFiles []string
	var isScanning bool

	scanButton = widget.NewButton("🔍 Scan for Files", func() {
		if isScanning {
			return
		}
		
		isScanning = true
		scanButton.SetText("🔄 Scanning...")
		scanButton.Disable()
		statusLabel.SetText("Scanning directories...")
		filesList.SetText("")
		fileCountLabel.SetText("Files found: 0")
		progressBar.Show()
		progressBar.SetValue(0)
		matchedFiles = []string{}

		// 异步扫描以避免阻塞UI
		go func() {
			var allFiles []string
			totalRules := len(rules)
			
			for i, rule := range rules {
				// 更新进度
				progressBar.SetValue(float64(i) / float64(totalRules))
				statusLabel.SetText(fmt.Sprintf("Scanning: %s", rule.Name))
				
				files := cleaner.MatchFiles(rule)
				allFiles = append(allFiles, files...)
				
				// 实时更新文件列表
				if len(files) > 0 {
					for _, f := range files {
						filesList.SetText(filesList.Text + f + "\n")
					}
					fileCountLabel.SetText(fmt.Sprintf("Files found: %d", len(allFiles)))
				}
				
				// 短暂暂停以允许UI更新
				time.Sleep(10 * time.Millisecond)
			}
			
			matchedFiles = allFiles
			progressBar.SetValue(1.0)
			statusLabel.SetText(fmt.Sprintf("Scan completed. Found %d files", len(allFiles)))
			
			// 短暂显示完成进度后隐藏
			time.Sleep(500 * time.Millisecond)
			progressBar.Hide()
			
			scanButton.SetText("🔍 Scan for Files")
			scanButton.Enable()
			isScanning = false
		}()
	})

	cleanButton = widget.NewButton("🗑️ Clean Files", func() {
		if len(matchedFiles) == 0 {
			dialog.ShowInformation("No Files", "No files to clean. Please scan first.", myWindow)
			return
		}
		
		confirmMsg := fmt.Sprintf("Delete %d files?\n\nThis action cannot be undone!", len(matchedFiles))
		dialog.ShowConfirm("⚠️ Confirm Deletion", confirmMsg, func(ok bool) {
			if ok {
				cleanButton.Disable()
				statusLabel.SetText("Deleting files...")
				progressBar.Show()
				progressBar.SetValue(0)
				
				go func() {
					deleted := 0
					failed := 0
					
					for i, file := range matchedFiles {
						progressBar.SetValue(float64(i) / float64(len(matchedFiles)))
						
						err := cleaner.DeleteFile(file)
						if err != nil {
							failed++
						} else {
							deleted++
						}
						
						time.Sleep(5 * time.Millisecond)
					}
					
					progressBar.SetValue(1.0)
					statusLabel.SetText(fmt.Sprintf("Cleanup finished: %d deleted, %d failed", deleted, failed))
					
					time.Sleep(500 * time.Millisecond)
					progressBar.Hide()
					cleanButton.Enable()
					
					// 清空文件列表
					filesList.SetText("")
					fileCountLabel.SetText("Files found: 0")
					matchedFiles = []string{}
					
					if failed > 0 {
						dialog.ShowInformation("Cleanup Complete", 
							fmt.Sprintf("Deleted %d files successfully.\n%d files failed to delete.", deleted, failed), myWindow)
					} else {
						dialog.ShowInformation("Cleanup Complete", 
							fmt.Sprintf("Successfully deleted %d files.", deleted), myWindow)
					}
				}()
			}
		}, myWindow)
	})

	aboutButton := widget.NewButton("ℹ️ About", func() {
		aboutContent := `Go-Cleaner v1.0.0

A powerful file cleanup utility written in Go.

Author: kukuqi666
Email: kukuqi666@gmail.com

Features:
• Smart file scanning with customizable rules
• Safe deletion with confirmation
• Real-time progress tracking
• Cross-platform compatibility`

		dialog.ShowInformation("About Go-Cleaner", aboutContent, myWindow)
	})

	// 创建更好的布局
	topContainer := container.NewHBox(scanButton, cleanButton, aboutButton)
	
	statusContainer := container.NewHBox(statusLabel, widget.NewSeparator(), fileCountLabel)
	
	mainContent := container.NewVBox(
		widget.NewCard("Actions", "", topContainer),
		widget.NewSeparator(),
		statusContainer,
		progressBar,
		widget.NewCard("Matched Files", "", filesScroll),
	)

	myWindow.SetContent(container.NewPadded(mainContent))
	myWindow.Resize(fyne.NewSize(800, 700))
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}

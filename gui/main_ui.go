package gui

import (
	"fmt"
	"time"

	"Go-Cleaner/cleaner"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var scanButton *widget.Button
var cleanButton *widget.Button

func RunUI(rules []cleaner.Rule) {
	myApp := app.New()
	myWindow := myApp.NewWindow("🧹 Go-Cleaner 文件清理工具")

	// 设置应用程序图标
	iconURI := storage.NewFileURI("assets/icon.png")
	icon, err := storage.LoadResourceFromURI(iconURI)
	if err == nil {
		myWindow.SetIcon(icon)
	}

	// 状态标签
	statusLabel := widget.NewLabel("准备就绪，点击扫描开始")

	// 文件计数标签
	fileCountLabel := widget.NewLabel("找到文件: 0")

	// 进度条
	progressBar := widget.NewProgressBar()
	progressBar.Hide()

	// 文件列表
	filesList := widget.NewMultiLineEntry()
	filesList.SetPlaceHolder("点击\"扫描文件\"按钮开始查找匹配清理规则的文件")
	filesList.Wrapping = fyne.TextWrapWord

	// 滚动容器包装文件列表
	filesScroll := container.NewScroll(filesList)
	filesScroll.SetMinSize(fyne.NewSize(650, 300))

	var matchedFiles []string
	var isScanning bool

	scanButton = widget.NewButton("🔍 扫描文件", func() {
		if isScanning {
			return
		}

		isScanning = true
		scanButton.SetText("🔄 扫描中...")
		scanButton.Disable()
		statusLabel.SetText("正在扫描目录...")
		filesList.SetText("")
		fileCountLabel.SetText("找到文件: 0")
		progressBar.Show()
		progressBar.SetValue(0)
		matchedFiles = []string{}

		// 异步扫描以避免阻塞UI
		go func() {
			var allFiles []string
			totalRules := len(rules)

			for i, rule := range rules {
				// 更新进度
				fyne.Do(func() {
					progressBar.SetValue(float64(i) / float64(totalRules))
					statusLabel.SetText(fmt.Sprintf("正在扫描: %s", rule.Name))
				})

				files := cleaner.MatchFiles(rule)
				allFiles = append(allFiles, files...)

				// 实时更新文件列表
				if len(files) > 0 {
					fyne.Do(func() {
						for _, f := range files {
							filesList.SetText(filesList.Text + f + "\n")
						}
						fileCountLabel.SetText(fmt.Sprintf("找到文件: %d", len(allFiles)))
					})
				}

				// 短暂暂停以允许UI更新
				time.Sleep(10 * time.Millisecond)
			}

			matchedFiles = allFiles
			fyne.Do(func() {
				progressBar.SetValue(1.0)
				statusLabel.SetText(fmt.Sprintf("扫描完成，找到 %d 个文件", len(allFiles)))
			})

			// 短暂显示完成进度后隐藏
			time.Sleep(500 * time.Millisecond)
			fyne.Do(func() {
				progressBar.Hide()
				scanButton.SetText("🔍 扫描文件")
				scanButton.Enable()
				isScanning = false
			})
		}()
	})

	cleanButton = widget.NewButton("🗑️ 清理文件", func() {
		if len(matchedFiles) == 0 {
			dialog.ShowInformation("没有文件", "没有找到可清理的文件。请先进行扫描。", myWindow)
			return
		}

		confirmMsg := fmt.Sprintf("确定要删除 %d 个文件吗？\n\n此操作无法撤销！", len(matchedFiles))
		dialog.ShowConfirm("⚠️ 确认删除", confirmMsg, func(ok bool) {
			if ok {
				cleanButton.Disable()
				statusLabel.SetText("正在删除文件...")
				progressBar.Show()
				progressBar.SetValue(0)

				go func() {
					deleted := 0
					failed := 0
					var failedFiles []string
					var failedReasons []string

					for i, file := range matchedFiles {
						fyne.Do(func() {
							progressBar.SetValue(float64(i) / float64(len(matchedFiles)))
						})

						result := cleaner.DeleteFile(file)
						if !result.Success {
							failed++
							failedFiles = append(failedFiles, file)
							failedReasons = append(failedReasons, result.Reason)
						} else {
							deleted++
						}

						time.Sleep(5 * time.Millisecond)
					}

					fyne.Do(func() {
						progressBar.SetValue(1.0)
						statusLabel.SetText(fmt.Sprintf("清理完成: 成功删除 %d 个文件，失败 %d 个", deleted, failed))
					})

					time.Sleep(500 * time.Millisecond)
					fyne.Do(func() {
						progressBar.Hide()
						cleanButton.Enable()

						// 清空文件列表
						filesList.SetText("")
						fileCountLabel.SetText("找到文件: 0")
						matchedFiles = []string{}
					})

					// 显示详细结果
					if failed > 0 {
						fyne.Do(func() {
							// 构建失败文件列表
							failedDetails := "删除失败的文件详情：\n\n"
							for i, file := range failedFiles {
								failedDetails += fmt.Sprintf("文件: %s\n原因: %s\n\n", file, failedReasons[i])
							}

							dialog.ShowInformation("清理完成",
								fmt.Sprintf("成功删除 %d 个文件。\n%d 个文件删除失败。\n\n%s", deleted, failed, failedDetails), myWindow)
						})
					} else {
						fyne.Do(func() {
							dialog.ShowInformation("清理完成",
								fmt.Sprintf("成功删除 %d 个文件。", deleted), myWindow)
						})
					}
				}()
			}
		}, myWindow)
	})

	aboutButton := widget.NewButton("ℹ️ 关于", func() {
		aboutContent := `Go-Cleaner v1.0.0

一个强大的文件清理实用程序，使用Go语言编写。

作者：kukuqi
邮箱：kukuqi666@gmail.com

主要特性：
• 基于可配置规则的智能文件扫描
• 安全删除并确认操作
• 实时进度跟踪
• 跨平台兼容性
• 现代化GUI界面
• 详细的错误提示`

		dialog.ShowInformation("关于清理工具", aboutContent, myWindow)
	})

	// 创建更好的布局
	topContainer := container.NewHBox(scanButton, cleanButton, aboutButton)

	statusContainer := container.NewHBox(statusLabel, widget.NewSeparator(), fileCountLabel)

	mainContent := container.NewVBox(
		widget.NewCard("操作面板", "", topContainer),
		widget.NewSeparator(),
		statusContainer,
		progressBar,
		widget.NewCard("匹配文件列表", "", filesScroll),
	)

	myWindow.SetContent(container.NewPadded(mainContent))
	myWindow.Resize(fyne.NewSize(800, 700))
	myWindow.CenterOnScreen()
	myWindow.ShowAndRun()
}

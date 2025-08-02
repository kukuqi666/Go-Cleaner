🗂️ **Go-Cleaner 目录结构**

```
Go-Cleaner/
 ├── main.go
 ├── go.mod
 ├── cleaner/
 │   ├── cleaner.go
 │   ├── rule.go
 ├── gui/
 │   ├── main_ui.go
 ├── rules/
 │   └── default_rules.json
 ├── assets/
 │   └── icon.png (占位，可用空文件)
 ├── build.bat
 ├── build.sh
 ├── README.md
```

------

下面我一份一份发👇

------

## ✅ 1️⃣ `go.mod`

```go
module Go-Cleaner

go 1.21

require fyne.io/fyne/v2 v2.4.3
```

------

## ✅ 2️⃣ `main.go`

```go
package main

import (
	"encoding/json"
	"os"

	"Go-Cleaner/cleaner"
	"Go-Cleaner/gui"
)

/*
  Go-Cleaner
  Author: kukuqi666
  Email: kukuqi666@gmail.com
  Version: 1.0.0
*/

func main() {
	file, _ := os.ReadFile("rules/default_rules.json")
	var rules []cleaner.Rule
	json.Unmarshal(file, &rules)

	gui.RunUI(rules)
}
```

------

## ✅ 3️⃣ `cleaner/rule.go`

```go
package cleaner

type Rule struct {
	Name        string   `json:"name"`
	Paths       []string `json:"paths"`
	Extensions  []string `json:"extensions"`
	MinSizeMB   int      `json:"min_size_mb"`
	MaxAgeDays  int      `json:"max_age_days"`
	Description string   `json:"description"`
}
```

------

## ✅ 4️⃣ `cleaner/cleaner.go`

```go
package cleaner

import (
	"io/fs"
	"os"
	"path/filepath"
	"time"
)

func MatchFiles(rule Rule) []string {
	var matched []string
	for _, path := range rule.Paths {
		filepath.Walk(path, func(p string, info fs.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}

			// 扩展名匹配
			matchedExt := false
			for _, ext := range rule.Extensions {
				if filepath.Ext(p) == ext {
					matchedExt = true
					break
				}
			}
			if !matchedExt {
				return nil
			}

			// 大小匹配
			sizeMB := info.Size() / (1024 * 1024)
			if sizeMB < int64(rule.MinSizeMB) {
				return nil
			}

			// 修改时间匹配
			ageDays := time.Since(info.ModTime()).Hours() / 24
			if int(ageDays) < rule.MaxAgeDays {
				return nil
			}

			matched = append(matched, p)
			return nil
		})
	}
	return matched
}

func DeleteFile(path string) error {
	return os.Remove(path)
}
```

------

## ✅ 5️⃣ `gui/main_ui.go`

```go
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
```

------

## ✅ 6️⃣ `rules/default_rules.json`

```json
[
  {
    "name": "Temp Files",
    "paths": ["C:\\\\Windows\\\\Temp", "/tmp"],
    "extensions": [".tmp", ".temp"],
    "min_size_mb": 0,
    "max_age_days": 0,
    "description": "Temporary files"
  },
  {
    "name": "Log Files",
    "paths": ["/var/log", "C:\\\\Logs"],
    "extensions": [".log"],
    "min_size_mb": 1,
    "max_age_days": 30,
    "description": "Old log files"
  }
]
```

------

## ✅ 7️⃣ `assets/icon.png`

可以放一个 128x128 PNG 图标，先留空也行，或者放个占位图片就好。

------

## ✅ 8️⃣ `build.bat`（Windows）

```bat
@echo off
echo Building Go-Cleaner.exe ...
go build -o Go-Cleaner.exe main.go
echo Done!
pause
```

------

## ✅ 9️⃣ `build.sh`（Linux/macOS）

```bash
#!/bin/bash
echo "Building Go-Cleaner ..."
go build -o Go-Cleaner main.go
echo "Done!"
```

------

## ✅  🔖 `README.md`

~~~md
# Go-Cleaner

**Author:** kukuqi666  
**Email:** kukuqi666@gmail.com  
**Version:** 1.0.0

## How to use

1. Install dependencies:  
   `go install fyne.io/fyne/v2/cmd/fyne@latest`

2. Run:
   ```bash
   go run main.go
~~~

1. Build:
   - Windows: double click `build.bat`
   - Linux/macOS: `chmod +x build.sh && ./build.sh`
2. Customize your rules in `rules/default_rules.json`.

Happy cleaning! 🧹

```
---

## ✅ 10️⃣ 全部保存后

- 把这些文件放到 `Go-Cleaner` 文件夹
- Windows 双击 `build.bat` 就会生成 `Go-Cleaner.exe`
- Linux/macOS 执行 `./build.sh` 就生成 `Go-Cleaner`

---

到这里，就是完整可运行、可打包的【增强版骨架】🎉  
需要我现在帮你做个**额外打包 ZIP 用法/压缩命令**，或者下一步做**在线规则库更新示例**吗？要的话直接说【要】！🚀
```
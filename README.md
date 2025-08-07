# 🧹 Go-Cleaner

> **强大的Windows系统文件清理工具**  
> 基于Go语言开发，具有现代化GUI界面

[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Platform](https://img.shields.io/badge/Platform-Windows-lightgrey.svg)](https://github.com/fyne-io/fyne)

## ✨ 特性

- 🔍 **智能文件扫描** - 基于可配置规则的文件匹配
- 🗑️ **安全删除** - 批量删除文件，带有确认对话框
- 📊 **实时进度跟踪** - 显示扫描和删除进度
- 🎨 **现代化GUI** - 基于Fyne框架的美观界面
- 🔧 **可自定义规则** - 支持JSON格式的清理规则配置
- 🛡️ **跨平台兼容** - 支持Windows系统
- ⚡ **高性能** - 异步操作，不阻塞界面

## 🚀 快速开始

### 环境要求
- Go 1.21 或更高版本
- Windows 操作系统

### 下载

- [Github Release](https://github.com/kukuqi666/Go-Cleaner/releases)
- [Github Actions](https://github.com/kukuqi666/Go-Cleaner/actions)


### 源码编译运行
```bash
go install fyne.io/fyne/v2/cmd/fyne@latest
```

### 运行程序
```bash
# 克隆项目
git clone https://github.com/your-username/Go-Cleaner.git
cd Go-Cleaner

# 运行程序
go run main.go
```

### 构建可执行文件
```bash
# 构建Windows可执行文件
go build -o go-cleaner.exe main.go

# 或使用fyne打包
fyne package -os windows
```

## 📋 清理规则

程序支持多种清理规则，包括：

| 规则类型 | 目标路径 | 文件类型 | 条件 |
|---------|---------|---------|------|
| **Windows临时文件** | `C:\Windows\Temp` | `.tmp`, `.temp` | 无限制 |
| **浏览器缓存** | Chrome/Edge缓存目录 | `.dat`, `.cache` | 7天以上 |
| **系统日志** | `C:\Windows\Logs` | `.log`, `.etl` | 30天以上，>1MB |
| **Windows更新缓存** | `C:\Windows\SoftwareDistribution` | `.cab`, `.msu` | 7天以上，>5MB |
| **回收站** | `C:\$Recycle.Bin` | 所有文件 | 7天以上 |
| **预取文件** | `C:\Windows\Prefetch` | `.pf` | 30天以上 |
| **缩略图缓存** | Explorer缓存目录 | `.db` | 30天以上，>1MB |
| **错误报告** | `C:\ProgramData\Microsoft\Windows\WER` | `.dmp`, `.tmp` | 30天以上，>1MB |

## 🛠️ 自定义规则

编辑 `rules/default_rules.json` 文件来自定义清理规则：

```json
{
  "name": "自定义规则",
  "paths": ["C:\\Users\\%USERNAME%\\Desktop\\测试文件夹"],
  "extensions": [".txt", ".log"],
  "min_size_mb": 1,
  "max_age_days": 30,
  "description": "清理桌面测试文件夹中的旧文件"
}
```

### 规则参数说明

- `name`: 规则名称
- `paths`: 目标路径数组（支持环境变量如 `%USERNAME%`）
- `extensions`: 文件扩展名数组（`.*` 匹配所有文件）
- `min_size_mb`: 最小文件大小（MB）
- `max_age_days`: 最大文件年龄（天）
- `description`: 规则描述

## 🎯 使用指南

1. **启动程序** - 运行 `go run main.go` 或双击可执行文件
2. **扫描文件** - 点击"🔍 扫描文件"按钮开始扫描
3. **查看结果** - 在文件列表中查看匹配的文件
4. **确认删除** - 点击"🗑️ 清理文件"按钮，确认删除操作
5. **完成清理** - 查看清理结果和统计信息

## ⚠️ 注意事项

- ⚠️ **备份重要文件** - 删除操作不可逆，请确保重要文件已备份
- 🔐 **管理员权限** - 某些系统目录可能需要管理员权限
- 📁 **路径验证** - 程序会自动跳过不存在的路径
- 🕒 **时间限制** - 大文件扫描可能需要较长时间

## 🏗️ 项目结构

```
Go-Cleaner/
├── main.go                 # 程序入口
├── gui/
│   └── main_ui.go         # GUI界面实现
├── cleaner/
│   ├── cleaner.go         # 核心清理逻辑
│   └── rule.go           # 规则结构定义
├── rules/
│   └── default_rules.json # 默认清理规则
├── assets/
│   └── icon.png          # 程序图标
├── build.bat             # 构建脚本
└── README.md            # 项目说明
```

## 🛡️ 安全特性

- ✅ **确认对话框** - 删除前必须确认
- ✅ **文件列表显示** - 清楚显示将要删除的文件
- ✅ **进度跟踪** - 实时显示操作进度
- ✅ **错误处理** - 优雅处理权限和路径错误
- ✅ **异步操作** - 不阻塞用户界面

## 🤝 贡献指南

欢迎提交Issue和Pull Request！

1. Fork 本项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 打开 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 👨‍💻 作者

**kukuqi666**  
📧 Email: kukuqi666@gmail.com  
🔗 GitHub: [@kukuqi666](https://github.com/kukuqi666)

## 🙏 致谢

- [Fyne](https://fyne.io/) - 优秀的Go GUI框架
- [Go](https://golang.org/) - 强大的编程语言
- 所有贡献者和用户的支持

---

<div align="center">

**⭐ 如果这个项目对你有帮助，请给它一个星标！**

</div>

package cleaner

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

// 自定义环境变量展开函数，支持Windows的%VAR%格式
func expandWindowsEnv(path string) string {
	// 匹配 %VAR% 格式的环境变量
	re := regexp.MustCompile(`%([^%]+)%`)
	return re.ReplaceAllStringFunc(path, func(match string) string {
		// 提取变量名
		varName := match[1 : len(match)-1]
		if value := os.Getenv(varName); value != "" {
			return value
		}
		return match // 如果环境变量不存在，返回原字符串
	})
}

func MatchFiles(rule Rule) []string {
	var matched []string
	for _, path := range rule.Paths {
		// 修复路径格式问题
		cleanPath := strings.ReplaceAll(path, "\\\\", "\\")
		// 展开环境变量（支持Windows格式）
		expandedPath := expandWindowsEnv(cleanPath)
		
		// 如果路径不存在，跳过
		if _, err := os.Stat(expandedPath); os.IsNotExist(err) {
			fmt.Printf("路径不存在: %s\n", expandedPath)
			continue
		}
		
		filepath.Walk(expandedPath, func(p string, info fs.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				return nil
			}

			// 扩展名匹配
			matchedExt := false
			if len(rule.Extensions) == 0 {
				matchedExt = true
			} else {
				fileExt := filepath.Ext(p)
				for _, ext := range rule.Extensions {
					// 支持 ".*" 匹配所有文件
					if ext == ".*" || fileExt == ext {
						matchedExt = true
						break
					}
				}
			}
			if !matchedExt {
				return nil
			}

			// 大小匹配
			sizeMB := info.Size() / (1024 * 1024)
			if rule.MinSizeMB > 0 && sizeMB < int64(rule.MinSizeMB) {
				return nil
			}

			// 修改时间匹配
			if rule.MaxAgeDays > 0 {
				ageDays := time.Since(info.ModTime()).Hours() / 24
				if int(ageDays) < rule.MaxAgeDays {
					return nil
				}
			}

			matched = append(matched, p)
			return nil
		})
	}
	return matched
}

// DeleteResult 删除结果结构
type DeleteResult struct {
	Success bool
	Error   string
	Reason  string
}

func DeleteFile(path string) DeleteResult {
	err := os.Remove(path)
	if err == nil {
		return DeleteResult{Success: true}
	}
	
	// 分析错误原因
	reason := getDeleteErrorReason(err)
	return DeleteResult{
		Success: false,
		Error:   err.Error(),
		Reason:  reason,
	}
}

// getDeleteErrorReason 分析删除失败的原因
func getDeleteErrorReason(err error) string {
	if err == nil {
		return ""
	}
	
	errStr := err.Error()
	
	// 检查是否是权限错误
	if strings.Contains(errStr, "Access is denied") || strings.Contains(errStr, "拒绝访问") {
		return "文件被系统保护或需要管理员权限"
	}
	
	// 检查是否是文件被占用
	if strings.Contains(errStr, "The process cannot access the file") || 
	   strings.Contains(errStr, "另一个程序正在使用此文件") {
		return "文件正在被其他程序使用，请关闭相关程序后重试"
	}
	
	// 检查是否是路径不存在
	if strings.Contains(errStr, "The system cannot find the path") || 
	   strings.Contains(errStr, "找不到指定的路径") {
		return "文件路径不存在或已被删除"
	}
	
	// 检查是否是只读文件
	if strings.Contains(errStr, "read-only") || strings.Contains(errStr, "只读") {
		return "文件为只读属性，无法删除"
	}
	
	// 检查是否是网络驱动器错误
	if strings.Contains(errStr, "network") || strings.Contains(errStr, "网络") {
		return "网络驱动器连接问题"
	}
	
	// 检查是否是磁盘空间不足
	if strings.Contains(errStr, "disk space") || strings.Contains(errStr, "磁盘空间") {
		return "磁盘空间不足"
	}
	
	// 检查是否是文件系统错误
	if strings.Contains(errStr, "file system") || strings.Contains(errStr, "文件系统") {
		return "文件系统错误"
	}
	
	// 默认错误信息
	return "未知错误，请检查文件权限和状态"
}

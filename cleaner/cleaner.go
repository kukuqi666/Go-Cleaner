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

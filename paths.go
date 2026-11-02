package main

import (
	"regexp"
	"strings"
)

var windowsPathRe = regexp.MustCompile(`[a-zA-Z]:\\`)

// parseDragPaths 解析拖拽输入的路径，跨平台兼容：
// - Windows 拖拽出的路径含盘符（如 C:\Users\...），按盘符位置切分
// - macOS/Linux 拖拽出的路径为引号包裹的绝对路径，按引号感知的空白切分
func parseDragPaths(input string) []string {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil
	}

	if windowsPathRe.MatchString(input) {
		return splitWindowsPaths(input)
	}
	return splitQuotedPaths(input)
}

func splitWindowsPaths(input string) []string {
	indices := windowsPathRe.FindAllStringIndex(input, -1)

	var paths []string
	for i := 0; i < len(indices); i++ {
		start := indices[i][0]
		end := len(input)
		if i+1 < len(indices) {
			end = indices[i+1][0]
		}

		raw := input[start:end]
		p := strings.TrimSpace(raw)
		p = strings.Trim(p, "\"")
		p = strings.TrimSpace(p)

		if p != "" {
			paths = append(paths, p)
		}
	}
	return paths
}

func splitQuotedPaths(input string) []string {
	var paths []string
	var sb strings.Builder
	inQuotes := false

	for i := 0; i < len(input); i++ {
		char := input[i]
		if char == '"' {
			inQuotes = !inQuotes
			continue
		}
		if char == ' ' && !inQuotes {
			if sb.Len() > 0 {
				paths = append(paths, sb.String())
				sb.Reset()
			}
			continue
		}
		sb.WriteByte(char)
	}
	if sb.Len() > 0 {
		paths = append(paths, sb.String())
	}
	return paths
}

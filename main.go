package main

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

var extensions = map[string]bool{
	".mp4": true, ".mov": true, ".nef": true,
	".jpg": true, ".png": true, ".mp3": true,
}

func main() {
	if len(os.Args) > 1 {
		runTask(os.Args[1:])
		fmt.Println("\n任务结束，按回车退出...")
		bufio.NewReader(os.Stdin).ReadString('\n')
		return
	}

	for {
		clearScreen()
		fmt.Println("\033[36m+------------------------------------------+\033[0m")
		fmt.Println("\033[36m|      媒体文件自动重命名工具     |\033[0m")
		fmt.Println("\033[36m+------------------------------------------+\033[0m")
		fmt.Println("请拖入【单个/多个 文件/文件夹】并按回车:")
		fmt.Print("\033[32m>> \033[0m")

		scanner := bufio.NewScanner(os.Stdin)
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		if strings.TrimSpace(input) == "" {
			continue
		}

		// 核心修复：更精准的路径拆分，确保每一项都不漏
		paths := splitWindowsPaths(input)
		if len(paths) > 0 {
			runTask(paths)
		} else {
			fmt.Println("\033[31m[ 错误 ] 未能识别有效路径，请确保拖入的是本地文件。\033[0m")
		}

		fmt.Println("\n\033[33m[ 任务结束 ] 按回车键 (Enter) 重新开始...\033[0m")
		bufio.NewReader(os.Stdin).ReadString('\n')
	}
}

// 修复后的拆分逻辑：通过盘符位置进行物理切割
func splitWindowsPaths(input string) []string {
	// 匹配 Windows 路径特征： X:\
	re := regexp.MustCompile(`[a-zA-Z]:\\`)
	indices := re.FindAllStringIndex(input, -1)

	var paths []string
	for i := 0; i < len(indices); i++ {
		start := indices[i][0]
		end := len(input)
		if i+1 < len(indices) {
			end = indices[i+1][0]
		}

		// 截取两个盘符之间的内容
		raw := input[start:end]
		// 清理：去掉首尾空格和双引号
		p := strings.TrimSpace(raw)
		p = strings.Trim(p, "\"")
		p = strings.TrimSpace(p)

		if p != "" {
			paths = append(paths, p)
		}
	}
	return paths
}

func runTask(paths []string) {
	start := time.Now()
	var allFiles []string

	for _, p := range paths {
		info, err := os.Stat(p)
		if err != nil {
			continue
		}

		if info.IsDir() {
			filepath.WalkDir(p, func(path string, d fs.DirEntry, err error) error {
				if err == nil && !d.IsDir() {
					ext := filepath.Ext(path)
					if extensions[strings.ToLower(ext)] {
						allFiles = append(allFiles, path)
					}
				}
				return nil
			})
		} else {
			if extensions[strings.ToLower(filepath.Ext(p))] {
				allFiles = append(allFiles, p)
			}
		}
	}

	total := len(allFiles)
	if total == 0 {
		fmt.Println("\033[31m[ 提示 ] 未发现匹配的媒体文件。\033[0m")
		return
	}

	sort.Strings(allFiles)
	success := 0
	for _, f := range allFiles {
		if err := processFile(f); err == nil {
			success++
		}
	}

	elapsed := time.Since(start).Round(time.Millisecond)
	printFinalReport(success, elapsed)
}

func processFile(oldPath string) error {
	originalExt := filepath.Ext(oldPath)
	fi, err := os.Stat(oldPath)
	if err != nil {
		return err
	}

	t := fi.ModTime()
	prefix := t.Format("20060102_150405_")
	msVal := t.Nanosecond() / 1000000
	if msVal == 0 {
		msVal = 1
	}

	dir := filepath.Dir(oldPath)
	oldName := filepath.Base(oldPath)

	for {
		newName := fmt.Sprintf("%s%03d%s", prefix, msVal, originalExt)
		newPath := filepath.Join(dir, newName)
		if oldName == newName {
			return nil
		}
		if _, err := os.Stat(newPath); os.IsNotExist(err) {
			return os.Rename(oldPath, newPath)
		}
		msVal++
		if msVal > 999 {
			msVal = 1
		}
		if msVal > 2000 {
			return fmt.Errorf("limit")
		}
	}
}

func printFinalReport(count int, duration time.Duration) {
	cyan, green, reset := "\033[36m", "\033[1;32m", "\033[0m"
	fmt.Printf("\n%s+------------------------------------------+%s\n", cyan, reset)
	fmt.Printf("%s|%s       [ 媒体文件处理任务完成 ]         %s%s|%s\n", cyan, green, reset, cyan, reset)
	fmt.Printf("%s+------------------------------------------+%s\n", cyan, reset)
	fmt.Printf(" >> 成功处理总数 : %s%d%s\n", green, count, reset)
	fmt.Printf(" >> 任务执行耗时 : %s%v%s\n", green, duration, reset)
	fmt.Printf("%s+------------------------------------------+%s\n", cyan, reset)
}

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

// EXIF Renamer — 根据文件修改时间批量重命名媒体文件
// 支持拖拽文件/文件夹，自动递归处理子目录
package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		runTask(os.Args[1:])
		fmt.Println("\n\033[33m[ 任务结束 ] 按回车键 (Enter) 继续到交互模式...\033[0m")
		bufio.NewReader(os.Stdin).ReadString('\n')
	}

	for {
		clearScreen()
		fmt.Println("\033[36m+------------------------------------------+\033[0m")
		fmt.Println("\033[36m|            媒体文件自动重命名工具        |\033[0m")
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

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

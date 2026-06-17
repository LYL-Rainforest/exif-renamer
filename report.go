package main

import (
	"fmt"
	"time"
)

func printFinalReport(count int, duration time.Duration) {
	cyan, green, reset := "\033[36m", "\033[1;32m", "\033[0m"
	fmt.Printf("\n%s+------------------------------------------+%s\n", cyan, reset)
	fmt.Printf("%s|%s       [ 媒体文件处理任务完成 ]         %s%s|%s\n", cyan, green, reset, cyan, reset)
	fmt.Printf("%s+------------------------------------------+%s\n", cyan, reset)
	fmt.Printf(" >> 成功处理总数 : %s%d%s\n", green, count, reset)
	fmt.Printf(" >> 任务执行耗时 : %s%v%s\n", green, duration, reset)
	fmt.Printf("%s+------------------------------------------+%s\n", cyan, reset)
}

package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"
)

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
					allFiles = append(allFiles, path)
				}
				return nil
			})
		} else {
			allFiles = append(allFiles, p)
		}
	}

	total := len(allFiles)
	if total == 0 {
		fmt.Println("\033[31m[ 提示 ] 未发现可处理的文件。\033[0m")
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

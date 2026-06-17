package main

import (
	"regexp"
	"strings"
)

func splitWindowsPaths(input string) []string {
	re := regexp.MustCompile(`[a-zA-Z]:\\`)
	indices := re.FindAllStringIndex(input, -1)

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

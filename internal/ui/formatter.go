// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"fmt"
	"strconv"
	"time"
)

func formatInt(i int) string {
	s := strconv.FormatInt(int64(i), 10)
	return addCommas(s)
}

func addCommas(s string) string {
	n := len(s)
	if n <= 3 {
		return s
	}

	// Start from the end and insert commas every three characters
	var result string
	for i, j := n-1, 1; i >= 0; i, j = i-1, j+1 {
		result = string(s[i]) + result
		if j%3 == 0 && i != 0 {
			result = "," + result
		}
	}
	return result
}

func formatBytes(bytes uint64) string {
	const (
		KB = 1024
		MB = KB * 1024
		GB = MB * 1024
		TB = GB * 1024
	)

	switch {
	case bytes >= TB:
		return fmt.Sprintf("%.2f TB", float64(bytes)/TB)
	case bytes >= GB:
		return fmt.Sprintf("%.2f GB", float64(bytes)/GB)
	case bytes >= MB:
		return fmt.Sprintf("%.2f MB", float64(bytes)/MB)
	case bytes >= KB:
		return fmt.Sprintf("%.2f KB", float64(bytes)/KB)
	default:
		return fmt.Sprintf("%d B", bytes)
	}
}

func formatDateTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

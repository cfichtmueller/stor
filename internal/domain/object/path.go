// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package object

import "strings"

func SplitPath(key, delimiter string) []string {
	if key == "" {
		return []string{}
	}
	parts := strings.Split(strings.TrimSuffix(key, delimiter), delimiter)
	return parts
}

func JoinPath(parts []string, delimiter string) string {
	if len(parts) == 0 {
		return ""
	}
	return strings.Join(parts, delimiter) + delimiter
}

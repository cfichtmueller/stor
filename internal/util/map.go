// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package util

func MapMany[T any, R any](x []T, m func(e T) R) []R {
	res := make([]R, 0, len(x))
	for _, e := range x {
		res = append(res, m(e))
	}
	return res
}

// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/goparts/e"

func FormLabel(htmlFor, text string) e.Node {
	return e.Label(
		e.For(htmlFor),
		srOnly(),
		e.Text(text),
	)
}

// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "io"

func RenderNotFoundPage(w io.Writer) error {
	return renderTemplate(w, "NotFoundPage", nil)
}

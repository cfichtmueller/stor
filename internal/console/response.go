// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"io"

	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/ui"
)

func nodeResponse(n e.Node) *srv.Response {
	if n == nil {
		return srv.Respond().Html("")
	}
	return srv.Respond().BodyFn("text/html", func(w io.Writer) error {
		return n.Render(w)
	})
}

func nodeResponseWithShell(c *srv.Context, n e.Node) *srv.Response {
	includeShell := !c.HxBoosted()
	return srv.Respond().BodyFn("text/html", func(w io.Writer) error {
		if includeShell {
			return ui.Shell("", n).Render(w)
		}
		return n.Render(w)
	})
}

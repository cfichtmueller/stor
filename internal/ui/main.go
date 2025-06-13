// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"embed"
	"log/slog"

	"github.com/cfichtmueller/srv"
)

var (
	//go:embed css/*
	css embed.FS
	//go:embed js/*
	js embed.FS
	//go:embed img/*
	img embed.FS
)

func RenderCss(name string) *srv.Response {
	return renderFile(css, "css/"+name, "text/css")
}

func RenderJs(name string) *srv.Response {
	return renderFile(js, "js/"+name, "application/javascript")
}

func RenderImg(name string) *srv.Response {
	return renderFile(img, "img/"+name, "image/png")
}

func renderFile(fs embed.FS, name, contentType string) *srv.Response {
	b, err := fs.ReadFile(name)
	if err != nil {
		slog.Error("unable to write file", "name", name, "error", err)
		return srv.Respond().InternalServerError()
	}
	return srv.Respond().Body(contentType, b)
}

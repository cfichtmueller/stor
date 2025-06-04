// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"embed"
	"log/slog"

	"github.com/cfichtmueller/jug"
)

var (
	//go:embed css/*
	css embed.FS
	//go:embed js/*
	js embed.FS
	//go:embed img/*
	img embed.FS
)

func RenderCss(ctx jug.Context, name string) {
	renderFile(ctx, css, "css/"+name, "text/css")
}

func RenderJs(ctx jug.Context, name string) {
	renderFile(ctx, js, "js/"+name, "application/javascript")
}

func RenderImg(ctx jug.Context, name string) {
	renderFile(ctx, img, "img/"+name, "image/png")
}

func renderFile(ctx jug.Context, fs embed.FS, name, contentType string) {
	b, err := fs.ReadFile(name)
	if err != nil {
		slog.Error("unable to write file", "name", name, "error", err)
		ctx.RespondInternalServerError(nil)
		return
	}
	ctx.SetHeader("Content-Type", contentType)
	ctx.Writer().Write(b)
}

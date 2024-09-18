// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"log"

	"github.com/cfichtmueller/jug"
)

var (
	//go:embed css/*
	css embed.FS
	//go:embed js/*
	js embed.FS
	//go:embed img/*
	img embed.FS
	//go:embed html/*
	htmlFiles embed.FS
	templates = template.Must(template.New("").ParseFS(htmlFiles, "html/*.html"))
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
		log.Printf("unable to write file %s: %v", name, err)
		ctx.RespondInternalServerError(nil)
		return
	}
	ctx.SetHeader("Content-Type", contentType)
	ctx.Writer().Write(b)
}

func RenderShellStart(w io.Writer) error {
	return renderTemplate(w, "ShellStart", nil)
}

func RenderShellEnd(w io.Writer) error {
	return renderTemplate(w, "ShellEnd", nil)
}

func renderTemplate(w io.Writer, name string, data any) error {
	if err := templates.ExecuteTemplate(w, name, data); err != nil {
		return fmt.Errorf("unable to render template %s: %v", name, err)
	}
	return nil
}

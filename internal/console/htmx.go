// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"encoding/json"
	"fmt"

	"github.com/cfichtmueller/jug"
)

func hxRefresh(c jug.Context) {
	if isHxRequest(c) {
		c.SetHeader("HX-Refresh", "true")
	}
}

func hxRedirect(c jug.Context, path string) {
	if isHxRequest(c) {
		c.SetHeader("HX-Redirect", path)
		return
	}
	c.SetHeader("Location", path)
	c.Status(302)
}

func hxReswap(c jug.Context, target string) {
	c.SetHeader("HX-Reswap", target)
}

type toast struct {
	Title   string
	Message string
}

func newToast(title, message string, args ...any) toast {
	return toast{
		Title:   title,
		Message: fmt.Sprintf(message, args...),
	}
}

type hxTriggerModel struct {
	Event string
	Toast toast
}

func hxTrigger(c jug.Context, m hxTriggerModel) {
	d := make(map[string]any)
	if m.Event != "" {
		d[m.Event] = map[string]any{}
	}
	if m.Toast.Title != "" {
		d["toast"] = map[string]any{
			"title":   m.Toast.Title,
			"message": m.Toast.Message,
		}
	}

	b, err := json.Marshal(d)
	if err != nil {
		panic(fmt.Errorf("unable to marshal htmx trigger: %v", err))
	}
	c.SetHeader("HX-Trigger", string(b))
}

func isHxRequest(c jug.Context) bool {
	return c.GetHeader("Hx-Request") == "true"
}

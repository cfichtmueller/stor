// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"encoding/json"
	"fmt"

	"github.com/cfichtmueller/srv"
)

func hxRedirectFn(path string) srv.Handler {
	return func(c *srv.Context) *srv.Response {
		return hxRedirect(c, path)
	}
}

func hxRedirect(c *srv.Context, path string) *srv.Response {
	if c.HxRequest() {
		return srv.Respond().HxRedirect(path)
	}
	return srv.Respond().Found(path)
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

func hxTrigger(m hxTriggerModel) string {
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
		panic(fmt.Errorf("unable to marshal htmx trigger: %w", err))
	}
	return string(b)
}

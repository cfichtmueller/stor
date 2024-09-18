// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"bytes"
	"testing"
)

func TestPageHea(t *testing.T) {
	w := &bytes.Buffer{}
	if err := renderTemplate(w, "PageHeader", pageHeaderModel{
		Title: "My Title",
	}); err != nil {
		t.Error(err)
	}
	t.Log(w.String())
}

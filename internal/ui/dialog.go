// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

type dialogModel struct {
	ID           string
	Title        string
	Description  string
	CancelTitle  string
	ConfirmTitle string
	Destructive  bool
	HxDelete     string
	HxPost       string
	HxTarget     string
	HxSwap       string
}

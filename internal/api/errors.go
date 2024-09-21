// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"

	"github.com/cfichtmueller/jug"
)

var (
	ErrNoSuchKey = &Error{StatusCode: 404, Code: "NoSuchKey", Message: "The specified key does not exist"}
)

type Error struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func internalError(e error) *Error {
	return &Error{
		StatusCode: 500,
		Code:       "InternalError",
		Message:    e.Error(),
	}
}

func handleError(ctx jug.Context, err error) {
	e, ok := err.(*Error)
	if !ok {
		ctx.HandleError(err)
		ctx.Abort()
		return
	}

	b, err := json.Marshal(e)
	if err != nil {
		ctx.HandleError(err)
		ctx.Abort()
		return
	}

	ctx.Data(e.StatusCode, "application/json", b)
}

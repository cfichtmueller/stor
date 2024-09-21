// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ec

var (
	NoSuchKey = &Error{StatusCode: 404, Code: "NoSuchKey", Message: "The specified key does not exist"}
)

type Error struct {
	StatusCode int    `json:"-"`
	Code       string `json:"code"`
	Message    string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func Internal(e error) *Error {
	return &Error{
		StatusCode: 500,
		Code:       "InternalError",
		Message:    e.Error(),
	}
}

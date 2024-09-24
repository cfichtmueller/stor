// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ec

var (
	AccountDisabled     = &Error{StatusCode: 401, Code: "AccountDisabled", Message: "The user account is disabled"}
	ArchiveNotPending   = &Error{StatusCode: 409, Code: "ArchiveNotPending", Message: "The archive is not pending"}
	BucketAlreadyExists = &Error{StatusCode: 409, Code: "BucketAlreadyExists", Message: "The requested bucket name is not available"}
	BucketNotEmpty      = &Error{StatusCode: 409, Code: "BucketNotEmpty", Message: "The bucket is not empty"}
	InvalidArgument     = &Error{StatusCode: 400, Code: "InvalidArgument", Message: "Invalid argument"}
	InvalidCredentials  = &Error{StatusCode: 401, Code: "InvalidCredentials", Message: "Invalid Credentials"}
	NoSuchArchive       = &Error{StatusCode: 404, Code: "NoSuchArchive", Message: "The specified archive does not exist"}
	NoSuchApiKey        = &Error{StatusCode: 404, Code: "NoSuchApiKey", Message: "The specified api key does not exist"}
	NoSuchBucket        = &Error{StatusCode: 404, Code: "NoSuchBucket", Message: "The specified bucket does not exist"}
	NoSuchKey           = &Error{StatusCode: 404, Code: "NoSuchKey", Message: "The specified key does not exist"}
	NoSuchUser          = &Error{StatusCode: 404, Code: "NoSuchUser", Message: "The specified user does not exist"}
	Unauthorized        = &Error{StatusCode: 401, Code: "Unauthorized", Message: "Unauthorized"}
	UserAlreadyExists   = &Error{StatusCode: 409, Code: "UserAlreadyExists", Message: "The requested user name is not available"}
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

func Wrap(e error) *Error {
	if ee, ok := e.(*Error); ok {
		return ee
	}
	return Internal(e)
}

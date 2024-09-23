// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package domain

import "github.com/cfichtmueller/stor/internal/util"

const idAlphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const etagAlphabet = "abcdef0123456789"

func RandomId() string {
	return util.RandomStringFromAlphabet(idAlphabet, 10)
}

func NewId(length int) string {
	return util.RandomStringFromAlphabet(idAlphabet, length)
}

func NewEtag() string {
	return util.RandomStringFromAlphabet(etagAlphabet, 64)
}

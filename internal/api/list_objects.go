// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/util"
)

type ListObjectsResponse struct {
	IsTruncated bool             `json:"isTruncated"`
	Contents    []ObjectResponse `json:"contents"`
	Name        string           `json:"name"`
	MaxKeys     int              `json:"maxKeys"`
	KeyCount    int              `json:"keyCount"`
	StartAfter  *string          `json:"startAfter,omitempty"`
}

func handleListObjects(c jug.Context) {
	startAfter := c.Query("start-after")
	maxKeys, err := c.DefaultIntQuery("max-keys", 1000)
	if err != nil {
		c.HandleError(err)
		return
	}
	if maxKeys > 1000 {
		maxKeys = 1000
	}
	b := contextGetBucket(c)
	contents, err := object.List(c, b.Name, startAfter, maxKeys)
	if err != nil {
		c.HandleError(err)
		return
	}
	keyCount := len(contents)
	totalKeys, err := object.Count(c, b.Name, startAfter)
	if err != nil {
		c.HandleError(err)
		return
	}
	var startAfterRes *string
	if startAfter != "" {
		startAfterRes = &startAfter
	}

	c.RespondOk(ListObjectsResponse{
		IsTruncated: totalKeys > keyCount,
		Contents:    util.MapMany(contents, newObjectResponse),
		Name:        b.Name,
		MaxKeys:     maxKeys,
		KeyCount:    keyCount,
		StartAfter:  startAfterRes,
	})
}

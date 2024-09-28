// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/uc"
	"github.com/cfichtmueller/stor/internal/util"
)

type ListObjectsResponse struct {
	IsTruncated    bool             `json:"isTruncated"`
	Objects        []ObjectResponse `json:"objects"`
	Name           string           `json:"name"`
	MaxKeys        int              `json:"maxKeys"`
	KeyCount       int              `json:"keyCount"`
	StartAfter     *string          `json:"startAfter,omitempty"`
	CommonPrefixes []string         `json:"commonPrefixes,omitempty"`
}

func handleListObjects(c jug.Context) {
	startAfter := c.Query("start-after")
	maxKeys, err := c.DefaultIntQuery("max-keys", 1000)
	if err != nil {
		handleError(c, err)
		return
	}
	if maxKeys > 1000 {
		maxKeys = 1000
	}
	delimiter := c.Query("delimiter")
	prefix := c.Query("prefix")
	b := contextGetBucket(c)

	if delimiter != "" {
		r, err := uc.ObjectPrefixSearch(c, b, delimiter, prefix, startAfter, maxKeys)
		if err != nil {
			handleError(c, err)
			return
		}
		c.RespondOk(ListObjectsResponse{
			IsTruncated:    r.IsTruncated,
			Objects:        util.MapMany(r.Objects, newObjectResponse),
			Name:           b.Name,
			MaxKeys:        maxKeys,
			KeyCount:       len(r.Objects),
			StartAfter:     &startAfter,
			CommonPrefixes: r.CommonPrefixes,
		})
		return
	}

	contents, err := object.List(c, b.Name, startAfter, maxKeys)
	if err != nil {
		handleError(c, err)
		return
	}
	keyCount := len(contents)
	totalKeys, err := object.Count(c, b.Name, startAfter)
	if err != nil {
		handleError(c, err)
		return
	}
	var startAfterRes *string
	if startAfter != "" {
		startAfterRes = &startAfter
	}

	c.RespondOk(ListObjectsResponse{
		IsTruncated: totalKeys > keyCount,
		Objects:     util.MapMany(contents, newObjectResponse),
		Name:        b.Name,
		MaxKeys:     maxKeys,
		KeyCount:    keyCount,
		StartAfter:  startAfterRes,
	})
}

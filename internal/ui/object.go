// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import "github.com/cfichtmueller/stor/internal/domain/object"

type objectModel struct {
	ID     string
	Bucket string
	Key    string
	Size   string
}

func newObjectModel(o *object.Object) objectModel {
	return objectModel{
		ID:     o.ID,
		Bucket: o.Bucket,
		Key:    o.Key,
		Size:   formatBytes(o.Size),
	}
}

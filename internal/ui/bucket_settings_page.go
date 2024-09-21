// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"io"

	"github.com/cfichtmueller/stor/internal/domain/bucket"
)

func RenderBucketSettingsPage(w io.Writer, b *bucket.Bucket) error {
	return renderTemplate(w, "BucketSettingsPage", newBucketPageModel(b, "settings"))
}

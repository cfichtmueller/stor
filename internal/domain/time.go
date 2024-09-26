// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package domain

import "time"

func TimeNow() time.Time {
	return time.Now().UTC().Truncate(time.Millisecond)
}

// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package disk

type Info struct {
	Total uint64
	Free  uint64
	Used  uint64
	Files uint64
	Ffree uint64
}

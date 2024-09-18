// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

//go:build linux && !s390x && !arm && !386
// +build linux,!s390x,!arm,!386

package disk

import (
	"syscall"
)

func GetInfo(path string) (Info, error) {
	s := syscall.Statfs_t{}
	if err := syscall.Statfs(path, &s); err != nil {
		return Info{}, err
	}

	reservedBlocks := s.Bfree - s.Bavail
	info := Info{
		Total: uint64(s.Bsize) * (s.Blocks - reservedBlocks),
		Free:  uint64(s.Bsize) * s.Bavail,
		Files: s.Files,
		Ffree: s.Ffree,
	}
	info.Used = info.Total - info.Free
	return info, nil
}

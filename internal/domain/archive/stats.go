// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package archive

import (
	"fmt"
	"time"
)

type Stats struct {
	start time.Time
	files int64
	bytes int64
}

func NewStats() *Stats {
	return &Stats{
		start: time.Now(),
	}
}

func (s *Stats) AddFiles(n int64) {
	s.files += n
}

func (s *Stats) AddBytes(n int64) {
	s.bytes += n
}

func (s *Stats) Summary() string {
	dur := time.Since(s.start)
	kb := s.bytes / 1024
	return fmt.Sprintf("%dms, %dkb, %.2fkb/s", dur.Milliseconds(), kb, float64(kb)/dur.Seconds())
}

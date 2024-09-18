// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"os"
	"path"
)

var (
	DataDir string
)

func Mkdir(name string) error {
	if err := os.Mkdir(path.Join(DataDir, name), 0700); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	return nil
}

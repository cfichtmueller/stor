// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package config

import (
	"os"
	"path"
)

var (
	DataDir      string
	ApiHost      string
	ApiPort      string
	ConsoleHost  string
	ConsolePort  string
	TrustProxies bool
)

func init() {
	DataDir = getEnv("DATA_DIR", "/var/stor")
	ApiHost = os.Getenv("API_HOST")
	ApiPort = getEnv("API_PORT", "8000")
	ConsoleHost = os.Getenv("CONSOLE_HOST")
	ConsolePort = getEnv("CONSOLE_PORT", "8001")
	TrustProxies = getEnv("TRUST_PROXIES", "false") == "true"
}

func Mkdir(name string) error {
	if err := os.Mkdir(path.Join(DataDir, name), 0700); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}
	return nil
}

func getEnv(key, fallback string) string {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	return v
}

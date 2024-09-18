// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/cfichtmueller/stor/internal/config"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "stor",
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&config.DataDir, "data", "d", "/var/stor", "location of the data directory")
	rootCmd.AddCommand(serveCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

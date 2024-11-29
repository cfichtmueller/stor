// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "stor",
}

func init() {
	rootCmd.AddCommand(serveCmd, checkCmd)
}

func Execute() error {
	return rootCmd.Execute()
}

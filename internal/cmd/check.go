// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"os"

	"github.com/cfichtmueller/stor/internal/shell"
	"github.com/spf13/cobra"
)

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Perform a system check",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Perform system check...\n")
		if !shell.Check() {
			os.Exit(1)
			return
		}
		fmt.Printf("Done\n")
	},
}

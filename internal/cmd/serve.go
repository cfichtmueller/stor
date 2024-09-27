// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cmd

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/cfichtmueller/stor/internal/api"
	"github.com/cfichtmueller/stor/internal/config"
	"github.com/cfichtmueller/stor/internal/console"
	"github.com/cfichtmueller/stor/internal/shell"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	g errgroup.Group
)

var serveCmd = &cobra.Command{
	Use: "serve",
	Run: func(cmd *cobra.Command, args []string) {
		shell.Configure()
		apiEngine := api.Configure()
		consoleEngine := console.Configure()

		apiAddr := config.ApiHost + ":" + config.ApiPort
		consoleAddr := config.ConsoleHost + ":" + config.ConsolePort

		engineServer := newServer(apiAddr, apiEngine)

		consoleServer := newServer(consoleAddr, consoleEngine)
		consoleServer.WriteTimeout = 10 * time.Second

		g.Go(func() error {
			return engineServer.ListenAndServe()
		})

		g.Go(func() error {
			return consoleServer.ListenAndServe()
		})

		fmt.Printf("Starting API on %s\n", apiAddr)
		fmt.Printf("Starting Console on %s\n", consoleAddr)

		if err := g.Wait(); err != nil {
			log.Fatal(err)
		}
	},
}

func newServer(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:        addr,
		Handler:     handler,
		ReadTimeout: 30 * time.Second,
	}
}

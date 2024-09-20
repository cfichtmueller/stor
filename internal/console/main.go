// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"fmt"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/ui"
)

func Configure() jug.Engine {
	console := jug.Default()

	console.GET("/css/style.css", func(c jug.Context) {
		ui.RenderCss(c, "style.css")
	})

	console.GET("/js/lib.js", func(c jug.Context) {
		ui.RenderJs(c, "lib.js")
	})

	console.GET("/js/htmx.min.js", func(c jug.Context) {
		ui.RenderJs(c, "htmx.min.js")
	})

	console.GET("/img/icon.png", func(c jug.Context) {
		c.SetHeader("Cache-Control", "max-age=31536000")
		ui.RenderImg(c, "icon.png")
	})

	console.GET("/img/empty.png", func(c jug.Context) {
		ui.RenderImg(c, "empty.png")
	})

	console.GET("/img/bucket-full.png", func(c jug.Context) {
		ui.RenderImg(c, "bucket-full.png")
	})

	console.GET("", handleHomePage)
	console.GET("/bootstrap", handleBootstrapPage)
	console.POST("/bootstrap", handleBootstrap)
	console.GET("/login", handleLoginPage)
	console.POST("/login", handleLogin)

	// c is for components
	componentsGroup := console.Group("/c", authenticatedFilter, requireHxRequest)

	componentsGroup.GET("/api-key-sheet", apiKeyFilter, handleRenderApiKeySheet)
	componentsGroup.GET("/api-key-delete-dialog", apiKeyFilter, handleRenderDeleteApiKeyDialog)
	componentsGroup.GET("/api-keys-table", handleRenderApiKeysTable)
	componentsGroup.GET("/buckets-table", handleRenderBucketsTable)
	componentsGroup.GET("/create-api-key-dialog", uiRenderFn("create api key dialog", ui.RenderCreateApiKeyDialog))
	componentsGroup.GET("/create-bucket-dialog", uiRenderFn("create bucket dialog", ui.RenderCreateBucketDialog))
	componentsGroup.GET("/dashboard-metrics", handleRenderDashboardMetrics)
	// /c/objects-table

	// r is for rpc
	r := console.Group("/r", authenticatedFilter, requireHxRequest)
	r.POST("/api-key", handleRpcCreateApiKey)
	r.DELETE("/api-key", apiKeyFilter, handleRpcDeleteApiKey)
	r.POST("/bucket", handleRpcCreateBucket)

	// DELETE /r/bucket
	// POST /r/invite-user
	// DELETE /r/user
	// PUT /r/profile

	// u is for user pages
	uGroup := console.Group("/u", authenticatedFilter, renderShell)
	uGroup.GET("", handleDashboardPage)

	uGroup.GET("/buckets", handleBucketsPage)

	uBucketGroup := uGroup.Group("/buckets/:bucketName", bucketFilter)
	uBucketGroup.GET("", handleBucketPage)
	uBucketGroup.GET("/objects", handleBucketObjectsPage)
	uBucketGroup.GET("/settings", handleBucketSettingsPage)

	uAdminGroup := uGroup.Group("/admin")
	uAdminGroup.GET("", handleAdminPage)
	uAdminGroup.GET("/api-keys", handleApiKeysPage)
	uAdminGroup.GET("/users", handleUsersPage)

	return console
}

func redirect(c jug.Context, to string) {
	c.SetHeader("Location", to)
	c.Status(301)
}

func requireHxRequest(c jug.Context) {
	hx := c.GetHeader("HX-Request")
	if hx == "true" {
		c.Next()
		return
	}
	c.Status(301)
	c.SetHeader("Location", "/u")
	c.Abort()
}

func must(what string, c jug.Context, err error) bool {
	if err == nil {
		return true
	}
	if _, ok := err.(*jug.ResponseStatusError); !ok {
		fmt.Printf("unable to %s: %v", what, err)
	}
	c.HandleError(err)
	return false
}

func renderShell(c jug.Context) {
	includeShell := c.GetHeader("Hx-Boosted") != "true"
	if includeShell {
		if !must("render shell start", c, ui.RenderShellStart(c.Writer())) {
			c.Abort()
			return
		}
	}
	c.Next()
	if includeShell {
		if !must("render shell end", c, ui.RenderShellEnd(c.Writer())) {
			c.Abort()
			return
		}
	}
}

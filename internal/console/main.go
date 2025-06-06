// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"fmt"

	"github.com/cfichtmueller/goparts/e"
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
	console.GET("/bootstrap", RequireNotBootstrapped, renderNode(withShell(handleBootstrapPage)))
	console.POST("/bootstrap", RequireNotBootstrapped, renderNode(handleBootstrap))
	console.GET("/login", renderNode(withShell(handleLoginPage)))
	console.POST("/login", renderNode(handleLogin))

	// c is for components
	componentsGroup := console.Group("/c", authenticatedFilter, requireHxRequest)

	componentsGroup.GET("/api-key-sheet", apiKeyFilter, renderNode(handleRenderApiKeySheet))
	componentsGroup.GET("/api-key-delete-dialog", apiKeyFilter, renderNode(handleRenderDeleteApiKeyDialog))
	componentsGroup.GET("/api-keys-table", renderNode(handleRenderApiKeysTable))
	componentsGroup.GET("/buckets-table", renderNode(handleRenderBucketsTable))
	componentsGroup.GET("/create-api-key-dialog", renderNodeFn(ui.CreateApiKeyDialog))
	componentsGroup.GET("/create-bucket-dialog", renderNodeFn(ui.CreateBucketDialog))
	componentsGroup.GET("/dashboard-metrics", renderNode(handleRenderDashboardMetrics))
	// /c/objects-table

	// r is for rpc
	r := console.Group("/r", authenticatedFilter, requireHxRequest)
	r.POST("/api-key", renderNode(handleRpcCreateApiKey))
	r.DELETE("/api-key", apiKeyFilter, handleRpcDeleteApiKey)
	r.POST("/bucket", handleRpcCreateBucket)

	console.GET("/open", authenticatedFilter, handleRpcOpenObject)
	console.GET("/download", authenticatedFilter, handleRpcDownloadObject)

	// DELETE /r/bucket
	// POST /r/invite-user
	// DELETE /r/user
	// PUT /r/profile

	// u is for user pages
	uGroup := console.Group("/u", authenticatedFilter)
	uGroup.GET("", renderNode(withShell(handleDashboardPage)))

	uGroup.GET("/buckets", renderNode(withShell(handleBucketsPage)))

	uBucketGroup := uGroup.Group("/buckets/:bucketName", bucketFilter)
	uBucketGroup.GET("", handleBucketPage)
	uBucketGroup.GET("/objects", renderNode(withShell(handleBucketObjectsPage)))
	uBucketGroup.GET("/object", renderNode(withShell(handleObjectPage)))
	uBucketGroup.GET("/properties", renderNode(withShell(handleBucketPropertiesPage)))
	uBucketGroup.GET("/settings", renderNode(withShell(handleBucketSettingsPage)))

	uAdminGroup := uGroup.Group("/admin")
	uAdminGroup.GET("", hxRedirectFn("/u/admin/users"))
	uAdminGroup.GET("/api-keys", renderNode(withShell(handleApiKeysPage)))
	uAdminGroup.GET("/users", renderNode(withShell(handleUsersPage)))

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

func withShell(h NodeHandler) NodeHandler {
	return func(c jug.Context) (e.Node, error) {
		includeShell := c.GetHeader("Hx-Boosted") != "true"
		n, err := h(c)
		if err != nil {
			return nil, err
		}
		if n == nil {
			return nil, nil
		}
		if includeShell {
			return ui.Shell("", n), nil
		}
		return n, nil
	}
}

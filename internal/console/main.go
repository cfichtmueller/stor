// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/config"
	"github.com/cfichtmueller/stor/internal/ui"
)

func Configure() *srv.Server {
	console := srv.NewServer().SetTrustRemoteIdHeaders(config.TrustProxies).Use(srv.LoggingMiddleware())

	console.GET("/css/style.css", func(c *srv.Context) *srv.Response {
		return ui.RenderCss("style.css")
	})

	console.GET("/js/lib.js", func(c *srv.Context) *srv.Response {
		return ui.RenderJs("lib.js")
	})

	console.GET("/js/htmx.min.js", func(c *srv.Context) *srv.Response {
		return ui.RenderJs("htmx.min.js")
	})

	console.GET("/img/icon.png", func(c *srv.Context) *srv.Response {
		return ui.RenderImg("icon.png").Header("Cache-Control", "max-age=31536000")
	})

	console.GET("/img/empty.png", func(c *srv.Context) *srv.Response {
		return ui.RenderImg("empty.png")
	})

	console.GET("/img/bucket-full.png", func(c *srv.Context) *srv.Response {
		return ui.RenderImg("bucket-full.png")
	})

	console.GET("", handleHomePage)
	console.GET("/bootstrap", handleBootstrapPage, RequireNotBootstrapped)
	console.POST("/bootstrap", handleBootstrap, RequireNotBootstrapped)
	console.GET("/login", handleLoginPage)
	console.POST("/login", handleLogin)

	// c is for components
	componentsGroup := console.Group("/c", authenticatedFilter, requireHxRequest)

	componentsGroup.GET("/api-key-sheet", renderNode(handleRenderApiKeySheet), apiKeyFilter)
	componentsGroup.GET("/api-key-delete-dialog", renderNode(handleRenderDeleteApiKeyDialog), apiKeyFilter)
	componentsGroup.GET("/api-keys-table", renderNode(handleRenderApiKeysTable))
	componentsGroup.GET("/buckets-table", renderNode(handleRenderBucketsTable))
	componentsGroup.GET("/create-api-key-dialog", renderNodeFn(ui.CreateApiKeyDialog))
	componentsGroup.GET("/create-bucket-dialog", renderNodeFn(ui.CreateBucketDialog))
	componentsGroup.GET("/delete-bucket-dialog", renderNode(handleRenderDeleteBucketDialog), withBucketFromQuery)
	componentsGroup.GET("/empty-bucket-dialog", renderNode(handleRenderEmptyBucketDialog), withBucketFromQuery)
	componentsGroup.GET("/dashboard-metrics", renderNode(handleRenderDashboardMetrics))
	// /c/objects-table

	// r is for rpc
	r := console.Group("/r", authenticatedFilter, requireHxRequest)
	r.POST("/api-key", handleRpcCreateApiKey)
	r.DELETE("/api-key", handleRpcDeleteApiKey, apiKeyFilter)
	r.POST("/bucket", handleRpcCreateBucket)
	r.DELETE("/bucket", handleRpcDeleteBucket, withBucketFromQuery)
	r.POST("/change-password", handleRpcChangePassword)
	r.POST("/logout-session", handleRpcLogoutSession)
	r.POST("/empty-bucket", handleRpcEmptyBucket, withBucketFromQuery)

	console.GET("/open", handleRpcOpenObject, authenticatedFilter)
	console.GET("/download", handleRpcDownloadObject, authenticatedFilter)

	// DELETE /r/bucket
	// POST /r/invite-user
	// DELETE /r/user
	// PUT /r/profile

	// u is for user pages
	uGroup := console.Group("/u", authenticatedFilter)
	uGroup.GET("", handleDashboardPage)

	uGroup.GET("/buckets", handleBucketsPage)

	uBucketGroup := uGroup.Group("/buckets/{bucketName}", bucketFilter)
	uBucketGroup.GET("", handleBucketPage)
	uBucketGroup.GET("/objects", handleBucketObjectsPage)
	uBucketGroup.GET("/object", handleObjectPage)
	uBucketGroup.GET("/properties", handleBucketPropertiesPage)
	uBucketGroup.GET("/settings", handleBucketSettingsPage)

	uAdminGroup := uGroup.Group("/admin")
	uAdminGroup.GET("", hxRedirectFn("/u/admin/users"))
	uAdminGroup.GET("/api-keys", handleApiKeysPage)
	uAdminGroup.GET("/users", handleUsersPage)

	uProfileGroup := uGroup.Group("/profile")
	uProfileGroup.GET("", handleProfilePage)

	return console
}

func requireHxRequest(c *srv.Context, next srv.Handler) *srv.Response {
	if c.HxRequest() {
		return next(c)
	}
	return srv.Respond().MovedPermanently("/u")
}

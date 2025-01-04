// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
)

type ApiKeysPageData struct {
	Keys []*apikey.ApiKey
}

func ApiKeysPage(data *ApiKeysPageData) e.Node {
	hasApiKeys := len(data.Keys) > 0
	return AdminPageLayout(admin_tab_active_api_keys,
		e.Div(
			e.HXTrigger("apiKeysUpdated from:body"),
			e.HXGet("/c/api-keys-table"),
			e.Iff(hasApiKeys, e.F(ApiKeysTable, data.Keys)),
			e.Iff(!hasApiKeys, ApiKeysEmptyState),
		),
	)
}

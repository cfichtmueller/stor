// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"io"

	"github.com/cfichtmueller/stor/internal/domain/apikey"
	"github.com/cfichtmueller/stor/internal/util"
)

type ApiKeysPageData struct {
	Keys []*apikey.ApiKey
}

type apiKeysPageModel struct {
	Layout adminPageModel
	Keys   []apiKeyModel
}

func RenderApiKeysPage(w io.Writer, data ApiKeysPageData) error {
	return renderTemplate(w, "ApiKeysPage", apiKeysPageModel{
		Layout: newAdminPageModel(admin_tab_active_api_keys),
		Keys:   util.MapMany(data.Keys, newApiKeyModel),
	})
}

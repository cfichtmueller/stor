// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package ui

import (
	"io"

	"github.com/cfichtmueller/stor/internal/domain/apikey"
)

type apiKeyModel struct {
	ID          string
	Prefix      string
	Description string
	CreatedAt   string
	CreatedBy   string
	ExpiresAt   string
}

func newApiKeyModel(k *apikey.ApiKey) apiKeyModel {
	return apiKeyModel{
		ID:          k.ID,
		Prefix:      k.Prefix,
		Description: k.Description,
		CreatedAt:   formatDateTime(k.CreatedAt),
		CreatedBy:   k.CreatedBy,
		ExpiresAt:   formatDateTime(k.ExpiresAt),
	}
}

type apiKeyCreatedData struct {
	Key   apiKeyModel
	Plain string
}

func RenderApiKeyCreatedDialog(w io.Writer, k *apikey.ApiKey, plain string) error {
	return renderTemplate(w, "ApiKeyCreatedDialog", apiKeyCreatedData{
		Key:   newApiKeyModel(k),
		Plain: plain,
	})
}

func RenderApiKeysEmptyState(w io.Writer) error {
	return renderTemplate(w, "ApiKeysEmptyState", nil)
}

type apiKeySheetModel struct {
	Sheet sheetModel
	Key   apiKeyModel
}

func RenderApiKeySheet(w io.Writer, key *apikey.ApiKey) error {
	return renderTemplate(w, "ApiKeySheet", apiKeySheetModel{
		Sheet: sheetModel{
			Title: "API Key",
			Lead:  "Make changes to the API key here.",
		},
		Key: newApiKeyModel(key),
	})
}

func RenderApiKeysTable(w io.Writer, keys []*apikey.ApiKey) error {
	return renderTemplate(w, "ApiKeysTable", map[string]any{
		"Keys": keys,
	})
}

func RenderCreateApiKeyDialog(w io.Writer) error {
	return renderTemplate(w, "CreateApiKeyDialog", nil)
}

func RenderDeleteApiKeyDialog(w io.Writer, key *apikey.ApiKey) error {
	return renderTemplate(w, "ConfirmationDialog", dialogModel{
		ID:           "delete-api-key-dialog",
		Title:        "Delete API Key",
		Description:  "U sure u want to delete key " + key.Description + "?",
		ConfirmTitle: "Delete",
		CancelTitle:  "Cancel",
		Destructive:  true,
		HxDelete:     "/r/api-key?key=" + key.ID,
	})
}

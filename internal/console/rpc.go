// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"time"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
	"github.com/cfichtmueller/stor/internal/uc"
	"github.com/cfichtmueller/stor/internal/ui"
)

//
// API Key
//

func handleRpcCreateApiKey(c jug.Context) {
	var description string
	if err := bindFormData(c, "description", &description); err != nil {
		c.HandleError(err)
		return
	}

	//TODO: get PRINCIPAL
	key, plain, err := apikey.Create(c, "aaabbbcc", apikey.CreateCommand{
		Description: description,
		TTL:         time.Hour * 24 * 360,
	})

	if err != nil {
		//TODO: give actual feedback
		c.HandleError(err)
		return
	}

	hxTrigger(c, hxTriggerModel{
		Event: "apiKeysUpdated",
		Toast: toast{
			Title:   "Success",
			Message: "API KEY " + key.Description + " created",
		},
	})
	hxReswap(c, "outerHTML")

	must("render api key created dialog", c, ui.RenderApiKeyCreatedDialog(c.Writer(), key, plain))
}

func handleRpcDeleteApiKey(c jug.Context) {
	key := contextGetApiKey(c)

	if err := apikey.Delete(c, key.ID); err != nil {
		c.HandleError(err)
		return
	}
	hxRefresh(c)
	hxTrigger(c, hxTriggerModel{
		Toast: newToast("Success", "API key deleted"),
	})
}

//
// Bucket
//

func handleRpcCreateBucket(c jug.Context) {
	if !must("parse form", c, c.Request().ParseForm()) {
		return
	}
	values := c.Request().Form
	name := values.Get("name")

	if _, err := uc.CreateBucket(c, name); err != nil {
		hxTrigger(c, hxTriggerModel{
			Toast: newToast("Error", "Failed to create bucket: %v", err),
		})
		return
	}

	hxTrigger(c, hxTriggerModel{
		Event: "bucketsUpdated",
		Toast: newToast("Success", "Bucket %s created", name),
	})
}

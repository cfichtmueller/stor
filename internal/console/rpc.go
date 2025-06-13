// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain/apikey"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/uc"
	"github.com/cfichtmueller/stor/internal/ui"
)

//
// API Key
//

func handleRpcCreateApiKey(c *srv.Context) *srv.Response {
	principal := contextMustGetPrincipal(c)
	var description string
	if err := bindFormData(c, "description", &description); err != nil {
		return responseFromError(err)
	}

	key, plain, err := apikey.Create(c, principal, apikey.CreateCommand{
		Description: description,
		TTL:         time.Hour * 24 * 360,
	})

	if err != nil {
		//TODO: give actual feedback
		return responseFromError(err)
	}

	return nodeResponse(ui.ApiKeyCreatedDialog(key, plain)).
		HxTrigger(hxTrigger(hxTriggerModel{
			Event: "apiKeysUpdated",
			Toast: newToast("Success", "API KEY %s created", key.Description),
		})).
		HxReswap("outerHTML")
}

func handleRpcDeleteApiKey(c *srv.Context) *srv.Response {
	key := contextGetApiKey(c)

	if err := apikey.Delete(c, key.ID); err != nil {
		return responseFromError(err)
	}
	return srv.Respond().
		HxRefresh().
		HxTrigger(hxTrigger(hxTriggerModel{
			Toast: newToast("Success", "API key deleted"),
		}))
}

//
// Bucket
//

func handleRpcCreateBucket(c *srv.Context) *srv.Response {
	values := c.FormValues()
	name := values.Get("name")

	if _, err := uc.CreateBucket(c, name); err != nil {
		return srv.Respond().
			HxTrigger(hxTrigger(hxTriggerModel{
				Toast: newToast("Error", "Failed to create bucket: %v", err),
			}))
	}

	return srv.Respond().
		HxTrigger(hxTrigger(hxTriggerModel{
			Event: "bucketsUpdated",
			Toast: newToast("Success", "Bucket %s created", name),
		}))
}

//
// Object
//

func handleRpcOpenObject(c *srv.Context) *srv.Response {
	bucketName := c.Query("bucket")
	key, r := c.StringQuery("key")
	if r != nil {
		return r
	}

	o, err := object.FindOne(c, bucketName, key, false)
	if err != nil {
		return responseFromError(err)
	}

	return srv.Respond().
		Header("Content-Disposition", "inline").
		BodyFn(o.ContentType, func(w io.Writer) error {
			return object.Write(c, o, w)
		})
}

func handleRpcDownloadObject(c *srv.Context) *srv.Response {
	bucketName := c.Query("bucket")
	key, r := c.StringQuery("key")
	if r != nil {
		return r
	}

	o, err := object.FindOne(c, bucketName, key, false)
	if err != nil {
		return responseFromError(err)
	}

	return srv.Respond().
		Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filepath.Base(o.Key))).
		BodyFn(o.ContentType, func(w io.Writer) error {
			return object.Write(c, o, w)
		})
}

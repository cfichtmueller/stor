// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"log"
	"strconv"
	"time"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/object"
	"github.com/cfichtmueller/stor/internal/uc"
	"github.com/cfichtmueller/stor/internal/util"
)

type ObjectResponse struct {
	Key         string    `json:"key"`
	ContentType string    `json:"contentType"`
	Size        uint64    `json:"size"`
	CreatedAt   time.Time `json:"createdAt"`
}

func newObjectResponse(o *object.Object) ObjectResponse {
	return ObjectResponse{
		Key:         o.Key,
		ContentType: o.ContentType,
		Size:        o.Size,
		CreatedAt:   o.CreatedAt,
	}
}

func handleListObjects(c jug.Context) {
	b := contextGetBucket(c)
	o, err := object.List(c, b.Name)
	if err != nil {
		log.Printf("Error: %v", err)
		c.HandleError(err)
		return
	}
	c.RespondOk(util.MapMany(o, newObjectResponse))
}

func handleGetObject(c jug.Context) {
	o := contextGetObject(c)
	c.Status(200)
	c.SetHeader("Content-Length", strconv.FormatInt(int64(o.Size), 10))
	c.SetHeader("Content-Type", o.ContentType)

	if err := object.Write(c, o, c.Writer()); err != nil {
		log.Printf("unable to write object: %v", err)
	}
}

func handleCreateObject(c jug.Context) {
	b := contextGetBucket(c)
	key := contextGetObjectKey(c)

	contentType := c.Request().Header.Get("Content-Type")
	if contentType == "" {
		contentType = "application/octet-stream"
	}
	d, err := c.GetRawData()
	if err != nil {
		c.HandleError(err)
		return
	}

	if err := uc.CreateObject(c, b, object.CreateCommand{
		Key:         key,
		ContentType: contentType,
		Data:        d,
	}); err != nil {
		log.Printf("unable to create object: %v", err)
		c.HandleError(err)
		return
	}

	c.RespondNoContent()
}

func handleDeleteObject(c jug.Context) {
	b := contextGetBucket(c)
	o := contextGetObject(c)

	if err := uc.DeleteObject(c, b, o); err != nil {
		c.HandleError(err)
		return
	}

	c.RespondNoContent()
}

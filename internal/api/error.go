// Copyright 2025 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package api

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/ec"
)

func responseFromError(err error) *srv.Response {
	if ve, ok := err.(*srv.ValidationError); ok {
		ve.Code = ec.InvalidArgument.Code
		return srv.Respond().Status(http.StatusBadRequest).Json(ve)
	}

	e, ok := err.(*ec.Error)
	if !ok {
		return srv.Respond().InternalServerError(err)
	}

	b, merr := json.Marshal(e)
	if merr != nil {
		slog.Error("unable to marshal error", "error", merr)
		return srv.Respond().Error(err)
	}

	return srv.Respond().Status(e.StatusCode).Body("application/json", b)
}

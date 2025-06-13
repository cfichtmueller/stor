// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"errors"

	"github.com/cfichtmueller/srv"
	"github.com/cfichtmueller/stor/internal/domain/session"
	"github.com/cfichtmueller/stor/internal/ec"
	"github.com/cfichtmueller/stor/internal/uc"
	"github.com/cfichtmueller/stor/internal/ui"
)

func handleBootstrapPage(c *srv.Context) *srv.Response {
	return nodeResponseWithShell(c, ui.BootstrapPage())
}

func RequireNotBootstrapped(c *srv.Context, next srv.Handler) *srv.Response {
	bootstrapped, err := uc.IsBootstrapped(c)
	if err != nil {
		return responseFromError(err)
	}
	if bootstrapped {
		return srv.Respond().MovedPermanently("/")
	}
	return next(c)
}

func handleBootstrap(c *srv.Context) *srv.Response {
	cmd := uc.BootstrapCommand{}
	if err := bindFormData(c, "email", &cmd.Email, "password", &cmd.Password, "passwordConfirmation", &cmd.PasswordConfirmation); err != nil {
		return responseFromError(err)
	}
	if err := cmd.Validate(); err != nil {
		return nodeResponse(ui.BootstrapForm(&ui.BootstrapFormData{Email: cmd.Email, ErrorMessage: err.Error()}))
	}
	if err := uc.Bootstrap(c, cmd); err != nil {
		return nodeResponse(ui.BootstrapForm(&ui.BootstrapFormData{Email: cmd.Email, ErrorMessage: err.Error()}))
	}

	return srv.Respond().HxRedirect("/login")
}

func handleLoginPage(c *srv.Context) *srv.Response {
	return nodeResponseWithShell(c, ui.LoginPage())
}

func handleLogin(c *srv.Context) *srv.Response {
	cmd := uc.LoginCommand{
		IpAddress: c.ClientIP(),
	}
	if err := bindFormData(c, "email", &cmd.Email, "password", &cmd.Password); err != nil {
		return responseFromError(err)
	}
	_, sid, err := uc.Login(c, cmd)
	if err != nil {
		if errors.Is(err, ec.InvalidCredentials) {
			return nodeResponse(ui.LoginForm(ui.LoginFormData{
				Email:        cmd.Email,
				ErrorMessage: "Invalid Credentials",
			}))
		}
		if errors.Is(err, ec.AccountDisabled) {
			return nodeResponse(ui.LoginForm(ui.LoginFormData{
				ErrorMessage: "Account is disabled",
			}))
		}
		return responseFromError(err)
	}

	return srv.Respond().
		Cookie("stor_auth", sid, int(session.TTL.Seconds()), "/", "", false, true).
		HxRedirect("/u")
}

func bindFormData(c *srv.Context, args ...any) error {
	req := c.Request()
	if err := req.ParseForm(); err != nil {
		return err
	}

	for i := 0; i < len(args); i += 2 {
		key := args[i].(string)
		val := args[i+1].(*string)

		*val = req.FormValue(key)
	}
	return nil
}

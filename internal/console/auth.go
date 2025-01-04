// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"errors"

	"github.com/cfichtmueller/goparts/e"
	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/session"
	"github.com/cfichtmueller/stor/internal/ec"
	"github.com/cfichtmueller/stor/internal/uc"
	"github.com/cfichtmueller/stor/internal/ui"
)

func handleBootstrapPage(c jug.Context) (e.Node, error) {
	return ui.BootstrapPage(), nil
}

func RequireNotBootstrapped(c jug.Context) {
	bootstrapped, err := uc.IsBootstrapped(c)
	if err != nil {
		c.HandleError(err)
		c.Abort()
		return
	}
	if bootstrapped {
		redirect(c, "/")
		c.Abort()
		return
	}
}

func handleBootstrap(c jug.Context) (e.Node, error) {
	cmd := uc.BootstrapCommand{}
	if err := bindFormData(c, "email", &cmd.Email, "password", &cmd.Password, "passwordConfirmation", &cmd.PasswordConfirmation); err != nil {
		return nil, err
	}
	if err := cmd.Validate(); err != nil {
		return ui.BootstrapForm(&ui.BootstrapFormData{Email: cmd.Email, ErrorMessage: err.Error()}), nil
	}
	if err := uc.Bootstrap(c, cmd); err != nil {
		return ui.BootstrapForm(&ui.BootstrapFormData{Email: cmd.Email, ErrorMessage: err.Error()}), nil
	}

	hxRedirect(c, "/login")
	return nil, nil
}

func handleLoginPage(c jug.Context) (e.Node, error) {
	return ui.LoginPage(), nil
}

func handleLogin(c jug.Context) (e.Node, error) {
	cmd := uc.LoginCommand{
		IpAddress: c.ClientIP(),
	}
	if err := bindFormData(c, "email", &cmd.Email, "password", &cmd.Password); err != nil {
		return nil, err
	}
	_, sid, err := uc.Login(c, cmd)
	if err != nil {
		if errors.Is(err, ec.InvalidCredentials) {
			return ui.LoginForm(ui.LoginFormData{
				Email:        cmd.Email,
				ErrorMessage: "Invalid Credentials",
			}), nil
		} else if errors.Is(err, ec.AccountDisabled) {
			return ui.LoginForm(ui.LoginFormData{
				ErrorMessage: "Account is disabled",
			}), nil
		} else {
			return nil, err
		}
	}
	c.SetCookie("stor_auth", sid, int(session.TTL.Seconds()), "/", "", true, true)
	hxRedirect(c, "/u")
	return nil, nil
}

func bindFormData(c jug.Context, args ...any) error {
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

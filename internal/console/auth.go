// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package console

import (
	"errors"

	"github.com/cfichtmueller/jug"
	"github.com/cfichtmueller/stor/internal/domain/session"
	"github.com/cfichtmueller/stor/internal/ec"
	"github.com/cfichtmueller/stor/internal/uc"
	"github.com/cfichtmueller/stor/internal/ui"
)

func handleBootstrapPage(c jug.Context) {
	bootstrapped, err := uc.IsBootstrapped(c)
	if err != nil {
		c.HandleError(err)
		return
	}
	if bootstrapped {
		redirect(c, "/")
		return
	}
	must("render bootstrap page", c, ui.RenderBootstrapPage(c.Writer()))
}

func handleBootstrap(c jug.Context) {
	cmd := uc.BootstrapCommand{}
	if err := bindFormData(c, "email", &cmd.Email, "password", &cmd.Password, "passwordConfirmation", &cmd.PasswordConfirmation); err != nil {
		c.HandleError(err)
		return
	}
	if err := cmd.Validate(); err != nil {
		must("render bootstrap form", c, ui.RenderBootstrapForm(c.Writer(), ui.BootstrapFormModel{Email: cmd.Email, ErrorMessage: err.Error()}))
		return
	}
	if err := uc.Bootstrap(c, cmd); err != nil {
		must("render bootstrap form", c, ui.RenderBootstrapForm(c.Writer(), ui.BootstrapFormModel{Email: cmd.Email, ErrorMessage: err.Error()}))
		return
	}

	hxRedirect(c, "/login")
}

func handleLoginPage(c jug.Context) {
	must("render login page", c, ui.RenderLoginPage(c.Writer()))
}

func handleLogin(c jug.Context) {
	cmd := uc.LoginCommand{
		IpAddress: c.ClientIP(),
	}
	if err := bindFormData(c, "email", &cmd.Email, "password", &cmd.Password); err != nil {
		c.HandleError(err)
		return
	}
	_, sid, err := uc.Login(c, cmd)
	if err != nil {
		if errors.Is(err, ec.InvalidCredentials) {
			ui.RenderLoginForm(c.Writer(), ui.LoginFormModel{
				Email:        cmd.Email,
				ErrorMessage: "Invalid Credentials",
			})
		} else if errors.Is(err, ec.AccountDisabled) {
			ui.RenderLoginForm(c.Writer(), ui.LoginFormModel{ErrorMessage: "Account is disabled"})
		} else {
			c.HandleError(err)
		}
		return
	}
	c.SetCookie("stor_auth", sid, int(session.TTL.Seconds()), "/", "", true, true)
	hxRedirect(c, "/u")
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

// Copyright 2024 Christoph Fichtmüller. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package user

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/cfichtmueller/stor/internal/db"
	"github.com/cfichtmueller/stor/internal/domain"
	"github.com/cfichtmueller/stor/internal/ec"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           string
	Email        string
	Enabled      bool
	passwordHash []byte
	CreatedAt    time.Time
	LastSeenAt   time.Time
}

func (u *User) SetPassword(pw string) error {
	if len(pw) < 8 || len(pw) > 70 {
		return ec.InvalidArgument
	}
	b, err := bcrypt.GenerateFromPassword([]byte(pw), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("unable to hash password: %v", err)
	}
	u.passwordHash = b
	return nil
}

func (u *User) PasswordMatches(pw string) bool {
	if err := bcrypt.CompareHashAndPassword(u.passwordHash, []byte(pw)); err != nil {
		return false
	}
	return true
}

type CreateCommand struct {
	Email    string
	Password string
}

var (
	createStmt *sql.Stmt
	listStmt   *sql.Stmt
	findStmt   *sql.Stmt
	getStmt    *sql.Stmt
	updateStmt *sql.Stmt
)

func Configure() {
	createStmt = db.Prepare("INSERT INTO users (id, email, enabled, password_hash, created_at, last_seen_at) VALUES ($1, $2, $3, $4, $5, $6)")
	listStmt = db.Prepare("SELECT * FROM users ORDER BY email")
	findStmt = db.Prepare("SELECT * FROM users WHERE email = $1 LIMIT 1")
	getStmt = db.Prepare("SELECT * FROM users WHERE id = $1 LIMIT 1")
	updateStmt = db.Prepare("UPDATE users SET email = $1, enabled = $2, password_hash = $3, last_seen_at = $3 WHERE id = $4")
}

func Urn(id string) string {
	return "user:" + id
}

func Create(ctx context.Context, cmd CreateCommand) (*User, error) {
	email := strings.ToLower(cmd.Email)
	existing, err := find(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ec.UserAlreadyExists
	}

	u := &User{
		ID:        domain.RandomId(),
		Email:     email,
		Enabled:   true,
		CreatedAt: time.Now(),
	}
	if err := u.SetPassword(cmd.Password); err != nil {
		return nil, err
	}

	if _, err := createStmt.ExecContext(
		ctx,
		u.ID,
		u.Email,
		u.Enabled,
		u.passwordHash,
		u.CreatedAt,
		u.LastSeenAt,
	); err != nil {
		return nil, fmt.Errorf("unable to save user: %v", err)
	}

	return u, nil
}

func List(ctx context.Context) ([]*User, error) {
	rows, err := listStmt.QueryContext(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to list users: %v", err)
	}
	users := make([]*User, 0)
	for rows.Next() {
		var u User
		if err := rows.Scan(
			&u.ID,
			&u.Email,
			&u.Enabled,
			&u.passwordHash,
			&u.CreatedAt,
			&u.LastSeenAt,
		); err != nil {
			return nil, fmt.Errorf("unable to decode user: %v", err)
		}
		users = append(users, &u)
	}
	return users, nil
}

func Get(ctx context.Context, id string) (*User, error) {
	u, err := decodeOne(getStmt.QueryRowContext(ctx, id))
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ec.NoSuchUser
	}
	return u, nil
}

func Update(ctx context.Context, u *User) error {
	if _, err := updateStmt.ExecContext(
		ctx,
		u.Email,
		u.Enabled,
		u.passwordHash,
		u.LastSeenAt,
		u.ID,
	); err != nil {
		return fmt.Errorf("unable to update user: %v", err)
	}
	return nil
}

func Login(ctx context.Context, email, password string) (*User, error) {
	u, err := find(ctx, strings.ToLower(email))
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, ec.InvalidCredentials
	}

	if !u.PasswordMatches(password) {
		return nil, ec.InvalidCredentials
	}

	if !u.Enabled {
		return nil, ec.AccountDisabled
	}

	return u, nil
}

func find(ctx context.Context, email string) (*User, error) {
	return decodeOne(findStmt.QueryRowContext(ctx, email))
}

func decodeOne(row *sql.Row) (*User, error) {
	var u User
	if err := row.Scan(
		&u.ID,
		&u.Email,
		&u.Enabled,
		&u.passwordHash,
		&u.CreatedAt,
		&u.LastSeenAt,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("unable to decode user: %v", err)
	}
	return &u, nil
}

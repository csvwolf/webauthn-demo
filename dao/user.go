package dao

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/csvwolf/goserver/db"
	"github.com/csvwolf/goserver/models"
)

type User struct {
}

func NewUser() *User {
	return &User{}
}

func (u *User) GetUser(db *sql.DB, username string) (user *models.User, err error) {
	user = &models.User{}
	query := "SELECT id, username, display_name, registered_at FROM users WHERE username = ?"
	row := db.QueryRow(query, username)
	if err = row.Scan(&user.ID, &user.Username, &user.DisplayName, &user.RegisteredAt); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		fmt.Println("[dao][GetUser] Error get user:", err)
		return
	}
	creds, err := NewPublicKeyCred().FindAllPublicKeyCred(db, user.Username)
	if err != nil {
		fmt.Println("[dao][GetUser] Error get public key cred:", err)
		return user, err
	}
	user.PublicKeyCreds = creds
	return user, err
}

func (u *User) CreateUser(execCtx db.ExecContext, user *models.User) error {
	query := "INSERT INTO users (id, username, display_name, registered_at) VALUES (?, ?, ?, ?)"
	_, err := execCtx.Exec(query, user.ID, user.Username, user.DisplayName, user.RegisteredAt)
	if err != nil {
		fmt.Printf("[dao][CreateUser] Error create user: %v", err)
		return err
	}
	return nil
}

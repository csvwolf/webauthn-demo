package service

import (
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/csvwolf/goserver/dao"
	"github.com/csvwolf/goserver/db"
	"github.com/csvwolf/goserver/models"
	"github.com/go-webauthn/webauthn/webauthn"
	"time"
)

func GetUser(username string) (*models.User, error) {
	return dao.NewUser().GetUser(db.GetClient().GetDB(), username)
}

func CreateUser(user *models.User, cred *webauthn.Credential) error {
	var (
		err      error
		credInfo []byte

		publicKeyCred = &models.PublicKeyCred{
			Username:     user.Username,
			CredentialID: base64.StdEncoding.EncodeToString(cred.ID),
			PublicKey:    base64.StdEncoding.EncodeToString(cred.PublicKey),
		}
	)

	if credInfo, err = json.Marshal(&cred); err != nil {
		fmt.Println("[service][CreateUser] marshal cred error: ", err)
		return err
	}

	publicKeyCred.CredentialInfo = string(credInfo)

	err = db.GetClient().DoTransaction(func(tx *sql.Tx) error {
		user.RegisteredAt = time.Now().UnixMilli()
		err = dao.NewUser().CreateUser(tx, user)
		if err != nil {
			return err
		}
		err = dao.NewPublicKeyCred().CreatePublicKey(tx, publicKeyCred)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		fmt.Println("[service][CreateUser] create user error: ", err)
		return err
	}
	return nil
}

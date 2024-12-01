package dao

import (
	"fmt"
	"github.com/csvwolf/goserver/db"
	"github.com/csvwolf/goserver/models"
)

type PublicKeyCred struct{}

func NewPublicKeyCred() *PublicKeyCred {
	return &PublicKeyCred{}
}

func (p *PublicKeyCred) CreatePublicKey(execCtx db.ExecContext, cred *models.PublicKeyCred) error {
	query := "INSERT INTO public_key_creds (username, credential_id, public_key, credential_info) VALUES (?, ?, ?, ?)"
	_, err := execCtx.Exec(query, cred.Username, cred.CredentialID, cred.PublicKey, cred.CredentialInfo)
	if err != nil {
		fmt.Println("[dao][CreatePublicKey] Error: ", err)
		return err
	}
	return nil
}

func (p *PublicKeyCred) FindAllPublicKeyCred(queryCtx db.QueryContext, username string) ([]models.PublicKeyCred, error) {
	query := "SELECT id, username, credential_id, public_key, credential_info  FROM public_key_creds WHERE username = ?"
	rows, err := queryCtx.Query(query, username)
	if err != nil {
		fmt.Println("[dao][FindAllPublicKeyCred] Error: ", err)
		return nil, err
	}
	defer rows.Close()

	var creds []models.PublicKeyCred
	for rows.Next() {
		var cred models.PublicKeyCred
		err := rows.Scan(&cred.ID, &cred.Username, &cred.CredentialID, &cred.PublicKey, &cred.CredentialInfo)
		if err != nil {
			fmt.Println("[dao][FindAllPublicKeyCred] Error: ", err)
			return nil, err
		}
		creds = append(creds, cred)
	}
	return creds, nil
}

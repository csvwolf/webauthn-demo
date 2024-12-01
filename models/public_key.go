package models

import (
	"encoding/json"
	"fmt"
	"github.com/go-webauthn/webauthn/webauthn"
)

type PublicKeyCred struct {
	ID             string `json:"id"`              // Primary key
	CredentialID   string `json:"credential_id"`   // 凭证 ID
	PublicKey      string `json:"public_key"`      // 公钥数据（经过 base64 编码的字符串）
	CredentialInfo string `json:"credential_info"` // JSON 格式的凭证信息
	Username       string `json:"username"`        // 用户名
}

func (cred *PublicKeyCred) ToWebAuthnCredential() (webauthn.Credential, error) {
	var (
		credential webauthn.Credential
		err        error
	)
	if err = json.Unmarshal([]byte(cred.CredentialInfo), &credential); err != nil {
		fmt.Println("[ToWebauthnCredential] json.Unmarshal error:", err)
	}
	return credential, err
}

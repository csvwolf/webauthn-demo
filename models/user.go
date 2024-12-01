package models

import (
	"encoding/binary"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
	"math/rand"
)

type User struct {
	ID             int64           `json:"id"`               // 用户唯一标识符
	Username       string          `json:"username"`         // 用户名
	DisplayName    string          `json:"display_name"`     // 用户显示名称
	CredentialIDs  []string        `json:"credential_ids"`   // 存储用户所有的 WebAuthn 凭证 ID（允许多个设备）
	PublicKeyCreds []PublicKeyCred `json:"public_key_creds"` // 存储 WebAuthn 公钥凭证信息
	RegisteredAt   int64           `json:"registered_at"`    // 注册时间戳
}

func (u *User) GenUserID() int64 {
	return rand.Int63()
}

func (u *User) WebAuthnID() []byte {
	if u == nil {
		return []byte{}
	}
	buf := make([]byte, binary.MaxVarintLen64)
	binary.PutUvarint(buf, uint64(u.ID))
	return buf
}

func (u *User) WebAuthnName() string {
	if u == nil {
		return ""
	}
	return u.Username
}

func (u *User) WebAuthnDisplayName() string {
	if u == nil {
		return ""
	}
	return u.DisplayName
}

func (u *User) WebAuthnCredentials() []webauthn.Credential {
	creds := make([]webauthn.Credential, 0)
	if u == nil {
		return creds
	}
	for _, cred := range u.PublicKeyCreds {
		credential, err := cred.ToWebAuthnCredential()
		if err != nil {
			continue
		}
		creds = append(creds, credential)
	}
	return creds
}

func (u *User) CredentialExcludeList() []protocol.CredentialDescriptor {
	var excludeList = make([]protocol.CredentialDescriptor, 0)
	if u == nil {
		return excludeList
	}
	for _, cred := range u.WebAuthnCredentials() {
		excludeList = append(excludeList, protocol.CredentialDescriptor{
			Type:         protocol.PublicKeyCredentialType,
			CredentialID: cred.ID,
		})
	}

	return excludeList
}

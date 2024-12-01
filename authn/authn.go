package authn

import (
	"github.com/go-webauthn/webauthn/webauthn"
)

var (
	authn *webauthn.WebAuthn
)

func NewAuthn() {
	var err error
	wconfig := &webauthn.Config{
		RPID:          "localhost",
		RPDisplayName: "WebAuthn Demo",
		RPOrigins:     []string{"http://localhost:8080"},
	}
	if authn, err = webauthn.New(wconfig); err != nil {
		panic("WebAuthn NewError: " + err.Error())
	}
	return
}

func GetAuthn() *webauthn.WebAuthn {
	return authn
}

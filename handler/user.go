package handler

import (
	"encoding/json"
	"fmt"
	"github.com/csvwolf/goserver/authn"
	"github.com/csvwolf/goserver/models"
	"github.com/csvwolf/goserver/service"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-webauthn/webauthn/protocol"
	"github.com/go-webauthn/webauthn/webauthn"
)

const (
	RegSessionDataKey    = "register_data"
	RegUserTempDataKey   = "temp_user"
	LoginSessionDataKey  = "login_data"
	LoginUserTempDataKey = "temp_user"
)

// https://webauthn.guide/#webauthn-api

type BeginRegisterReq struct {
	Username string `json:"username" binding:"required"`
	Nickname string `json:"nickname" binding:"required"`
}

type BeginLoginReq struct {
	Username string `json:"username" binding:"required"`
}

func BeginRegister(c *gin.Context) {
	var req BeginRegisterReq
	session := sessions.Default(c)
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}
	user, err := service.GetUser(req.Username)
	if err != nil {
		fmt.Println("[handler][BeginRegister] service.GetUser error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if user != nil {
		c.JSON(400, gin.H{"error": "user already exists"})
		return
	}
	user = &models.User{}
	user.ID = user.GenUserID()
	user.Username = req.Username
	user.DisplayName = req.Nickname
	registerOptions := func(credCreationOpts *protocol.PublicKeyCredentialCreationOptions) {
		credCreationOpts.CredentialExcludeList = user.CredentialExcludeList()
	}
	options, sessionData, err := authn.GetAuthn().BeginRegistration(user, registerOptions)
	if err != nil {
		fmt.Println("[handler][BeginRegister] authn.BeginRegistration error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if sessionDataStr, err := json.Marshal(&sessionData); err != nil {
		fmt.Println("[handler][BeginRegister] json.Marshal error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		session.Set(RegSessionDataKey, sessionDataStr)
	}
	if userStr, err := json.Marshal(&user); err != nil {
		fmt.Println("[handler][BeginRegister] json.Marshal error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		session.Set(RegUserTempDataKey, userStr)
	}
	err = session.Save()
	if err != nil {
		fmt.Println("[handler][BeginRegister] session.Save error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": options})
}

func FinishRegister(c *gin.Context) {
	var (
		//req             = protocol.ParsedPublicKeyCredential{}
		sessionData     = webauthn.SessionData{}
		user            = models.User{}
		credential      *webauthn.Credential
		err             error
		ok              bool
		registerDataStr []byte
		userStr         []byte
	)
	//if err = c.ShouldBindJSON(&req); err != nil {
	//	c.JSON(400, gin.H{"error": err.Error()})
	//	return
	//}
	session := sessions.Default(c)
	if registerDataStr, ok = session.Get(RegSessionDataKey).([]byte); !ok {
		c.JSON(400, gin.H{"error": "register_data not found"})
		return
	}
	if userStr, ok = session.Get(RegUserTempDataKey).([]byte); !ok {
		c.JSON(400, gin.H{"error": "temp_user not found"})
		return
	}
	if err = json.Unmarshal(registerDataStr, &sessionData); err != nil {
		fmt.Println("[handler][FinishRegister] json.Unmarshal error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err = json.Unmarshal(userStr, &user); err != nil {
		fmt.Println("[handler][FinishRegister] json.Unmarshal error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	if credential, err = authn.GetAuthn().FinishRegistration(&user, sessionData, c.Request); err != nil {
		fmt.Printf("[handler][FinishRegister] authn.FinishRegistration error: %+v\n", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	if err = service.CreateUser(&user, credential); err != nil {
		fmt.Println("[handler][FinishRegister] service.CreateUser error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	session.Delete(RegSessionDataKey)
	session.Delete(RegUserTempDataKey)
	if err = session.Save(); err != nil {
		fmt.Println("[handler][FinishRegister] session.Save error", err)
	}
	c.JSON(200, gin.H{"success": true})
}

func BeginLogin(c *gin.Context) {
	var req = BeginLoginReq{}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	user, err := service.GetUser(req.Username)
	if err != nil {
		fmt.Println("[handler][BeginLogin] service.GetUser error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	options, sessionData, err := authn.GetAuthn().BeginLogin(user)
	if err != nil {
		fmt.Println("[handler][BeginLogin] authn.BeginLogin error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	session := sessions.Default(c)
	if sessionDataStr, err := json.Marshal(&sessionData); err != nil {
		fmt.Println("[handler][BeginLogin] json.Marshal error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		session.Set(LoginSessionDataKey, sessionDataStr)
	}
	if userStr, err := json.Marshal(&user); err != nil {
		fmt.Println("[handler][BeginRegister] json.Marshal error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	} else {
		session.Set(LoginUserTempDataKey, userStr)
	}
	if err = session.Save(); err != nil {
		fmt.Println("[handler][BeginLogin] session.Save error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": options})
}

func FinishLogin(c *gin.Context) {
	var (
		sessionData = webauthn.SessionData{}
		user        = models.User{}
		err         error
	)
	session := sessions.Default(c)
	if sessionDataStr, ok := session.Get(LoginSessionDataKey).([]byte); !ok {
		c.JSON(400, gin.H{"error": "login_data not found"})
		return
	} else {
		if err = json.Unmarshal(sessionDataStr, &sessionData); err != nil {
			fmt.Println("[handler][FinishLogin] json.Unmarshal error", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	if userStr, ok := session.Get(LoginUserTempDataKey).([]byte); !ok {
		c.JSON(400, gin.H{"error": "temp_user not found"})
		return
	} else {
		if err = json.Unmarshal(userStr, &user); err != nil {
			fmt.Println("[handler][FinishLogin] json.Unmarshal error", err)
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
	}
	_, err = authn.GetAuthn().FinishLogin(&user, sessionData, c.Request)
	if err != nil {
		fmt.Println("[handler][FinishLogin] authn.FinishLogin error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	session.Delete(LoginSessionDataKey)
	session.Delete(LoginUserTempDataKey)
	if err = session.Save(); err != nil {
		fmt.Println("[handler][FinishLogin] session.Save error", err)
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"success": true, "data": user})
}

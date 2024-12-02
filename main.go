package main

import (
	"github.com/csvwolf/goserver/authn"
	"github.com/csvwolf/goserver/db"
	"github.com/csvwolf/goserver/handler"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

func main() {
	authn.NewAuthn()
	dbclient := db.NewClient()
	defer dbclient.Close()

	r := gin.Default()
	store, _ := redis.NewStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	r.Use(sessions.Sessions("mysession", store))
	r.LoadHTMLFiles("./index.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	r.GET("/ping", handler.Ping)
	r.GET("/user", handler.GetCurrentUser)
	r.POST("logout", handler.Logout)
	r.POST("/register/begin", handler.BeginRegister)
	r.POST("/register/finish", handler.FinishRegister)
	r.POST("/login/begin", handler.BeginLogin)
	r.POST("/login/finish", handler.FinishLogin)
	r.Run() // 监听并在 0.0.0.0:8080 上启动服务
}

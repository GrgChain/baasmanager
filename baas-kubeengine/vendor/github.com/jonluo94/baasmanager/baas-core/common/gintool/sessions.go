package gintool

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/memstore"
)

var store = memstore.NewStore([]byte("secret"))

func UseSession(router *gin.Engine) {
	router.Use(sessions.Sessions("mysession", store))
}

func SetSession(ctx *gin.Context, k string, o interface{}) {
	session := sessions.Default(ctx)
	session.Set(k, o)
	session.Save()
}

func GetSession(ctx *gin.Context, k string) interface{} {
	session := sessions.Default(ctx)
	return session.Get(k)
}

func RemoveSession(ctx *gin.Context, k string) {
	session := sessions.Default(ctx)
	session.Delete(k)
}

func ClearAllSession(ctx *gin.Context) {
	session := sessions.Default(ctx)
	session.Clear()
	session.Save()
	return
}

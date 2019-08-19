package controller

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"github.com/jonluo94/baasmanager/baas-gateway/entity"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
	"github.com/jonluo94/baasmanager/baas-gateway/model"
	"github.com/jonluo94/baasmanager/baas-core/common/password"
)

func (a *ApiController) UserAdd(ctx *gin.Context) {

	user := new(entity.User)

	if err := ctx.ShouldBindJSON(user); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.userService.Add(user)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}
func (a *ApiController) UserAddAuth(ctx *gin.Context) {

	ur := new(entity.UserRole)

	if err := ctx.ShouldBindJSON(ur); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.userService.AddAuth(ur)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) UserDelAuth(ctx *gin.Context) {

	ur := new(entity.UserRole)

	if err := ctx.ShouldBindJSON(ur); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.userService.DelAuth(ur)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) UserUpdate(ctx *gin.Context) {

	user := new(entity.User)

	if err := ctx.ShouldBindJSON(user); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.userService.Update(user)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) UserDelete(ctx *gin.Context) {

	user := new(entity.User)

	if err := ctx.ShouldBindJSON(user); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.userService.Delete(user.Id)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) UserLogin(ctx *gin.Context) {

	login := new(model.LoginForm)
	if err := ctx.ShouldBind(&login); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	user := &entity.User{
		Account: login.UserName,
	}
	has, u := a.userService.GetByUser(user)
	if !has {
		gintool.ResultFail(ctx, "username error")
		return
	}
	vali := password.Validate(login.Password, u.Password)
	if !vali {
		gintool.ResultFail(ctx, "password error")
		return
	}

	type UserInfo map[string]interface{}

	token := a.userService.GetToken(u)
	//保存session
	gintool.SetSession(ctx, token.Token, u.Id)
	gintool.ResultOk(ctx, token)

}

func (a *ApiController) UserLogout(ctx *gin.Context) {
	token := ctx.GetHeader("X-Token")
	gintool.RemoveSession(ctx, token)
	gintool.ResultMsg(ctx, "logout success")
}

func (a *ApiController) UserInfo(ctx *gin.Context) {

	token := ctx.Query("token")

	session := gintool.GetSession(ctx, token)
	if nil == session {
		gintool.ResultFail(ctx, "token不存在")
		return
	}
	user, err := a.userService.CheckToken(token, &entity.User{Id: session.(int)})

	if err != nil {
		if err.Error() == "token已过期" || err.Error() == "token无效" {
			m := make(map[string]interface{})
			m["code"] = 2
			m["msg"] = err.Error()
			gintool.ResultMap(ctx, m)
		}
		gintool.ResultFail(ctx, err.Error())
	} else {
		gintool.ResultOk(ctx, user)
	}
}

func (a *ApiController) UserAuthorize(ctx *gin.Context) {
	var token string
	var err error
	token = ctx.GetHeader("X-Token")
	if token == "" {
		token, err = ctx.Cookie("Admin-Token")
		if err != nil {
			gintool.ResultFail(ctx, err.Error())
			ctx.Abort()
			return
		}
	}

	session := gintool.GetSession(ctx, token)
	if nil == session {
		gintool.ResultFail(ctx, "token不存在")
		return
	}
	_, err = a.userService.CheckToken(token, &entity.User{Id: session.(int)})

	if err != nil {
		if err.Error() == "token已过期" || err.Error() == "token无效" {
			m := make(map[string]interface{})
			m["code"] = 2
			m["msg"] = err.Error()
			gintool.ResultMap(ctx, m)
		} else {
			gintool.ResultFail(ctx, err.Error())
		}
		ctx.Abort()
		return
	} else {
		ctx.Next()
	}
}

func (a *ApiController) UserList(ctx *gin.Context) {

	page, err := strconv.Atoi(ctx.Query("page"))
	if err != nil {
		gintool.ResultFail(ctx, "page error")
		return
	}
	limit, err := strconv.Atoi(ctx.Query("limit"))
	if err != nil {
		gintool.ResultFail(ctx, "limit error")
		return
	}
	name := ctx.Query("name")

	b, list, total := a.userService.GetList(&entity.User{Name: name}, page, limit)

	if b {
		gintool.ResultList(ctx, list, total)

	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

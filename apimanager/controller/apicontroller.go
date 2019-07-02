package controller

import (
	"gitee.com/jonluo/baasmanager/apimanager/service"
	"github.com/gin-gonic/gin"
	"gitee.com/jonluo/baasmanager/apimanager/entity"
	"github.com/jonluo94/commontools/gintool"
	"strconv"
	"github.com/jonluo94/commontools/password"
	"time"
	"net/http"
	"fmt"
)

type LoginForm struct {
	UserName string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type ApiController struct {
	chainService     *service.ChainService
	channelService   *service.ChannelService
	chaincodeService *service.ChaincodeService
	dashboardService *service.DashboardService
	userService      *service.UserService
	roleService      *service.RoleService
}

func NewApiController(userService *service.UserService, roleService *service.RoleService, chainService *service.ChainService, channelService *service.ChannelService, chaincodeService *service.ChaincodeService, dashboardService *service.DashboardService) *ApiController {
	return &ApiController{
		userService:      userService,
		roleService:      roleService,
		chainService:     chainService,
		channelService:   channelService,
		chaincodeService: chaincodeService,
		dashboardService: dashboardService,
	}
}

func (a *ApiController) Upload(ctx *gin.Context) {
	// single file
	file, _ := ctx.FormFile("file")
	path := fmt.Sprintf("/tmp/%d", time.Now().UnixNano())
	ctx.SaveUploadedFile(file, path)
	ctx.String(http.StatusOK, path)

}

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

	login := new(LoginForm)
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
		token,err = ctx.Cookie("Admin-Token")
		if err != nil{
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
		}else {
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

func (a *ApiController) RoleList(ctx *gin.Context) {

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

	b, list, total := a.roleService.GetList(&entity.Role{Name: name}, page, limit)

	if b {
		gintool.ResultList(ctx, list, total)

	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) RoleAllList(ctx *gin.Context) {

	b, list := a.roleService.GetAll()
	if b {
		gintool.ResultOk(ctx, list)

	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) RoleAdd(ctx *gin.Context) {

	role := new(entity.Role)

	if err := ctx.ShouldBindJSON(role); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.roleService.Add(role)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) RoleUpdate(ctx *gin.Context) {

	role := new(entity.Role)

	if err := ctx.ShouldBindJSON(role); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.roleService.Update(role)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) RoleDelete(ctx *gin.Context) {

	role := new(entity.Role)

	if err := ctx.ShouldBindJSON(role); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.roleService.Delete(role.Rkey)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainAdd(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	chain.Created = time.Now().Unix()
	isSuccess, msg := a.chainService.Add(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainGet(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, chain := a.chainService.GetByChain(chain)
	if isSuccess {
		gintool.ResultOk(ctx, chain)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) ChainList(ctx *gin.Context) {

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
	userAccount := ctx.Query("userAccount")

	if a.userService.HasAdminRole(userAccount) {
		//admin 可看所有
		userAccount = ""
	}

	b, list, total := a.chainService.GetList(&entity.Chain{
		Name:        name,
		UserAccount: userAccount,
	}, page, limit)

	if b {
		gintool.ResultList(ctx, list, total)

	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) ChainUpdate(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.chainService.Update(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainDeleted(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.chainService.Delete(chain.Id)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainBuild(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := a.chainService.BuildChain(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}

}

func (a *ApiController) ChainRun(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := a.chainService.RunChain(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}

}

func (a *ApiController) ChainStop(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.chainService.StopChain(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChainRelease(ctx *gin.Context) {

	chain := new(entity.Chain)

	if err := ctx.ShouldBindJSON(chain); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.chainService.ReleaseChain(chain)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}



func (a *ApiController) ChainDownload(ctx *gin.Context) {

	chainId, err := strconv.Atoi(ctx.Query("chainId"))
	if err != nil {
		gintool.ResultFail(ctx, "chainId error")
		return
	}

	chain := new(entity.Chain)
	chain.Id = chainId
	isSuccess, chain := a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	reader, contentLength, name := a.chainService.DownloadChainArtifacts(chain)
	extraHeaders := map[string]string{
		"Content-Disposition": fmt.Sprintf(`attachment; filename="%s"`, name),
	}

	ctx.DataFromReader(http.StatusOK, contentLength, "application/x-tar", reader, extraHeaders)

}

func (a *ApiController) ChannelAdd(ctx *gin.Context) {

	channel := new(entity.Channel)

	if err := ctx.ShouldBindJSON(channel); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain := a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	isSuccess, msg := a.channelService.AddChannel(chain, channel)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChannelGet(ctx *gin.Context) {

	chn := new(entity.Channel)

	if err := ctx.ShouldBindJSON(chn); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, chn := a.channelService.GetByChannel(chn)
	if isSuccess {
		gintool.ResultOk(ctx, chn)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) ChannelAll(ctx *gin.Context) {

	chainId, err := strconv.Atoi(ctx.Query("chainId"))
	if err != nil {
		gintool.ResultFail(ctx, "chainId error")
		return
	}
	isSuccess, data := a.channelService.GetAllList(chainId)
	if isSuccess {
		gintool.ResultOk(ctx, data)
	} else {
		gintool.ResultFail(ctx, data)
	}
}

func (a *ApiController) ChaincodeAdd(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, "channel 不存在")
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	isSuccess, msg := a.chaincodeService.AddChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeDeploy(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, "channel 不存在")
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	isSuccess, msg := a.chaincodeService.DeployChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeUpgrade(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, "channel 不存在")
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	isSuccess, msg := a.chaincodeService.UpgradeChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeQuery(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, "channel 不存在")
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	isSuccess, msg := a.chaincodeService.QueryChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultOk(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeInvoke(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	channel := new(entity.Channel)
	channel.Id = cc.ChannelId
	isSuccess, channel := a.channelService.GetByChannel(channel)
	if !isSuccess {
		gintool.ResultFail(ctx, "channel 不存在")
		return
	}

	chain := new(entity.Chain)
	chain.Id = channel.ChainId
	isSuccess, chain = a.chainService.GetByChain(chain)
	if !isSuccess {
		gintool.ResultFail(ctx, "chain 不存在")
		return
	}

	isSuccess, msg := a.chaincodeService.InvokeChaincode(chain, channel, cc)
	if isSuccess {
		gintool.ResultOk(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeGet(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, chain := a.chaincodeService.GetByChaincode(cc)
	if isSuccess {
		gintool.ResultOk(ctx, chain)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) ChaincodeList(ctx *gin.Context) {

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
	name := ctx.Query("chaincodeName")
	channelId, err := strconv.Atoi(ctx.Query("channelId"))
	if err != nil {
		gintool.ResultFail(ctx, "channelId error")
		return
	}
	b, list, total := a.chaincodeService.GetList(&entity.Chaincode{
		ChaincodeName: name,
		ChannelId:     channelId,
	}, page, limit)

	if b {
		gintool.ResultList(ctx, list, total)

	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) ChaincodeUpdate(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.chaincodeService.Update(cc)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) ChaincodeDeleted(ctx *gin.Context) {

	cc := new(entity.Chaincode)

	if err := ctx.ShouldBindJSON(cc); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	isSuccess, msg := a.chaincodeService.Delete(cc.Id)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (a *ApiController) DashboardCounts(ctx *gin.Context) {

	userAccount := ctx.Query("userAccount")
	if a.userService.HasAdminRole(userAccount) {
		//admin 可看所有
		userAccount = ""
	}
	isSuccess, ash := a.dashboardService.Counts(userAccount)
	if isSuccess {
		gintool.ResultOk(ctx, ash)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

func (a *ApiController) DashboardSevenDays(ctx *gin.Context) {

	start, err := strconv.Atoi(ctx.Query("start"))
	if err != nil {
		gintool.ResultFail(ctx, "start error")
		return
	}
	end, err := strconv.Atoi(ctx.Query("end"))
	if err != nil {
		gintool.ResultFail(ctx, "end error")
		return
	}

	userAccount := ctx.Query("userAccount")
	if a.userService.HasAdminRole(userAccount) {
		//admin 可看所有
		userAccount = ""
	}
	isSuccess, ash := a.dashboardService.SevenDays(userAccount, start, end)
	if isSuccess {
		gintool.ResultOk(ctx, ash)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

package controller

import (
	"github.com/gin-gonic/gin"
	"time"
	"net/http"
	"fmt"
	"github.com/jonluo94/baasmanager/baas-gateway/service"
)

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

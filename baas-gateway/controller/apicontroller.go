package controller

import (
	"data/baas-gateway/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"
)

type ApiController struct {
	chainService     *service.ChainService
	channelService   *service.ChannelService
	chaincodeService *service.ChaincodeService
	dashboardService *service.DashboardService
	userService      *service.UserService
	roleService      *service.RoleService
	orgService       *service.OrgService
}

func NewApiController(userService *service.UserService, roleService *service.RoleService, chainService *service.ChainService, channelService *service.ChannelService, chaincodeService *service.ChaincodeService, dashboardService *service.DashboardService, orgService *service.OrgService) *ApiController {
	return &ApiController{
		userService:      userService,
		roleService:      roleService,
		chainService:     chainService,
		channelService:   channelService,
		chaincodeService: chaincodeService,
		dashboardService: dashboardService,

		orgService: orgService,
	}
}

func (a *ApiController) UploadToNano(ctx *gin.Context) {
	// single file
	file, _ := ctx.FormFile("file")
	path := fmt.Sprintf("/tmp/%d", time.Now().UnixNano())
	ctx.SaveUploadedFile(file, path)
	ctx.String(http.StatusOK, path)
}

type UploadForm struct {
	Filename     string                `form:"filename"`
	RelativePath string                `form:"relativePath"`
	File         *multipart.FileHeader `form:"file"`
}

func uploadDir(prefix string, form UploadForm) string {
	var uploadir string
	// 切分保存路径
	idx := strings.LastIndex(form.RelativePath, form.Filename) // 切掉文件名的部分
	uploadir = form.RelativePath[0 : idx-1]
	// 创建保存目录
	uploadir = prefix + uploadir
	return uploadir
}

func (a *ApiController) Upload(ctx *gin.Context) {
	// single file with relative path
	// 接收参数
	var form UploadForm
	err := ctx.ShouldBind(&form)

	if err != nil {
		fmt.Println("..bind UploadForm error..", err)
		return
	}

	//fmt.Println(form.RelativePath, form.Filename, form.File)
	// 创建目录
	prefixPath := "/tmp/chaincode/"
	path := uploadDir(prefixPath, form)

	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
	}
	//创建文件
	//fmt.Println("..before save file..", path)
	err = ctx.SaveUploadedFile(form.File, path+string(os.PathSeparator)+form.Filename)
	if err != nil {
		fmt.Println("..save file error..", err)
		return
	}

	fmt.Printf("..save file success, path = <%s> filename = <%s>\n", path, form.Filename)
	ctx.String(http.StatusOK, prefixPath)
}

package controller

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"

	"data/baas-gateway/entity"
	"data/baas-gateway/model"
	"github.com/gin-gonic/gin"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
	uuid "github.com/satori/go.uuid"
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
	if u.IdCard != "激活" {
		gintool.ResultFail(ctx, "用户无登录权限，请联系管理员")
		return
	}
	if login.Password != u.Password {
		// }
		// vali := password.Validate(login.Password, u.Password)
		// if !vali {
		gintool.ResultFail(ctx, "password error")
		return
	}

	type UserInfo map[string]interface{}
	// orgid to channel, chain, chaincode
	org := &entity.Org{
		Id: login.Org,
	}
	has, org = a.orgService.GetByOrg(org)
	if has {
		fmt.Println("has", has)
		fmt.Printf(
			"org.Id=%d\n org.Chain=%d\n org.Channels=%s\n org.Name=%s\n org.Publickey=%s\n",
			org.Id,
			org.Chain,
			org.Channels,
			org.Name,
			org.Publickey,
		)
	}

	token := a.userService.GetToken(u)
	// 保存session
	gintool.SetSession(ctx, token.Token, u.Id)
	gintool.ResultOk(ctx, token)
}

type RegisterRequest struct {
	Account         string `json:"account" binding:"required"`
	Password        string `json:"password" binding:"required"`
	Phone           string `json:"phone" binding:"required"`
	Org             string `json:"org" binding:"required"`
	BusinessLicense string `json:"businessLicense" `
}

func (a *ApiController) UserRegister(ctx *gin.Context) {
	user := new(RegisterRequest)
	if err := ctx.ShouldBindJSON(user); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	// // form获取文件
	// fileReader, err := ctx.FormFile("businessLicense")
	// if err != nil {
	// 	gintool.ResultFail(ctx, err)
	// 	return
	// }
	//
	// fi, err := fileReader.Open()
	// if err != nil {
	// 	gintool.ResultFail(ctx, err)
	// 	return
	// }
	//
	// //// ctx.Request.Body
	// body, err := io.ReadAll(fi)
	// if err != nil {
	// 	gintool.ResultFail(ctx, err)
	// 	return
	// }
	//
	// // 获取uuid
	// u1 := uuid.NewV4()
	// // 创建文件
	// file, err := os.Create("./file/" + u1.String() + ".png")
	// if err != nil {
	// 	gintool.ResultFail(ctx, err)
	// 	return
	// }
	// defer file.Close()
	//
	// // 写入文件
	// _, err = file.Write(body)
	// if err != nil {
	// 	gintool.ResultFail(ctx, err)
	// 	return
	// }
	//
	// if user.Org != "farmer" && user.Org != "maker" {
	// 	gintool.ResultFail(ctx, "请选择身份注册")
	// 	return
	// }

	utype, _ := strconv.Atoi(user.Org)
	isSuccess, msg := a.userService.Register(&entity.User{
		Account:  user.Account,
		Password: user.Password,
		Phone:    user.Phone,
		IdCard:   "",
		Type:     utype,
		Img:      user.BusinessLicense,
	})
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, "注册失败，用户名可能已存在")
	}
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

func (a *ApiController) Cors(c *gin.Context) {
	//c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	//c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
	//c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	//c.Writer.Header().
	//	Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	//
	//method := c.Request.Method

	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header(
		"Access-Control-Expose-Headers",
		"Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type",
	)
	c.Header("Access-Control-Allow-Credentials", "true")
	if c.Request.Method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}
	//c.AbortWithStatus(http.StatusNoContent)
	//} else {
	c.Next()
	//}
}

func (a *ApiController) UserAuthorize(ctx *gin.Context) {
	var token string
	var err error
	m := make(map[string]interface{})
	m["code"] = 2

	token = ctx.GetHeader("X-Token")
	if token == "" {
		token, err = ctx.Cookie("Admin-Token")
		if err != nil {
			m["msg"] = err.Error()
			gintool.ResultMap(ctx, m)
			ctx.Abort()
			return
		}
	}

	session := gintool.GetSession(ctx, token)
	if nil == session {
		m["msg"] = "token不存在"
		gintool.ResultMap(ctx, m)
		ctx.Abort()
		return
	}
	user, err := a.userService.CheckToken(token, &entity.User{Id: session.(int)})
	if err != nil {
		if err.Error() == "token已过期" || err.Error() == "token无效" {
			m["msg"] = err.Error()
			gintool.ResultMap(ctx, m)
		} else {
			gintool.ResultFail(ctx, err.Error())
		}
		ctx.Abort()
		return
	} else {
		ctx.Set("userid", user.Id)
		ctx.Set("userName", user.Name)
		ctx.Set("userType", user.Type)
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

// 上传文件
func (f *ApiController) UpFile(ctx *gin.Context) {
	// form获取文件
	fileReader, err := ctx.FormFile("file")
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	fi, err := fileReader.Open()
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	//// ctx.Request.Body
	body, err := io.ReadAll(fi)
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	// 获取uuid
	u1 := uuid.NewV4()
	// 创建文件
	file, err := os.Create("./file/" + u1.String() + ".png")
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	defer file.Close()

	// 写入文件
	_, err = file.Write(body)
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	gintool.ResultList(ctx, u1.String(), 0)
}

// 下载文件
func (f *ApiController) DownFile(ctx *gin.Context) {
	fid := ctx.Query("fileId")
	// 创建文件
	// file, err := os.Open()
	// if err != nil {
	// 	gintool.ResultFail(ctx, err)
	// 	return
	// }
	// defer file.Close()
	ctx.File("./file/" + fid + ".png")
}

// FindUserList
func (f *ApiController) FindUserList(ctx *gin.Context) {
	userid, isOk := ctx.Get("userid")
	if !isOk {
		gintool.ResultFail(ctx, "fail")
		return
	}
	activet := ctx.Query("isok")
	userId := userid.(int)
	userList, err := f.userService.FindUserListById(userId)
	if err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	userList2 := make([]*entity.User, 0)
	if activet != "" {
		for _, u := range userList {
			if u.IsOk == "通过" {
				userList2 = append(userList2, u)
			}
		}
		gintool.ResultList(ctx, userList2, 1)
		return
	}
	gintool.ResultList(ctx, userList, 1)
}

// UpdateUser
func (f *ApiController) UpdateUser(ctx *gin.Context) {
	var user entity.User
	if err := ctx.ShouldBindJSON(&user); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}
	fmt.Println("user:", user)
	// if user.Password != "" {
	// 	user.Password = password.Encode(user.Password, 12, "default")
	// }

	if err := f.userService.UpdateUserById(&user); err != nil {
		gintool.ResultFail(ctx, err.Error())
		return
	}
	gintool.ResultOk(ctx, nil)
}

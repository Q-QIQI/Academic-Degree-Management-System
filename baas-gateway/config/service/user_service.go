package service

import (
	"fmt"
	"time"

	"data/baas-gateway/entity"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-xorm/xorm"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
	jwttool "github.com/jonluo94/baasmanager/baas-core/common/jwt"
	"github.com/jonluo94/baasmanager/baas-core/common/password"
)

const TokenKey = "baas user secret"

type UserService struct {
	DbEngine *xorm.Engine
}

func (l *UserService) Add(user *entity.User) (bool, string) {
	if user.Password != "" {
		user.Password = password.Encode(user.Password, 12, "default")
	}
	i, err := l.DbEngine.Insert(user)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "add success"
	}
	return false, "add fail"
}

func (l *UserService) Update(user *entity.User) (bool, string) {
	if user.Password != "" {
		user.Password = password.Encode(user.Password, 12, "default")
	}
	i, err := l.DbEngine.Where("id = ?", user.Id).Update(user)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "update success"
	}
	return false, "update fail"
}

func (l *UserService) Delete(id int) (bool, string) {
	i, err := l.DbEngine.Where("id = ?", id).Delete(&entity.User{})
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "delete success"
	}
	return false, "delete fail"
}

func (l *UserService) GetByUser(user *entity.User) (bool, *entity.User) {
	has, err := l.DbEngine.Get(user)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, user
}

func (l *UserService) GetList(user *entity.User, page, size int) (bool, []entity.UserDetail, int64) {
	pager := gintool.CreatePager(page, size)

	users := make([]*entity.User, 0)

	values := make([]interface{}, 0)
	where := "1=1"
	if user.Account != "" {
		where += " and account = ? "
		values = append(values, user.Account)
	}
	if user.Name != "" {
		where += " and name like ? "
		values = append(values, "%"+user.Name+"%")
	}

	err := l.DbEngine.Where(where, values...).Limit(pager.PageSize, pager.NumStart).Find(&users)
	if err != nil {
		logger.Error(err.Error())
	}

	total, err := l.DbEngine.Where(where, values...).Count(new(entity.User))
	if err != nil {
		logger.Error(err.Error())
	}

	userIds := make([]int, len(users))
	userDatas := make([]entity.UserDetail, len(users))
	for i, u := range users {
		userIds[i] = u.Id
		userDatas[i].Id = u.Id
		userDatas[i].Account = u.Account
		userDatas[i].Password = u.Password
		// userDatas[i].Avatar = u.Avatar
		userDatas[i].Name = u.Name
		userDatas[i].Created = u.Created
	}

	roles := make([]entity.UserRole, 0)
	err = l.DbEngine.In("user_id", userIds).Find(&roles)
	if err != nil {
		logger.Error(err.Error())
	}

	for i, d := range userDatas {
		keys := make([]string, 0)
		for _, r := range roles {
			if r.UserId == d.Id {
				keys = append(keys, r.RoleKey)
			}
		}
		d.Roles = keys
		userDatas[i] = d
	}

	return true, userDatas, total
}

// 生成 JWT 令牌。
func (l *UserService) GetToken(user *entity.User) *entity.JwtToken {
	info := make(map[string]interface{})
	now := time.Now()
	info["userId"] = user.Id
	info["exp"] = now.Add(time.Hour * 1).Unix() // 1 小时过期
	info["iat"] = now.Unix()
	tokenString := jwttool.CreateToken(TokenKey, info)

	return &entity.JwtToken{
		Token: tokenString,
	}
}

// 验证 JWT 令牌，并返回用户信息。
func (l *UserService) CheckToken(token string, user *entity.User) (*entity.UserInfo, error) {
	info, ok := jwttool.ParseToken(token, TokenKey)
	infoMap := info.(jwt.MapClaims)
	if ok {
		expTime := infoMap["exp"].(float64)
		if float64(time.Now().Unix()) >= expTime {
			return nil, fmt.Errorf("%s", "token已过期")
		} else {
			l.DbEngine.Get(user)
			ur := make([]entity.UserRole, 0)
			err := l.DbEngine.Where("user_id = ?", user.Id).Find(&ur)
			if err != nil {
				logger.Error(err.Error())
			}
			roles := make([]string, len(ur))
			for i, m := range ur {
				roles[i] = m.RoleKey
			}
			info := &entity.UserInfo{
				Id: user.Id,
				// Avatar:  user.Avatar,
				Roles:   roles,
				Name:    user.Name,
				Account: user.Account,
				Type:    user.Type,
			}
			return info, nil
		}
	} else {
		return nil, fmt.Errorf("%s", "token无效")
	}
}

func (l *UserService) AddAuth(ur *entity.UserRole) (bool, string) {
	i, err := l.DbEngine.Insert(ur)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "add success"
	}
	return false, "add fail"
}

func (l *UserService) DelAuth(ur *entity.UserRole) (bool, string) {
	i, err := l.DbEngine.Delete(ur)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "del success"
	}
	return false, "del fail"
}

// 检查用户是否具有管理员角色。
func (l *UserService) HasAdminRole(account string) bool {
	user := &entity.User{Account: account}
	_, user = l.GetByUser(user)

	ur := make([]entity.UserRole, 0)
	err := l.DbEngine.Where("user_id = ?", user.Id).Find(&ur)
	if err != nil {
		logger.Error(err.Error())
	}
	for _, m := range ur {
		if m.RoleKey == "admin" {
			return true
		}
	}
	return false
}

func NewUserService(engine *xorm.Engine) *UserService {
	return &UserService{
		DbEngine: engine,
	}
}

func (l *UserService) Register(user *entity.User) (bool, string) {
	// if user.Password != "" {
	// 	user.Password = password.Encode(user.Password, 12, "default")
	// }
	_, err := l.DbEngine.Insert(user)
	if err != nil {
		logger.Error(err.Error())
		return false, err.Error()
	}

	// if user.Org == "farmer" {
	// 	user.Type = 1
	// } else {
	// 	user.Type = 2
	// }

	return true, "success"
}

// FindUserListById
func (f *UserService) FindUserListById(id int) ([]*entity.User, error) {
	user := make([]*entity.User, 0)
	err := f.DbEngine.Find(&user)
	if err != nil {
		logger.Error("FindUserListById failed:", err.Error())
		return nil, err
	}

	resp := make([]*entity.User, 0)
	for k, u := range user {
		if u.Type == -1 {
			continue
		}
		user[k].Img = "http://127.0.0.1:6991/api/file/download?fileId=" + u.Img
		resp = append(resp, u)
	}
	return resp, nil
}

// FindUserListByIds
func (f *UserService) FindUserListByIds(ids []int) ([]*entity.User, error) {
	user := make([]*entity.User, 0)
	err := f.DbEngine.In("id", ids).Find(&user)
	if err != nil {
		logger.Error("FindUserListByIds failed:", err.Error())
		return nil, err
	}
	return user, nil
}

// UpdateUserById
func (f *UserService) UpdateUserById(user *entity.User) error {
	// 根据account查询用户
	userPre := &entity.User{}
	_, err := f.DbEngine.Where("id = ?", user.Id).Get(userPre)
	if err != nil {
		return err
	}
	// if userPre.Type != 0 {
	// 	return errors.New("已分配过角色")
	// }

	fmt.Println("userPre.Type:", userPre.IsOk)
	user.Updated = time.Now().Format("2006-01-02 15:04:05")
	// uid := uuid.NewV4()
	// user.IdCard = uid.String()

	_, err = f.DbEngine.ID(user.Id).Update(user)
	if err != nil {
		logger.Error("UpdateUserById failed:", err.Error())
		return err
	}

	if user.IdCard != "" {
		return nil
	}

	// 根据account查询用户
	_, err = f.DbEngine.Where("id = ?", user.Id).Get(user)
	if err != nil {
		return err
	}

	if user.IsOk != "通过" {
		return nil
	}

	switch user.Type {
	case 1:
		user.Org = "farmer"
	case 2:
		user.Org = "maker"
	case 3:
		user.Org = "dealer"
	case 4:
		user.Org = "consumer"
	case 5:
		user.Org = "regulation"
	case 6:
		user.Org = "student"
	case 7:
		user.Org = "enterprise"
	}

	f.DelAuth(&entity.UserRole{
		UserId: user.Id,
	})

	f.AddAuth(&entity.UserRole{
		UserId:  user.Id,
		RoleKey: user.Org,
	})

	return nil
}

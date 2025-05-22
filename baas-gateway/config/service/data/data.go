package data

import (
	"encoding/json"
	"fmt"
	"strings"

	"data/baas-gateway/blockchian"
	entity "data/baas-gateway/entity"
	data "data/baas-gateway/entity/data"

	"github.com/go-xorm/xorm"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
	"github.com/jonluo94/baasmanager/baas-core/common/log"
	uuid "github.com/satori/go.uuid"

	"data/baas-gateway/service"
)

// 初始化一个日志记录器，用于记录错误和调试信息。
var logger = log.GetLogger("service", log.ERROR)

type DataService struct {
	DbEngine      *xorm.Engine           //初始化一个日志记录器，用于记录错误和调试信息。
	FabircService *service.FabricService //区块链服务，用于与区块链网络交互。
}

// 初始化DataService实例。
func NewDataService(engine *xorm.Engine, fabircService *service.FabricService) *DataService {
	createTable(engine, new(data.EducationalInfor))
	createTable(engine, new(data.EducationalApplication))
	createTable(engine, new(data.Company))
	createTable(engine, new(entity.User))
	createTable(engine, new(entity.UserRole))

	return &DataService{engine, fabircService}
}

// createTable: 检查表是否存在，如果不存在则创建表。
func createTable(engine *xorm.Engine, tabel interface{}) {
	if flag, _ := engine.IsTableExist(tabel); flag {
		logger.Info("table exist")
	} else {
		if err := engine.CreateTables(tabel); err != nil {
			logger.Error("table create failed:", err.Error())
		}
	}
}

// FindEducationalInforList
func (f *DataService) FindEducationalInforList(userId int64, name string, page, size int, adminstatus int) (bool, []*data.EducationalInfor, int64) {
	pager := gintool.CreatePager(page, size)

	resp := make([]*data.EducationalInfor, 0)

	conditions := []string{"1 = 1"}
	args := []interface{}{}

	if name != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+name+"%")
	}
	if userId != 0 {
		conditions = append(conditions, "user_id = ?")
		args = append(args, userId)
	}
	if adminstatus != 0 {
		conditions = append(conditions, "admin_status = ?")
		args = append(args, adminstatus)
	}

	where := strings.Join(conditions, " AND ")

	err := f.DbEngine.Limit(pager.PageSize, pager.NumStart).Where(where, args...).Find(&resp)
	if err != nil {
		logger.Error(err.Error())
		return false, nil, 0
	}

	total, err := f.DbEngine.Where(where, args...).Count(new(data.EducationalInfor))
	if err != nil {
		logger.Error(err.Error())
		return false, nil, 0
	}

	return true, resp, total
}

// GetEducationalInforByUserId
// 根据用户ID查询单个教育信息。
func (f *DataService) GetEducationalInforByUserId(userId int64) (bool, *data.EducationalInfor) {
	resp := new(data.EducationalInfor)
	_, err := f.DbEngine.Where("user_id = ?", userId).Get(resp)
	if err != nil {
		logger.Error(err.Error())
		return false, nil
	}

	return true, resp
}

func (f *DataService) GetEducationalInforByUserTraceCode(traceCode string) (bool, *data.EducationalInfor) {
	resp := new(data.EducationalInfor)
	_, err := f.DbEngine.Where("trace_code = ?", traceCode).Get(resp)
	if err != nil {
		logger.Error(err.Error())
		return false, nil
	}

	return true, resp
}

// UpdateEducationalInforById
func (f *DataService) UpdateEducationalInforById(educationalInfor *data.EducationalInfor) (bool, string) {

	i, err := f.DbEngine.Id(educationalInfor.Id).Update(educationalInfor)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		f.InvokeChainCode(fmt.Sprint("EducationalInfor_", educationalInfor.Id), educationalInfor)
		return true, "update success"
	}
	return false, "update fail"
}

// CreateEducationalInfor
func (f *DataService) CreateEducationalInfor(educationalInfor *data.EducationalInfor) (bool, string) {
	// uuid
	educationalInfor.TraceCode = uuid.NewV4().String()
	i, err := f.DbEngine.Insert(educationalInfor)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		f.InvokeChainCode(fmt.Sprint("EducationalInfor_", educationalInfor.Id), educationalInfor)
		return true, "add success"
	}
	return false, "add fail"
}

// CreateEducationalApplication
func (f *DataService) CreateEducationalApplication(educationalApplication *data.EducationalApplication) (bool, string) {
	i, err := f.DbEngine.Insert(educationalApplication)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		f.InvokeChainCode(fmt.Sprint("EducationalApplication_", educationalApplication.Id), educationalApplication)
		return true, "add success"
	}
	return false, "add fail"
}

// FindEducationalApplicationList
func (f *DataService) FindEducationalApplicationList(name string, userid, studentid int64, page, size int, studentids []int64) (bool, []*data.EducationalApplication, int64) {
	pager := gintool.CreatePager(page, size)

	resp := make([]*data.EducationalApplication, 0)

	conditions := []string{"1 = 1"}
	args := []interface{}{}

	if userid != 0 {
		conditions = append(conditions, "user_id = ?")
		args = append(args, userid)
	}
	if studentid != 0 {
		conditions = append(conditions, "student_id = ?")
		args = append(args, studentid)
	}
	if name != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+name+"%")
	}
	//if len(studentids) > 0 {
	//	conditions = append(conditions, "student_id IN ?")
	//	args = append(args, studentids)
	//}

	where := strings.Join(conditions, " AND ")

	err := f.DbEngine.Limit(pager.PageSize, pager.NumStart).Where(where, args...).Find(&resp)
	if err != nil {
		logger.Error(err.Error())
		return false, nil, 0
	}

	total, err := f.DbEngine.Where(where, args...).Count(new(data.EducationalApplication))
	if err != nil {
		logger.Error(err.Error())
		return false, nil, 0
	}

	return true, resp, total
}

// UpdateEducationalApplicationById
func (f *DataService) UpdateEducationalApplicationById(educationalApplication *data.EducationalApplication) (bool, string) {
	i, err := f.DbEngine.Id(educationalApplication.Id).Update(educationalApplication)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		f.InvokeChainCode(fmt.Sprint("EducationalApplication_", educationalApplication.Id), educationalApplication)
		return true, "update success"
	}
	return false, "update fail"
}

// FindCompanyList
func (f *DataService) FindCompanyList(name string, page, size int) (bool, []*data.Company, int64) {
	pager := gintool.CreatePager(page, size)

	resp := make([]*data.Company, 0)

	conditions := []string{"1 = 1"}
	args := []interface{}{}

	if name != "" {
		conditions = append(conditions, "name LIKE ?")
		args = append(args, "%"+name+"%")
	}

	where := strings.Join(conditions, " AND ")

	err := f.DbEngine.Limit(pager.PageSize, pager.NumStart).Where(where, args...).Find(&resp)
	if err != nil {
		logger.Error(err.Error())
		return false, nil, 0
	}

	total, err := f.DbEngine.Where(where, args...).Count(new(data.Company))
	if err != nil {
		logger.Error(err.Error())
		return false, nil, 0
	}

	return true, resp, total
}

func (f *DataService) FindCompanyListByUserIds(userIds []int64) (bool, []*data.Company) {
	resp := make([]*data.Company, 0)

	err := f.DbEngine.In("user_id", userIds).Find(&resp)
	if err != nil {
		logger.Error(err.Error())
		return false, nil
	}

	return true, resp
}

// GetCompanyByUserId
func (f *DataService) GetCompanyByUserId(userId int64) (bool, *data.Company) {
	resp := new(data.Company)
	_, err := f.DbEngine.Where("user_id = ?", userId).Get(resp)
	if err != nil {
		logger.Error(err.Error())
	}

	return true, resp
}

// UpdateCompanyById
func (f *DataService) UpdateCompanyById(company *data.Company) (bool, string) {
	if company.Id == 0 {
		return f.CreateCompany(company)
	}

	i, err := f.DbEngine.Id(company.Id).Update(company)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		f.InvokeChainCode(fmt.Sprint("Company_", company.Id), company)
		return true, "update success"
	}
	return false, "update fail"
}

// CreateCompany
func (f *DataService) CreateCompany(company *data.Company) (bool, string) {
	i, err := f.DbEngine.Insert(company)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		f.InvokeChainCode(fmt.Sprint("Company_", company.Id), company)
		return true, "add success"
	}
	return false, "add fail"
}

// 调用区块链智能合约，将数据写入区块链。
func (f *DataService) InvokeChainCode(key string, value interface{}) (bool, string) {
	data, _ := json.Marshal(value)
	orgClient := blockchian.New(
		"./config/rock.yaml",
		"Org1",
		"Admin",
		"User1",
	)
	_, err := orgClient.InvokeCC(
		[]string{
			"peer0.org1.example.com",
		},
		"set",
		[]string{
			key,
			string(data),
		},
	)
	if err != nil {
		return false, err.Error()
	}

	return true, ""
}

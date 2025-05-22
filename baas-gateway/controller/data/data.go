package data

import (
	"data/baas-gateway/entity/data"
	"strconv"

	dataService "data/baas-gateway/service/data"

	"github.com/gin-gonic/gin"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
)

type DataControl struct {
	dataService *dataService.DataService
}

func NewDataControl(fs *dataService.DataService) *DataControl {
	return &DataControl{fs}
}

func (f *DataControl) FindEducationalInforListSelf(ctx *gin.Context) {
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

	uid, uType, _ := GetUserIdAndType(ctx)
	if uType == -1 {
		uid = 0
	}

	b, list, total := f.dataService.FindEducationalInforList(uid, name, page, limit, 0)

	ids := []int64{}
	for _, v := range list {
		ids = append(ids, v.UserId)
	}

	if b {
		gintool.ResultList(ctx, list, total)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

// FindEducationalInforList
func (f *DataControl) FindEducationalInforList(ctx *gin.Context) {
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

	uid, _, _ := GetUserIdAndType(ctx)

	b, list, total := f.dataService.FindEducationalInforList(0, name, page, limit, 2)

	ids := []int64{}
	for _, v := range list {
		ids = append(ids, v.UserId)
	}

	_, aps, _ := f.dataService.FindEducationalApplicationList("", uid, 0, 1, 1000, ids)

	for _, v := range list {
		for _, a := range aps {
			if v.UserId == a.StudentId && v.Id == a.EducationalInforId {
				v.ApplicationStatus = a.Status
				if v.ApplicationStatus != 2 {
					v.Gender = ""
					v.IDCard = ""
					v.NativePlace = ""
					v.BirthDate = ""
					v.Ethnicity = ""
					v.EnrollmentDate = ""
					v.GraduationDate = ""
					v.SchoolName = ""
					v.Major = ""
					v.EducationType = ""
					v.StudyMode = ""
					v.Level = ""
					v.Duration = ""
					v.Graduated = ""
					v.CertificateNo = ""
					v.TraceCode = ""
				}
			}
		}
	}

	if len(aps) == 0 {
		for _, v := range list {
			v.Gender = ""
			v.IDCard = ""
			v.NativePlace = ""
			v.BirthDate = ""
			v.Ethnicity = ""
			v.EnrollmentDate = ""
			v.GraduationDate = ""
			v.SchoolName = ""
			v.Major = ""
			v.EducationType = ""
			v.StudyMode = ""
			v.Level = ""
			v.Duration = ""
			v.Graduated = ""
			v.CertificateNo = ""
			v.TraceCode = ""
		}

	}

	if b {
		gintool.ResultList(ctx, list, total)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

// GetEducationalInforByUserId
func (f *DataControl) GetEducationalInforByUserId(ctx *gin.Context) {
	userId, _, _ := GetUserIdAndType(ctx)
	b, resp := f.dataService.GetEducationalInforByUserId(int64(userId))
	if b {
		gintool.ResultList(ctx, resp, 0)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

// UpdateEducationalInforById
func (f *DataControl) UpdateEducationalInforById(ctx *gin.Context) {
	educationalInfor := new(data.EducationalInfor)
	if err := ctx.ShouldBindJSON(educationalInfor); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	userId, userType, _ := GetUserIdAndType(ctx)
	educationalInfor.UserId = int64(userId)
	id := educationalInfor.Id
	adminstatus := educationalInfor.AdminStatus
	if userType == -1 {
		educationalInfor = new(data.EducationalInfor)
		educationalInfor.Id = id
		educationalInfor.AdminStatus = adminstatus
	}

	isSuccess, msg := f.dataService.UpdateEducationalInforById(educationalInfor)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

// CreateEducationalInfor
func (f *DataControl) CreateEducationalInfor(ctx *gin.Context) {
	educationalInfor := new(data.EducationalInfor)
	if err := ctx.ShouldBindJSON(educationalInfor); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	userId, _, _ := GetUserIdAndType(ctx)
	educationalInfor.UserId = int64(userId)
	// _, resp := f.dataService.GetEducationalInforByUserId(int64(userId))
	// if resp.Id != 0 {
	// 	gintool.ResultFail(ctx, "already exists, please to update")
	// 	return
	// }

	isSuccess, msg := f.dataService.CreateEducationalInfor(educationalInfor)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

// FindCompanyList
func (f *DataControl) FindCompanyList(ctx *gin.Context) {
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

	b, list, total := f.dataService.FindCompanyList(name, page, limit)
	if b {
		gintool.ResultList(ctx, list, total)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

// GetCompanyByUserId
func (f *DataControl) GetCompanyByUserId(ctx *gin.Context) {
	userId, _, _ := GetUserIdAndType(ctx)
	b, resp := f.dataService.GetCompanyByUserId(int64(userId))
	if b {
		gintool.ResultList(ctx, resp, 0)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

// UpdateCompanyById
func (f *DataControl) UpdateCompanyById(ctx *gin.Context) {
	company := new(data.Company)
	if err := ctx.ShouldBindJSON(company); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	userId, _, _ := GetUserIdAndType(ctx)
	company.UserId = int64(userId)

	isSuccess, msg := f.dataService.UpdateCompanyById(company)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

// CreateCompany
func (f *DataControl) CreateCompany(ctx *gin.Context) {
	company := new(data.Company)
	if err := ctx.ShouldBindJSON(company); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	userId, _, _ := GetUserIdAndType(ctx)
	company.UserId = int64(userId)

	b, _ := f.dataService.GetCompanyByUserId(int64(userId))
	if b {
		gintool.ResultFail(ctx, "already exists, please to update")
		return
	}

	isSuccess, msg := f.dataService.CreateCompany(company)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

// CreateEducationalApplication
func (f *DataControl) CreateEducationalApplication(ctx *gin.Context) {
	educationalApplication := new(data.EducationalApplication)
	if err := ctx.ShouldBindJSON(educationalApplication); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	userId, _, _ := GetUserIdAndType(ctx)
	educationalApplication.UserId = int64(userId)

	// 查询公司
	_, company := f.dataService.GetCompanyByUserId(int64(userId))
	if company.Id == 0 {
		gintool.ResultFail(ctx, "company not exists")
		return
	}
	educationalApplication.CompanyName = company.Name
	educationalApplication.CompanyId = company.Id

	isSuccess, msg := f.dataService.CreateEducationalApplication(educationalApplication)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func (f *DataControl) FindEducationalApplicationList(ctx *gin.Context) {
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

	userId, _, _ := GetUserIdAndType(ctx)

	b, list, total := f.dataService.FindEducationalApplicationList("", 0, userId, page, limit, nil)

	ids := []int64{}
	for _, v := range list {
		ids = append(ids, v.UserId)
	}

	_, coms := f.dataService.FindCompanyListByUserIds(ids)
	for _, v := range list {
		for _, c := range coms {
			if v.UserId == c.UserId {
				v.CompanyName = c.Name
			}
		}
	}

	if b {
		gintool.ResultList(ctx, list, total)
	} else {
		gintool.ResultFail(ctx, "fail")
	}
}

// UpdateEducationalApplicationById
func (f *DataControl) UpdateEducationalApplicationById(ctx *gin.Context) {
	educationalApplication := new(data.EducationalApplication)
	if err := ctx.ShouldBindJSON(educationalApplication); err != nil {
		gintool.ResultFail(ctx, err)
		return
	}

	isSuccess, msg := f.dataService.UpdateEducationalApplicationById(educationalApplication)
	if isSuccess {
		gintool.ResultMsg(ctx, msg)
	} else {
		gintool.ResultFail(ctx, msg)
	}
}

func GetUserIdAndType(ctx *gin.Context) (int64, int, string) {
	userid, isOk := ctx.Get("userid")
	if !isOk {
		gintool.ResultFail(ctx, "fail")
		return 0, 0, ""
	}
	userId := userid.(int)

	username, isOk := ctx.Get("userName")
	if !isOk {
		gintool.ResultFail(ctx, "fail")
		return 0, 0, ""
	}
	userName := username.(string)

	// 1-农民 2 - 制造商 3 - 经销商 4 -消费者 5 物流
	usertype, isOk := ctx.Get("userType")
	if !isOk {
		gintool.ResultFail(ctx, "fail")
		return 0, 0, ""
	}
	userType := usertype.(int)
	return int64(userId), userType, userName
}

func (f *DataControl) FindTrace(ctx *gin.Context) {
	traceCode := ctx.Query("trace_code")
	isOk, ed := f.dataService.GetEducationalInforByUserTraceCode(traceCode)
	if !isOk || ed.Id == 0 {
		gintool.ResultFail(ctx, "the traceability information does not exist")
	} else {
		gintool.ResultList(ctx, ed, 0)
	}
}

package data

type EducationalInfor struct {
	Id                int64  `json:"id"          xorm:"pk autoincr INT(11)"`
	Name              string `json:"name"           xorm:"VARCHAR(200)"`
	Gender            string `json:"gender"         xorm:"VARCHAR(200)"`
	IDCard            string `json:"id_card"        xorm:"VARCHAR(200)"`
	NativePlace       string `json:"native_place"   xorm:"VARCHAR(200)"`
	BirthDate         string `json:"birth_date"     xorm:"VARCHAR(200)"`
	Ethnicity         string `json:"ethnicity"      xorm:"VARCHAR(200)"`
	EnrollmentDate    string `json:"enrollment_date" xorm:"VARCHAR(200)"`
	GraduationDate    string `json:"graduation_date" xorm:"VARCHAR(200)"`
	SchoolName        string `json:"school_name"    xorm:"VARCHAR(200)"`
	Major             string `json:"major"          xorm:"VARCHAR(200)"`
	EducationType     string `json:"education_type" xorm:"VARCHAR(200)"`
	StudyMode         string `json:"study_mode"     xorm:"VARCHAR(200)"`
	Level             string `json:"level"          xorm:"VARCHAR(200)"`
	Duration          string `json:"duration"       xorm:"VARCHAR(200)"`
	Graduated         string `json:"graduated"      xorm:"VARCHAR(200)"`
	CertificateNo     string `json:"certificate_no" xorm:"VARCHAR(200)"`
	UserId            int64  `json:"user_id"     xorm:"INT(11) unique"`
	TraceCode         string `json:"trace_code" xorm:"VARCHAR(200)"`
	ApplicationStatus int    `json:"status"         xorm:"-"`
	AdminStatus       int    `json:"adminstatus" xorm:"INT(11)"` // 1. 申请中 2. 申请通过 3. 申请拒绝
}

type EducationalApplication struct {
	Id int64 `json:"id"          xorm:"pk autoincr INT(11)"` // 预警ID，自增主键，唯一标识每次预警记录
	// 申请企业id
	CompanyId int64 `json:"company_id" xorm:"INT(11)"`
	// 企业名称
	CompanyName string `json:"company_name" xorm:"-"`
	// 学历信息id
	EducationalInforId int64 `json:"educational_infor_id" xorm:"INT(11)"`
	// 学历名称
	EducationalInforName string `json:"educational_infor_name" xorm:"-"`
	// 申请状态
	Status    int   `json:"status" xorm:"INT(11)"` // 1. 申请中 2. 申请通过 3. 申请拒绝
	StudentId int64 `json:"student_id" xorm:"INT(11)"`
	UserId    int64 `json:"user_id"     xorm:"INT(11)"`
}

// 企业信息
type Company struct {
	Id int64 `json:"id"          xorm:"pk autoincr INT(11)"` // 预警ID，自增主键，唯一标识每次预警记录
	// 企业名称
	Name string `json:"name" xorm:"not null VARCHAR(64)"`
	// 企业类型
	Type string `json:"type" xorm:"VARCHAR(64)"`
	// 企业地址
	Address string `json:"address" xorm:"VARCHAR(64)"`
	// 企业联系方式
	Contact string `json:"contact" xorm:"VARCHAR(64)"`
	// 企业法人
	LegalPerson string `json:"legal_person" xorm:"VARCHAR(64)"`
	// 企业注册资金
	RegisteredCapital string `json:"registered_capital" xorm:"VARCHAR(64)"`
	// 企业经营项目
	BusinessItems string `json:"business_items" xorm:"VARCHAR(64)"`
	// 企业经营状态
	BusinessStatus string `json:"business_status" xorm:"VARCHAR(64)"`
	// 企业注册号
	RegistrationNumber string `json:"registration_number" xorm:"VARCHAR(64)"`
	// 企业统一社会信用代码
	UnifiedSocialCreditCode string `json:"unified_social_credit_code" xorm:"VARCHAR(64)"`
	UserId                  int64  `json:"user_id"     xorm:"INT(11) unique"`
}

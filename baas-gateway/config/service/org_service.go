package service

import (
	"data/baas-gateway/entity"
	"fmt"

	"github.com/go-xorm/xorm"
)

type OrgService struct {
	DbEngine      *xorm.Engine
	FabricService *FabricService
}

// 向数据库中插入一条新的组织记录。
func (l *OrgService) Add(org *entity.Org) (bool, string) {

	//调用 l.DbEngine.Insert(org) 方法，将 org 对象插入到数据库中。
	i, err := l.DbEngine.Insert(org)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "add success"
	}
	return false, "add fail"
}

// 根据组织的名称和所属链删除数据库中的组织记录。
func (l *OrgService) Delete(org *entity.Org) (bool, string) {
	i, err := l.DbEngine.Where("name = ?", org.Name).And("chain = ?", org.Chain).Delete(&entity.Org{})
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "delete success"
	}
	return false, "delete fail"
}

// 根据组织的 ID 删除数据库中的组织记录。
func (l *OrgService) DeleteById(id int) (bool, string) {
	i, err := l.DbEngine.Where("id = ?", id).Delete(&entity.Org{})
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "delete success"
	}
	return false, "delete fail"
}

// 根据组织的名称和所属链更新数据库中的组织记录。
func (l *OrgService) Update(org *entity.Org) (bool, string) {

	i, err := l.DbEngine.Where("name = ?", org.Name).And("chain = ?", org.Chain).Update(org)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "update success"
	}
	return false, "update fail"
}

func (l *OrgService) UpdateById(org *entity.Org) (bool, string) {

	i, err := l.DbEngine.Where("id = ?", org.Id).Update(org)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "update success"
	}
	return false, "update fail"
}

// 根据传入的 org 对象查询数据库中的组织记录。
func (l *OrgService) GetByOrg(org *entity.Org) (bool, *entity.Org) {
	fmt.Printf("GetByOrg org.id = %d\n", org.Id)
	has, err := l.DbEngine.Get(org)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, org
}

func (l *OrgService) GetAllList() (bool, []*entity.Org) {
	orgs := make([]*entity.Org, 0)
	err := l.DbEngine.Where("1=1").Find(&orgs)
	if err != nil {
		logger.Error(err.Error())
	}

	return true, orgs
}

func NewOrgService(engine *xorm.Engine, FabricService *FabricService) *OrgService {
	return &OrgService{
		DbEngine:      engine,
		FabricService: FabricService,
	}
}

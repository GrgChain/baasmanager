package service

import (
	"github.com/go-xorm/xorm"
	"github.com/jonluo94/baasmanager/baas-gateway/entity"
	"github.com/jonluo94/baasmanager/baas-core/common/gintool"
)

type RoleService struct {
	DbEngine *xorm.Engine
}

func (l *RoleService) Add(role *entity.Role) (bool, string) {

	i, err := l.DbEngine.Insert(role)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "add success"
	}
	return false, "add fail"
}

func (l *RoleService) Update(role *entity.Role) (bool, string) {
	i, err := l.DbEngine.Where("rkey = ?", role.Rkey).Update(role)
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "update success"
	}
	return false, "update fail"
}

func (l *RoleService) Delete(key string) (bool, string) {
	i, err := l.DbEngine.Where("rkey = ?", key).Delete(&entity.Role{})
	if err != nil {
		logger.Error(err.Error())
	}
	if i > 0 {
		return true, "delete success"
	}
	return false, "delete fail"
}

func (l *RoleService) GetByRole(role *entity.Role) (bool, *entity.Role) {
	has, err := l.DbEngine.Get(role)
	if err != nil {
		logger.Error(err.Error())
	}
	return has, role
}

func (l *RoleService) GetList(role *entity.Role, page, size int) (bool, []*entity.Role, int64) {

	pager := gintool.CreatePager(page, size)

	roles := make([]*entity.Role, 0)

	values := make([]interface{}, 0)
	where := "1=1"
	if role.Name != "" {
		where += " and name like ? "
		values = append(values, "%"+role.Name+"%")
	}

	err := l.DbEngine.Where(where, values...).Limit(pager.PageSize, pager.NumStart).Find(&roles)
	if err != nil {
		logger.Error(err.Error())
	}
	total, err := l.DbEngine.Where(where, values...).Count(new(entity.Role))
	if err != nil {
		logger.Error(err.Error())
	}

	return true, roles, total
}

func (l *RoleService) GetAll() (bool, []*entity.Role) {

	roles := make([]*entity.Role, 0)

	err := l.DbEngine.Find(&roles)
	if err != nil {
		logger.Error(err.Error())
	}
	return true, roles
}

func NewRoleService(engine *xorm.Engine) *RoleService {
	return &RoleService{
		DbEngine: engine,
	}
}

package service

import (
	"errors"

	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/util"

	"github.com/yockii/quick-system/internal/model"
)

var ApplicationConfigService = new(applicationConfigService)

type applicationConfigService struct{}

func (s *applicationConfigService) Add(instance *model.ApplicationConfig) (isDuplicated bool, success bool, err error) {
	if instance.ApplicationId == "" {
		return false, false, errors.New("应用ID不能为空")
	}
	var c int64 = 0
	c, err = database.DB.Count(&model.ApplicationConfig{
		ApplicationId: instance.ApplicationId,
	})
	if err != nil {
		return
	}
	if c > 0 {
		isDuplicated = true
		return
	}
	instance.Id = model.ApplicationConfigIdPrefix + util.GenerateDatabaseID()
	_, err = database.DB.Insert(instance)
	success = err == nil
	return
}

func (s *applicationConfigService) Remove(instance *model.ApplicationConfig) (bool, error) {
	if instance.Id == "" {
		return false, errors.New("id不能为空")
	}
	c, err := database.DB.Delete(instance)
	if err != nil {
		return false, err
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}

func (s *applicationConfigService) Update(instance *model.ApplicationConfig) (bool, error) {
	if instance.Id == "" {
		return false, errors.New("ID不能为空")
	}
	// 不允许更改的字段

	c, err := database.DB.ID(instance.Id).Update(&model.ApplicationConfig{
		// 允许更改的字段
		Pc:     instance.Pc,
		Mobile: instance.Mobile,
	})
	if err != nil {
		return false, err
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}

func (s *applicationConfigService) Get(instance *model.ApplicationConfig) (*model.ApplicationConfig, error) {
	if instance.Id == "" {
		return nil, errors.New("ID不能为空")
	}
	has, err := database.DB.Get(instance)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, nil
	}
	return instance, nil
}

func (s *applicationConfigService) Paginate(condition *model.ApplicationConfig, limit, offset int, orderBy string) (int, []*model.ApplicationConfig, error) {
	return s.PaginateBetweenTimes(condition, limit, offset, orderBy, nil)
}

func (s *applicationConfigService) PaginateBetweenTimes(condition *model.ApplicationConfig, limit, offset int, orderBy string, tcList map[string]*domain.TimeCondition) (int, []*model.ApplicationConfig, error) {
	// 处理不允许查询的字段

	// 处理sql
	session := database.DB.NewSession()
	if limit > -1 && offset > -1 {
		session.Limit(limit, offset)
	}

	if orderBy != "" {
		session.OrderBy(orderBy)
	}
	session.Desc("create_time")

	// 处理时间字段，在某段时间之间
	for tc, tr := range tcList {
		if tc != "" {
			if !tr.Start.IsZero() && !tr.End.IsZero() {
				session.Where(tc+" between ? and ?", tr.Start, tr.End)
			} else if tr.Start.IsZero() {
				session.Where(tc+" <= ?", tr.End)
			} else if tr.End.IsZero() {
				session.Where(tc+" > ?", tr.Start)
			}
		}
	}

	// 模糊查找

	var list []*model.ApplicationConfig
	total, err := session.FindAndCount(&list, condition)
	if err != nil {
		return 0, nil, err
	}
	return int(total), list, nil
}

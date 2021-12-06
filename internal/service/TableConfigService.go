package service

import (
	"errors"

	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/util"

	"github.com/yockii/quick-system/internal/model"
)

var TableConfigService = new(tableConfigService)

type tableConfigService struct{}

func (s *tableConfigService) Add(instance *model.TableConfig) (isDuplicated bool, success bool, err error) {
	if instance.ApplicationId == "" {
		return false, false, errors.New("应用ID不能为空")
	}
	if instance.TableName == "" {
		return false, false, errors.New("表名不能为空")
	}
	var c int64 = 0
	c, err = database.DB.Count(&model.TableConfig{
		ApplicationId: instance.ApplicationId,
		TableName:     instance.TableName,
	})
	if err != nil {
		return
	}
	if c > 0 {
		isDuplicated = true
		return
	}

	instance.Id = model.TableConfigIdPrefix + util.GenerateDatabaseID()
	if instance.RecordType == 0 {
		instance.RecordType = 1
	}
	_, err = database.DB.Insert(instance)
	success = err == nil
	return
}

func (s *tableConfigService) Remove(instance *model.TableConfig) (bool, error) {
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

func (s *tableConfigService) Update(instance *model.TableConfig) (bool, error) {
	if instance.Id == "" {
		return false, errors.New("ID不能为空")
	}
	// 不允许更改的字段

	c, err := database.DB.ID(instance.Id).Update(&model.TableConfig{
		// 允许更改的字段
		TableName:    instance.TableName,
		TableComment: instance.TableComment,
		RecordType:   instance.RecordType,
	})
	if err != nil {
		return false, err
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}

func (s *tableConfigService) Get(instance *model.TableConfig) (*model.TableConfig, error) {
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

func (s *tableConfigService) Paginate(condition *model.TableConfig, limit, offset int, orderBy string) (int, []*model.TableConfig, error) {
	return s.PaginateBetweenTimes(condition, limit, offset, orderBy, nil)
}

func (s *tableConfigService) PaginateBetweenTimes(condition *model.TableConfig, limit, offset int, orderBy string, tcList map[string]*domain.TimeCondition) (int, []*model.TableConfig, error) {
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
	if condition.TableName != "" {
		session.Where("table_name like ?", condition.TableName+"%")
		condition.TableName = ""
	}
	if condition.TableComment != "" {
		session.Where("table_comment like ?", condition.TableComment+"%")
		condition.TableComment = ""
	}
	var list []*model.TableConfig
	total, err := session.FindAndCount(&list, condition)
	if err != nil {
		return 0, nil, err
	}
	return int(total), list, nil
}

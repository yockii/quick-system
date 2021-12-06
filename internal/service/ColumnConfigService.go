package service

import (
	"errors"

	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/util"

	"github.com/yockii/quick-system/internal/model"
)

var ColumnConfigService = new(columnConfigService)

type columnConfigService struct{}

func (s *columnConfigService) Add(instance *model.ColumnConfig) (isDuplicated bool, success bool, err error) {
	if instance.ApplicationId == "" {
		return false, false, errors.New("应用ID不能为空")
	}
	if instance.TableId == "" {
		return false, false, errors.New("表ID不能为空")
	}
	if instance.ColumnName == "" {
		return false, false, errors.New("字段名不能为空")
	}
	var c int64 = 0
	c, err = database.DB.Count(&model.ColumnConfig{
		ApplicationId: instance.ApplicationId,
		TableId:       instance.TableId,
		ColumnName:    instance.ColumnName,
	})
	if err != nil {
		return
	}
	if c > 0 {
		isDuplicated = true
		return
	}

	instance.Id = model.ColumnConfigIdPrefix + util.GenerateDatabaseID()
	if instance.ColumnType == 0 {
		instance.ColumnType = 1
	}
	if instance.UpdateType == 0 {
		instance.UpdateType = 7
	}
	if instance.StringType == 0 {
		instance.StringType = 1
	}
	if instance.StringSearch == 0 {
		instance.StringSearch = 1
	}
	_, err = database.DB.Insert(instance)
	success = err == nil
	return
}

func (s *columnConfigService) Remove(instance *model.ColumnConfig) (bool, error) {
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

func (s *columnConfigService) Update(instance *model.ColumnConfig) (bool, error) {
	if instance.Id == "" {
		return false, errors.New("ID不能为空")
	}
	// 不允许更改的字段

	c, err := database.DB.ID(instance.Id).Update(&model.ColumnConfig{
		// 允许更改的字段
		ApplicationId: instance.ApplicationId,
		TableId:       instance.TableId,
		ColumnName:    instance.ColumnName,
		ColumnComment: instance.ColumnComment,
		ColumnType:    instance.ColumnType,
		UpdateType:    instance.UpdateType,
		UpdateAlone:   instance.UpdateAlone,
		ZeroValue:     instance.ZeroValue,
		UniqueCheck:   instance.UniqueCheck,
		StringType:    instance.StringType,
		StringSearch:  instance.StringSearch,
		EnumJson:      instance.EnumJson,
		ColumnLength:  instance.ColumnLength,
		DecimalLength: instance.DecimalLength,
	})
	if err != nil {
		return false, err
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}

func (s *columnConfigService) Get(instance *model.ColumnConfig) (*model.ColumnConfig, error) {
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

func (s *columnConfigService) Paginate(condition *model.ColumnConfig, limit, offset int, orderBy string) (int, []*model.ColumnConfig, error) {
	return s.PaginateBetweenTimes(condition, limit, offset, orderBy, nil)
}

func (s *columnConfigService) PaginateBetweenTimes(condition *model.ColumnConfig, limit, offset int, orderBy string, tcList map[string]*domain.TimeCondition) (int, []*model.ColumnConfig, error) {
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
	if condition.ColumnName != "" {
		session.Where("column_name like ?", condition.ColumnName+"%")
		condition.ColumnName = ""
	}
	if condition.ColumnComment != "" {
		session.Where("column_comment like ?", condition.ColumnComment+"%")
		condition.ColumnComment = ""
	}
	if condition.ZeroValue != "" {
		session.Where("zero_value like ?", condition.ZeroValue+"%")
		condition.ZeroValue = ""
	}
	if condition.EnumJson != "" {
		session.Where("enum_json like ?", condition.EnumJson+"%")
		condition.EnumJson = ""
	}
	var list []*model.ColumnConfig
	total, err := session.FindAndCount(&list, condition)
	if err != nil {
		return 0, nil, err
	}
	return int(total), list, nil
}

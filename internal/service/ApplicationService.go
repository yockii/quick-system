package service

import (
	"errors"
	"fmt"

	"github.com/yockii/qscore/pkg/database"
	"github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/logger"
	"github.com/yockii/qscore/pkg/util"

	gDomain "github.com/yockii/qs-code-generator/pkg/domain"
	"github.com/yockii/qs-code-generator/pkg/generator"

	"github.com/yockii/quick-system/internal/model"
)

var ApplicationService = new(applicationService)

type applicationService struct{}

func (s *applicationService) Add(instance *model.Application) (isDuplicated bool, success bool, err error) {
	if instance.AppName == "" {
		return false, false, errors.New("字典键名不能为空")
	}
	var c int64 = 0
	c, err = database.DB.Count(&model.Application{
		AppName: instance.AppName,
	})
	if err != nil {
		return
	}
	if c > 0 {
		isDuplicated = true
		return
	}
	instance.Id = model.ApplicationIdPrefix + util.GenerateDatabaseID()
	_, err = database.DB.Insert(instance)
	success = err == nil
	return
}

func (s *applicationService) Remove(instance *model.Application) (bool, error) {
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

func (s *applicationService) Update(instance *model.Application) (bool, error) {
	if instance.Id == "" {
		return false, errors.New("ID不能为空")
	}
	// 不允许更改的字段

	c, err := database.DB.ID(instance.Id).Update(&model.Application{
		// 允许更改的字段
		AppName: instance.AppName,
		AppDesc: instance.AppDesc,
		OwnerId: instance.OwnerId,
	})
	if err != nil {
		return false, err
	}
	if c == 0 {
		return false, nil
	}
	return true, nil
}

func (s *applicationService) Get(instance *model.Application) (*model.Application, error) {
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

func (s *applicationService) Paginate(condition *model.Application, limit, offset int, orderBy string) (int, []*model.Application, error) {
	return s.PaginateBetweenTimes(condition, limit, offset, orderBy, nil)
}

func (s *applicationService) PaginateBetweenTimes(condition *model.Application, limit, offset int, orderBy string, tcList map[string]*domain.TimeCondition) (int, []*model.Application, error) {
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
	if condition.AppName != "" {
		session.Where("app_name like ?", condition.AppName+"%")
		condition.AppName = ""
	}
	if condition.AppDesc != "" {
		session.Where("app_desc like ?", condition.AppDesc+"%")
		condition.AppDesc = ""
	}
	var list []*model.Application
	total, err := session.FindAndCount(&list, condition)
	if err != nil {
		return 0, nil, err
	}
	return int(total), list, nil
}

func (s *applicationService) GenerateCode(id string) (bool, error) {
	if id == "" {
		return false, errors.New("ID不能为空")
	}
	application := new(model.Application)
	if exist, err := database.DB.ID(id).Get(application); err != nil {
		return false, err
	} else if !exist {
		return false, errors.New("ID所指向的应用不存在")
	}
	app := new(gDomain.Application)
	app.Package = application.Package
	var tables []*model.TableConfig
	if err := database.DB.Find(&tables, &model.TableConfig{ApplicationId: id}); err != nil {
		return false, err
	}
	var gtables []*gDomain.Table
	for _, table := range tables {
		gtable := &gDomain.Table{
			Name:             table.TableName,
			RecordCreateTime: table.RecordType&1 == 1,
			RecordUpdateTime: table.RecordType&2 == 2,
			RecordDeleteTime: table.RecordType&4 == 4,
		}
		uniqueCheckColumns := make(map[int][]string)
		var cs []*gDomain.Column
		var columns []*model.ColumnConfig
		if err := database.DB.Find(&columns, &model.ColumnConfig{TableId: table.Id}); err != nil {
			return false, err
		}
		for _, column := range columns {
			c := &gDomain.Column{
				Name:        column.ColumnName,
				DisplayName: column.DisplayName,
				Type:        column.ColumnType,
				Comment:     column.ColumnComment,
				Nullable:    column.ZeroValue != "!NIL",
				Updatable:   column.UpdateType&2 == 2,
				Searchable:  column.UpdateType&4 == 4,
			}
			if c.Type == 0 {
				c.Type = 1
			}
			if c.Type == 1 {
				if column.StringType == 2 {
					c.ColumnType = "longtext"
				} else if column.ColumnLength > 0 {
					c.ColumnType = fmt.Sprintf("varchar(%d)", column.ColumnLength)
				}
			} else if c.Type == 2 {
				if column.DecimalLength > 0 {
					c.ColumnType = fmt.Sprintf("decimal(%d,%d)", column.ColumnLength, column.DecimalLength)
				} else if column.ColumnLength > 0 {
					c.ColumnType = fmt.Sprintf("int(%d)", column.ColumnLength)
				}
			}
			cs = append(cs, c)
			if column.UniqueCheck > 0 {
				_, ok := uniqueCheckColumns[column.UniqueCheck]
				if !ok {
					uniqueCheckColumns[column.UniqueCheck] = []string{column.ColumnName}
				} else {
					uniqueCheckColumns[column.UniqueCheck] = append(uniqueCheckColumns[column.UniqueCheck], column.ColumnName)
				}
			}
		}
		gtable.Columns = cs
		for _, v := range uniqueCheckColumns {
			gtable.UniqueCheckColumnNames = append(gtable.UniqueCheckColumnNames, v)
		}
		gtables = append(gtables, gtable)
	}
	app.Tables = gtables

	bs, err := generator.GenerateApplicationSource(app)
	if err != nil {
		return false, err
	}
	if bs != nil {
		logger.Debug("代码生成成功!")
		// 入库
		_, err = database.DB.Insert(&model.ApplicationSource{
			Id:            util.GenerateDatabaseID(),
			ApplicationId: id,
			Source:        bs,
		})
		if err != nil {
			logger.Error(err)
		}
	}
	return true, nil
}

package controller

import (
	"github.com/gofiber/fiber/v2"
	"github.com/yockii/qscore/pkg/constant"
	"github.com/yockii/qscore/pkg/domain"
	"github.com/yockii/qscore/pkg/logger"

	"github.com/yockii/quick-system/internal/model"
	"github.com/yockii/quick-system/internal/service"
)

var ApplicationController = new(applicationController)

type applicationController struct{}

func (c *applicationController) Add(ctx *fiber.Ctx) error {
	instance := new(model.Application)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}

	// 处理必填
	if instance.AppName == "" || instance.Package == "" {
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "应用名/包名必须提供",
		})
	}

	if instance.OwnerId == "" {
		uidPtr := ctx.Locals("userId")
		if uidPtr != nil {
			instance.OwnerId = uidPtr.(string)
		}
	}
	if instance.OwnerId == "" {
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeNotFound,
			Msg:  "所属用户未能获取到",
		})
	}

	duplicated, success, err := service.ApplicationService.Add(instance)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	if duplicated {
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeDuplicate,
			Msg:  "有重复记录",
		})
	}
	if success {
		return ctx.JSON(&domain.CommonResponse{Data: instance})
	}
	return ctx.JSON(&domain.CommonResponse{
		Code: constant.ErrorCodeUnknown,
		Msg:  "服务出现异常",
	})
}

func (c *applicationController) Delete(ctx *fiber.Ctx) error {
	instance := new(model.Application)
	if err := ctx.QueryParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	if instance.Id == "" {
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "ID必须提供",
		})
	}
	deleted, err := service.ApplicationService.Remove(instance)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	if deleted {
		return ctx.JSON(&domain.CommonResponse{})
	}
	return ctx.JSON(&domain.CommonResponse{
		Msg:  "无数据被删除",
		Data: false,
	})
}

func (c *applicationController) Update(ctx *fiber.Ctx) error {
	instance := new(model.Application)
	if err := ctx.BodyParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	if instance.Id == "" {
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "ID必须提供",
		})
	}
	updated, err := service.ApplicationService.Update(instance)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	if updated {
		return ctx.JSON(&domain.CommonResponse{})
	}
	return ctx.JSON(&domain.CommonResponse{
		Msg:  "无数据被更新",
		Data: false,
	})
}

func (c *applicationController) Paginate(ctx *fiber.Ctx) error {
	pr := new(model.ApplicationRequest)
	if err := ctx.QueryParser(pr); err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	limit, offset, orderBy, err := parsePaginationInfoFromQuery(ctx)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}

	timeRangeMap := make(map[string]*domain.TimeCondition)
	//if pr.CreateTimeRange != nil {
	//	timeRangeMap["create_time"] = &domain.TimeCondition{
	//		Start: pr.CreateTimeRange.Start,
	//		End:   pr.CreateTimeRange.End,
	//	}
	//}

	total, list, err := service.ApplicationService.PaginateBetweenTimes(pr.Application, limit, offset, orderBy, timeRangeMap)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	return ctx.JSON(&domain.CommonResponse{Data: &domain.Paginate{
		Total:  total,
		Offset: offset,
		Limit:  limit,
		Items:  list,
	}})
}

func (c *applicationController) Get(ctx *fiber.Ctx) error {
	instance := new(model.Application)
	var err error
	if err = ctx.QueryParser(instance); err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeBodyParse,
			Msg:  "参数解析失败!",
		})
	}
	instance, err = service.ApplicationService.Get(instance)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "服务出现异常",
		})
	}
	return ctx.JSON(&domain.CommonResponse{
		Data: instance,
	})
}

func (c *applicationController) GenerateCode(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	if id == "" {
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeLackOfField,
			Msg:  "ID必须提供",
		})
	}
	ok, err := service.ApplicationService.GenerateCode(id)
	if err != nil {
		logger.Error(err)
		return ctx.JSON(&domain.CommonResponse{
			Code: constant.ErrorCodeService,
			Msg:  "代码生成出错!",
		})
	}
	return ctx.JSON(&domain.CommonResponse{Data: ok})
}

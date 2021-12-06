package model

import (
	"github.com/yockii/qscore/pkg/domain"
)

const (
	ApplicationIdPrefix       = "application"
	ApplicationConfigIdPrefix = "applicationConfig"
)

type Application struct {
	Id         string          `json:"id,omitempty" xorm:"pk varchar(50)"`
	AppName    string          `json:"appName,omitempty" xorm:"comment('应用名称')"`
	Package    string          `json:"package,omitempty" xorm:"comment('应用包名')"`
	AppDesc    string          `json:"appDesc,omitempty" xorm:"varchar(500) comment('应用说明')"`
	OwnerId    string          `json:"ownerId,omitempty" xorm:"varchar(50) comment('创建人/所有人ID')"`
	CreateTime domain.DateTime `json:"createTime" xorm:"created"`
}

type ApplicationSource struct {
	Id            string          `json:"id,omitempty" xorm:"pk varchar(50)"`
	ApplicationId string          `json:"applicationId,omitempty" xorm:"index varchar(50)"`
	Source        []byte          `json:"source,omitempty" xorm:"blob"`
	CreateTime    domain.DateTime `json:"createTime" xorm:"created"`
}

type ApplicationConfig struct {
	Id               string `json:"id,omitempty" xorm:"pk varchar(50)"`
	ApplicationId    string `json:"applicationId,omitempty" xorm:"index varchar(50)"`
	Pc               int    `json:"pc,omitempty" xorm:"int(1) comment('生成PC界面，0未配置 1生成 2不生成')"`
	Mobile           int    `json:"mobile,omitempty" xorm:"int(1) comment('生成手机端界面，0未配置 1生成 2不生成')"`
	TokenExpireHours int    `json:"tokenExpireHours,omitempty" xorm:"comment('token失效时长')"`
}

func init() {
	SyncModels = append(SyncModels, Application{}, ApplicationConfig{}, ApplicationSource{})
}

type ApplicationRequest struct {
	*Application
	CreateTimeRange *domain.TimeCondition `json:"createTimeRange,omitempty"`
}

type ApplicationConfigRequest struct {
	*ApplicationConfig
}

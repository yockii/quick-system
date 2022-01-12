package model

const (
	TableConfigIdPrefix  = "tableConfig"
	ColumnConfigIdPrefix = "columnConfig"
)

type TableConfig struct {
	Id            string `json:"id,omitempty" xorm:"pk varchar(50)"`
	ApplicationId string `json:"applicationId,omitempty" xorm:"index varchar(50)"`
	TableName     string `json:"tableName,omitempty" xorm:"varchar(50) comment('表名')"`
	TableComment  string `json:"tableComment,omitempty" xorm:"comment('表说明')"`
	RecordType    int    `json:"recordType,omitempty" xorm:"comment('表记录类型 0-未设置 1-记录创建时间 2-记录更新时间 4-记录删除时间')"`
}

type ColumnConfig struct {
	Id            string `json:"id,omitempty" xorm:"pk varchar(50)"`
	ApplicationId string `json:"applicationId,omitempty" xorm:"index varchar(50)"`
	TableId       string `json:"tableId,omitempty" xorm:"index varchar(50)"`
	ColumnName    string `json:"columnName,omitempty" xorm:"varchar(50) comment('字段名')"`
	DisplayName   string `json:"displayName,omitempty" xorm:"comment('字段显示名')"`
	ColumnComment string `json:"columnComment,omitempty" xorm:"comment('字段说明')"`
	DisplayType   int    `json:"displayType,omitempty" xorm:"comment('显示类型 0-未设置 1-添加显示 2-更新显示 4-列表显示 8-详情显示')"`
	ColumnType    int    `json:"columnType,omitempty" xorm:"comment('字段类型 0-未知 1-string 2-int 3-DateTime 4-decimal')"`
	UpdateType    int    `json:"updateType,omitempty" xorm:"comment('字段更新方式 0-未设置 1-允许新增 2-允许更改 4-允许作为查询条件')"`
	UpdateAlone   int    `json:"updateAlone,omitempty" xorm:"comment('独立更改 0-未设置 1-只能独立更改 将会生成单独的接口进行更改')"`
	ZeroValue     string `json:"zeroValue,omitempty" xorm:"comment('空时默认值 !NIL-表示不允许为空，如有其他值表示为空时采用该默认值')"`
	UniqueCheck   int    `json:"uniqueCheck,omitempty" xorm:"comment('参与唯一性校验，非0值的字段将都参与唯一性校验，相同的值一起参与唯一性')"`
	StringType    int    `json:"stringType,omitempty" xorm:"comment('字符串数据库存储类型 0-未设置 1-varchar 2-longtext')"`
	StringSearch  int    `json:"stringSearch,omitempty" xorm:"comment('字符串类搜索方式 0-未设置 1-开头模糊匹配 2-全量模糊匹配 3-精确匹配')"`
	EnumJson      string `json:"enumJson,omitempty" xorm:"comment('枚举值，如果是字符串类，以,分割，如果是int，则需存入json，格式为[{key:1,value:''}]')"`
	ColumnLength  int    `json:"columnLength,omitempty" xorm:"comment('字段长度')"`
	DecimalLength int    `json:"decimalLength,omitempty" xorm:"comment('小数部分长度')"`
}

func init() {
	SyncModels = append(SyncModels, TableConfig{}, ColumnConfig{})
}

type TableConfigRequest struct {
	*TableConfig
}
type ColumnConfigRequest struct {
	*ColumnConfig
}

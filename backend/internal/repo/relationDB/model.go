package relationDB

type Example struct {
	ID    int64 `gorm:"column:id;type:bigint;primary_key;AUTO_INCREMENT"`    // id编号
	alias int64 `gorm:"column:alias;type:bigint;primary_key;AUTO_INCREMENT"` // 别名
}

func (m *Example) TableName() string {
	return "example"
}

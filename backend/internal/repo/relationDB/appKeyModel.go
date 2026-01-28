package relationDB

import "time"

const TableNameAppKey = "resource_app_key"

// AppKeyModel
type AppKeyModel struct {
	ID             int64     `gorm:"column:id;primaryKey" json:"id"`
	AppSecretKey   string    `gorm:"column:app_secret_key;not null" json:"appSecretKey"`
	AppSecretValue string    `gorm:"column:app_secret_value;not null" json:"appSecretValue"`
	Status         int32     `gorm:"column:status;not null" json:"status"`
	CreateTime     time.Time `gorm:"column:create_time;default:now()" json:"createTime"`
}

// TableName 返回表名
func (AppKeyModel) TableName() string {
	return TableNameAppKey
}

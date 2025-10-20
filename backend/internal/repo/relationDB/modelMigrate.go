package relationDB

import (
	"context"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
)

var NeedInitColumn bool

func Migrate(c conf.Database, schema string) error {
	// return nil
	//if c.IsInitTable == false {
	//	return nil
	//}
	db := stores.GetCommonConn(context.TODO())
	if !db.Migrator().HasTable(&UnsNamespace{}) {
		//需要初始化表
		NeedInitColumn = true
	}
	err := db.AutoMigrate(
		// &UnsNamespace{},
		&UnsLabel{},
	)
	if c.DBType == "pgsql" && schema != "" {
		db.Exec("create schema if not exists " + schema)
		db.Exec("SET search_path TO " + schema)
	}
	if err != nil {
		return err
	}

	return err
}

func migrateTableColumn() error {
	return nil
}

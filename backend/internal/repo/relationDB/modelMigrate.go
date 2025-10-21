package relationDB

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"strings"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
)

var NeedInitColumn bool

//go:embed migrations_sqls/*
var sqlFiles embed.FS

func Migrate(c conf.Database, schema string) error {
	// return nil
	//if c.IsInitTable == false {
	//	return nil
	//}
	db := stores.GetCommonConn(context.TODO())
	fs.WalkDir(sqlFiles, "migrations_sqls", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".sql") && !d.IsDir() {
			if bs, er := sqlFiles.ReadFile(path); er == nil && len(bs) > 0 {
				sqlContent := strings.Split(string(bs), ";")
				if len(sqlContent) > 0 {
					for _, sql := range sqlContent {
						err := db.Exec(sql).Error
						if err != nil {
							log.Println("SQL WARN:", err, sql)
						} else {
							log.Println("SQL: ", strings.TrimSpace(sql))
						}
					}
				}
			} else {
				log.Println("NOT FOUND: ", path, er)
			}
		}
		return nil
	})
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

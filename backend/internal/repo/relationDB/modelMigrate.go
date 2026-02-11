package relationDB

import (
	"context"
	"embed"
	"io/fs"
	"log"
	"strings"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"gorm.io/gorm"
)

var NeedInitColumn bool

//go:embed migrations_sqls/*
var sqlFiles embed.FS
var dbConfig conf.Database

func InitDbConfig(c conf.Database) {
	dbConfig = c
}
func Migrate(c conf.Database) error {
	dbConfig = c
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
	addPrimaryKeyForFlowUnsLink(db)
	var unsDao UnsNamespaceRepo
	unsDao.migrate(db)
	if !db.Migrator().HasTable(&UnsNamespace{}) {
		//需要初始化表
		NeedInitColumn = true
	}
	err := db.AutoMigrate(
		// &UnsNamespace{},
		&UnsLabel{},
	)

	return err
}
func addPrimaryKeyForFlowUnsLink(db *gorm.DB) {
	sql := `-- 检查是否已存在主键
BEGIN;
DO $$
DECLARE
    has_primary_key BOOLEAN;
BEGIN
    -- 检查主键是否存在
    SELECT EXISTS (
        SELECT 1
        FROM pg_constraint pc
        JOIN pg_class c ON c.oid = pc.conrelid
        JOIN pg_namespace n ON n.oid = c.relnamespace
        WHERE c.relname = 'supos_node_flow_models'
          AND n.nspname = current_schema()
          AND pc.contype = 'p'
    ) INTO has_primary_key;
    
    -- 如果没有主键，则执行修改
    IF NOT has_primary_key THEN
        -- 1. 更新现有数据
        UPDATE supos_node_flow_models 
        SET alias = '' 
        WHERE alias IS NULL;

        -- 使用临时表记录需要删除的行
       CREATE TEMP TABLE duplicate_rows_to_delete AS
        SELECT "ctid" as dupid
        FROM (
            SELECT "ctid",
                   ROW_NUMBER() OVER (
                       PARTITION BY parent_id, alias 
                       ORDER BY create_time
                   ) as row_num
            FROM supos_node_flow_models
        ) ranked
        WHERE row_num > 1;
        -- 删除重复行
        DELETE FROM supos_node_flow_models
        WHERE ctid IN (SELECT dupid FROM duplicate_rows_to_delete);
        
        -- 删除临时表
        DROP TABLE duplicate_rows_to_delete;
        
        -- 2. 修改列约束
        ALTER TABLE supos_node_flow_models 
        ALTER COLUMN alias SET NOT NULL,
        ALTER COLUMN alias SET DEFAULT '';
        
        -- 3. 添加联合主键
        ALTER TABLE supos_node_flow_models 
        ADD PRIMARY KEY (parent_id, alias);
        
        RAISE NOTICE '表结构已修改完成';
    ELSE
        RAISE NOTICE '表已有主键，跳过修改';
    END IF;
END $$;
-- 提交事务
COMMIT;
`
	err := db.Exec(sql).Error
	if err != nil {
		log.Println("SQL WARN:", err, sql)
	}
}
func migrateTableColumn() error {
	return nil
}

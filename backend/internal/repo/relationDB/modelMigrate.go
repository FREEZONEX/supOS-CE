package relationDB

import (
	"backend/internal/common/enums"
	"context"
	"embed"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"

	"gitee.com/unitedrhino/share/conf"
	"gitee.com/unitedrhino/share/stores"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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
	fs.WalkDir(sqlFiles, "migrations_sqls", func(path string, d fs.DirEntry, err error) error {
		if strings.HasSuffix(path, ".sql") && !d.IsDir() {
			if bs, er := sqlFiles.ReadFile(path); er == nil && len(bs) > 0 {
				sqlContent := strings.Split(string(bs), ";")
				if len(sqlContent) > 0 {
					for _, sql := range sqlContent {
						ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
						db := stores.GetCommonConn(ctx)
						err := db.Exec(sql).Error
						cancel()
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
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	db := stores.GetCommonConn(ctx)
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
		&IamUser{},
		&IamRole{},
		&IamUserRole{},
		&IamRoleResource{},
		&IamSession{},
		&IamOAuthClient{},
		&IamOAuthAuthorizationCode{},
		&IamOAuthAccessToken{},
	)
	if err == nil {
		seedIAMDefaults(db)
		seedSuposDefaultResources(db)
		seedSuperAdminRoleResources(db)
	}
	cancel()

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

func seedIAMDefaults(db *gorm.DB) {
	if db == nil {
		return
	}
	defaultRoles := []IamRole{
		{
			ID:          enums.RoleSuperAdmin.ID,
			RoleKey:     "admin",
			RoleName:    "admin",
			Description: "builtin admin role",
			Builtin:     true,
			Status:      1,
		},
		{
			ID:          enums.RoleNormalUser.ID,
			RoleKey:     "user",
			RoleName:    "user",
			Description: "builtin user role",
			Builtin:     true,
			Status:      1,
		},
	}
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"role_key", "role_name", "description", "builtin", "status", "updated_at"}),
	}).Create(&defaultRoles).Error; err != nil {
		log.Println("seed supos roles WARN:", err)
	}

	defaultClient := IamOAuthClient{
		ClientID:     firstNonEmptyEnv("IAM_PORTAINER_CLIENT_ID", "portainer"),
		ClientSecret: firstNonEmptyEnv("IAM_PORTAINER_CLIENT_SECRET", "Tier0PortainerSecret@1304"),
		ClientName:   "portainer",
		Scopes:       "openid",
		GrantTypes:   "authorization_code",
		Enabled:      true,
		Trusted:      true,
		RedirectURIs: "",
	}
	if err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "client_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"client_secret", "client_name", "scopes", "grant_types", "enabled", "trusted", "updated_at"}),
	}).Create(&defaultClient).Error; err != nil {
		log.Println("seed supos oauth client WARN:", err)
	}

	seedIAMBootstrapUser(db)
}

func seedSuposDefaultResources(db *gorm.DB) {
	if db == nil {
		return
	}

	now := time.Now()
	resources := []SuposResource{
		{
			ID:          50,
			Type:        1,
			Source:      stringPtr("platform"),
			Code:        "menu.tag.devtools",
			NameCode:    stringPtr("menu.tag.devtools"),
			RouteSource: intPtr(2),
			URLType:     intPtr(1),
			OpenType:    intPtr(0),
			Icon:        stringPtr("menu.tag.devtools.svg"),
			Sort:        intPtr(3),
			EditEnable:  boolPtr(true),
			HomeEnable:  boolPtr(true),
			Fixed:       boolPtr(false),
			Enable:      boolPtr(true),
			UpdateAt:    now,
			CreateAt:    now,
		},
		{
			ID:          60,
			Type:        1,
			Source:      stringPtr("platform"),
			Code:        "menu.tag.system",
			NameCode:    stringPtr("menu.tag.system"),
			RouteSource: intPtr(2),
			URLType:     intPtr(1),
			OpenType:    intPtr(0),
			Icon:        stringPtr("menu.tag.system.svg"),
			Sort:        intPtr(5),
			EditEnable:  boolPtr(false),
			HomeEnable:  boolPtr(true),
			Fixed:       boolPtr(true),
			Enable:      boolPtr(true),
			UpdateAt:    now,
			CreateAt:    now,
		},
		{
			ID:              61,
			ParentID:        int64Ptr(60),
			Type:            2,
			Source:          stringPtr("platform"),
			Code:            "RoutingManagement",
			NameCode:        stringPtr("RoutingManagement"),
			RouteSource:     intPtr(2),
			URL:             stringPtr("/konga/home/"),
			URLType:         intPtr(2),
			OpenType:        intPtr(0),
			Icon:            stringPtr("RoutingManagement.svg"),
			DescriptionCode: stringPtr("menu.desc.konga"),
			Sort:            intPtr(1),
			EditEnable:      boolPtr(false),
			HomeEnable:      boolPtr(true),
			Fixed:           boolPtr(true),
			Enable:          boolPtr(true),
			UpdateAt:        now,
			CreateAt:        now,
		},
		{
			ID:              63,
			ParentID:        int64Ptr(60),
			Type:            2,
			Source:          stringPtr("platform"),
			Code:            "UserManagement",
			NameCode:        stringPtr("UserManagement"),
			RouteSource:     intPtr(2),
			URL:             stringPtr("/account-management"),
			URLType:         intPtr(1),
			OpenType:        intPtr(0),
			Icon:            stringPtr("UserManagement.svg"),
			DescriptionCode: stringPtr("menu.desc.account"),
			Sort:            intPtr(3),
			EditEnable:      boolPtr(false),
			HomeEnable:      boolPtr(true),
			Fixed:           boolPtr(true),
			Enable:          boolPtr(true),
			UpdateAt:        now,
			CreateAt:        now,
		},
		{
			ID:              65,
			ParentID:        int64Ptr(50),
			Type:            2,
			Source:          stringPtr("platform"),
			Code:            "ContainerManagement",
			NameCode:        stringPtr("ContainerManagement"),
			RouteSource:     intPtr(2),
			URL:             stringPtr("/portainer/home/"),
			URLType:         intPtr(2),
			OpenType:        intPtr(0),
			Icon:            stringPtr("ContainerManagement.svg"),
			DescriptionCode: stringPtr("menu.desc.dockerMgmt"),
			Sort:            intPtr(5),
			EditEnable:      boolPtr(true),
			HomeEnable:      boolPtr(true),
			Fixed:           boolPtr(false),
			Enable:          boolPtr(true),
			UpdateAt:        now,
			CreateAt:        now,
		},
		{
			ID:              66,
			ParentID:        int64Ptr(60),
			Type:            2,
			Source:          stringPtr("platform"),
			Code:            "AboutUs",
			NameCode:        stringPtr("AboutUs"),
			RouteSource:     intPtr(2),
			URL:             stringPtr("/aboutus"),
			URLType:         intPtr(1),
			OpenType:        intPtr(0),
			Icon:            stringPtr("AboutUs.svg"),
			DescriptionCode: stringPtr("menu.desc.aboutus"),
			Sort:            intPtr(6),
			EditEnable:      boolPtr(false),
			HomeEnable:      boolPtr(true),
			Fixed:           boolPtr(true),
			Enable:          boolPtr(true),
			UpdateAt:        now,
			CreateAt:        now,
		},
		{
			ID:              67,
			ParentID:        int64Ptr(50),
			Type:            2,
			Source:          stringPtr("platform"),
			Code:            "AdvancedUse",
			NameCode:        stringPtr("AdvancedUse"),
			RouteSource:     intPtr(2),
			URL:             stringPtr("/advanced-use"),
			URLType:         intPtr(1),
			OpenType:        intPtr(0),
			Icon:            stringPtr("AdvancedUse.svg"),
			DescriptionCode: stringPtr("menu.desc.advanceUse"),
			Sort:            intPtr(7),
			EditEnable:      boolPtr(true),
			HomeEnable:      boolPtr(true),
			Fixed:           boolPtr(false),
			Enable:          boolPtr(true),
			UpdateAt:        now,
			CreateAt:        now,
		},
		{
			ID:              300,
			ParentID:        int64Ptr(60),
			Type:            2,
			Source:          stringPtr("platform"),
			Code:            "MenuConfiguration",
			NameCode:        stringPtr("MenuConfiguration"),
			RouteSource:     intPtr(2),
			URL:             stringPtr("/MenuConfiguration"),
			URLType:         intPtr(1),
			OpenType:        intPtr(0),
			Icon:            stringPtr("MenuConfiguration.svg"),
			DescriptionCode: stringPtr("menu.desc.MenuConfiguration"),
			Sort:            intPtr(18),
			EditEnable:      boolPtr(false),
			HomeEnable:      boolPtr(true),
			Fixed:           boolPtr(true),
			Enable:          boolPtr(true),
			UpdateAt:        now,
			CreateAt:        now,
		},
	}

	if err := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"parent_id", "type", "source", "code", "name_code", "route_source", "url",
			"url_type", "open_type", "icon", "description_code", "sort", "edit_enable",
			"home_enable", "fixed", "enable", "update_at",
		}),
	}).Create(&resources).Error; err != nil {
		log.Println("seed supos default resources WARN:", err)
	}
}

func seedSuperAdminRoleResources(db *gorm.DB) {
	if db == nil {
		return
	}

	var resourceIDs []int64
	if err := db.Model(&SuposResource{}).
		Where("type IN ?", []int{2, 3, 5}).
		Where("COALESCE(enable, true) = true").
		Order("id ASC").
		Pluck("id", &resourceIDs).Error; err != nil {
		log.Println("query super admin resources WARN:", err)
		return
	}
	if len(resourceIDs) == 0 {
		return
	}

	now := time.Now()
	assignments := make([]IamRoleResource, 0, len(resourceIDs))
	for _, resourceID := range resourceIDs {
		assignments = append(assignments, IamRoleResource{
			RoleID:     enums.RoleSuperAdmin.ID,
			ResourceID: resourceID,
			CreatedAt:  now,
			CreatedBy:  "system",
		})
	}

	if err := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&assignments).Error; err != nil {
		log.Println("seed super admin role resources WARN:", err)
	}
}

func seedIAMBootstrapUser(db *gorm.DB) {
	if db == nil {
		return
	}

	username := firstNonEmptyEnv("IAM_BOOTSTRAP_USERNAME", "tier0")
	password := firstNonEmptyEnv("IAM_BOOTSTRAP_PASSWORD", "tier0")
	email := strings.TrimSpace(os.Getenv("IAM_BOOTSTRAP_EMAIL"))
	userID := username

	var user IamUser
	err := db.Where("LOWER(username) = LOWER(?)", username).Take(&user).Error
	switch {
	case err == nil:
		userID = user.ID
	case err == gorm.ErrRecordNotFound:
		now := time.Now()
		user = IamUser{
			ID:             userID,
			Username:       username,
			DisplayName:    username,
			Email:          email,
			Enabled:        true,
			Source:         "local",
			HomePage:       "/uns",
			FirstTimeLogin: 0,
			TipsEnable:     1,
			CreatedAt:      now,
			UpdatedAt:      now,
		}
		if createErr := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&user).Error; createErr != nil {
			log.Println("seed supos bootstrap user WARN:", createErr)
			return
		}
	default:
		log.Println("query supos bootstrap user WARN:", err)
		return
	}

	if strings.TrimSpace(user.Password) == "" {
		hash, hashErr := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if hashErr != nil {
			log.Println("seed iam bootstrap password WARN:", hashErr)
		} else {
			updates := map[string]any{
				"email":            firstNonEmptyString(user.Email, email),
				"home_page":        firstNonEmptyString(user.HomePage, "/uns"),
				"first_time_login": 0,
				"tips_enable":      1,
				"password":         string(hash),
				"updated_at":       time.Now(),
			}
			if createErr := db.Model(&IamUser{}).Where("id = ?", userID).Updates(updates).Error; createErr != nil {
				log.Println("seed supos bootstrap password WARN:", createErr)
			}
		}
	}

	assignment := IamUserRole{
		UserID: userID,
		RoleID: enums.RoleSuperAdmin.ID,
	}
	if assignErr := db.Clauses(clause.OnConflict{DoNothing: true}).Create(&assignment).Error; assignErr != nil {
		log.Println("seed supos bootstrap role WARN:", assignErr)
	}
}

func firstNonEmptyEnv(key, fallback string) string {
	if value := strings.TrimSpace(os.Getenv(key)); value != "" {
		return value
	}
	return fallback
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if trimmed := strings.TrimSpace(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func stringPtr(value string) *string {
	return &value
}

func intPtr(value int) *int {
	return &value
}

func int64Ptr(value int64) *int64 {
	return &value
}

func boolPtr(value bool) *bool {
	return &value
}

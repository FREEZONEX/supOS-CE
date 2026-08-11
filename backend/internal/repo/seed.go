package repo

import (
	"context"
	"errors"
	"time"

	"backend/internal/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func (s *Store) Seed(ctx context.Context, security config.SecurityConf, gateway config.GatewayConf) error {
	now := time.Now().UTC().UnixMilli()
	ts := repoTimeFromMilli(now)
	builderPassword := security.InitialAdminPassword
	if builderPassword == "" {
		return errors.New("initial admin password is required")
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(builderPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	return s.CommonDB().WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := migrateLegacyBuilderRole(tx, ts); err != nil {
			return err
		}

		builderRole := sysRoleInfo{ID: SystemRoleAdminID, Name: "Admin", Code: BuiltinRoleAdmin, Description: "built-in admin", Type: 1, Status: 1, DefaultHomePage: DefaultHomePage}
		builderRole.CreatedTime = ts
		builderRole.UpdatedTime = ts
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}, {Name: "deleted_time"}},
			DoUpdates: clause.Assignments(map[string]any{"default_home_page": DefaultHomePage, "updated_time": ts}),
		}).Create(&builderRole).Error; err != nil {
			return err
		}
		builderRoleID := builderRole.ID

		builderUser := sysUserInfo{UserName: "tier0", NickName: "Admin", Password: string(passwordHash), Email: "tier0@example.com", Status: 1, IsRandomPwd: false}
		builderUser.CreatedTime = ts
		builderUser.UpdatedTime = ts
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "user_name"}, {Name: "deleted_time"}},
			DoUpdates: clause.Assignments(map[string]any{"updated_time": ts}),
		}).Create(&builderUser).Error; err != nil {
			return err
		}
		builderUserID := builderUser.ID

		builderWorkspaceUser := sysWorkspaceUser{WorkspaceID: 1, UserID: builderUserID, RoleCode: BuiltinRoleAdmin, RoleID: builderRoleID}
		builderWorkspaceUser.CreatedTime = ts
		builderWorkspaceUser.UpdatedTime = ts
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "workspace_id"}, {Name: "user_id"}, {Name: "deleted_time"}},
			DoUpdates: clause.Assignments(map[string]any{"role_code": BuiltinRoleAdmin, "role_id": builderRoleID, "updated_time": ts}),
		}).Create(&builderWorkspaceUser).Error; err != nil {
			return err
		}

		operatorRole := sysRoleInfo{ID: SystemRoleOperatorID, Name: "Operator", Code: BuiltinRoleOperator, Description: "built-in operator", Type: 1, Status: 1, DefaultHomePage: DefaultOperatorHomePage}
		operatorRole.CreatedTime = ts
		operatorRole.UpdatedTime = ts
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "code"}, {Name: "deleted_time"}},
			DoUpdates: clause.Assignments(map[string]any{"default_home_page": DefaultOperatorHomePage, "updated_time": ts}),
		}).Create(&operatorRole).Error; err != nil {
			return err
		}
		operatorRoleID := operatorRole.ID

		resources := seedResources()
		for _, r := range resources {
			row := sysResourceInfo{ID: r.ID, ParentID: r.ParentID, ResourceKey: r.ResourceKey, Name: r.Name, Type: r.Type, RoutePath: r.RoutePath, UrlType: r.UrlType, OpenType: r.OpenType, Icon: r.Icon, Sort: r.Sort, Enabled: 1, SystemGenerated: 1}
			row.CreatedTime = ts
			row.UpdatedTime = ts
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "id"}},
				DoUpdates: clause.AssignmentColumns([]string{"parent_id", "resource_key", "name", "type", "route_path", "url_type", "open_type", "icon", "sort", "enabled", "system_generated", "updated_time"}),
			}).Create(&row).Error; err != nil {
				return err
			}
			builderRoleResource := roleResourceWithTime(builderRoleID, r.ID, now)
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "role_id"}, {Name: "resource_id"}, {Name: "deleted_time"}},
				DoNothing: true,
			}).Create(builderRoleResource).Error; err != nil {
				return err
			}
		}
		if err := tx.Model(&sysResourceInfo{}).Where("resource_key = ?", "system.localization").
			Updates(map[string]any{"enabled": 0, "updated_time": ts}).Error; err != nil {
			return err
		}
		if err := tx.Model(&sysResourceAction{}).Where("action_value LIKE ?", "Localization.%").
			Updates(map[string]any{"enabled": 0, "updated_time": ts}).Error; err != nil {
			return err
		}
		operatorResources := operatorResourceIDs()
		if err := tx.Where("role_id = ? AND deleted_time = 0 AND resource_id NOT IN ?", operatorRoleID, operatorResources).
			Delete(&sysRoleResource{}).Error; err != nil {
			return err
		}
		for _, resourceID := range operatorResources {
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "role_id"}, {Name: "resource_id"}, {Name: "deleted_time"}},
				DoNothing: true,
			}).Create(roleResourceWithTime(operatorRoleID, resourceID, now)).Error; err != nil {
				return err
			}
		}

		for _, a := range seedActions() {
			row := sysResourceAction{ResourceID: a.ResourceID, ActionType: a.ActionType, ActionValue: a.ActionValue, Methods: a.Methods, Enabled: 1, SystemGenerated: 1}
			row.CreatedTime = ts
			row.UpdatedTime = ts
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "action_type"}, {Name: "action_value"}, {Name: "methods"}, {Name: "deleted_time"}},
				DoUpdates: clause.AssignmentColumns([]string{"resource_id", "enabled", "system_generated", "updated_time"}),
			}).Create(&row).Error; err != nil {
				return err
			}
		}
		if err := disableRetiredSeedEntries(tx, now); err != nil {
			return err
		}
		for _, route := range builtinRoutes(gateway) {
			if err := upsertGatewayRoute(tx, route, now); err != nil {
				return err
			}
		}

		return nil
	})
}

// migrateLegacyBuilderRole upgrades the old built-in "builder" role to the new
// "admin" role code/name so that existing deployments continue to work after the
// built-in role rename. It also rewrites workspace role codes that still point to
// the legacy code.
func migrateLegacyBuilderRole(tx *gorm.DB, ts time.Time) error {
	result := tx.Model(&sysRoleInfo{}).
		Where("id = ? AND code = ? AND deleted_time = 0", SystemRoleAdminID, legacyBuiltinRoleBuilder).
		Updates(map[string]any{
			"name":              "Admin",
			"code":              BuiltinRoleAdmin,
			"default_home_page": DefaultHomePage,
			"updated_time":      ts,
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected > 0 {
		if err := tx.Model(&sysWorkspaceUser{}).
			Where("role_id = ? AND role_code = ? AND deleted_time = 0", SystemRoleAdminID, legacyBuiltinRoleBuilder).
			Updates(map[string]any{
				"role_code":    BuiltinRoleAdmin,
				"updated_time": ts,
			}).Error; err != nil {
			return err
		}
	}
	return nil
}

func roleResourceWithTime(roleID, resourceID, now int64) *sysRoleResource {
	ts := repoTimeFromMilli(now)
	row := sysRoleResource{RoleID: roleID, ResourceID: resourceID}
	row.ID = systemRoleResourceID(roleID, resourceID)
	row.CreatedTime = ts
	row.UpdatedTime = ts
	return &row
}

func systemRoleResourceID(roleID, resourceID int64) int64 {
	return roleID*100000 + resourceID
}

func resourceKeys(resources []Resource) []string {
	out := make([]string, 0, len(resources))
	for _, r := range resources {
		out = append(out, r.ResourceKey)
	}
	return out
}

func disableRetiredSeedEntries(tx *gorm.DB, now int64) error {
	ts := repoTimeFromMilli(now)
	retiredResourceKeys := retiredSeedResourceKeys()
	retiredActionValues := []string{
		"button:MenuConfiguration.addOperation",
		"button:MenuConfiguration.editOperation",
		"button:MenuConfiguration.deleteOperation",
		"/sourceflow/**",
		"/nodered/home/**",
		"/eventflow/home/**",
		"/emqx/**",
		"advancedUse.advancedUse",
		"menu.emqx",
	}

	var resourceIDs []int64
	if err := tx.Model(&sysResourceInfo{}).
		Where("resource_key IN ?", retiredResourceKeys).
		Pluck("id", &resourceIDs).Error; err != nil {
		return err
	}
	if len(resourceIDs) > 0 {
		if err := tx.Model(&sysResourceInfo{}).
			Where("id IN ?", resourceIDs).
			Updates(map[string]any{"enabled": 0, "updated_time": ts}).Error; err != nil {
			return err
		}
		if err := tx.Model(&sysResourceAction{}).
			Where("resource_id IN ?", resourceIDs).
			Updates(map[string]any{"enabled": 0, "updated_time": ts}).Error; err != nil {
			return err
		}
	}
	if err := tx.Model(&sysResourceAction{}).
		Where("action_value IN ?", retiredActionValues).
		Updates(map[string]any{"enabled": 0, "updated_time": ts}).Error; err != nil {
		return err
	}
	for _, action := range retiredSeedActionMethodPairs() {
		if err := tx.Model(&sysResourceAction{}).
			Where(
				"action_type = ? AND action_value = ? AND methods = ? AND deleted_time = 0",
				action.ActionType,
				action.ActionValue,
				action.Methods,
			).
			Updates(map[string]any{"enabled": 0, "updated_time": ts}).Error; err != nil {
			return err
		}
	}
	if err := tx.Table("sys_gateway_route").
		Where("route_key IN ? AND deleted_time = 0", retiredGatewayRouteKeys()).
		Updates(map[string]any{"enabled": false, "updated_time": ts}).Error; err != nil {
		return err
	}
	return nil
}

func retiredGatewayRouteKeys() []string {
	return []string{
		"proxy.nodered.home",
		"proxy.eventflow.home",
		"proxy.nodered.tenant",
		"proxy.eventflow.tenant",
		"proxy.sourceflow",
		"proxy.eventflow",
		"proxy.grafana",
		"proxy.portainer.home",
		"proxy.portainer",
		"proxy.emqx.home",
		"proxy.emqx",
	}
}

func retiredSeedActionMethodPairs() []ResourceAction {
	return []ResourceAction{
		{ActionType: "api", ActionValue: "/api/core/common/api-keys", Methods: "GET,POST"},
		{ActionType: "api", ActionValue: "/api/core/common/api-keys/**", Methods: "PUT,DELETE"},
		{ActionType: "openapi", ActionValue: "/openapi/v1/mock/ping", Methods: "GET"},
		{ActionType: "openapi", ActionValue: "/openapi/v1/uns/browse", Methods: "GET"},
		{ActionType: "openapi", ActionValue: "/openapi/v1/uns/nodes/**", Methods: "GET,POST,PUT,DELETE"},
		{ActionType: "openapi", ActionValue: "/openapi/v1/uns/labels", Methods: "GET,POST"},
		{ActionType: "openapi", ActionValue: "/openapi/v1/uns/labels/**", Methods: "GET,PUT,DELETE"},
		{ActionType: "openapi", ActionValue: "/openapi/v1/uns/**", Methods: "GET,POST"},
		{ActionType: "openapi", ActionValue: "/openapi/v1/flows", Methods: "GET"},
		{ActionType: "openapi", ActionValue: "/openapi/v1/flows/**", Methods: "GET"},
		{ActionType: "openapi", ActionValue: "/openapi/v1/assets", Methods: "GET"},
		{ActionType: "openapi", ActionValue: "/openapi/v1/assets/**", Methods: "GET"},
		{ActionType: "openapi", ActionValue: "/openapi/v1/uns/:unsId/attachments", Methods: "GET"},
	}
}

func retiredSeedResourceKeys() []string {
	return []string{
		"app.management",
		"plugin.management",
		"app.space",
		"app.space.view",
		"appspace.view",
		"app.store",
		"app.market",
		"uns.browse",
		"uns.search",
		"uns.create",
		"uns.update",
		"uns.restore",
		"uns.attachment.read",
		"uns.attachment.manage",
		"uns.label.manage",
		"flow",
		"flow.list",
		"flow.write",
		"flow.event.view",
		"flow.event.edit",
		"dashboard.view",
		"openapi.mock.ping",
		"mock.permission",
		"system.portainer.view",
		"dev.tools",
		"dev.advanced.use",
		"gateway.route.write",
		"iam.role.write",
		"apikey.write",
		"system.menu.write",
		"iam.resource.write",
		"app.container.view",
		"app.container.manage",
		"system.emqx.view",
	}
}

func operatorResourceIDs() []int64 {
	return nil
}

func seedResources() []Resource {
	return []Resource{
		{ID: 198, ParentID: 0, ResourceKey: "uns", Name: "UNS", Type: 1, Icon: "menu.tag.uns", Sort: 10, UrlType: 1, OpenType: 0},
		{ID: 199, ParentID: 198, ResourceKey: "uns.page", Name: "Namespace", Type: 2, RoutePath: "/uns", Icon: "Namespace", Sort: 10, UrlType: 1, OpenType: 0},
		{ID: 212, ParentID: 199, ResourceKey: "uns.read", Name: "UNS 读取", Type: 3, Sort: 10, UrlType: 1, OpenType: 0},
		{ID: 214, ParentID: 199, ResourceKey: "uns.write", Name: "UNS 写入", Type: 3, Sort: 20, UrlType: 1, OpenType: 0},
		{ID: 215, ParentID: 199, ResourceKey: "uns.delete", Name: "UNS 删除", Type: 3, Sort: 30, UrlType: 1, OpenType: 0},
		{ID: 210, ParentID: 199, ResourceKey: "uns.manage", Name: "UNS 管理", Type: 3, Sort: 80, UrlType: 1, OpenType: 0},
		{ID: 298, ParentID: 0, ResourceKey: "flow", Name: "Flow", Type: 1, Icon: "menu.tag.connections", Sort: 15, UrlType: 1, OpenType: 0},
		{ID: 299, ParentID: 298, ResourceKey: "flow.collection.page", Name: "Flow", Type: 2, RoutePath: "/flow", Icon: "SourceFlow", Sort: 10, UrlType: 1, OpenType: 0},
		{ID: 319, ParentID: 298, ResourceKey: "flow.event.page", Name: "Event Flow", Type: 2, RoutePath: "/flow", Icon: "EventFlow", Sort: 20, UrlType: 1, OpenType: 0},
		{ID: 300, ParentID: 298, ResourceKey: "flow.read", Name: "Flow 查看", Type: 3, Sort: 30, UrlType: 1, OpenType: 0},
		{ID: 310, ParentID: 298, ResourceKey: "flow.manage", Name: "Flow 管理", Type: 3, Sort: 40, UrlType: 1, OpenType: 0},
		{ID: 100, ParentID: 0, ResourceKey: "system", Name: "System", Type: 1, Icon: "menu.tag.system", Sort: 60, UrlType: 1, OpenType: 0},
		{ID: 110, ParentID: 100, ResourceKey: "iam.user.view", Name: "User Management", Type: 2, RoutePath: "/account-management", Icon: "UserManagement", Sort: 20, UrlType: 1, OpenType: 0},
		{ID: 120, ParentID: 110, ResourceKey: "iam.role.view", Name: "Role Settings", Type: 3, Sort: 10, UrlType: 1, OpenType: 0},
		{ID: 121, ParentID: 120, ResourceKey: "iam.role.create", Name: "Role Create", Type: 3, Sort: 20, UrlType: 1, OpenType: 0},
		{ID: 123, ParentID: 120, ResourceKey: "iam.role.update", Name: "Role Update", Type: 3, Sort: 30, UrlType: 1, OpenType: 0},
		{ID: 122, ParentID: 120, ResourceKey: "iam.role.delete", Name: "Role Delete", Type: 3, Sort: 40, UrlType: 1, OpenType: 0},
		{ID: 150, ParentID: 100, ResourceKey: "apikey.manage", Name: "API Key", Type: 2, RoutePath: "/OpenData", Icon: "OpenData", Sort: 30, UrlType: 1, OpenType: 0},
		{ID: 151, ParentID: 150, ResourceKey: "apikey.create", Name: "API Key Create", Type: 3, Sort: 10, UrlType: 1, OpenType: 0},
		{ID: 153, ParentID: 150, ResourceKey: "apikey.update", Name: "API Key Update", Type: 3, Sort: 20, UrlType: 1, OpenType: 0},
		{ID: 152, ParentID: 150, ResourceKey: "apikey.delete", Name: "API Key Delete", Type: 3, Sort: 30, UrlType: 1, OpenType: 0},
		{ID: 500, ParentID: 0, ResourceKey: "openapi.base", Name: "OpenAPI Base", Type: 3, Sort: 50, UrlType: 1, OpenType: 0},
		{ID: 700, ParentID: 0, ResourceKey: "asset.file.upload", Name: "文件上传", Type: 3, Sort: 70, UrlType: 1, OpenType: 0},
		{ID: 710, ParentID: 0, ResourceKey: "asset.file.read", Name: "文件查看", Type: 3, Sort: 71, UrlType: 1, OpenType: 0},
		{ID: 720, ParentID: 0, ResourceKey: "asset.file.delete", Name: "文件删除", Type: 3, Sort: 72, UrlType: 1, OpenType: 0},
		{ID: 730, ParentID: 0, ResourceKey: "asset.file.download", Name: "文件下载", Type: 3, Sort: 73, UrlType: 1, OpenType: 0},
	}
}

func seedActions() []ResourceAction {
	return []ResourceAction{
		{ResourceID: 110, ActionType: "ui", ActionValue: "menu.system.users"},
		{ResourceID: 110, ActionType: "ui", ActionValue: "button:UserManagement.add"},
		{ResourceID: 110, ActionType: "ui", ActionValue: "button:UserManagement.edit"},
		{ResourceID: 110, ActionType: "ui", ActionValue: "button:UserManagement.enable"},
		{ResourceID: 110, ActionType: "ui", ActionValue: "button:UserManagement.disable"},
		{ResourceID: 110, ActionType: "ui", ActionValue: "button:UserManagement.resetPassword"},
		{ResourceID: 110, ActionType: "ui", ActionValue: "button:UserManagement.delete"},
		{ResourceID: 120, ActionType: "ui", ActionValue: "button:UserManagement.role_setting"},
		{ResourceID: 110, ActionType: "api", ActionValue: "/api/core/iam/users", Methods: "GET"},
		{ResourceID: 110, ActionType: "api", ActionValue: "/api/core/iam/users", Methods: "POST"},
		{ResourceID: 110, ActionType: "api", ActionValue: "/api/core/iam/users/**", Methods: "PUT"},
		{ResourceID: 110, ActionType: "api", ActionValue: "/api/core/iam/users/**", Methods: "POST"},
		{ResourceID: 110, ActionType: "api", ActionValue: "/api/core/iam/users/**", Methods: "DELETE"},
		{ResourceID: 120, ActionType: "ui", ActionValue: "menu.system.roles"},
		{ResourceID: 120, ActionType: "api", ActionValue: "/api/core/iam/roles", Methods: "GET"},
		{ResourceID: 121, ActionType: "api", ActionValue: "/api/core/iam/roles", Methods: "POST"},
		{ResourceID: 123, ActionType: "api", ActionValue: "/api/core/iam/roles/**", Methods: "PUT"},
		{ResourceID: 122, ActionType: "api", ActionValue: "/api/core/iam/roles/**", Methods: "DELETE"},
		{ResourceID: 150, ActionType: "ui", ActionValue: "menu.system.apikey"},
		{ResourceID: 150, ActionType: "api", ActionValue: "/api/core/common/api-keys", Methods: "GET"},
		{ResourceID: 151, ActionType: "ui", ActionValue: "button:dataopen.addKey"},
		{ResourceID: 151, ActionType: "api", ActionValue: "/api/core/common/api-keys", Methods: "POST"},
		{ResourceID: 153, ActionType: "ui", ActionValue: "button:dataopen.enable"},
		{ResourceID: 153, ActionType: "ui", ActionValue: "button:dataopen.disable"},
		{ResourceID: 153, ActionType: "api", ActionValue: "/api/core/common/api-keys/**", Methods: "PUT"},
		{ResourceID: 152, ActionType: "ui", ActionValue: "button:dataopen.delete"},
		{ResourceID: 152, ActionType: "api", ActionValue: "/api/core/common/api-keys/**", Methods: "DELETE"},
		{ResourceID: 199, ActionType: "ui", ActionValue: "menu.uns"},
		{ResourceID: 199, ActionType: "api", ActionValue: "/api/core/uns/nodes", Methods: "GET"},
		{ResourceID: 199, ActionType: "api", ActionValue: "/api/core/uns/newMsg", Methods: "GET"},
		{ResourceID: 199, ActionType: "api", ActionValue: "/api/core/uns/dashboard", Methods: "GET"},
		{ResourceID: 210, ActionType: "ui", ActionValue: "button:Namespace.uns_import"},
		{ResourceID: 210, ActionType: "ui", ActionValue: "button:Namespace.uns_export"},
		{ResourceID: 210, ActionType: "api", ActionValue: "/api/core/uns/import", Methods: "POST"},
		{ResourceID: 210, ActionType: "api", ActionValue: "/api/core/uns/export", Methods: "POST"},
		{ResourceID: 210, ActionType: "api", ActionValue: "/api/core/uns/export/global", Methods: "POST"},
		{ResourceID: 210, ActionType: "api", ActionValue: "/api/core/uns/import-jobs", Methods: "POST"},
		{ResourceID: 210, ActionType: "api", ActionValue: "/api/core/uns/import-jobs/**", Methods: "GET"},
		{ResourceID: 210, ActionType: "api", ActionValue: "/api/core/uns/export-jobs", Methods: "POST"},
		{ResourceID: 210, ActionType: "api", ActionValue: "/api/core/uns/export-jobs/**", Methods: "GET"},
		{ResourceID: 212, ActionType: "openapi", ActionValue: "/openapi/v1/uns/read", Methods: "POST"},
		{ResourceID: 212, ActionType: "openapi", ActionValue: "/openapi/v1/uns/browse", Methods: "POST"},
		{ResourceID: 212, ActionType: "openapi", ActionValue: "/openapi/v1/uns/search", Methods: "POST"},
		{ResourceID: 212, ActionType: "openapi", ActionValue: "/openapi/v1/uns/history", Methods: "POST"},
		{ResourceID: 212, ActionType: "openapi", ActionValue: "/openapi/v1/uns/nodes", Methods: "POST"},
		{ResourceID: 210, ActionType: "ui", ActionValue: "button:Namespace.folder_detail"},
		{ResourceID: 210, ActionType: "ui", ActionValue: "button:Namespace.file_detail"},
		{ResourceID: 212, ActionType: "api", ActionValue: "/api/core/uns/nodes/**", Methods: "GET"},
		{ResourceID: 212, ActionType: "api", ActionValue: "/api/core/uns/recycle", Methods: "GET"},
		{ResourceID: 212, ActionType: "api", ActionValue: "/api/core/uns/labels", Methods: "GET"},
		{ResourceID: 212, ActionType: "api", ActionValue: "/api/core/uns/labels/**", Methods: "GET"},
		{ResourceID: 214, ActionType: "openapi", ActionValue: "/openapi/v1/uns/write", Methods: "POST"},
		{ResourceID: 210, ActionType: "openapi", ActionValue: "/openapi/v1/uns/create", Methods: "POST"},
		{ResourceID: 210, ActionType: "openapi", ActionValue: "/openapi/v1/uns/update", Methods: "POST"},
		{ResourceID: 210, ActionType: "openapi", ActionValue: "/openapi/v1/uns/restore", Methods: "POST"},
		{ResourceID: 215, ActionType: "openapi", ActionValue: "/openapi/v1/uns/delete", Methods: "POST"},
		{ResourceID: 222, ActionType: "api", ActionValue: "/api/core/uns/labels", Methods: "POST"},
		{ResourceID: 222, ActionType: "api", ActionValue: "/api/core/uns/labels/**", Methods: "PUT,DELETE"},
		{ResourceID: 210, ActionType: "api", ActionValue: "/api/core/uns/nodes", Methods: "POST"},
		{ResourceID: 210, ActionType: "ui", ActionValue: "button:Namespace.folder_add"},
		{ResourceID: 210, ActionType: "ui", ActionValue: "button:Namespace.file_add"},
		{ResourceID: 210, ActionType: "api", ActionValue: "/api/core/uns/nodes/**", Methods: "PUT"},
		{ResourceID: 210, ActionType: "ui", ActionValue: "button:Namespace.folder_copy"},
		{ResourceID: 210, ActionType: "ui", ActionValue: "button:Namespace.folder_paste"},
		{ResourceID: 210, ActionType: "ui", ActionValue: "button:Namespace.file_copy"},
		{ResourceID: 210, ActionType: "ui", ActionValue: "button:Namespace.file_paste"},
		{ResourceID: 215, ActionType: "api", ActionValue: "/api/core/uns/nodes/**", Methods: "DELETE"},
		{ResourceID: 215, ActionType: "ui", ActionValue: "button:Namespace.folder_delete"},
		{ResourceID: 215, ActionType: "ui", ActionValue: "button:Namespace.file_delete"},
		{ResourceID: 210, ActionType: "api", ActionValue: "/api/core/uns/nodes/**", Methods: "POST"},
		{ResourceID: 299, ActionType: "ui", ActionValue: "menu.flow.source"},
		{ResourceID: 319, ActionType: "ui", ActionValue: "menu.flow.event"},
		{ResourceID: 300, ActionType: "ui", ActionValue: "button:SourceFlow.node_management"},
		{ResourceID: 300, ActionType: "ui", ActionValue: "button:EventFlow.node_management"},
		{ResourceID: 300, ActionType: "api", ActionValue: "/api/core/flows", Methods: "GET"},
		{ResourceID: 300, ActionType: "api", ActionValue: "/api/core/flows/**", Methods: "GET"},
		{ResourceID: 300, ActionType: "openapi", ActionValue: "/openapi/v1/flow/list", Methods: "POST"},
		{ResourceID: 300, ActionType: "openapi", ActionValue: "/openapi/v1/flow/get", Methods: "POST"},
		{ResourceID: 300, ActionType: "openapi", ActionValue: "/openapi/v1/flow/flowdata", Methods: "POST"},
		{ResourceID: 300, ActionType: "openapi", ActionValue: "/openapi/v1/flow/nodes", Methods: "POST"},
		{ResourceID: 300, ActionType: "gateway", ActionValue: "/nodered/flows", Methods: "GET,HEAD"},
		{ResourceID: 300, ActionType: "gateway", ActionValue: "/eventflow/flows", Methods: "GET,HEAD"},
		{ResourceID: 300, ActionType: "gateway", ActionValue: "/nodered/home/flows", Methods: "GET,HEAD"},
		{ResourceID: 300, ActionType: "gateway", ActionValue: "/eventflow/home/flows", Methods: "GET,HEAD"},
		{ResourceID: 300, ActionType: "gateway", ActionValue: "/nodered/**", Methods: "GET"},
		{ResourceID: 300, ActionType: "gateway", ActionValue: "/eventflow/**", Methods: "GET"},
		{ResourceID: 310, ActionType: "api", ActionValue: "/api/core/flows", Methods: "POST"},
		{ResourceID: 310, ActionType: "api", ActionValue: "/api/core/flows/**", Methods: "PUT,DELETE,POST"},
		{ResourceID: 310, ActionType: "openapi", ActionValue: "/openapi/v1/flow/create", Methods: "POST"},
		{ResourceID: 310, ActionType: "openapi", ActionValue: "/openapi/v1/flow/update", Methods: "POST"},
		{ResourceID: 310, ActionType: "openapi", ActionValue: "/openapi/v1/flow/delete", Methods: "POST"},
		{ResourceID: 310, ActionType: "openapi", ActionValue: "/openapi/v1/flow/deploy", Methods: "POST"},
		{ResourceID: 310, ActionType: "openapi", ActionValue: "/openapi/v1/uns/unsBindFlow", Methods: "POST"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "menu.flow.editor"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:SourceFlow.add"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:SourceFlow.edit"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:SourceFlow.delete"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:SourceFlow.design"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:SourceFlow.save"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:SourceFlow.deploy"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:SourceFlow.import"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:SourceFlow.export"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:SourceFlow.copy"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:SourceFlow.moveToGroup"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:EventFlow.add"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:EventFlow.edit"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:EventFlow.delete"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:EventFlow.design"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:EventFlow.save"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:EventFlow.deploy"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:EventFlow.import"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:EventFlow.export"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:EventFlow.copy"},
		{ResourceID: 310, ActionType: "ui", ActionValue: "button:EventFlow.moveToGroup"},
		{ResourceID: 310, ActionType: "gateway", ActionValue: "/nodered/**", Methods: "POST,PUT,DELETE,PATCH"},
		{ResourceID: 310, ActionType: "gateway", ActionValue: "/eventflow/**", Methods: "POST,PUT,DELETE,PATCH"},
		{ResourceID: 500, ActionType: "openapi", ActionValue: "/openapi/v1/auth/whoami", Methods: "POST"},
		{ResourceID: 500, ActionType: "openapi", ActionValue: "/openapi/v1/info", Methods: "POST"},
		{ResourceID: 700, ActionType: "api", ActionValue: "/api/core/assets", Methods: "POST"},
		{ResourceID: 700, ActionType: "api", ActionValue: "/api/core/assets/multipart/init", Methods: "POST"},
		{ResourceID: 700, ActionType: "api", ActionValue: "/api/core/assets/multipart/part-urls", Methods: "POST"},
		{ResourceID: 700, ActionType: "api", ActionValue: "/api/core/assets/multipart/complete", Methods: "POST"},
		{ResourceID: 700, ActionType: "api", ActionValue: "/api/core/assets/multipart/abort", Methods: "POST"},
		{ResourceID: 700, ActionType: "api", ActionValue: "/api/core/asset-bindings", Methods: "POST"},
		{ResourceID: 700, ActionType: "openapi", ActionValue: "/openapi/v1/assets/upload", Methods: "POST"},
		{ResourceID: 700, ActionType: "openapi", ActionValue: "/openapi/v1/uns/attachments", Methods: "POST"},
		{ResourceID: 700, ActionType: "openapi", ActionValue: "/openapi/v1/assets/files", Methods: "POST"},
		{ResourceID: 700, ActionType: "openapi", ActionValue: "/openapi/v1/assets/files/multipart/init", Methods: "POST"},
		{ResourceID: 700, ActionType: "openapi", ActionValue: "/openapi/v1/assets/files/multipart/part-urls", Methods: "POST"},
		{ResourceID: 700, ActionType: "openapi", ActionValue: "/openapi/v1/assets/files/multipart/complete", Methods: "POST"},
		{ResourceID: 700, ActionType: "openapi", ActionValue: "/openapi/v1/assets/files/multipart/abort", Methods: "POST"},
		{ResourceID: 710, ActionType: "ui", ActionValue: "menu.assets"},
		{ResourceID: 710, ActionType: "api", ActionValue: "/api/core/assets", Methods: "GET"},
		{ResourceID: 710, ActionType: "api", ActionValue: "/api/core/assets/**", Methods: "GET"},
		{ResourceID: 710, ActionType: "openapi", ActionValue: "/openapi/v1/uns/attachments/list", Methods: "POST"},
		{ResourceID: 720, ActionType: "api", ActionValue: "/api/core/assets/**", Methods: "DELETE"},
		{ResourceID: 720, ActionType: "api", ActionValue: "/api/core/asset-bindings/**", Methods: "DELETE"},
		{ResourceID: 720, ActionType: "openapi", ActionValue: "/openapi/v1/assets/files/delete", Methods: "POST"},
		{ResourceID: 730, ActionType: "openapi", ActionValue: "/openapi/v1/assets/files/download", Methods: "GET"},
		{ResourceID: 730, ActionType: "openapi", ActionValue: "/openapi/v1/assets/files/url", Methods: "GET"},
	}
}

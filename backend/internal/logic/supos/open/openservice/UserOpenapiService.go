// Code scaffolded by goctl. Safe to edit.

package openservice

import (
	"context"
	"strings"

	"backend/internal/common/constants"
	"backend/internal/svc"
	"backend/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserOpenapiService struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// NewUserOpenapiService creates a new UserOpenapiService instance
func NewUserOpenapiService(ctx context.Context, svcCtx *svc.ServiceContext) *UserOpenapiService {
	return &UserOpenapiService{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// UserPageQueryDto 用户分页查询参数
type UserPageQueryDto struct {
	PageNo      int
	PageSize    int
	Username    string // 用户名，精准查询
	DisplayName string // 显示名称，支持模糊查询
	Email       string // 邮箱，支持模糊查询
	Phone       string // 手机号，支持模糊查询
	Enabled     *bool  // 是否启用
	// ExactUsername 精确用户名查询（用于用户详情接口）
	ExactUsername string
}

// UserManageResult 用户管理查询结果
type UserManageResult struct {
	Users []types.OpenUserInfo
	Total int
}

// UserManageList 获取用户列表（供列表和详情接口共用）
func (s *UserOpenapiService) UserManageList(params UserPageQueryDto) (*UserManageResult, error) {
	roleService := NewRoleService(s.ctx, s.svcCtx)

	// 构建查询参数
	searchRequest := SearchUsersRequest{}

	// 如果是精确用户名查询（用户详情接口）
	if params.ExactUsername != "" {
		searchRequest.Username = params.ExactUsername
		searchRequest.Exact = true
	} else {
		// 分页查询（用户列表接口）
		searchRequest.First = (params.PageNo - 1) * params.PageSize
		searchRequest.Max = params.PageSize
		searchRequest.Username = params.Username     // 精准查询
		searchRequest.FirstName = params.DisplayName // 模糊查询
		searchRequest.Email = params.Email           // 模糊查询
		searchRequest.Enabled = params.Enabled
	}

	// 通过 RoleService 获取用户列表
	users, err := roleService.SearchUsers(searchRequest)
	if err != nil {
		s.Errorf("Failed to search users: %v", err)
		return nil, err
	}

	resultList := make([]types.OpenUserInfo, 0, len(users))

	for _, user := range users {
		// 过滤服务账户
		if roleService.IsServiceAccount(user.Attributes) {
			continue
		}

		// 按手机号过滤（需要在获取用户后进行属性过滤）
		if params.Phone != "" {
			userPhone := roleService.GetStringAttribute(user.Attributes, "phone")
			if userPhone == "" || !strings.Contains(userPhone, params.Phone) {
				continue
			}
		}

		// 获取用户角色
		roles, err := roleService.GetRoleListByUserId(user.ID)
		if err != nil {
			s.Errorf("Failed to get roles for user %s: %v", user.ID, err)
			continue
		}

		// 转换角色列表
		roleList := make([]types.RoleSummary, 0, len(roles))
		for _, r := range roles {
			roleList = append(roleList, types.RoleSummary{
				RoleID:          r.RoleID,
				RoleName:        r.RoleName,
				RoleDescription: r.RoleDescription,
				ClientRole:      r.ClientRole,
			})
		}

		// 解析用户属性并构建 OpenUserInfo
		phone := roleService.GetStringAttribute(user.Attributes, "phone")
		homePage := roleService.GetStringAttribute(user.Attributes, "homePage")
		if homePage == "" {
			homePage = constants.DefaultHomepage
		}
		source := roleService.GetStringAttribute(user.Attributes, "source")
		firstTimeLogin := roleService.GetIntAttribute(user.Attributes, "firstTimeLogin", 1)
		tipsEnable := roleService.GetIntAttribute(user.Attributes, "tipsEnable", 1)

		userInfo := types.OpenUserInfo{
			ID:                user.ID,
			Email:             user.Email,
			EmailVerified:     false,
			FirstName:         user.FirstName,
			PreferredUsername: user.PreferredUsername,
			Sub:               user.ID,
			Enabled:           user.Enabled,
			RoleList:          roleList,
			FirstTimeLogin:    firstTimeLogin,
			TipsEnable:        tipsEnable,
			HomePage:          homePage,
			Phone:             phone,
			Source:            source,
		}
		resultList = append(resultList, userInfo)
	}

	return &UserManageResult{
		Users: resultList,
		Total: len(resultList),
	}, nil
}

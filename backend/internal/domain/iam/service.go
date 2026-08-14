package iam

import (
	"context"
	"errors"
	"strings"

	"backend/internal/repo"

	"golang.org/x/crypto/bcrypt"
)

type Service struct {
	repo *repo.IAMRepo
}

type UserRoleCommand struct {
	RoleID   int64
	RoleName string
}

type UserSaveCommand struct {
	Username  string
	Password  string
	FirstName string
	Email     string
	Phone     string
	Enabled   bool
	RoleList  []UserRoleCommand
	UserID    int64
}

type UserUpdateCommand struct {
	ID        int64
	Username  string
	FirstName string
	Email     string
	Phone     string
	Enabled   bool
	RoleList  []UserRoleCommand
	UserID    int64
}

type UserPasswordResetCommand struct {
	UserID   int64
	Password string
	ActorID  int64
}

type CurrentUserPasswordCommand struct {
	UserID      int64
	OldPassword string
	NewPassword string
}

type ResourceSaveCommand struct {
	ResourceID  int64
	ParentID    int64
	ResourceKey string
	Name        string
	Type        int64
	RoutePath   string
	UrlType     int64
	OpenType    int64
	Icon        string
	Sort        int64
	Enabled     int64
}

func New(ctx context.Context) *Service {
	return &Service{
		repo: repo.NewIAMRepo(ctx),
	}
}

func (s *Service) IsAdmin(ctx context.Context, userID int64) bool {
	if s == nil || s.repo == nil {
		return false
	}
	ok, err := s.repo.IsAdmin(ctx, userID)
	return err == nil && ok
}

func (s *Service) UserRoleCode(ctx context.Context, userID int64) string {
	if s == nil || s.repo == nil {
		return ""
	}
	code, _ := s.repo.UserRoleCode(ctx, userID)
	return code
}

func (s *Service) Users(ctx context.Context) ([]map[string]any, error) {
	return s.repo.ListUsers(ctx)
}

func (s *Service) CreateUser(ctx context.Context, cmd UserSaveCommand) (map[string]any, error) {
	username := strings.TrimSpace(cmd.Username)
	password := strings.TrimSpace(cmd.Password)
	roleID := repo.SystemRoleAdminID
	if len(cmd.RoleList) > 0 {
		roleID = cmd.RoleList[0].RoleID
	}
	if username == "" || password == "" || roleID <= 0 {
		return nil, ErrInvalidUser
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	status := int64(0)
	if cmd.Enabled {
		status = 1
		if err := s.checkUserQuota(ctx); err != nil {
			return nil, err
		}
	}
	return s.repo.CreateUser(ctx, repo.UserCreate{
		UserName: username,
		NickName: strings.TrimSpace(cmd.FirstName),
		Email:    strings.TrimSpace(cmd.Email),
		Phone:    strings.TrimSpace(cmd.Phone),
		Password: string(hash),
		Status:   status,
		RoleID:   roleID,
		UserID:   cmd.UserID,
	})
}

func (s *Service) UpdateUser(ctx context.Context, cmd UserUpdateCommand) (map[string]any, error) {
	roleID := int64(0)
	hasRole := len(cmd.RoleList) > 0
	if hasRole {
		roleID = cmd.RoleList[0].RoleID
	}
	status := int64(0)
	if cmd.Enabled {
		status = 1
		current, err := s.repo.GetUserByID(ctx, cmd.ID)
		if err != nil {
			return nil, err
		}
		if current.Status != 1 {
			if err := s.checkUserQuota(ctx); err != nil {
				return nil, err
			}
		}
	}
	return s.repo.UpdateUser(ctx, repo.UserUpdate{
		UserID:   cmd.ID,
		UserName: strings.TrimSpace(cmd.Username),
		NickName: strings.TrimSpace(cmd.FirstName),
		Email:    strings.TrimSpace(cmd.Email),
		Phone:    strings.TrimSpace(cmd.Phone),
		Status:   status,
		RoleID:   roleID,
		HasRole:  hasRole,
		ActorID:  cmd.UserID,
	})
}

func (s *Service) DeleteUser(ctx context.Context, userID, actorID int64) error {
	return s.repo.DeleteUser(ctx, userID, actorID)
}

func (s *Service) ResetUserPassword(ctx context.Context, cmd UserPasswordResetCommand) error {
	password := strings.TrimSpace(cmd.Password)
	if cmd.UserID <= 0 || password == "" {
		return ErrInvalidUser
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.ResetUserPassword(ctx, cmd.UserID, string(hash), cmd.ActorID)
}

func (s *Service) ChangeCurrentUserPassword(ctx context.Context, cmd CurrentUserPasswordCommand) error {
	oldPassword := strings.TrimSpace(cmd.OldPassword)
	newPassword := strings.TrimSpace(cmd.NewPassword)
	if cmd.UserID <= 0 || oldPassword == "" || newPassword == "" {
		return ErrInvalidUser
	}
	user, err := s.repo.GetUserByID(ctx, cmd.UserID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return ErrInvalidPassword
	}
	hash, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return s.repo.ChangeUserPassword(ctx, cmd.UserID, string(hash))
}

func (s *Service) CountEnabledUsers(ctx context.Context) (int64, error) {
	return s.repo.CountEnabledUsers(ctx)
}

func (s *Service) checkUserQuota(ctx context.Context) error {
	return nil
}

func (s *Service) Roles(ctx context.Context) ([]repo.Role, error) {
	return s.repo.ListRoles(ctx)
}

func (s *Service) CreateRole(ctx context.Context, name string, userID int64) (repo.Role, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return repo.Role{}, ErrInvalidRole
	}
	return s.repo.CreateRole(ctx, name, "", userID)
}

func (s *Service) UpdateRole(ctx context.Context, roleID int64, name string, allowURIs []string, defaultHomePage string, userID int64) error {
	if roleID <= 0 {
		return ErrInvalidRole
	}
	return s.repo.UpdateRole(ctx, roleID, name, allowURIs, defaultHomePage, userID)
}

func (s *Service) DeleteRole(ctx context.Context, roleID, userID int64) error {
	if roleID <= 0 {
		return ErrInvalidRole
	}
	return s.repo.DeleteRole(ctx, roleID, userID)
}

func (s *Service) Resources(ctx context.Context) ([]repo.Resource, error) {
	return s.repo.ListResources(ctx)
}

func (s *Service) Menus(ctx context.Context) ([]repo.Resource, error) {
	return s.repo.ListMenuResources(ctx)
}

func (s *Service) SaveResource(ctx context.Context, cmd ResourceSaveCommand, userID int64) (repo.Resource, error) {
	if strings.TrimSpace(cmd.ResourceKey) == "" || strings.TrimSpace(cmd.Name) == "" || cmd.Type <= 0 {
		return repo.Resource{}, ErrInvalidRole
	}
	return s.repo.SaveResource(ctx, repo.Resource{
		ID:          cmd.ResourceID,
		ParentID:    cmd.ParentID,
		ResourceKey: strings.TrimSpace(cmd.ResourceKey),
		Name:        strings.TrimSpace(cmd.Name),
		Type:        cmd.Type,
		RoutePath:   strings.TrimSpace(cmd.RoutePath),
		UrlType:     cmd.UrlType,
		OpenType:    cmd.OpenType,
		Icon:        strings.TrimSpace(cmd.Icon),
		Sort:        cmd.Sort,
	}, userID, cmd.Enabled)
}

func (s *Service) ResourceActions(ctx context.Context, resourceID int64) ([]repo.ResourceAction, error) {
	if resourceID <= 0 {
		return nil, ErrInvalidRole
	}
	return s.repo.ListResourceActions(ctx, resourceID)
}

func (s *Service) ResolveUserHomePage(ctx context.Context, userID int64, preferred string) (string, error) {
	keys, err := s.repo.ResourceKeysForUser(ctx, userID)
	if err != nil {
		return "", err
	}
	resourceMap, err := s.repo.ResourceMapByKeys(ctx, keys)
	if err != nil {
		return "", err
	}
	allowedPages := map[string]string{}
	for _, item := range resourceMap {
		page := strings.TrimSpace(item.RoutePath)
		if page == "" || !strings.HasPrefix(page, "/") || strings.HasPrefix(page, "//") {
			continue
		}
		allowedPages[strings.ToLower(page)] = page
	}
	orderedPages, err := s.repo.HomePagesForUser(ctx, userID)
	if err != nil {
		return "", err
	}
	orderedAllowedPages := orderedAllowedHomePages(allowedPages, orderedPages)
	preferred = strings.TrimSpace(preferred)
	if preferred != "" {
		if page, ok := allowedPages[strings.ToLower(preferred)]; ok {
			return page, nil
		}
	}
	defaultPage, err := s.repo.DefaultHomePageForUser(ctx, userID)
	if err != nil {
		return "", err
	}
	if page, ok := allowedPages[strings.ToLower(strings.TrimSpace(defaultPage))]; ok {
		return page, nil
	}
	return defaultAllowedHomePageInOrder(allowedPages, orderedAllowedPages), nil
}

func (s *Service) UpdateUserHomePage(ctx context.Context, userID int64, preferred string) (string, error) {
	homePage, err := s.ResolveUserHomePage(ctx, userID, preferred)
	if err != nil {
		return "", err
	}
	if _, err := repo.NewUserConfigRepo(ctx).UpdateUserConfig(ctx, userID, &homePage, nil); err != nil {
		return "", err
	}
	return homePage, nil
}

func normalizeRoleHomePage(value string, resourceList []map[string]string) string {
	allowed := map[string]string{}
	ordered := make([]string, 0, len(resourceList))
	for _, item := range resourceList {
		page := strings.TrimSpace(item["routePath"])
		if page == "" || !strings.HasPrefix(page, "/") || strings.HasPrefix(page, "//") {
			continue
		}
		key := strings.ToLower(page)
		if _, exists := allowed[key]; exists {
			continue
		}
		allowed[key] = page
		ordered = append(ordered, page)
	}
	if page, ok := allowed[strings.ToLower(strings.TrimSpace(value))]; ok {
		return page
	}
	return defaultAllowedHomePageInOrder(allowed, ordered)
}

func orderedAllowedHomePages(allowed map[string]string, ordered []string) []string {
	out := make([]string, 0, len(ordered))
	seen := map[string]struct{}{}
	for _, page := range ordered {
		key := strings.ToLower(strings.TrimSpace(page))
		if key == "" {
			continue
		}
		value, ok := allowed[key]
		if !ok {
			continue
		}
		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}
		out = append(out, value)
	}
	return out
}

func defaultAllowedHomePageInOrder(allowed map[string]string, ordered []string) string {
	if page, ok := allowed[strings.ToLower(repo.DefaultHomePage)]; ok {
		return page
	}
	if page, ok := allowed[strings.ToLower(repo.DefaultOperatorHomePage)]; ok {
		return page
	}
	if len(ordered) > 0 {
		return ordered[0]
	}
	return repo.DefaultOperatorHomePage
}

func (s *Service) SaveResourceAction(ctx context.Context, action repo.ResourceAction, userID int64) (repo.ResourceAction, error) {
	if action.ResourceID <= 0 {
		return repo.ResourceAction{}, ErrInvalidRole
	}
	action.ActionType = strings.TrimSpace(strings.ToLower(action.ActionType))
	action.ActionValue = strings.TrimSpace(action.ActionValue)
	action.Methods = strings.TrimSpace(strings.ToUpper(action.Methods))
	if action.ActionType == "" || action.ActionValue == "" {
		return repo.ResourceAction{}, ErrInvalidRole
	}
	return s.repo.SaveResourceAction(ctx, action, userID)
}

func (s *Service) DeleteResourceAction(ctx context.Context, actionID, userID int64) error {
	if actionID <= 0 {
		return ErrInvalidRole
	}
	return s.repo.DeleteResourceAction(ctx, actionID, userID)
}

var ErrInvalidRole = errors.New("invalid role")
var ErrInvalidUser = errors.New("invalid user")
var ErrInvalidPassword = errors.New("invalid password")

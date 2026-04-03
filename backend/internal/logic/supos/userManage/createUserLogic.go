package userManage

import (
	"context"
	"strings"
	"time"

	"backend/internal/common/constants"
	"backend/internal/repo/relationDB"
	"backend/internal/svc"
	"backend/internal/types"

	"gitee.com/unitedrhino/share/errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CreateUserLogic struct {
	baseUserManageLogic
}

// Create user
func NewCreateUserLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateUserLogic {
	return &CreateUserLogic{
		baseUserManageLogic: newBaseUserManageLogic(ctx, svcCtx),
	}
}

func (l *CreateUserLogic) CreateUser(req *types.UserCreateReq) (*types.OperationResult, error) {
	if req == nil {
		return nil, errors.Parameter.WithMsg("request body is empty")
	}
	username := strings.TrimSpace(req.Username)
	if username == "" {
		return nil, errors.Parameter.WithMsg("user.username.empty")
	}
	if len(username) < 3 || len(username) > 30 {
		return nil, errors.Parameter.WithMsg("user.username.invalid")
	}
	password := strings.TrimSpace(req.Password)
	if password == "" {
		return nil, errors.Parameter.WithMsg("user.password.empty")
	}

	db, err := l.db()
	if err != nil {
		return nil, err
	}

	if existing, err := l.getIAMUserByUsername(db, username); err != nil {
		return nil, err
	} else if existing != nil {
		return nil, errors.Parameter.WithMsg("user.username.already.exists")
	}

	email := strings.TrimSpace(req.Email)
	if email != "" {
		if existing, err := l.getIAMUserByEmail(db, email); err != nil {
			return nil, err
		} else if existing != nil {
			return nil, errors.Parameter.WithMsg("user.email.already.exists")
		}
	}

	roleIDs, err := l.resolveRoleIDs(db, req.RoleList)
	if err != nil {
		return nil, err
	}
	if len(roleIDs) == 0 {
		defaultRoleID, roleErr := l.getDefaultUserRoleID(db)
		if roleErr != nil {
			return nil, roleErr
		}
		roleIDs = []string{defaultRoleID}
	}

	userID := strings.TrimSpace(req.ID)
	if userID == "" {
		userID = uuid.NewString()
	}
	now := time.Now()
	user := relationDB.IamUser{
		ID:             userID,
		Username:       username,
		DisplayName:    firstNonEmpty(req.FirstName, username),
		Email:          email,
		Enabled:        boolValue(req.Enabled, true),
		Source:         strings.TrimSpace(req.Source),
		Phone:          strings.TrimSpace(req.Phone),
		HomePage:       constants.DefaultHomepage,
		FirstTimeLogin: 1,
		TipsEnable:     1,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	if err := db.Transaction(func(tx *gorm.DB) error {
		if txErr := tx.Create(&user).Error; txErr != nil {
			l.Errorf("create user failed: %v", txErr)
			return errors.System.WithMsg("failed to create user")
		}
		if passwordErr := l.upsertUserPassword(tx, userID, password); passwordErr != nil {
			return passwordErr
		}
		return l.replaceUserRoles(tx, userID, roleIDs)
	}); err != nil {
		return nil, err
	}

	l.invalidateUserCache(userID)
	return l.success(), nil
}

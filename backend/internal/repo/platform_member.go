package repo

import (
	"context"
	"time"
)

type PlatformMemberRow struct {
	UserID            int64     `gorm:"column:user_id"`
	UserName          string    `gorm:"column:user_name"`
	NickName          string    `gorm:"column:nick_name"`
	Email             string    `gorm:"column:email"`
	Status            int64     `gorm:"column:status"`
	RoleID            int64     `gorm:"column:role_id"`
	RoleCode          string    `gorm:"column:role_code"`
	RoleName          string    `gorm:"column:role_name"`
	RoleDescription   string    `gorm:"column:role_description"`
	UserUpdatedTime   time.Time `gorm:"column:user_updated_time"`
	MemberUpdatedTime time.Time `gorm:"column:member_updated_time"`
	RoleUpdatedTime   time.Time `gorm:"column:role_updated_time"`
}

func (r *IAMRepo) ListPlatformMemberRows(ctx context.Context, workspaceID int64) ([]PlatformMemberRow, error) {
	var rows []PlatformMemberRow
	err := r.db.WithContext(ctx).Table("sys_workspace_user AS wu").
		Select(`wu.user_id,
			u.user_name,
			u.nick_name,
			u.email,
			u.status,
			ri.id AS role_id,
			ri.code AS role_code,
			ri.name AS role_name,
			ri.desc AS role_description,
			u.updated_time AS user_updated_time,
			wu.updated_time AS member_updated_time,
			ri.updated_time AS role_updated_time`).
		Joins("JOIN sys_user_info AS u ON u.user_id = wu.user_id AND u.deleted_time = 0").
		Joins("JOIN sys_role_info AS ri ON ri.id = wu.role_id AND ri.deleted_time = 0").
		Where("wu.workspace_id = ? AND wu.deleted_time = 0", workspaceID).
		Order("wu.user_id ASC, ri.id ASC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for index := range rows {
		email, _, err := decryptUserContacts(rows[index].UserID, rows[index].Email, "")
		if err != nil {
			return nil, err
		}
		rows[index].Email = email
	}
	return rows, err
}

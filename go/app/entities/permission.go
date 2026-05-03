// app/entities/permission.go

package entities

import "time"

type Action string

const (
	ActionView    Action = "view"
	ActionEdit    Action = "edit"
	ActionDelete  Action = "delete"
	ActionCreate  Action = "create"
	ActionShare   Action = "share"
	ActionApprove Action = "approve"
)

// Optional: Validation helper
func (a Action) Valid() bool {
	switch a {
	case ActionView, ActionEdit, ActionDelete, ActionCreate, ActionShare, ActionApprove:
		return true
	}
	return false
}

type Permission struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       string    `gorm:"column:user_id;not null;index:idx_user_resource" json:"user_id"`
	Action       Action    `gorm:"column:action;not null;index:idx_user_resource" json:"action"`
	ResourceType string    `gorm:"column:resource_type;not null;index:idx_user_resource" json:"resource_type"`
	ResourceID   *string   `gorm:"column:resource_id" json:"resource_id,omitempty"`
	GrantedAt    time.Time `gorm:"column:granted_at;default:now()" json:"granted_at"`
	CreatedBy    *string   `gorm:"column:created_by" json:"created_by,omitempty"`
}

func (Permission) TableName() string {
	return "permissions"
}

// app/entities/decision.go

package entities

import "time"

type Decision struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	UserID       string    `gorm:"column:user_id;not null;index" json:"user_id"`
	Action       string    `gorm:"column:action;not null" json:"action"`
	ResourceType string    `gorm:"column:resource_type;not null" json:"resource_type"`
	ResourceID   *string   `gorm:"column:resource_id" json:"resource_id,omitempty"`
	Allowed      bool      `gorm:"column:allowed;not null" json:"allowed"`
	Reason       string    `gorm:"column:reason;type:text" json:"reason"`
	Timestamp    time.Time `gorm:"column:timestamp;default:now();index" json:"timestamp"`
}

func (Decision) TableName() string {
	return "decisions"
}

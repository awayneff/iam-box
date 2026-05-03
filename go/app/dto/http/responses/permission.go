// app/dto/http/requests/permission.go

package responses

import "time"

type Permission struct {
	ID           uint      `json:"id"`
	UserID       string    `json:"user_id"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   *string   `json:"resource_id"`
	GrantedAt    time.Time `json:"granted_at"`
	CreatedBy    *string   `json:"created_by"`
}

type GetPermissionsResponse struct {
	Permissions []Permission `json:"permissions"`
	Count       int          `json:"count"`
}

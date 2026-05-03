// app/dto/http/requests/permission.go

package requests

type UnifiedPermissionRequest struct {
	UserID       string  `json:"user_id" validate:"required"`
	Action       string  `json:"action" validate:"required"`
	ResourceType string  `json:"resource_type" validate:"required"`
	ResourceID   *string `json:"resource_id"`
}

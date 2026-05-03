// app/dto/http/responses/decision.go

package responses

import "time"

type Decision struct {
	ID           uint      `json:"id"`
	UserID       string    `json:"user_id"`
	Action       string    `json:"action"`
	ResourceType string    `json:"resource_type"`
	ResourceID   *string   `json:"resource_id"`
	Allowed      bool      `json:"allowed"`
	Reason       string    `json:"reason"`
	Timestamp    time.Time `json:"timestamp"`
}

type DecisiosnResponse struct {
	Decisions []Decision `json:"decision"`
	Count     int        `json:"count"`
}

// app/controllers/permission.go

package controllers

import (
	"context"
	"encoding/json"
	"iam-box/app/api"
	"iam-box/app/dto/http/requests"
	"iam-box/app/dto/http/responses"
	"iam-box/app/service"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

type permissionController struct {
	permissionService service.PermissionService
	validator         *validator.Validate
}

func NewPermissionController(permissionService service.PermissionService) *permissionController {
	validator := validator.New()

	return &permissionController{
		permissionService: permissionService,
		validator:         validator,
	}
}

func (c *permissionController) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.RespondWithBadMethod(w, map[string]string{"message": "Expected method POST"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req requests.UnifiedPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.validator.Struct(req); err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.permissionService.Create(ctx, req.UserID, req.Action, req.ResourceType, req.ResourceID); err != nil {
		api.RespondWithInternalError(w, err)
		return
	}

	api.RespondWithCreated(w, map[string]string{"message": "permission created"})
}

func (c *permissionController) GetByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.RespondWithBadMethod(w, map[string]string{"message": "Expected method GET"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	user_id := chi.URLParam(r, "user_id")

	permissions, err := c.permissionService.GetByUser(ctx, user_id)
	if err != nil {
		api.RespondWithInternalError(w, err)
		return
	} else if permissions == nil {
		api.RespondWithJSON(w, http.StatusOK, responses.GetPermissionsResponse{
			Permissions: []responses.Permission{},
			Count:       0,
		})
	}

	perms := make([]responses.Permission, len(*permissions))
	for i, p := range *permissions {
		perms[i] = responses.Permission{
			ID:           p.ID,
			UserID:       p.UserID,
			Action:       string(p.Action),
			ResourceType: p.ResourceType,
			ResourceID:   p.ResourceID,
			GrantedAt:    p.GrantedAt,
			CreatedBy:    p.CreatedBy,
		}
	}

	api.RespondWithJSON(w, http.StatusOK, responses.GetPermissionsResponse{
		Permissions: perms,
		Count:       len(perms),
	})
}

func (c *permissionController) Check(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		api.RespondWithBadMethod(w, map[string]string{"message": "Expected method POST"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req requests.UnifiedPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.validator.Struct(req); err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	can, err := c.permissionService.Check(ctx, req.UserID, req.Action, req.ResourceType, req.ResourceID)
	if err != nil {
		api.RespondWithInternalError(w, err)
		return
	}

	if can {
		api.RespondWithJSON(w, http.StatusOK, map[string]string{"message": "access granted"})
		return
	}

	api.RespondWithJSON(w, http.StatusForbidden, map[string]string{"message": "access rejected"})
}

func (c *permissionController) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		api.RespondWithBadMethod(w, map[string]string{"message": "Expected method DELETE"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var req requests.UnifiedPermissionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.validator.Struct(req); err != nil {
		api.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := c.permissionService.Delete(ctx, req.UserID, req.Action, req.ResourceType, req.ResourceID); err != nil {
		api.RespondWithInternalError(w, err)
		return
	}

	api.RespondWithOK(w, map[string]string{"message": "permission deleted"})
}

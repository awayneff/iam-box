// app/controllers/decision.go

package controllers

import (
	"context"
	"iam-box/app/api"
	"iam-box/app/dto/http/responses"
	errs "iam-box/app/errors"
	"iam-box/app/service"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
)

type decisionController struct {
	decisionService service.DecisionService
	validator       *validator.Validate
}

func NewDecisionController(decisionService service.DecisionService) *decisionController {
	validator := validator.New()

	return &decisionController{
		decisionService: decisionService,
		validator:       validator,
	}
}

func (c *decisionController) List(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		api.RespondWithBadMethod(w, map[string]string{"message": "Expected method GET"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	limit := api.ParseIntQuery(r, "limit", 100)
	offset := api.ParseIntQuery(r, "offset", 0)

	if limit > 1000 {
		api.RespondWithError(w, http.StatusBadRequest, errs.ErrLimitViolation.Error())
		return
	}

	decisions, err := c.decisionService.List(ctx, limit, offset)
	if err != nil {
		api.RespondWithInternalError(w, err)
		return
	}

	d := make([]responses.Decision, len(decisions))
	for i, decision := range decisions {
		d[i] = responses.Decision{
			ID:           decision.ID,
			UserID:       decision.UserID,
			Action:       decision.Action,
			ResourceType: decision.ResourceType,
			ResourceID:   decision.ResourceID,
			Allowed:      decision.Allowed,
			Reason:       decision.Reason,
			Timestamp:    decision.Timestamp,
		}
	}

	api.RespondWithCreated(w, responses.DecisiosnResponse{
		Decisions: d,
		Count:     len(d),
	})
}

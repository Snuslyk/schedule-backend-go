package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/service"
)

type OrganizationHandler struct {
	orgService *service.OrganizationService
}

func NewOrganizationHandler(orgService *service.OrganizationService) *OrganizationHandler {
	return &OrganizationHandler{orgService: orgService}
}

// Create godoc
// @Summary      Create a new organization
// @Tags         organization
// @Accept       json
// @Produce      json
// @Param        input  body      domain.CreateOrganizationDTO  true  "Organization payload"
// @Success      201    {object}  domain.Organization
// @Failure      400    {object}  domain.ErrorResponse
// @Router       /api/organization [post]
func (h *OrganizationHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto domain.CreateOrganizationDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if dto.Name == "" {
		BadRequest(w, "Organization name is required")
		return
	}

	org, err := h.orgService.Create(r.Context(), dto)
	if err != nil {
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, org)
}

// CreateMany godoc
// @Summary      Create multiple organizations
// @Tags         organization
// @Accept       json
// @Produce      json
// @Param        input  body      []domain.CreateOrganizationDTO  true  "List of organizations"
// @Success      201    {object}  map[string]int
// @Failure      400    {object}  domain.ErrorResponse
// @Router       /api/organization/bulk [post]
func (h *OrganizationHandler) CreateMany(w http.ResponseWriter, r *http.Request) {
	var dtos []domain.CreateOrganizationDTO
	if err := json.NewDecoder(r.Body).Decode(&dtos); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	count, err := h.orgService.CreateMany(r.Context(), dtos)
	if err != nil {
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, map[string]int{"count": count})
}

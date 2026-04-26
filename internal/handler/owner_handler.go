package handler

import (
	"net/http"
	"strconv"

	"github.com/Snuslyk/schedule-backend-go/internal/service"
)

type OwnerHandler struct {
	ownerService *service.OwnerService
}

func NewOwnerHandler(ownerService *service.OwnerService) *OwnerHandler {
	return &OwnerHandler{ownerService: ownerService}
}

// FindByName godoc
// @Summary      Search groups and teachers by name
// @Tags         owner
// @Produce      json
// @Param        name      query     string  false  "Search query"
// @Param        quantity  query     int     false  "Max items (default 5, max 10)"
// @Success      200       {array}   domain.Owner
// @Failure      400       {object}  domain.ErrorResponse
// @Router       /api/owner [get]
func (h *OwnerHandler) FindByName(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	quantityStr := r.URL.Query().Get("quantity")

	quantity := 5
	if quantityStr != "" {
		if q, err := strconv.Atoi(quantityStr); err == nil {
			quantity = q
		}
	}

	if quantity > 10 {
		BadRequest(w, "Quantity must be less than or equal to 10")
		return
	}

	owners, err := h.ownerService.FindByName(r.Context(), name, quantity)
	if err != nil {
		BadRequest(w, err.Error())
		return
	}

	JSON(w, http.StatusOK, owners)
}

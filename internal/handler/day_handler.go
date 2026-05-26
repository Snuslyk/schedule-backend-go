package handler

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Snuslyk/schedule-backend-go/internal/service"
	"github.com/go-chi/chi/v5"
)

type DayHandler struct {
	dayService *service.DayService
}

func NewDayHandler(dayService *service.DayService) *DayHandler {
	return &DayHandler{dayService: dayService}
}

// GetByID godoc
// @Summary      Get day by ID
// @Tags         day
// @Produce      json
// @Param        id   path      int  true  "Day ID"
// @Success      200  {object}  domain.Day
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Router       /api/day/{id} [get]
func (h *DayHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(w, "Invalid day id")
		return
	}

	day, err := h.dayService.GetByID(r.Context(), id)
	if err != nil {
		InternalError(w, err.Error())
		return
	}
	if day == nil {
		NotFound(w, fmt.Sprintf("Day with id %d not found", id))
		return
	}

	JSON(w, http.StatusOK, day)
}

// DeleteByID godoc
// @Summary      Delete day by ID
// @Tags         day
// @Param        id   path      int  true  "Day ID"
// @Success      204  "No Content"
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Router       /api/day/{id} [delete]
func (h *DayHandler) DeleteByID(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(w, "Invalid day id")
		return
	}

	deleted, err := h.dayService.DeleteByID(r.Context(), id)
	if err != nil {
		InternalError(w, err.Error())
		return
	}
	if !deleted {
		NotFound(w, fmt.Sprintf("Day with id %d not found", id))
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

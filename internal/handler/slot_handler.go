package handler

import (
	"encoding/json"
	"net/http"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/repository/postgres"
)

type SlotHandler struct {
	slotRepo *postgres.SlotRepository
}

func NewSlotHandler(slotRepo *postgres.SlotRepository) *SlotHandler {
	return &SlotHandler{slotRepo: slotRepo}
}

// FindAll godoc
// @Summary      List bell timetable slots
// @Tags         slot
// @Produce      json
// @Success      200  {array}   domain.Slot
// @Failure      500  {object}  domain.ErrorResponse
// @Router       /api/slot [get]
func (h *SlotHandler) FindAll(w http.ResponseWriter, r *http.Request) {
	slots, err := h.slotRepo.FindAll(r.Context())
	if err != nil {
		InternalError(w, err.Error())
		return
	}
	JSON(w, http.StatusOK, slots)
}

// Create godoc
// @Summary      Create bell slot
// @Tags         slot
// @Accept       json
// @Produce      json
// @Param        input  body      domain.CreateSlotDTO  true  "Slot payload"
// @Success      201    {object}  domain.Slot
// @Failure      400    {object}  domain.ErrorResponse
// @Router       /api/slot [post]
func (h *SlotHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto domain.CreateSlotDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if dto.SlotNumber <= 0 || dto.StartMinute < 0 || dto.EndMinute <= dto.StartMinute {
		BadRequest(w, "Invalid slot number or start/end minutes")
		return
	}

	slot := &domain.Slot{
		SlotNumber:  dto.SlotNumber,
		StartMinute: dto.StartMinute,
		EndMinute:   dto.EndMinute,
		DayOfWeek:   dto.DayOfWeek,
	}

	if err := h.slotRepo.Create(r.Context(), slot); err != nil {
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, slot)
}

// CreateMany godoc
// @Summary      Bulk create bell slots
// @Tags         slot
// @Accept       json
// @Produce      json
// @Param        input  body      []domain.CreateSlotDTO  true  "List of slots"
// @Success      201    {object}  map[string]int
// @Failure      400    {object}  domain.ErrorResponse
// @Router       /api/slot/bulk [post]
func (h *SlotHandler) CreateMany(w http.ResponseWriter, r *http.Request) {
	var dtos []domain.CreateSlotDTO
	if err := json.NewDecoder(r.Body).Decode(&dtos); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if err := h.slotRepo.CreateMany(r.Context(), dtos); err != nil {
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, map[string]int{"count": len(dtos)})
}

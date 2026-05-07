package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/service"
	"github.com/Snuslyk/schedule-backend-go/internal/utils"
)

type ScheduleHandler struct {
	scheduleService *service.ScheduleService
}

func NewScheduleHandler(scheduleService *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{scheduleService: scheduleService}
}

// GetGroupDay godoc
// @Summary      Get group schedule for a day
// @Tags         schedule
// @Produce      json
// @Param        id   query     int     true   "Group ID"
// @Param        day  query     string  false  "Date YYYY-MM-DD"
// @Success      200  {object}  domain.DayScheduleResponse
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Router       /api/schedule/day/group [get]
func (h *ScheduleHandler) GetGroupDay(w http.ResponseWriter, r *http.Request) {
	dayStr := r.URL.Query().Get("day")
	idStr := r.URL.Query().Get("id")

	date, err := utils.ParseDate(dayStr)
	if err != nil {
		date = time.Now()
	}

	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(w, "Invalid or missing group id")
		return
	}

	resp, err := h.scheduleService.GetGroupDay(r.Context(), date, groupID)
	if err != nil {
		NotFound(w, err.Error())
		return
	}

	JSON(w, http.StatusOK, resp)
}

// GetGroupWeek godoc
// @Summary      Get group schedule for a week (cached in Redis)
// @Tags         schedule
// @Produce      json
// @Param        id    query     int     true   "Group ID"
// @Param        week  query     string  false  "Date YYYY-MM-DD"
// @Success      200   {object}  domain.WeekScheduleResponse
// @Failure      400   {object}  domain.ErrorResponse
// @Failure      404   {object}  domain.ErrorResponse
// @Router       /api/schedule/week/group [get]
func (h *ScheduleHandler) GetGroupWeek(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")
	idStr := r.URL.Query().Get("id")

	date, err := utils.ParseDate(weekStr)
	if err != nil {
		date = time.Now()
	}

	groupID, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(w, "Invalid or missing group id")
		return
	}

	resp, err := h.scheduleService.GetGroupWeek(r.Context(), date, groupID)
	if err != nil {
		NotFound(w, err.Error())
		return
	}

	JSON(w, http.StatusOK, resp)
}

// GetTeacherDay godoc
// @Summary      Get teacher schedule for a day (direct SELECT)
// @Tags         schedule
// @Produce      json
// @Param        id   query     int     true   "Teacher ID"
// @Param        day  query     string  false  "Date YYYY-MM-DD"
// @Success      200  {object}  domain.DayScheduleResponse
// @Failure      400  {object}  domain.ErrorResponse
// @Failure      404  {object}  domain.ErrorResponse
// @Router       /api/schedule/day/teacher [get]
func (h *ScheduleHandler) GetTeacherDay(w http.ResponseWriter, r *http.Request) {
	dayStr := r.URL.Query().Get("day")
	idStr := r.URL.Query().Get("id")

	date, err := utils.ParseDate(dayStr)
	if err != nil {
		date = time.Now()
	}

	teacherID, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(w, "Invalid or missing teacher id")
		return
	}

	resp, err := h.scheduleService.GetTeacherDay(r.Context(), date, teacherID)
	if err != nil {
		NotFound(w, err.Error())
		return
	}

	JSON(w, http.StatusOK, resp)
}

// GetTeacherWeek godoc
// @Summary      Get teacher schedule for a week (cached in Redis)
// @Tags         schedule
// @Produce      json
// @Param        id    query     int     true   "Teacher ID"
// @Param        week  query     string  false  "Date YYYY-MM-DD"
// @Success      200   {object}  domain.WeekScheduleResponse
// @Failure      400   {object}  domain.ErrorResponse
// @Failure      404   {object}  domain.ErrorResponse
// @Router       /api/schedule/week/teacher [get]
func (h *ScheduleHandler) GetTeacherWeek(w http.ResponseWriter, r *http.Request) {
	weekStr := r.URL.Query().Get("week")
	idStr := r.URL.Query().Get("id")

	date, err := utils.ParseDate(weekStr)
	if err != nil {
		date = time.Now()
	}

	teacherID, err := strconv.Atoi(idStr)
	if err != nil {
		BadRequest(w, "Invalid or missing teacher id")
		return
	}

	resp, err := h.scheduleService.GetTeacherWeek(r.Context(), date, teacherID)
	if err != nil {
		NotFound(w, err.Error())
		return
	}

	JSON(w, http.StatusOK, resp)
}

// Create godoc
// @Summary      Create academic period
// @Tags         schedule
// @Accept       json
// @Produce      json
// @Param        input  body      domain.CreateScheduleDTO  true  "Period payload"
// @Success      201    {object}  domain.CreateScheduleDTO
// @Failure      400    {object}  domain.ErrorResponse
// @Router       /api/schedule [post]
func (h *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto domain.CreateScheduleDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	if err := h.scheduleService.CreatePeriod(r.Context(), dto); err != nil {
		InternalError(w, err.Error())
		return
	}

	JSON(w, http.StatusCreated, dto)
}

// CreateMany godoc
// @Summary      Bulk create academic periods
// @Tags         schedule
// @Accept       json
// @Produce      json
// @Param        input  body      []domain.CreateScheduleDTO  true  "List of periods"
// @Success      201    {object}  map[string]int
// @Failure      400    {object}  domain.ErrorResponse
// @Router       /api/schedule/bulk [post]
func (h *ScheduleHandler) CreateMany(w http.ResponseWriter, r *http.Request) {
	var dtos []domain.CreateScheduleDTO
	if err := json.NewDecoder(r.Body).Decode(&dtos); err != nil {
		BadRequest(w, "Invalid request body")
		return
	}

	for _, dto := range dtos {
		if err := h.scheduleService.CreatePeriod(r.Context(), dto); err != nil {
			InternalError(w, err.Error())
			return
		}
	}

	JSON(w, http.StatusCreated, map[string]int{"count": len(dtos)})
}

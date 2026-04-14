package service

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Snuslyk/schedule-backend-go/internal/cache"
	"github.com/Snuslyk/schedule-backend-go/internal/domain"
	"github.com/Snuslyk/schedule-backend-go/internal/repository/postgres"
	"github.com/Snuslyk/schedule-backend-go/internal/utils"
)

type ScheduleService struct {
	lessonRepo  *postgres.LessonRepository
	replaceRepo *postgres.ReplaceRepository
	slotRepo    *postgres.SlotRepository
	periodRepo  *postgres.PeriodRepository
	groupRepo   *postgres.GroupRepository
	teacherRepo *postgres.TeacherRepository
	cache       *cache.ScheduleCache
}

func NewScheduleService(
	lessonRepo *postgres.LessonRepository,
	replaceRepo *postgres.ReplaceRepository,
	slotRepo *postgres.SlotRepository,
	periodRepo *postgres.PeriodRepository,
	groupRepo *postgres.GroupRepository,
	teacherRepo *postgres.TeacherRepository,
	cache *cache.ScheduleCache,
) *ScheduleService {
	return &ScheduleService{
		lessonRepo:  lessonRepo,
		replaceRepo: replaceRepo,
		slotRepo:    slotRepo,
		periodRepo:  periodRepo,
		groupRepo:   groupRepo,
		teacherRepo: teacherRepo,
		cache:       cache,
	}
}

func (s *ScheduleService) GetGroupDay(ctx context.Context, date time.Time, groupID int) (*domain.DayScheduleResponse, error) {
	weekResp, err := s.GetGroupWeek(ctx, date, groupID)
	if err != nil {
		return nil, err
	}

	dayIdx := utils.GetWeekDayIndex(date)
	if dayIdx >= 0 && dayIdx < len(weekResp.Days) {
		return &weekResp.Days[dayIdx], nil
	}

	return &domain.DayScheduleResponse{
		Slots:   []domain.SlotResponse{},
		Lessons: []domain.LessonResponse{},
	}, nil
}

func (s *ScheduleService) GetGroupWeek(ctx context.Context, weekDate time.Time, groupID int) (*domain.WeekScheduleResponse, error) {
	startOfWeek := utils.StartOfWeek(weekDate)
	cacheKey := fmt.Sprintf("schedule:group:week:%d:%s", groupID, startOfWeek.Format("2006-01-02"))

	var cached domain.WeekScheduleResponse
	if utils.IsSameWeek(time.Now(), weekDate) && s.cache.Get(ctx, cacheKey, &cached) {
		return &cached, nil
	}

	group, err := s.groupRepo.FindByID(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if group == nil {
		return nil, fmt.Errorf("group with id %d not found", groupID)
	}

	activePeriod, _ := s.periodRepo.FindActive(ctx)
	parity := utils.GetWeekParity(startOfWeek)
	mode := "PARITY"
	if activePeriod != nil && activePeriod.Mode != "" {
		mode = activePeriod.Mode
	}

	weekParity := parity
	if mode == "OTHER" {
		weekParity = "OTHER"
	}

	days := make([]domain.DayScheduleResponse, 7)
	for i := range 7 {
		currentDayDate := startOfWeek.AddDate(0, 0, i)

		if activePeriod != nil && !utils.IsDateInRange(currentDayDate, activePeriod.StartDate, activePeriod.EndDate) {
			days[i] = domain.DayScheduleResponse{
				Slots:   []domain.SlotResponse{},
				Lessons: []domain.LessonResponse{},
			}
			continue
		}

		daySlots, err := s.slotRepo.FindForDay(ctx, i)
		if err != nil {
			return nil, err
		}

		slotsResp := make([]domain.SlotResponse, len(daySlots))
		for idx, sl := range daySlots {
			slotsResp[idx] = domain.SlotResponse{
				Start: sl.StartMinute,
				End:   sl.EndMinute,
			}
		}

		lessons, err := s.lessonRepo.FindByGroupAndDay(ctx, groupID, i, weekParity)
		if err != nil {
			return nil, err
		}

		replaces, err := s.replaceRepo.FindByGroupAndDate(ctx, groupID, currentDayDate)
		if err != nil {
			return nil, err
		}

		lessonsMap := make(map[int]domain.LessonResponse)
		for _, l := range lessons {
			lessonsMap[l.SlotNumber] = domain.LessonResponse{
				ID:          l.ID,
				Classroom:   l.Classroom,
				SlotNumber:  l.SlotNumber,
				SlotLength:  l.SlotLength,
				IsAvailable: l.IsAvailable,
				Subject:     &domain.NamedObject{Name: l.SubjectName},
				Teacher:     &domain.NamedObject{Name: l.TeacherName},
			}
		}

		for _, r := range replaces {
			lessonsMap[r.SlotNumber] = domain.LessonResponse{
				ID:          r.ID,
				Classroom:   r.Classroom,
				SlotNumber:  r.SlotNumber,
				SlotLength:  1,
				IsAvailable: r.IsAvailable,
				Subject:     &domain.NamedObject{Name: r.SubjectName},
				Teacher:     &domain.NamedObject{Name: r.TeacherName},
			}
		}

		lessonsResp := make([]domain.LessonResponse, 0, len(lessonsMap))
		for s := range 20 {
			slotNum := s + 1
			if item, exists := lessonsMap[slotNum]; exists {
				lessonsResp = append(lessonsResp, item)
			}
		}

		days[i] = domain.DayScheduleResponse{
			Slots:   slotsResp,
			Lessons: lessonsResp,
		}
	}

	resp := &domain.WeekScheduleResponse{
		Type:     parity,
		WeekType: strings.ToLower(parity),
		Days:     days,
	}

	if utils.IsSameWeek(time.Now(), weekDate) {
		_ = s.cache.Set(ctx, cacheKey, resp, utils.DurationUntilEndOfDay(time.Now()))
		_ = s.cache.Set(ctx, fmt.Sprintf("group_week:%d", groupID), resp, utils.DurationUntilEndOfDay(time.Now()))
	}

	return resp, nil
}

func (s *ScheduleService) GetTeacherDay(ctx context.Context, date time.Time, teacherID int) (*domain.DayScheduleResponse, error) {
	weekResp, err := s.GetTeacherWeek(ctx, date, teacherID)
	if err != nil {
		return nil, err
	}

	dayIdx := utils.GetWeekDayIndex(date)
	if dayIdx >= 0 && dayIdx < len(weekResp.Days) {
		return &weekResp.Days[dayIdx], nil
	}

	return &domain.DayScheduleResponse{
		Slots:   []domain.SlotResponse{},
		Lessons: []domain.LessonResponse{},
	}, nil
}

func (s *ScheduleService) GetTeacherWeek(ctx context.Context, weekDate time.Time, teacherID int) (*domain.WeekScheduleResponse, error) {
	startOfWeek := utils.StartOfWeek(weekDate)
	cacheKey := fmt.Sprintf("schedule:teacher:week:%d:%s", teacherID, startOfWeek.Format("2006-01-02"))

	var cached domain.WeekScheduleResponse
	if utils.IsSameWeek(time.Now(), weekDate) && s.cache.Get(ctx, cacheKey, &cached) {
		return &cached, nil
	}

	teacher, err := s.teacherRepo.FindByID(ctx, teacherID)
	if err != nil {
		return nil, err
	}
	if teacher == nil {
		return nil, fmt.Errorf("teacher with id %d not found", teacherID)
	}

	activePeriod, _ := s.periodRepo.FindActive(ctx)
	parity := utils.GetWeekParity(startOfWeek)

	days := make([]domain.DayScheduleResponse, 7)
	for i := range 7 {
		currentDayDate := startOfWeek.AddDate(0, 0, i)

		if activePeriod != nil && !utils.IsDateInRange(currentDayDate, activePeriod.StartDate, activePeriod.EndDate) {
			days[i] = domain.DayScheduleResponse{
				Slots:   []domain.SlotResponse{},
				Lessons: []domain.LessonResponse{},
			}
			continue
		}

		daySlots, err := s.slotRepo.FindForDay(ctx, i)
		if err != nil {
			return nil, err
		}

		slotsResp := make([]domain.SlotResponse, len(daySlots))
		for idx, sl := range daySlots {
			slotsResp[idx] = domain.SlotResponse{
				Start: sl.StartMinute,
				End:   sl.EndMinute,
			}
		}

		lessons, err := s.lessonRepo.FindByTeacherAndDay(ctx, teacherID, i, parity)
		if err != nil {
			return nil, err
		}

		replaces, err := s.replaceRepo.FindByTeacherAndDate(ctx, teacherID, currentDayDate)
		if err != nil {
			return nil, err
		}

		lessonsMap := make(map[int]domain.LessonResponse)
		for _, l := range lessons {
			lessonsMap[l.SlotNumber] = domain.LessonResponse{
				ID:          l.ID,
				Classroom:   l.Classroom,
				SlotNumber:  l.SlotNumber,
				SlotLength:  l.SlotLength,
				IsAvailable: l.IsAvailable,
				Subject:     &domain.NamedObject{Name: l.SubjectName},
				GroupName:   l.GroupName,
			}
		}

		for _, r := range replaces {
			lessonsMap[r.SlotNumber] = domain.LessonResponse{
				ID:          r.ID,
				Classroom:   r.Classroom,
				SlotNumber:  r.SlotNumber,
				SlotLength:  1,
				IsAvailable: r.IsAvailable,
				Subject:     &domain.NamedObject{Name: r.SubjectName},
				GroupName:   r.GroupName,
			}
		}

		lessonsResp := make([]domain.LessonResponse, 0, len(lessonsMap))
		for s := range 20 {
			slotNum := s + 1
			if item, exists := lessonsMap[slotNum]; exists {
				lessonsResp = append(lessonsResp, item)
			}
		}

		days[i] = domain.DayScheduleResponse{
			Slots:   slotsResp,
			Lessons: lessonsResp,
		}
	}

	resp := &domain.WeekScheduleResponse{
		Type:     parity,
		WeekType: strings.ToLower(parity),
		Days:     days,
	}

	if utils.IsSameWeek(time.Now(), weekDate) {
		_ = s.cache.Set(ctx, cacheKey, resp, utils.DurationUntilEndOfDay(time.Now()))
		_ = s.cache.Set(ctx, fmt.Sprintf("teacher_week:%d", teacherID), resp, utils.DurationUntilEndOfDay(time.Now()))
	}

	return resp, nil
}

func (s *ScheduleService) CreatePeriod(ctx context.Context, dto domain.CreateScheduleDTO) error {
	name := dto.Name
	if name == "" {
		name = "Semester"
	}
	mode := dto.Mode
	if mode == "" {
		mode = "PARITY"
	}

	p := &domain.AcademicPeriod{
		Name:      name,
		StartDate: dto.Start,
		EndDate:   dto.End,
		Mode:      mode,
		IsActive:  true,
	}
	return s.periodRepo.Create(ctx, p)
}

package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Snuslyk/schedule-backend-go/internal/cache"
	"github.com/Snuslyk/schedule-backend-go/internal/config"
	"github.com/Snuslyk/schedule-backend-go/internal/handler"
	"github.com/Snuslyk/schedule-backend-go/internal/repository/postgres"
	"github.com/Snuslyk/schedule-backend-go/internal/service"
	"github.com/Snuslyk/schedule-backend-go/internal/storage"
)

// @title           Schedule API
// @version         1.0
// @description     Educational timetable and schedule API with flat lessons, Redis caching, and S3 storage.

// @host      localhost:4242
// @BasePath  /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and the JWT token.

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load configuration: %v", err)
	}

	db, err := postgres.NewDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to postgres: %v", err)
	}
	defer db.Close()
	log.Println("[INFO] connected to postgres successfully")

	redisClient, err := cache.NewRedisClient(cfg.RedisURL, cfg.RedisPassword)
	if err != nil {
		log.Printf("[WARNING] could not connect to redis: %v (cache will degrade gracefully)", err)
	} else {
		defer redisClient.Client().Close()
		log.Println("[INFO] connected to redis successfully")
	}

	scheduleCache := cache.NewScheduleCache(redisClient)

	var s3Storage storage.Storage
	s3Inst, err := storage.NewS3Storage(cfg)
	if err != nil {
		log.Printf("[WARNING] could not initialize s3 storage: %v", err)
	} else {
		s3Storage = s3Inst
	}
	avatarService := storage.NewAvatarService(s3Storage)

	userRepo := postgres.NewUserRepository(db)
	groupRepo := postgres.NewGroupRepository(db)
	teacherRepo := postgres.NewTeacherRepository(db)
	subjectRepo := postgres.NewSubjectRepository(db)
	slotRepo := postgres.NewSlotRepository(db)
	periodRepo := postgres.NewPeriodRepository(db)
	lessonRepo := postgres.NewLessonRepository(db)
	replaceRepo := postgres.NewReplaceRepository(db)
	ownerRepo := postgres.NewOwnerRepository(db)
	orgRepo := postgres.NewOrganizationRepository(db)
	wtRepo := postgres.NewWeekTemplateRepository(db)
	dayRepo := postgres.NewDayRepository(db)

	authService := service.NewAuthService(userRepo, avatarService, cfg)
	scheduleService := service.NewScheduleService(lessonRepo, replaceRepo, slotRepo, periodRepo, groupRepo, teacherRepo, scheduleCache)
	replaceService := service.NewReplaceService(replaceRepo, groupRepo, scheduleCache)
	ownerService := service.NewOwnerService(ownerRepo)
	orgService := service.NewOrganizationService(orgRepo)
	wtService := service.NewWeekTemplateService(wtRepo, groupRepo, scheduleCache)
	dayService := service.NewDayService(dayRepo)

	router := handler.SetupRouter(handler.RouterParams{
		Config:              cfg,
		AuthHandler:         handler.NewAuthHandler(authService, avatarService),
		GroupHandler:        handler.NewGroupHandler(groupRepo),
		TeacherHandler:      handler.NewTeacherHandler(teacherRepo),
		SubjectHandler:      handler.NewSubjectHandler(subjectRepo),
		SlotHandler:         handler.NewSlotHandler(slotRepo),
		LessonHandler:       handler.NewLessonHandler(lessonRepo),
		ReplaceHandler:      handler.NewReplaceHandler(replaceService),
		ScheduleHandler:     handler.NewScheduleHandler(scheduleService),
		OwnerHandler:        handler.NewOwnerHandler(ownerService),
		OrganizationHandler: handler.NewOrganizationHandler(orgService),
		WeekTemplateHandler: handler.NewWeekTemplateHandler(wtService),
		DayHandler:          handler.NewDayHandler(dayService),
	})

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("[INFO] server listening on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("[INFO] shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("server forced to shutdown: %v", err)
	}

	log.Println("[INFO] server exited gracefully")
}

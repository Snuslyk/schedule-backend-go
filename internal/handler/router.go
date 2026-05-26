package handler

import (
	"net/http"

	"github.com/Snuslyk/schedule-backend-go/internal/config"
	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

type RouterParams struct {
	Config          *config.Config
	AuthHandler     *AuthHandler
	GroupHandler    *GroupHandler
	TeacherHandler  *TeacherHandler
	SubjectHandler  *SubjectHandler
	SlotHandler     *SlotHandler
	LessonHandler   *LessonHandler
	ReplaceHandler      *ReplaceHandler
	ScheduleHandler     *ScheduleHandler
	OwnerHandler        *OwnerHandler
	OrganizationHandler *OrganizationHandler
	WeekTemplateHandler *WeekTemplateHandler
	DayHandler          *DayHandler
}

func SetupRouter(p RouterParams) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.Logger)
	r.Use(chimiddleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{p.Config.CORSOrigin, "http://localhost:3000", "http://127.0.0.1:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(AuthMiddleware(p.Config))

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		JSON(w, http.StatusOK, map[string]string{"message": "pong"})
	})

	r.Get("/reference", ScalarHandler)
	r.Get("/reference/*", ScalarHandler)
	r.Get("/docs", ScalarHandler)
	r.Get("/openapi.json", OpenAPIHandler)
	r.Get("/swagger.json", OpenAPIHandler)

	r.Route("/api", func(api chi.Router) {
		api.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
			JSON(w, http.StatusOK, map[string]string{"message": "pong"})
		})
		api.Get("/reference", ScalarHandler)
		api.Get("/reference/*", ScalarHandler)
		api.Get("/docs", ScalarHandler)
		api.Get("/openapi.json", OpenAPIHandler)
		api.Get("/swagger.json", OpenAPIHandler)

		registerAuthRoutes := func(sub chi.Router) {
			sub.Post("/register", p.AuthHandler.Register)
			sub.Post("/login", p.AuthHandler.Login)
			sub.Post("/refresh", p.AuthHandler.Refresh)
			sub.Post("/logout", p.AuthHandler.Logout)
			sub.Get("/profile", RequireAuth(p.AuthHandler.GetProfile))
			sub.Post("/role", RequireAdmin(p.AuthHandler.AddRole))
			sub.Post("/avatar", RequireAuth(p.AuthHandler.UploadAvatar))
		}

		api.Route("/auth", registerAuthRoutes)
		api.Route("/user", registerAuthRoutes)

		api.Route("/group", func(g chi.Router) {
			g.Get("/id/{id}", p.GroupHandler.GetByID)
			g.Get("/name/{name}", p.GroupHandler.GetByName)
			g.Post("/", p.GroupHandler.Create)
			g.Post("/bulk", p.GroupHandler.CreateMany)
			g.Delete("/id/{id}", p.GroupHandler.DeleteByID)
			g.Delete("/name/{name}", p.GroupHandler.DeleteByName)
		})

		api.Route("/teacher", func(t chi.Router) {
			t.Get("/", p.TeacherHandler.FindAll)
			t.Get("/{id}", p.TeacherHandler.FindOne)
			t.Post("/", p.TeacherHandler.Create)
			t.Post("/bulk", p.TeacherHandler.CreateMany)
			t.Delete("/{id}", p.TeacherHandler.Remove)
		})

		api.Route("/subject", func(s chi.Router) {
			s.Get("/", p.SubjectHandler.FindAll)
			s.Post("/", p.SubjectHandler.Create)
			s.Post("/bulk", p.SubjectHandler.CreateMany)
		})

		api.Route("/slot", func(sl chi.Router) {
			sl.Get("/", p.SlotHandler.FindAll)
			sl.Post("/", p.SlotHandler.Create)
			sl.Post("/bulk", p.SlotHandler.CreateMany)
		})

		api.Route("/lesson", func(l chi.Router) {
			l.Post("/", p.LessonHandler.Create)
			l.Post("/bulk", p.LessonHandler.CreateMany)
			l.Delete("/{id}", p.LessonHandler.DeleteByID)
		})

		api.Route("/replace", func(rep chi.Router) {
			rep.Post("/{group}", p.ReplaceHandler.Create)
			rep.Post("/{group}/bulk", p.ReplaceHandler.CreateMany)
		})

		api.Route("/schedule", func(sc chi.Router) {
			sc.Get("/day/group", p.ScheduleHandler.GetGroupDay)
			sc.Get("/week/group", p.ScheduleHandler.GetGroupWeek)
			sc.Get("/day/teacher", p.ScheduleHandler.GetTeacherDay)
			sc.Get("/week/teacher", p.ScheduleHandler.GetTeacherWeek)
			sc.Post("/", p.ScheduleHandler.Create)
			sc.Post("/bulk", p.ScheduleHandler.CreateMany)
		})

		api.Route("/owner", func(o chi.Router) {
			o.Get("/", p.OwnerHandler.FindByName)
		})

		api.Route("/organization", func(org chi.Router) {
			org.Post("/", p.OrganizationHandler.Create)
			org.Post("/bulk", p.OrganizationHandler.CreateMany)
		})

		api.Route("/week-template", func(wt chi.Router) {
			wt.Post("/{id}", p.WeekTemplateHandler.Create)
			wt.Post("/{id}/bulk", p.WeekTemplateHandler.CreateMany)
			wt.Delete("/{id}", p.WeekTemplateHandler.DeleteByID)
		})

		api.Route("/day", func(d chi.Router) {
			d.Get("/{id}", p.DayHandler.GetByID)
			d.Delete("/{id}", p.DayHandler.DeleteByID)
		})
	})

	return r
}

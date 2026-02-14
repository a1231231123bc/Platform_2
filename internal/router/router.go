package router

import (
	"encoding/json"
	"net/http"

	"platform/internal/handler"
	mw "platform/internal/middleware"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Deps struct {
	AuthHandler          *handler.AuthHandler
	UsersHandler         *handler.UsersHandler
	OrganizationsHandler *handler.OrganizationsHandler
	ContractorsHandler   *handler.ContractorsHandler
	JobsHandler          *handler.JobsHandler
	MatchingHandler      *handler.MatchingHandler
	ResponsesHandler     *handler.ResponsesHandler
	AssignmentsHandler   *handler.AssignmentsHandler
	RatingsHandler       *handler.RatingsHandler
	JWTSecret            string
}

func New(deps Deps) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(corsMiddleware)

	r.Get("/", healthCheck)

	// Auth routes (public)
	r.Post("/auth/register", deps.AuthHandler.Register)
	r.Post("/auth/login", deps.AuthHandler.Login)
	r.Get("/responses/{token}", deps.ResponsesHandler.GetByToken)
	r.Post("/responses/{token}/accept", deps.ResponsesHandler.Accept)
	r.Post("/responses/{token}/decline", deps.ResponsesHandler.Decline)

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(mw.JWTAuth(deps.JWTSecret))

		r.Get("/auth/me", deps.AuthHandler.GetMe)
		r.Get("/users", deps.UsersHandler.List)
		r.Get("/users/{id}", deps.UsersHandler.GetByID)
		r.Get("/organizations/{id}", deps.OrganizationsHandler.GetByID)
		r.Get("/contractors", deps.ContractorsHandler.List)
		r.Get("/contractors/{id}", deps.ContractorsHandler.GetByID)
		r.Get("/contractors/{id}/history", deps.ContractorsHandler.History)
		r.Get("/jobs", deps.JobsHandler.List)
		r.Get("/jobs/{id}", deps.JobsHandler.GetByID)
		r.Get("/matching/jobs/{jobId}/contractors", deps.MatchingHandler.MatchByJobID)
		r.Get("/assignments/job/{jobId}", deps.AssignmentsHandler.ListByJob)
		r.Get("/ratings/contractor/{contractorId}", deps.RatingsHandler.ListByContractor)
		r.Get("/ratings/job/{jobId}", deps.RatingsHandler.ListByJob)

		r.Group(func(r chi.Router) {
			r.Use(mw.RequireRoles("ADMIN", "MANAGER"))
			r.Patch("/users/{id}", deps.UsersHandler.Update)
			r.Patch("/organizations/{id}", deps.OrganizationsHandler.Update)
			r.Post("/contractors", deps.ContractorsHandler.Create)
			r.Patch("/contractors/{id}", deps.ContractorsHandler.Update)
			r.Post("/contractors/{id}/blacklist", deps.ContractorsHandler.AddToBlacklist)
			r.Post("/jobs", deps.JobsHandler.Create)
			r.Patch("/jobs/{id}", deps.JobsHandler.Update)
			r.Post("/jobs/{id}/publish", deps.JobsHandler.Publish)
			r.Post("/jobs/{id}/cancel", deps.JobsHandler.Cancel)
			r.Post("/jobs/{id}/complete", deps.JobsHandler.Complete)
			r.Post("/jobs/{id}/duplicate", deps.JobsHandler.Duplicate)
			r.Post("/assignments", deps.AssignmentsHandler.Create)
			r.Patch("/assignments/{id}/status", deps.AssignmentsHandler.UpdateStatus)
			r.Post("/ratings", deps.RatingsHandler.Create)
		})

		r.Group(func(r chi.Router) {
			r.Use(mw.RequireRoles("ADMIN"))
			r.Delete("/users/{id}", deps.UsersHandler.Delete)
		})
	})

	return r
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

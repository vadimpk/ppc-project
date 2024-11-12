package controller

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter(h *Handlers, authMiddleware, corsMiddleware func(http.Handler) http.Handler) *chi.Mux {
	r := chi.NewRouter()

	// Global middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(corsMiddleware)

	// API routes
	r.Route("/api/v1", func(r chi.Router) {
		// Auth routes - no authentication required
		r.Group(func(r chi.Router) {
			r.Post("/auth/login", h.User.Login)
			r.Post("/auth/register", h.User.Register)
		})

		// Routes requiring authentication
		r.Group(func(r chi.Router) {
			r.Use(authMiddleware)

			// Business routes
			r.Route("/businesses", func(r chi.Router) {
				r.Post("/", h.Business.Create)

				// Business-specific routes
				r.Route("/{businessID}", func(r chi.Router) {
					r.Get("/", h.Business.Get)
					r.Put("/", h.Business.Update)
					r.Patch("/appearance", h.Business.UpdateAppearance)

					// Service routes
					r.Route("/services", func(r chi.Router) {
						r.Get("/", h.Service.List)
						r.Post("/", h.Service.Create)
						r.Get("/{serviceID}", h.Service.Get)
						r.Put("/{serviceID}", h.Service.Update)
						r.Delete("/{serviceID}", h.Service.Delete)
					})

					// Employee routes
					r.Route("/employees", func(r chi.Router) {
						r.Get("/", h.Employee.List)
						r.Post("/", h.Employee.Create)

						r.Route("/{employeeID}", func(r chi.Router) {
							r.Get("/", h.Employee.Get)
							r.Put("/", h.Employee.Update)

							// Employee services
							r.Get("/services", h.Employee.ListServices)
							r.Post("/services", h.Employee.AssignServices)
							r.Delete("/services", h.Employee.RemoveServices)

							// Employee schedule
							r.Route("/schedule", func(r chi.Router) {
								// Templates
								r.Get("/templates", h.Schedule.ListTemplates)
								r.Post("/templates", h.Schedule.CreateTemplate)
								r.Put("/templates/{templateID}", h.Schedule.UpdateTemplate)
								r.Delete("/templates/{templateID}", h.Schedule.DeleteTemplate)

								// Overrides
								r.Get("/overrides", h.Schedule.ListOverrides)
								r.Post("/overrides", h.Schedule.CreateOverride)
								r.Put("/overrides/{overrideID}", h.Schedule.UpdateOverride)
								r.Delete("/overrides/{overrideID}", h.Schedule.DeleteOverride)
							})
						})
					})

					// Appointment routes
					r.Route("/appointments", func(r chi.Router) {
						r.Get("/", h.Appointment.ListByBusiness)
						r.Post("/", h.Appointment.Create)
						r.Get("/slots", h.Appointment.GetAvailableSlots)

						r.Route("/{appointmentID}", func(r chi.Router) {
							r.Get("/", h.Appointment.Get)
							r.Put("/", h.Appointment.Update)
							r.Delete("/", h.Appointment.Cancel)
						})
					})
				})
			})

			// User routes
			r.Route("/users", func(r chi.Router) {
				r.Route("/{userID}", func(r chi.Router) {
					r.Get("/", h.User.Get)
					r.Put("/", h.User.Update)
					r.Get("/appointments", h.Appointment.ListByClient)
				})
			})
		})
	})

	return r
}

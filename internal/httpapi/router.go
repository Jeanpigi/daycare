package httpapi

import (
	"net/http"

	"daycare/internal/httpapi/handlers"
	"daycare/internal/httpapi/middleware"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func NewRouter(
	// public
	authH *handlers.AuthHandler,
	bootstrapH *handlers.BootstrapHandler,

	// protected
	childrenH *handlers.ChildrenHandler,
	attH *handlers.AttendanceHandler,

	// admin protected
	adminPricingH *handlers.AdminPricingHandler,
	adminPromosH *handlers.AdminPromotionsHandler,
	adminUsersH *handlers.AdminUsersHandler,

	// middleware
	authMW *middleware.AuthMiddleware,
	adminMW *middleware.RequireAdminMiddleware,
) http.Handler {
	r := chi.NewRouter()

	// middlewares
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// health
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"ok":true}`))
	})

	// auth
	r.Post("/auth/login", authH.Login)

	// bootstrap: PUBLIC pero controlado por APP_ENV dentro del handler
	r.Post("/admin/bootstrap", bootstrapH.CreateFirstAdmin)

	// protected group (staff + admin)
	r.Group(func(pr chi.Router) {
		pr.Use(authMW.RequireAuth)

		pr.Route("/children", func(r chi.Router) {
			r.Post("/", childrenH.Create)
			r.Get("/by-document/{document}", childrenH.GetByDocument)
			r.Get("/{id}", childrenH.GetByID)
		})

		pr.Route("/attendances", func(r chi.Router) {
			r.Post("/check-in", attH.CheckInByDocument)
			r.Post("/check-out", attH.CheckOutByDocument)
		})

		// admin group (solo ADMIN)
		pr.Route("/admin", func(ar chi.Router) {
			ar.Use(adminMW.RequireAdmin)

			ar.Route("/pricing", func(r chi.Router) {
				r.Get("/", adminPricingH.GetActive)
				r.Post("/", adminPricingH.CreateAndActivate)
			})

			ar.Route("/promotions", func(r chi.Router) {
				r.Post("/", adminPromosH.Create)
				r.Post("/{id}/activate", adminPromosH.Activate)
				r.Post("/{id}/deactivate", adminPromosH.Deactivate)
			})

			ar.Post("/users", adminUsersH.Create)
		})
	})

	return r
}

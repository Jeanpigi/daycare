package main

import (
	"log"
	"net/http"

	"daycare/internal/config"
	"daycare/internal/db"
	"daycare/internal/httpapi"
	"daycare/internal/httpapi/handlers"
	"daycare/internal/httpapi/middleware"
	mysqlrepo "daycare/internal/repository/mysql"
	"daycare/internal/service"
)

func main() {
	// config
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}

	// db (usa DSN del config)
	conn, err := db.Open(cfg.DSN())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// repos
	userRepo := mysqlrepo.NewUserRepo(conn)
	childrenRepo := mysqlrepo.NewChildrenRepo(conn)
	attendanceRepo := mysqlrepo.NewAttendanceRepo(conn)
	pricingRepo := mysqlrepo.NewPricingRepo(conn)
	promoRepo := mysqlrepo.NewPromotionRepo(conn)
	auditRepo := mysqlrepo.NewAuditRepo(conn)

	// services (según tus firmas reales)
	authSvc := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTTTLMinutes)
	pricingSvc := service.NewPricingService()
	childrenSvc := service.NewChildrenService(childrenRepo)
	attendanceSvc := service.NewAttendanceService(childrenRepo, attendanceRepo, pricingRepo, promoRepo, pricingSvc)
	adminSvc := service.NewAdminService(pricingRepo, promoRepo, auditRepo)
	bootstrapSvc := service.NewBootstrapService(userRepo, pricingRepo, promoRepo)
	adminUsersSvc := service.NewAdminUsersService(userRepo)

	// handlers
	authH := handlers.NewAuthHandler(authSvc)
	childrenH := handlers.NewChildrenHandler(childrenSvc)
	attH := handlers.NewAttendanceHandler(attendanceSvc)

	// AdminPricingHandler pide (adminSvc, pricingRepo) según tu firma real
	adminPricingH := handlers.NewAdminPricingHandler(adminSvc, pricingRepo)

	// AdminPromotionsHandler: mantengo como lo tenías.
	// Si te da error de argumentos aquí, sacamos su firma y lo ajustamos igual que pricing.
	adminPromosH := handlers.NewAdminPromotionsHandler(adminSvc)

	bootstrapH := handlers.NewBootstrapHandler(bootstrapSvc, cfg.AppEnv)
	adminUsersH := handlers.NewAdminUsersHandler(adminUsersSvc)

	// middleware
	authMW := middleware.NewAuthMiddleware(cfg.JWTSecret)
	adminMW := middleware.NewRequireAdminMiddleware()

	// router (tu firma actual)
	router := httpapi.NewRouter(
		authH,
		bootstrapH,
		childrenH,
		attH,
		adminPricingH,
		adminPromosH,
		adminUsersH,
		authMW,
		adminMW,
	)

	log.Printf("listening on %s (env=%s)", cfg.HTTPAddr, cfg.AppEnv)
	log.Fatal(http.ListenAndServe(cfg.HTTPAddr, router))
}

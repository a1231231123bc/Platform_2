package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"platform/internal/config"
	"platform/internal/database"
	"platform/internal/handler"
	"platform/internal/repository"
	"platform/internal/router"
	"platform/internal/service"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	if err := database.RunMigrations(cfg.DatabaseURL); err != nil {
		log.Fatalf("run migrations: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	pool, err := database.NewPool(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("connect to database: %v", err)
	}
	defer pool.Close()

	userRepo := repository.NewUserRepository(pool)
	orgRepo := repository.NewOrganizationRepository(pool)
	contractorRepo := repository.NewContractorRepository(pool)
	blacklistRepo := repository.NewContractorBlacklistRepository(pool)
	jobRepo := repository.NewJobRepository(pool)
	jobOfferRepo := repository.NewJobOfferRepository(pool)
	assignmentRepo := repository.NewAssignmentRepository(pool)

	authService := service.NewAuthService(userRepo, orgRepo, cfg)
	usersService := service.NewUsersService(userRepo)
	orgService := service.NewOrganizationsService(orgRepo)
	contractorsService := service.NewContractorsService(contractorRepo, blacklistRepo)
	jobsService := service.NewJobsService(jobRepo)
	matchingService := service.NewMatchingService(jobRepo, contractorRepo)
	responsesService := service.NewResponsesService(jobOfferRepo)
	assignmentsService := service.NewAssignmentsService(assignmentRepo, jobRepo, contractorRepo)

	authHandler := handler.NewAuthHandler(authService)
	usersHandler := handler.NewUsersHandler(usersService)
	orgHandler := handler.NewOrganizationsHandler(orgService)
	contractorsHandler := handler.NewContractorsHandler(contractorsService)
	jobsHandler := handler.NewJobsHandler(jobsService)
	matchingHandler := handler.NewMatchingHandler(matchingService)
	responsesHandler := handler.NewResponsesHandler(responsesService)
	assignmentsHandler := handler.NewAssignmentsHandler(assignmentsService)

	r := router.New(router.Deps{
		AuthHandler:          authHandler,
		UsersHandler:         usersHandler,
		OrganizationsHandler: orgHandler,
		ContractorsHandler:   contractorsHandler,
		JobsHandler:          jobsHandler,
		MatchingHandler:      matchingHandler,
		ResponsesHandler:     responsesHandler,
		AssignmentsHandler:   assignmentsHandler,
		JWTSecret:            cfg.JWTSecret,
	})

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.Port),
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("server starting on port %s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}

	log.Println("server stopped")
}

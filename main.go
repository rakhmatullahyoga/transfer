package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"transfer/transfer/handler"
	"transfer/transfer/repository"
	"transfer/transfer/usecase"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// Setup dependencies
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%s/%s", os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_PASSWORD"), os.Getenv("POSTGRES_HOST"), os.Getenv("POSTGRES_PORT"), os.Getenv("POSTGRES_DB"))
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		fmt.Fprintf(os.Stderr, "unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer pool.Close()

	// Setup repositories
	bankRepo := repository.NewBankRepository(os.Getenv("BANK_HOST"))
	dbRepo := repository.NewDbRepository(pool)

	// Setup usecase
	uc := usecase.NewUsecase(bankRepo, dbRepo)

	// Setup handler
	h := handler.NewTransferHttpHandler(uc)

	// Setup http server
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.SetHeader("Content-Type", "application/json"))
	r.Mount("/transfer", h.Route())

	srvPort := os.Getenv("SERVER_PORT")
	if srvPort == "" {
		srvPort = "3000"
	}

	s := &http.Server{
		Addr:    fmt.Sprintf(":%s", srvPort),
		Handler: r,
	}
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
	go func(s *http.Server) {
		log.Printf("transfer API now available at %s\n", s.Addr)
		if serr := s.ListenAndServe(); serr != http.ErrServerClosed {
			log.Fatal(serr)
		}
	}(s)

	<-sigChan

	err = s.Shutdown(context.Background())
	if err != nil {
		log.Fatal("something wrong when stopping server : ", err)
		return
	}

	log.Printf("server stopped")
}

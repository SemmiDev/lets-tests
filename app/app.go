package app

import (
	"context"
	"github.com/SemmiDev/lets-tests/domain"
	chimiddleware "github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/go-chi/httprate"
	"github.com/joho/godotenv"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var router = chi.NewRouter()

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("sad .env file found")
	}
}

func StartApp() {
	dbdriver := os.Getenv("DBDRIVER")
	username := os.Getenv("USERNAME")
	password := os.Getenv("PASSWORD")
	host := os.Getenv("HOST")
	database := os.Getenv("DATABASE")
	port := os.Getenv("PORT")
	httpRateLimitRequest := os.Getenv("HTTP_RATE_LIMIT_REQUEST")
	httpRateLimitTime := os.Getenv("HTTP_RATE_LIMIT_TIME")

	limitRequest, _ := strconv.Atoi(httpRateLimitRequest)
	limitTime, _ := time.ParseDuration(httpRateLimitTime)

	domain.ChatRepo.Initialize(dbdriver, username, password, port, host, database)

	router.Use(httprate.LimitByIP(limitRequest, limitTime))
	router.Use(cors.AllowAll().Handler)
	router.Use(chimiddleware.Logger)
	router.Use(chimiddleware.Recoverer)

	routes(router)
	run()
}

func run() {
	ctx, cancel := context.WithCancel(context.Background())
	httpServer := &http.Server{
		Addr:        ":3333",
		Handler:     router,
		BaseContext: func(_ net.Listener) context.Context { return ctx },
	}

	// if your BaseContext is more complex you might want to use this instead of doing it manually
	// httpServer.RegisterOnShutdown(cancel)

	// Run server
	go func() {
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			// it is fine to use Fatal here because it is not main gorutine
			log.Fatalf("HTTP server ListenAndServe: %v", err)
		}
	}()

	signalChan := make(chan os.Signal, 1)

	signal.Notify(
		signalChan,
		syscall.SIGHUP,  // kill -SIGHUP XXXX
		syscall.SIGINT,  // kill -SIGINT XXXX or Ctrl+c
		syscall.SIGQUIT, // kill -SIGQUIT XXXX
	)

	<-signalChan
	log.Print("os.Interrupt - shutting down...\n")

	go func() {
		<-signalChan
		log.Fatal("os.Kill - terminating...\n")
	}()

	gracefullCtx, cancelShutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelShutdown()

	if err := httpServer.Shutdown(gracefullCtx); err != nil {
		log.Printf("shutdown error: %v\n", err)
		defer os.Exit(1)
		return
	} else {
		log.Printf("gracefully stopped\n")
	}

	// manually cancel context if not using httpServer.RegisterOnShutdown(cancel)
	cancel()

	defer os.Exit(0)
	return
}

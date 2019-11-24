package main

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Smet1/bmstu-gis/internal/pathfinding"

	"github.com/Smet1/bmstu-gis/internal/config"
	"github.com/Smet1/bmstu-gis/internal/db"
	"github.com/Smet1/bmstu-gis/internal/handlers"
	"github.com/Smet1/bmstu-gis/internal/logger"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
)

func main() {
	configPath := flag.String(
		"config",
		"./config.yaml",
		"path to config",
	)
	flag.Parse()

	filenameHook := filename.NewHook()
	filenameHook.Field = "sourcelog"

	log := logrus.New()
	log.AddHook(filenameHook)
	log.SetFormatter(&logrus.JSONFormatter{})

	cfg := &config.Config{}
	err := config.ReadConfig(*configPath, &cfg)
	if err != nil {
		log.WithError(err).Fatal("can't read config")
	}
	log.WithField("config", cfg).Info("started with config")

	conn, err := db.EnsureDBConn(cfg)
	if err != nil {
		log.WithError(err).Fatal("can't create db connection")
	}

	mapBMSTU, err := pathfinding.InitMap(conn)
	if err != nil {
		log.WithError(err).Fatal("can't create map")
	}

	mux := chi.NewRouter()
	mux.Use(middleware.NoCache)
	mux.Use(logger.GetLoggerMiddleware(log))

	mux.Route("/api", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.Post("/", handlers.GetCreateUserHandler(conn))
			r.Get("/{userID}", handlers.GetGetUserHandler(conn))
			r.Post("/login", handlers.GetLoginUserHandler(conn))
		})

		r.Route("/news", func(r chi.Router) {
			r.Get("/", handlers.GetGetNewsHandler(conn))
			r.Post("/", handlers.GetCreateNewsHandler(conn))
		})

		r.Route("/history", func(r chi.Router) {
			r.Post("/", handlers.GetCreateHistoryHandler(conn))
			r.Get("/{userID}", handlers.GetGetHistoryHandler(conn))
			r.Delete("/{userID}", handlers.GetClearHistoryHandler(conn))
		})

		r.Route("/map", func(r chi.Router) {
			r.Route("/point", func(r chi.Router) {
				r.Post("/", handlers.GetCreatePointHandler(conn))
				r.Get("/", handlers.GetGetPointsHandler(conn))
			})

			r.Route("/search", func(r chi.Router) {
				r.Get("/", handlers.GetPathFindingHandler(conn, mapBMSTU))
			})
		})
	})

	server := http.Server{
		Handler: mux,
		Addr:    cfg.ServeAddr,
	}

	go func() {
		log.Infof("syncing orders service started on port %s", cfg.ServeAddr)
		if err = server.ListenAndServe(); err != nil {
			if err == http.ErrServerClosed {
				log.Info("graceful shutdown")
			} else {
				log.WithError(err).Fatal("sync service")
			}
		}
	}()

	sgnl := make(chan os.Signal, 1)
	signal.Notify(sgnl,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)
	stop := <-sgnl

	if err = server.Shutdown(context.Background()); err != nil {
		log.WithError(err).Error("error on shutdown")
	}

	log.WithField("signal", stop).Info("stopping")
}

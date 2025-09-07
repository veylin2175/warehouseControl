package main

import (
	"WarehouseControl/internal/config"
	"WarehouseControl/internal/http-server/handlers"
	authMiddleware "WarehouseControl/internal/http-server/handlers/middleware"
	"WarehouseControl/internal/http-server/middleware/mwlogger"
	"WarehouseControl/internal/lib/logger/handlers/slogpretty"
	"WarehouseControl/internal/lib/logger/sl"
	"WarehouseControl/internal/storage/postgres"
	"errors"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	cfg := config.MustLoad()

	log := setupLogger(cfg.Env)

	log.Info("Starting warehouse control system", slog.String("env", cfg.Env))
	log.Debug("Debug messages are enabled")

	storage, err := postgres.InitDB(&cfg.Database)
	if err != nil {
		log.Error("failed to init storage", sl.Err(err))
		os.Exit(1)
	}

	// Инициализация репозиториев
	userStorage := postgres.NewUserStorage(storage.DB)
	itemStorage := postgres.NewItemStorage(storage.DB)
	historyStorage := postgres.NewHistoryStorage(storage.DB)

	// Инициализация хендлеров
	authHandler := handlers.NewAuthHandler(userStorage, "secret-key", log)
	itemsHandler := handlers.NewItemsHandler(itemStorage, log)
	historyHandler := handlers.NewHistoryHandler(historyStorage, log)

	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(mwlogger.New(log))
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)

	// Раздача статических файлов
	fs := http.FileServer(http.Dir("./static/"))
	router.Handle("/static/*", http.StripPrefix("/static/", fs))

	// Главная страница - проверяем авторизацию через JS
	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/index.html")
	})

	// Страница логина
	router.Get("/login.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./static/login.html")
	})

	// Публичные API маршруты
	router.Post("/login", authHandler.Login)

	// Защищенные API маршруты - применяем middleware
	router.Group(func(r chi.Router) {
		r.Use(authMiddleware.AuthMiddleware("secret-key", log))

		r.Post("/items", itemsHandler.CreateItem)
		r.Get("/items", itemsHandler.GetAllItems)
		r.Get("/items/{id}", itemsHandler.GetItemByID)
		r.Put("/items/{id}", itemsHandler.UpdateItem)
		r.Delete("/items/{id}", itemsHandler.DeleteItem)
		r.Get("/history", historyHandler.GetAllHistory)
		r.Get("/history/{id}", historyHandler.GetHistoryByItemID)
	})

	log.Info("starting server", slog.String("address", cfg.HTTPServer.Address))

	srv := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)

	go func() {
		if err = srv.ListenAndServe(); err != nil && !errors.Is(http.ErrServerClosed, err) {
			log.Error("failed to start server", sl.Err(err))
			stop <- syscall.SIGTERM
		}
	}()

	sign := <-stop

	log.Info("application stopping", slog.String("signal", sign.String()))

	if err = srv.Shutdown(nil); err != nil {
		log.Error("failed to shutdown server", sl.Err(err))
	}

	log.Info("application stopped")

	if err = storage.Close(); err != nil {
		log.Error("failed to close postgres connection", sl.Err(err))
	}

	log.Info("postgres connection closed")
}

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger

	switch env {
	case envLocal:
		log = setupPrettySlog()
	case envDev:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		log = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	h := opts.NewPrettyHandler(os.Stdout)

	return slog.New(h)
}

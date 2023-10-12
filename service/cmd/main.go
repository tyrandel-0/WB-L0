package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"service/natsStreaming"
	"service/server/cache"
	"service/server/config"
	"service/server/http-server/handlers"
	"service/server/storage/postgres"
)

func main() {
	cfg := InitConfig("config/local.yml")

	var wg sync.WaitGroup
	defer wg.Wait()

	storage := InitStorage(cfg)

	router := chi.NewRouter()

	InitMiddleware(router)

	initCache := InitCache(storage, &wg)

	InitHandlers(router, storage, initCache)

	InitNatsStreaming(&wg, storage, initCache)

	RunServer(cfg, router)
}

func InitConfig(configPath string) *config.Config {
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal(err)
	}
	return cfg
}

func InitStorage(cfg *config.Config) *postgres.Storage {
	connectionString := fmt.Sprintf(
		"host=%s port=%d user=%s dbname=%s sslmode=%s",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User,
		cfg.Database.DBName, cfg.Database.SSLMode,
	)

	storage, err := postgres.New(connectionString)
	if err != nil {
		log.Fatal(err)
	}

	return storage
}

func InitMiddleware(router chi.Router) {
	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.URLFormat)
}

func InitCache(storage *postgres.Storage, wg *sync.WaitGroup) *cache.Cache {
	newCache, err := cache.New(storage, wg)
	if err != nil {
		log.Println("Can not init newCache")
	}
	return newCache
}

func InitHandlers(router *chi.Mux, storage *postgres.Storage, cache *cache.Cache) {
	// home
	router.Get("/", handlers.HomePage)

	// add
	router.Get("/add", handlers.AddOrderPage(""))
	router.Post("/add", handlers.AddOrder(storage, cache))

	// find
	router.Get("/find", handlers.FindOrderByIDPage(""))
	router.Post("/find", handlers.FindOrderByID(storage, cache))

	// order
	router.Get("/order", handlers.OrderDetailsPage(nil))
}

func InitNatsStreaming(wg *sync.WaitGroup, storage *postgres.Storage, cache *cache.Cache) {
	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		defer fmt.Println("Shutting down...")
		defer wg.Done()

		if err := natsStreaming.RunNatsStreaming(storage, cache); err != nil {
			log.Println(err)
		}
	}(wg)
}

func RunServer(cfg *config.Config, router *chi.Mux) {
	log.Printf("starting server\naddress: %s", cfg.HttpServer.Address)

	srv := &http.Server{
		Addr:         cfg.HttpServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HttpServer.Timeout,
		WriteTimeout: cfg.HttpServer.Timeout,
		IdleTimeout:  cfg.HttpServer.IdleTimeout,
	}
	if err := srv.ListenAndServe(); err != nil {
		log.Println("failed to start server")
	}

	log.Println("server stopped")
}

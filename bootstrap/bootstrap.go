package bootstrap

import (
	"fmt"
	"log"
	"time"

	"kong-go-assignment/config"
	"kong-go-assignment/handlers"
	"kong-go-assignment/models"
	"kong-go-assignment/repositories"
	"kong-go-assignment/services"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type App struct {
	Config *config.Config
	DB     *gorm.DB
	Router *gin.Engine
}

func InitializeApp() (*App, error) {
	cfg := initConfig()
	db, err := initDatabase(cfg)
	if err != nil {
		return nil, err
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	router := initRouter(initServices(db))

	return &App{
		Config: cfg,
		DB:     db,
		Router: router,
	}, nil
}

func initConfig() *config.Config {
	cfg := config.Load()
	return &cfg
}

func initDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxOpenConns(25)
	sqlDB.SetMaxIdleConns(25)
	sqlDB.SetConnMaxLifetime(5 * time.Minute)

	log.Println("Database connection established")
	return db, nil
}

func runMigrations(db *gorm.DB) error {
	return db.AutoMigrate(&models.Service{}, &models.ServiceVersion{})
}

func initServices(db *gorm.DB) *handlers.ServiceHandler {
	serviceRepo := repositories.NewServiceRepository(db)
	serviceSvc := services.NewServiceService(serviceRepo)
	return handlers.NewServiceHandler(serviceSvc)
}

func initRouter(serviceHandler *handlers.ServiceHandler) *gin.Engine {
	return handlers.NewRouter(serviceHandler)
}

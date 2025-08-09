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

	cfg := config.Load()

	db, err := connectDatabase(&cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	if err := db.AutoMigrate(&models.Service{}, &models.ServiceVersion{}); err != nil {
		return nil, fmt.Errorf("failed to migrate DB: %w", err)
	}

	serviceRepo := repositories.NewServiceRepository(db)

	serviceSvc := services.NewServiceService(serviceRepo)

	serviceHandler := handlers.NewServiceHandler(serviceSvc)
	router := handlers.NewRouter(serviceHandler)

	return &App{
		Config: &cfg,
		DB:     db,
		Router: router,
	}, nil
}

func connectDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&loc=Local",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
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

package http

import (
	"context"
	"excel-service/internal/configs"
	"excel-service/internal/repository"
	"excel-service/internal/service"
	"excel-service/internal/transport/http/handler"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	container "github.com/vielendanke/go-db-lb"
)

func StartHTTPServer(ctx context.Context, errCh chan<- error) {
	app := echo.New()

	if envErr := godotenv.Load(); envErr != nil {
		log.Warnf("failed to find env file: %v", envErr)
	}

	cfg := newConfig()

	pool, poolErr := InitDBX(ctx, fmt.Sprintf("user=%s host=%s port=%s password=%s dbname=%s sslmode=%s", cfg.DB.User, cfg.DB.Host, cfg.DB.Port, cfg.DB.Password, cfg.DB.DBName, cfg.DB.SslMode))
	if poolErr != nil {
		log.Errorf("failed to initialize database: %s", poolErr)
		errCh <- poolErr
		return
	}

	lb, lbErr := container.NewLoadBalancer(ctx, 2, 2)
	if lbErr != nil {
		log.Errorf("failed to create container: %v", lbErr)
		errCh <- lbErr
	}

	prErr := lb.AddPGxPoolPrimaryNode(ctx, pool)
	if prErr != nil {
		log.Errorf("failed add primary node: %v", prErr)
		errCh <- prErr
	}

	srErr := lb.AddPGxPoolNode(ctx, pool)
	if srErr != nil {
		log.Errorf("failed add secondary node: %v", srErr)
		errCh <- srErr
	}

	excelRepo := repository.NewExcelRepository(lb)
	excelService := service.NewExcelService(excelRepo, lb, cfg)

	cron := gocron.NewScheduler(time.UTC)

	_, err := cron.Every(5).Minute().Do(func() { fmt.Println("раз кроно два кроно") })
	if err != nil {
		fmt.Println("crono error: ", err)
		errCh <- err
	}

	cron.StartAsync()

	srvHandler := handler.NewHandler(excelService)

	app.POST("api/v1/upload/excel", srvHandler.SaveExcelFile)
	app.POST("api/v1/upload/mtr", srvHandler.SaveMtr)
	app.POST("api/v1/upload/category", srvHandler.NewCategory)
	app.POST("api/v1/upload/company", srvHandler.NewCompany)
	app.POST("api/v1/upload/orgNomenclature", srvHandler.SaveOrganizerNomenclature)
	app.POST("api/v1/upload/bank", srvHandler.SaveBanks)
	app.POST("api/v1/upload/aws/object", srvHandler.GetExcelFromAwsByFileId)
	app.GET("dimeken", dimeken)

	errCh <- app.Start(cfg.Port)
}
func dimeken(c echo.Context) error {
	return c.JSON(http.StatusOK, nil)
}

func dimeken(c echo.Context) error {
	return c.JSON(http.StatusCreated, nil)
}

func InitDBX(ctx context.Context, url string) (*pgxpool.Pool, error) {
	conf, cfgErr := pgxpool.ParseConfig(url)
	if cfgErr != nil {
		return nil, cfgErr
	}

	conf.MaxConns = 20
	conf.MinConns = 10
	conf.MaxConnIdleTime = 10 * time.Second

	pool, poolErr := pgxpool.ConnectConfig(ctx, conf)
	if poolErr != nil {
		return nil, poolErr
	}

	if pingErr := pool.Ping(ctx); pingErr != nil {
		return nil, pingErr
	}
	return pool, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func newConfig() *configs.Configs {
	return &configs.Configs{
		Port: ":" + getEnv("port", "9090"),
		DB: &configs.DBCfg{
			User:     getEnv("DB_USERNAME", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			Host:     getEnv("DB_HOST", "postrges"),
			Port:     getEnv("DB_PORT", "5432"),
			DBName:   getEnv("DB_DATABASE", "postrges"),
			SslMode:  getEnv("DB_SSLMODE", "disable"),
		},
		Aws: &configs.AwsConfig{
			Host:      getEnv("XCLOUD_DIRECTUS_S3_HOST", "postgres"),
			AccessKey: getEnv("XCLOUD_DIRECTUS_S3_KEY", "postgres"),
			SecretKey: getEnv("XCLOUD_DIRECTUS_S3_SECRET", "postgres"),
			Bucket:    getEnv("XCLOUD_DIRECTUS_S3_BUCKET", "postgres"),
		},
	}
}

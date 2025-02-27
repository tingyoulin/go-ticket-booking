package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/tingyoulin/go-ticket-booking/booking"
	"github.com/tingyoulin/go-ticket-booking/flight"
	mysqlRepo "github.com/tingyoulin/go-ticket-booking/internal/repository/mysql"
	redisRepo "github.com/tingyoulin/go-ticket-booking/internal/repository/redis"
	"github.com/tingyoulin/go-ticket-booking/internal/rest/api"
	"github.com/tingyoulin/go-ticket-booking/internal/rest/middleware"
	"github.com/tingyoulin/go-ticket-booking/passenger"
)

const (
	defaultTimeout = 30
	defaultAddress = ":9090"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func main() {
	//prepare database
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("DATABASE_USER")
	dbPass := os.Getenv("DATABASE_PASS")
	dbName := os.Getenv("DATABASE_NAME")
	connection := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUser, dbPass, dbHost, dbPort, dbName)
	val := url.Values{}
	val.Add("parseTime", "1")
	val.Add("loc", "Asia/Jakarta")
	dsn := fmt.Sprintf("%s?%s", connection, val.Encode())

	gormDB, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatal("failed to open connection to database", err)
	}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// prepare echo
	e := echo.New()
	e.Use(middleware.CORS)
	timeoutStr := os.Getenv("CONTEXT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutStr)
	if err != nil {
		log.Println("failed to parse timeout, using default timeout")
		timeout = defaultTimeout
	}
	timeoutContext := time.Duration(timeout) * time.Second
	e.Use(middleware.SetRequestContextWithTimeout(timeoutContext))

	// Prepare Repository
	bookingRepo := mysqlRepo.NewBookingRepository(gormDB)
	flightRepo := mysqlRepo.NewFlightRepository(gormDB)
	passengerRepo := mysqlRepo.NewPassengerRepository(gormDB)
	tokenRedisRepo := redisRepo.NewTokenRedisRepository(redisClient)

	// Auth Middleware
	authMiddleware := middleware.Auth(tokenRedisRepo)

	// Build service Layer
	bookingSvc := booking.NewService(bookingRepo, flightRepo)
	flightSvc := flight.NewService(flightRepo)
	passengerSvc := passenger.NewService(passengerRepo, tokenRedisRepo)

	api.NewBookingHandler(e, bookingSvc, authMiddleware)
	api.NewFlightHandler(e, flightSvc, authMiddleware)
	api.NewPassengerHandler(e, passengerSvc, authMiddleware)

	// Start Server
	address := os.Getenv("SERVER_ADDRESS")
	if address == "" {
		address = defaultAddress
	}
	log.Fatal(e.Start(address)) //nolint
}

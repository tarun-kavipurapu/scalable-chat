package internal

import (
	"context"
	"log"
	"tarun-kavipurapu/test-go-chat/config"
	"tarun-kavipurapu/test-go-chat/db/adapters"
	db "tarun-kavipurapu/test-go-chat/db/sqlc"
	"tarun-kavipurapu/test-go-chat/internal/middlewares"

	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

type Server struct {
	store db.Store
	// store  *db.SQLStore
	router      *gin.Engine
	redisClient *redis.Client
}

// runs on a specific address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func NewHTTPServer() *Server {
	ctx := context.Background()
	router := gin.Default()

	router.Use(middlewares.CORSMiddleware())

	dbSource := config.EnvVars.DBSource
	pgConn := adapters.InitDb(dbSource)

	// redisURL := config.EnvVars.REDIS_URL

	redisClient := adapters.CreateRedisClient("redis://localhost:6379")
	pong, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	} else {
		log.Println("Successfully connected to Redis:", pong)
	}
	err = redisClient.Set(ctx, "mykey", "Hello from Redis!", 0).Err()
	if err != nil {
		log.Fatal("Failed to set key in Redis:", err)
	}

	// if pgConn != nil {
	// 	log.Fatal("Connected to database")
	// }

	store := db.NewStore(pgConn)

	// ctx := context.Background()

	server := &Server{

		router:      router,
		store:       store,
		redisClient: redisClient,
	}

	log.Println("Setting up server...")

	SetupRouter(server)

	return server
}

package internal

import (
	"log"
	"os"
	"path/filepath"
	"tarun-kavipurapu/test-go-chat/config"
	"tarun-kavipurapu/test-go-chat/db/adapters"
	db "tarun-kavipurapu/test-go-chat/db/sqlc"
	"tarun-kavipurapu/test-go-chat/internal/middlewares"

	"github.com/gin-gonic/gin"
)

type Server struct {
	store db.Store
	// store  *db.SQLStore
	router *gin.Engine
}

// runs on a specific address
func (s *Server) Start(address string) error {
	return s.router.Run(address)
}

func NewHTTPServer() *Server {
	router := gin.Default()

	router.Use(middlewares.CORSMiddleware())
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current working directory:", err)
	}

	// Navigate up to the root directory of your project
	rootDir := filepath.Dir(filepath.Dir(cwd))

	// Load configuration from the root directory
	cfg, err := config.LoadConfig(rootDir)
	if err != nil {
		log.Fatal("Failed to load configuration:", err)
	}

	dbSource := cfg.DBSource
	pgConn := adapters.InitDb(dbSource)

	// if pgConn != nil {
	// 	log.Fatal("Connected to database")
	// }

	store := db.NewStore(pgConn)

	// ctx := context.Background()

	server := &Server{

		router: router,
		store:  store,
	}

	log.Println("Setting up server...")

	SetupRouter(server)

	return server
}

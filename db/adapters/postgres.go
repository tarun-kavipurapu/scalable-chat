package adapters

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func InitDb(connString string) *pgxpool.Pool {
	// config, err := config.LoadConfig(".") // Load configuration
	// if err != nil {
	// 	log.Fatalf("cannot load configuration: %v", err)
	// }

	pool, err := pgxpool.New(context.Background(), connString) // Connect to database
	if err != nil {
		log.Fatal("Unable to connect to database", err)
	}

	return pool
}

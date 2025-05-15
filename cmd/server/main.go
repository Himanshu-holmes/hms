package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"golang.org/x/net/context"

	"github.com/himanshu-holmes/hms/internal/db"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	DB db.Queries
}

func main() {
	godotenv.Load("../../.env")
	portEnv := os.Getenv("PORT")
	if portEnv == "" {
		portEnv = "3000"
	}
	dbUrl := os.Getenv("DB_URL")
	if dbUrl == "" {
		log.Fatal("DB_URL is not found in the env")
	}
	// Create a connection pool
	dbpool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v\n", err)
	}
	defer dbpool.Close()

	// Verify connection
	if err := dbpool.Ping(context.Background()); err != nil {
		log.Fatalf("Unable to ping database: %v\n", err)
	}

	fmt.Println("Successfully connected to PostgreSQL!")

	apiCfg := apiConfig{
		DB: *db.New(dbpool),
	}
	
}

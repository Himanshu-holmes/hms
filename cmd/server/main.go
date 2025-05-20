package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"golang.org/x/net/context"
	"github.com/himanshu-holmes/hms/internal/db"
	"github.com/himanshu-holmes/hms/internal/handler"
	"github.com/himanshu-holmes/hms/internal/repository"
	"github.com/himanshu-holmes/hms/internal/service"
	"github.com/joho/godotenv"
)


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


	// Initialize the repositories
	userRepo := repository.NewUserRepo(db.New(dbpool))
	patientRepo := repository.NewPatientRepo(db.New(dbpool))
    patientVisitRepo := repository.NewPatientVisitRepo(db.New(dbpool))


	// Initialize the services
	userService := service.NewAuthService(userRepo)
	patientService := service.NewPatientService(patientRepo)
	patientVisitService := service.NewPatientVisitService(patientVisitRepo, patientRepo)

	// Initialize the handlers
	userHandler := handler.NewAuthHandler(userService)
	patientHandler := handler.NewPatientHandler(patientService)
	patientVisitHandler := handler.NewPatientVisitHandler(patientVisitService)

	// Initialize the router	
	r := gin.Default()

	api := r.Group("/api/v1")
	{
		api.POST("/login", userHandler.Login)
		api.POST("/register", userHandler.CreateUser)
		api.POST("/create-patient", patientHandler.RegisterPatient)
		api.GET("/patients/:id", patientHandler.GetPatient)
		api.GET("/patients", patientHandler.ListPatients)
		api.GET("/patients/:id/visits", patientVisitHandler.ListPatientVisits)
		api.POST("/patients/:id/visits", patientVisitHandler.RecordPatientVisit)
	}



	r.Run(":" + portEnv)
       
	


}

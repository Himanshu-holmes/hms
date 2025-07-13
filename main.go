// @title HMS API
// @version 1.0
// @description Hospital Management System API.
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @host localhost:3000
// @BasePath /api/v1
package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/himanshu-holmes/hms/docs"
	"github.com/himanshu-holmes/hms/internal/db"
	"github.com/himanshu-holmes/hms/internal/handler"
	"github.com/himanshu-holmes/hms/internal/middleware"
	"github.com/himanshu-holmes/hms/internal/repository"
	"github.com/himanshu-holmes/hms/internal/service"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"golang.org/x/net/context"
	"github.com/bugsnag/bugsnag-go/v2"
	"github.com/bugsnag/bugsnag-go-gin"
)

func main() {
	godotenv.Load(".env")
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
	 r.Use(bugsnaggin.AutoNotify(bugsnag.Configuration{
        // Your Bugsnag project API key, required unless set as environment
        // variable $BUGSNAG_API_KEY
        APIKey:        os.Getenv("BUGSNAG_API_KEY"),
        // The import paths for the Go packages containing your source files
        ProjectPackages: []string{"main", "github.com/org/myapp"},
    }))
	r.GET("/healthz", func(c *gin.Context) { c.Status(200) })

	api := r.Group("/api/v1")

	// Swagger documentation
	docs.SwaggerInfo.BasePath = "/api/v1"
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
	{
		//auth
		api.POST("/auth/login", userHandler.Login)
		api.POST("/auth/register", userHandler.CreateUser)
		// patient
		api.POST("/patients/create", middleware.AuthMiddleware(), patientHandler.RegisterPatient)
		api.GET("/patients/:id", middleware.AuthMiddleware(), patientHandler.GetPatient)
		api.GET("/patients", middleware.AuthMiddleware(), patientHandler.ListPatients)
		api.PATCH("/patients/:id", middleware.AuthMiddleware(), patientHandler.UpdatePatient)
		api.DELETE("/patients/:id", middleware.AuthMiddleware(), patientHandler.DeletePatient)
		// visit
		api.POST("/visits/create", middleware.AuthMiddleware(), patientVisitHandler.RecordPatientVisit)
		api.GET("/visits/:id", middleware.AuthMiddleware(), patientVisitHandler.GetPatientVisitDetails)
		api.GET("/visits/:id/list", middleware.AuthMiddleware(), patientVisitHandler.ListPatientVisits)
		api.PATCH("/visits/:id", middleware.AuthMiddleware(), patientVisitHandler.UpdatePatientVisit)
		
	}
	r.Run(":" + portEnv)
}

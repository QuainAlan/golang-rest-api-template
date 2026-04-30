package main

import (
	"context"
	"golang-rest-api-template/pkg/api"
	"golang-rest-api-template/pkg/auth"
	"golang-rest-api-template/pkg/cache"
	"golang-rest-api-template/pkg/database"
	"golang-rest-api-template/pkg/middleware"
	"log"
	"os"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8001
// @BasePath  /api/v1

// @securityDefinitions.apikey JwtAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name X-API-Key

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	if err := auth.SetJWTSigningKey([]byte(os.Getenv("JWT_SECRET_KEY"))); err != nil {
		log.Fatalf("invalid JWT_SECRET_KEY: %v", err)
	}
	if err := middleware.SetAPISecretKey([]byte(os.Getenv("API_SECRET_KEY"))); err != nil {
		log.Fatalf("invalid API_SECRET_KEY: %v", err)
	}

	redisClient := cache.NewRedisClient()
	db := database.NewDatabase()
	dbWrapper := &database.GormDatabase{DB: db}
	mongo := database.SetupMongoDB()
	ctx := context.Background()
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	//gin.SetMode(gin.ReleaseMode)
	gin.SetMode(gin.DebugMode)

	r := api.NewRouter(logger, mongo, dbWrapper, redisClient, &ctx)

	if err := r.Run(":8001"); err != nil {
		log.Fatal(err)
	}
}

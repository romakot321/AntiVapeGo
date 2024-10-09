package main

import (
  "log"
  "os"

  "antivape/db"
  "antivape/services"
  "antivape/repositories"
  "antivape/handlers"
  _ "antivape/docs"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
  "github.com/gofiber/swagger"
)

func InitApp() *fiber.App {
  dsn := os.Getenv("DB_CONFIG")
  if len(dsn) == 0 {
    dsn = "host=localhost port=5432 user=postgres dbname=db password=postgres sslmode=disable"
  }

  dbConnection, err := db.InitDatabase(dsn)
  if err != nil {
    log.Fatal(err)
  }
  redisConnection := db.InitRedis()

  sensorService := services.NewSensorService(dbConnection)
  roomService := services.NewRoomService(dbConnection)
  userRepository := repositories.NewUserRepository(dbConnection)
  authService := services.NewAuthService(userRepository)
  zoneService := services.NewZoneService(dbConnection)
  externalService := services.NewExternalService(redisConnection, dbConnection)
  userService := services.NewUserService(dbConnection)

  authHandler := handlers.NewAuthHandler(authService)
  zoneHandler := handlers.NewZoneHandler(zoneService, authService)
  sensorHandler := handlers.NewSensorHandler(sensorService, authService)
  roomHandler := handlers.NewRoomHandler(roomService, authService)
  externalHandler := handlers.NewExternalHandler(externalService)
  userHandler := handlers.NewUserHandler(userService, authService)

  app := fiber.New()
  app.Get("/swagger/*", swagger.HandlerDefault) // default
  app.Use(cors.New())
  authHandler.Register(app)
  zoneHandler.Register(app)
  sensorHandler.Register(app)
  roomHandler.Register(app)
  userHandler.Register(app)
  externalHandler.Register(app)
  go externalService.RunTransferingCycle()

  return app
}

// @title AntiVape API in golang
// @version 1.0
// @description AntiVape API
// @termsOfService http://swagger.io/terms/
// @host localhost:8000
// @BasePath /
//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
func main() {
  app := InitApp()
  app.Listen(":8080")
}

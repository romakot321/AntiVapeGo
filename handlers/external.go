package handlers

import (
  "antivape/schemas"
  "antivape/services"
	"github.com/gofiber/fiber/v2"
)

type ExternalHandler interface {
  Register(app *fiber.App)
}

type externalHandler struct {
  externalService services.ExternalService  
}

// Store sensordata godoc
//
//	@Summary		store sensordata
//	@Description	store sensordata
//	@Tags			External
//	@Accept			json
//	@Produce		json
//	@Param			account	body		schemas.ExternalSensorDataSchema true	"Create room"
//	@Success		200		{object}	nil
//	@Router			/external/sensors_data [post]
func (h externalHandler) handleStore(c *fiber.Ctx) error {
  var schema schemas.ExternalSensorDataSchema
  if err := c.BodyParser(&schema); err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  h.externalService.Store(schema)
  return nil
}

func (h externalHandler) Register(app *fiber.App) {
  router := app.Group("/external")

  router.Post("/sensors_data", h.handleStore)
}

func NewExternalHandler(externalService services.ExternalService) ExternalHandler {
  return externalHandler{externalService: externalService}
}

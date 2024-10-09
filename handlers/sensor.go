package handlers

import (
  "strconv"

  "antivape/services"
  "antivape/schemas"
  "antivape/middlewares"
	"github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/logger"
)

type SensorHandler interface {
  Register(app *fiber.App)
}

type sensorHandler struct {
  sensorService services.SensorService
  authService services.AuthService
}

// Create sensor godoc
//
//	@Summary		Create sensor
//	@Description	create sensor
//	@Tags			Sensor
//	@Accept			json
//	@Produce		json
//	@Param			account	body		schemas.SensorCreateSchema true	"Create sensor"
//	@Success		200		{object}	schemas.SensorSchema
//	@Router			/sensor [post]
//	@Security ApiKeyAuth
func (h sensorHandler) handleCreate(c *fiber.Ctx) error {
  if !h.authService.IsSuperuser(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }

  var schema schemas.SensorCreateSchema
  if err := c.BodyParser(&schema); err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  resp := h.sensorService.Create(schema)
  return c.JSON(resp)
}

// Get sensor godoc
//
//	@Summary		Get sensor
//	@Description	Get sensor
//	@Tags			Sensor
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Sensor ID"
//	@Success		200		{object}	schemas.SensorSchema
//	@Router			/sensor/{id} [get]
//	@Security ApiKeyAuth
func (h sensorHandler) handleTake(c *fiber.Ctx) error {
  sensorID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  sensor := h.sensorService.Take(uint(sensorID))
  if sensor.OwnerID != h.authService.CurrentUserID(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }
  return c.JSON(sensor)
}

// Find sensors godoc
//
//	@Summary		Find sensors
//	@Description	Find sensors
//	@Tags			Sensor
//	@Accept			json
//	@Produce		json
//	@Param			q	query		schemas.SensorFindSchema false	"find filters"
//	@Success		200		{array}	schemas.SensorSchema
//	@Router			/sensor/ [get]
// @Security ApiKeyAuth
func (h sensorHandler) handleFind(c *fiber.Ctx) error {
  var schema schemas.SensorFindSchema
  if err := c.BodyParser(&schema); err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  sensors := h.sensorService.Find(schema)
  filteredSensors := h.sensorService.FilterByOwnerID(h.authService.CurrentUserID(c), sensors...)
  return c.JSON(filteredSensors)
}

// UpdateSensor godoc
//
//	@Summary		Update an sensor
//	@Description	Update by json sensor
//	@Tags			Sensor
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Sensor ID"
//	@Param			sensor	body		schemas.SensorUpdateSchema true	"Update sensor"
//	@Success		204		{object}	nil
//	@Router			/sensor/{id} [patch]
//	@Security ApiKeyAuth
func (h sensorHandler) handleUpdate(c *fiber.Ctx) error {
  var schema schemas.SensorUpdateSchema
  sensorID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }
  if err := c.BodyParser(&schema); err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  sensor := h.sensorService.Take(uint(sensorID))
  if sensor.OwnerID != h.authService.CurrentUserID(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }
  h.sensorService.Update(uint(sensorID), schema)
  c.Status(204)
  return nil
}

// DeleteSensor godoc
//
//	@Summary		Delete an sensor
//	@Description	Delete by id sensor
//	@Tags			Sensor
//	@Param			id		path		int					true	"Sensor ID"
//	@Success		204		{object}	nil
//	@Router			/sensor/{id} [delete]
//	@Security ApiKeyAuth
func (h sensorHandler) handleDelete(c *fiber.Ctx) error {
  sensorID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  sensor := h.sensorService.Take(uint(sensorID))
  if sensor.OwnerID != h.authService.CurrentUserID(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }
  h.sensorService.Delete(uint(sensorID))
  c.Status(204)
  return nil
}

func (h sensorHandler) Register(app *fiber.App) {
  router := app.Group("/sensor", middlewares.Protected(), logger.New())

  router.Post("/", h.handleCreate)
  router.Get("/:id<int>/", h.handleTake)
  router.Get("/", h.handleFind)
  router.Patch("/:id", h.handleUpdate)
  router.Delete("/:id", h.handleDelete)
}

func NewSensorHandler(sensorService services.SensorService, authService services.AuthService) SensorHandler {
  return sensorHandler{sensorService: sensorService, authService: authService}
}

package handlers

import (
  "strconv"

  "antivape/services"
  "antivape/schemas"
  "antivape/middlewares"
	"github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/logger"
)

type ZoneHandler interface {
  Register(app *fiber.App)
}

type zoneHandler struct {
  zoneService services.ZoneService
  authService services.AuthService
}

// Create zone godoc
//
//	@Summary		Create zone
//	@Description	create zone
//	@Tags			Zone
//	@Accept			json
//	@Produce		json
//	@Param			account	body		schemas.ZoneCreateSchema true	"Create zone"
//	@Success		201		{object}	schemas.ZoneSchema
//	@Router			/zone [post]
//	@Security ApiKeyAuth
func (h zoneHandler) handleCreate(c *fiber.Ctx) error {
  if !h.authService.IsSuperuser(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }

  var schema schemas.ZoneCreateSchema
  if err := c.BodyParser(&schema); err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  resp := h.zoneService.Create(schema)
  return c.Status(201).JSON(resp)
}

// Get zone godoc
//
//	@Summary		Get zone
//	@Description	Get zone
//	@Tags			Zone
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Zone ID"
//	@Success		200		{object}	schemas.ZoneSchema
//	@Router			/zone/{id} [get]
//	@Security ApiKeyAuth
func (h zoneHandler) handleTake(c *fiber.Ctx) error {
  zoneID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  zone := h.zoneService.Take(uint(zoneID))
  if zone.OwnerID != h.authService.CurrentUserID(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }
  return c.JSON(zone)
}

// Find zones godoc
//
//	@Summary		Find zones
//	@Description	Find zones
//	@Tags			Zone
//	@Accept			json
//	@Produce		json
//	@Param			q	query		schemas.ZoneFindSchema false	"find filters"
//	@Success		200		{array}	schemas.ZoneSchema
//	@Router			/zone/ [get]
// @Security ApiKeyAuth
func (h zoneHandler) handleFind(c *fiber.Ctx) error {
  var schema schemas.ZoneFindSchema
  c.BodyParser(&schema)

  zones := h.zoneService.Find(schema)
  filteredZones := h.zoneService.FilterByOwnerID(h.authService.CurrentUserID(c), zones...)
  return c.JSON(filteredZones)
}

// UpdateZone godoc
//
//	@Summary		Update an zone
//	@Description	Update by json zone
//	@Tags			Zone
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Zone ID"
//	@Param			zone	body		schemas.ZoneUpdateSchema true	"Update zone"
//	@Success		204		{object}	nil
//	@Router			/zone/{id} [patch]
//	@Security ApiKeyAuth
func (h zoneHandler) handleUpdate(c *fiber.Ctx) error {
  var schema schemas.ZoneUpdateSchema
  zoneID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }
  if err := c.BodyParser(&schema); err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  zone := h.zoneService.Take(uint(zoneID))
  if zone.OwnerID != h.authService.CurrentUserID(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }
  h.zoneService.Update(uint(zoneID), schema)
  c.Status(204)
  return nil
}

// DeleteZone godoc
//
//	@Summary		Delete an zone
//	@Description	Delete by id zone
//	@Tags			Zone
//	@Param			id		path		int					true	"Zone ID"
//	@Success		204		{object}	nil
//	@Router			/zone/{id} [delete]
//	@Security ApiKeyAuth
func (h zoneHandler) handleDelete(c *fiber.Ctx) error {
  zoneID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  zone := h.zoneService.Take(uint(zoneID))
  if zone.OwnerID != h.authService.CurrentUserID(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }
  h.zoneService.Delete(uint(zoneID))
  c.Status(204)
  return nil
}

// Get zone statistic godoc
//
//	@Summary		Get zone statistic
//	@Description	Get zone statistic
//	@Tags			Zone
//	@Produce		json
//	@Param			id	path		int	true	"Zone ID"
//	@Success		200		{object}	schemas.SensorDataZoneSchema
//	@Router			/zone/{id}/statistic [get]
//	@Security ApiKeyAuth
func (h zoneHandler) handleStatistic(c *fiber.Ctx) error {
  zoneID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  zone := h.zoneService.Take(uint(zoneID))
  if zone.OwnerID != h.authService.CurrentUserID(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }
  statistic := h.zoneService.GetStatistic(uint(zoneID))
  return c.JSON(statistic)
}

func (h zoneHandler) Register(app *fiber.App) {
  router := app.Group("/zone", middlewares.Protected(), logger.New())

  router.Post("/", h.handleCreate)
  router.Get("/:id<int>/", h.handleTake)
  router.Get("/:id<int>/statistic", h.handleStatistic)
  router.Get("/", h.handleFind)
  router.Patch("/:id<int>", h.handleUpdate)
  router.Delete("/:id<int>", h.handleDelete)
}

func NewZoneHandler(zoneService services.ZoneService, authService services.AuthService) ZoneHandler {
  return zoneHandler{zoneService: zoneService, authService: authService}
}

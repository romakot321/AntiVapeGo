package handlers

import (
  "strconv"

  "antivape/services"
  "antivape/schemas"
  "antivape/middlewares"
	"github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/logger"
)

type RoomHandler interface {
  Register(app *fiber.App)
}

type roomHandler struct {
  roomService services.RoomService
  authService services.AuthService
}

// Create room godoc
//
//	@Summary		Create room
//	@Description	create room
//	@Tags			Room
//	@Accept			json
//	@Produce		json
//	@Param			account	body		schemas.RoomCreateSchema true	"Create room"
//	@Success		201		{object}	schemas.RoomSchema
//	@Router			/room [post]
//	@Security ApiKeyAuth
func (h roomHandler) handleCreate(c *fiber.Ctx) error {
  if !h.authService.IsSuperuser(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }

  var schema schemas.RoomCreateSchema
  if err := c.BodyParser(&schema); err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  resp := h.roomService.Create(schema)
  return c.Status(201).JSON(resp)
}

// Get room godoc
//
//	@Summary		Get room
//	@Description	Get room
//	@Tags			Room
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"Room ID"
//	@Success		200		{object}	schemas.RoomSchema
//	@Router			/room/{id} [get]
//	@Security ApiKeyAuth
func (h roomHandler) handleTake(c *fiber.Ctx) error {
  roomID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  room := h.roomService.Take(uint(roomID))
  if room.OwnerID != h.authService.CurrentUserID(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }
  return c.JSON(room)
}

// Find rooms godoc
//
//	@Summary		Find rooms
//	@Description	Find rooms
//	@Tags			Room
//	@Accept			json
//	@Produce		json
//	@Param			q	query		schemas.RoomFindSchema false	"find filters"
//	@Success		200		{array}	schemas.RoomSchema
//	@Router			/room/ [get]
// @Security ApiKeyAuth
func (h roomHandler) handleFind(c *fiber.Ctx) error {
  var schema schemas.RoomFindSchema
  if err := c.BodyParser(&schema); err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  rooms := h.roomService.Find(schema)
  filteredRooms := h.roomService.FilterByOwnerID(h.authService.CurrentUserID(c), rooms...)
  return c.JSON(filteredRooms)
}

// UpdateRoom godoc
//
//	@Summary		Update an room
//	@Description	Update by json room
//	@Tags			Room
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"Room ID"
//	@Param			room	body		schemas.RoomUpdateSchema true	"Update room"
//	@Success		204		{object}	nil
//	@Router			/room/{id} [patch]
//	@Security ApiKeyAuth
func (h roomHandler) handleUpdate(c *fiber.Ctx) error {
  var schema schemas.RoomUpdateSchema
  roomID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }
  if err := c.BodyParser(&schema); err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  room := h.roomService.Take(uint(roomID))
  if room.OwnerID != h.authService.CurrentUserID(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }
  h.roomService.Update(uint(roomID), schema)
  c.Status(204)
  return nil
}

// DeleteRoom godoc
//
//	@Summary		Delete an room
//	@Description	Delete by id room
//	@Tags			Room
//	@Param			id		path		int					true	"Room ID"
//	@Success		204		{object}	nil
//	@Router			/room/{id} [delete]
//	@Security ApiKeyAuth
func (h roomHandler) handleDelete(c *fiber.Ctx) error {
  roomID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  room := h.roomService.Take(uint(roomID))
  if room.OwnerID != h.authService.CurrentUserID(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }
  h.roomService.Delete(uint(roomID))
  c.Status(204)
  return nil
}

func (h roomHandler) Register(app *fiber.App) {
  router := app.Group("/room", middlewares.Protected(), logger.New())

  router.Post("/", h.handleCreate)
  router.Get("/:id<int>/", h.handleTake)
  router.Get("/", h.handleFind)
  router.Patch("/:id", h.handleUpdate)
  router.Delete("/:id", h.handleDelete)
}

func NewRoomHandler(roomService services.RoomService, authService services.AuthService) RoomHandler {
  return roomHandler{roomService: roomService, authService: authService}
}

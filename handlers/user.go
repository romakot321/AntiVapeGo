package handlers

import (
  "strconv"

  "antivape/services"
  "antivape/schemas"
  "antivape/middlewares"
	"github.com/gofiber/fiber/v2"
  "github.com/gofiber/fiber/v2/middleware/logger"
)

type UserHandler interface {
  Register(app *fiber.App)
}

type userHandler struct {
  userService services.UserService
  authService services.AuthService
}

// Get user godoc
//
//	@Summary		Get user
//	@Description	Get user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id	path		int	true	"user ID"
//	@Success		200		{object}	schemas.UserSchema
//	@Router			/user/{id} [get]
//	@Security ApiKeyAuth
func (h userHandler) handleTake(c *fiber.Ctx) error {
  userID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  user := h.userService.TakeByID(uint(userID))
  if user.ID != h.authService.CurrentUserID(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }
  return c.JSON(user)
}

// Find users godoc
//
//	@Summary		Find users
//	@Description	Find users
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Success		200		{array}	schemas.UserSchema
//	@Router			/user/ [get]
// @Security ApiKeyAuth
func (h userHandler) handleFind(c *fiber.Ctx) error {
  actor := h.authService.ParseToken(c)
  if !actor.IsSuperuser {
    return c.Status(401).SendString("Not enough rights for this request")
  }

  users := h.userService.Find(schemas.UserFindSchema{})
  return c.JSON(users)
}

// Updateuser godoc
//
//	@Summary		Update an user
//	@Description	Update by json user
//	@Tags			user
//	@Accept			json
//	@Produce		json
//	@Param			id		path		int					true	"user ID"
//	@Param			user	body		schemas.UserUpdateSchema true	"Update user"
//	@Success		204		{object}	nil
//	@Router			/user/{id} [patch]
//	@Security ApiKeyAuth
func (h userHandler) handleUpdate(c *fiber.Ctx) error {
  var schema schemas.UserUpdateSchema
  userID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }
  if err := c.BodyParser(&schema); err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  user := h.userService.TakeByID(uint(userID))
  if user.ID != h.authService.CurrentUserID(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }
  h.userService.Update(uint(userID), schema)
  c.Status(204)
  return nil
}

// Deleteuser godoc
//
//	@Summary		Delete an user
//	@Description	Delete by id user
//	@Tags			user
//	@Param			id		path		int					true	"user ID"
//	@Success		204		{object}	nil
//	@Router			/user/{id} [delete]
//	@Security ApiKeyAuth
func (h userHandler) handleDelete(c *fiber.Ctx) error {
  userID, err := strconv.Atoi(c.Params("id"))
  if err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }

  user := h.userService.TakeByID(uint(userID))
  if user.ID != h.authService.CurrentUserID(c) {
    return c.Status(401).SendString("Not enough rights for this request")
  }
  h.userService.Delete(uint(userID))
  c.Status(204)
  return nil
}

func (h userHandler) Register(app *fiber.App) {
  router := app.Group("/user", middlewares.Protected(), logger.New())

  router.Get("/:id<int>/", h.handleTake)
  router.Get("/", h.handleFind)
  router.Patch("/:id", h.handleUpdate)
  router.Delete("/:id", h.handleDelete)
}

func NewUserHandler(userService services.UserService, authService services.AuthService) UserHandler {
  return userHandler{userService: userService, authService: authService}
}

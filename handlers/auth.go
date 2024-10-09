package handlers

import (
  "antivape/services"
  "antivape/schemas"
  "antivape/middlewares"
	"github.com/gofiber/fiber/v2"
)

type AuthHandler interface {
  Register(app *fiber.App)
  HandleLogin(c *fiber.Ctx) error
  HandleRegister(c *fiber.Ctx) error
}

type authHandler struct {
  authService services.AuthService
}

//  Get me godoc
//
//	@Summary		Get me
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Success		200		{object}	schemas.UserSchema
//	@Router			/auth/me [get]
//	@Security ApiKeyAuth
func (h authHandler) HandleGetMe(c *fiber.Ctx) error {
  userID := h.authService.CurrentUserID(c)
  if userID == 0 {
    return c.Status(401).SendString("Invalid token")
  }
  resp := h.authService.GetMe(userID)
  return c.JSON(resp)
}

//  Login godoc
//
//	@Summary		Login
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			account	body		schemas.LoginSchema true	"login"
//	@Success		200		{object}	schemas.TokenSchema
//	@Router			/auth/login [post]
func (h authHandler) HandleLogin(c *fiber.Ctx) error {
  var schema schemas.LoginSchema
  if err := c.BodyParser(&schema); err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }
  resp, err2 := h.authService.Login(schema)
  if err2 != nil {
    return c.Status(401).JSON(fiber.Map{"status": "error", "data": err2})
  }

  return c.JSON(resp)
}

//  Register godoc
//
//	@Summary		Register
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			account	body		schemas.RegisterSchema true	"register"
//	@Success		200		{object}	schemas.UserSchema
//	@Router			/auth/register [post]
func (h authHandler) HandleRegister(c *fiber.Ctx) error {
  var schema schemas.RegisterSchema
  if err := c.BodyParser(&schema); err != nil {
    return c.Status(422).JSON(fiber.Map{"status": "error", "data": err})
  }
  resp := h.authService.Register(schema)

  return c.JSON(resp)
}

func (h authHandler) Register(app *fiber.App) {
  router := app.Group("/auth")

  router.Post("/login", h.HandleLogin)
  router.Post("/register", h.HandleRegister)
  router.Get("/me", middlewares.Protected(), h.HandleGetMe)
}

func NewAuthHandler(authService services.AuthService) AuthHandler {
  return authHandler{authService: authService}
}


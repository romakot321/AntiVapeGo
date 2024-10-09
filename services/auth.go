package services

import (
	"time"
  "errors"

  models "antivape/db"
  repositories "antivape/repositories"
  schemas "antivape/schemas"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

var secretKey = []byte("secret")

type AuthService struct {
  userRepository repositories.UserRepository
}

func (s AuthService) GetMe(userID uint) schemas.UserSchema {
  user := s.userRepository.TakeByID(userID)
  return schemas.UserSchema{
    ID: user.ID,
    Name: user.Name,
    IsSuperuser: user.IsSuperuser,
  }
}

func (s AuthService) Login(schema schemas.LoginSchema) (schemas.TokenSchema, error) {
  user, err := s.validateLogin(schema.Username, schema.Password)
  if err != nil {
    return schemas.TokenSchema{}, err
  }

  token, err := createToken(user.ID, user.IsSuperuser)
  if err != nil {
    return schemas.TokenSchema{}, err
  }

  return schemas.TokenSchema{Token: token}, nil
}

func (s AuthService) Register(schema schemas.RegisterSchema) schemas.UserSchema {
  passwordHash := HashPassword(schema.Password)
  user := s.userRepository.Create(schema.Name, passwordHash)
  return schemas.UserSchema{
    ID: user.ID,
    Name: user.Name,
    IsSuperuser: user.IsSuperuser,
  }
}

func (s AuthService) ParseToken(ctx *fiber.Ctx) models.User {
  userID := parseToken(ctx)
  return s.userRepository.TakeByID(userID)
}

func (s AuthService) validateLogin(username string, password string) (models.User, error) {
  hashedPassword := HashPassword(password)
  user := s.userRepository.TakeByName(username)
  if user.PasswordHash != hashedPassword {
    return models.User{}, errors.New("Invalid username or password")
  }
  return user, nil
}

func (s AuthService) IsSuperuser(ctx *fiber.Ctx) bool {
	token := ctx.Locals("user").(*jwt.Token)
	claims := token.Claims.(jwt.MapClaims)
	status := claims["is_superuser"].(bool)
  return status
}

func (s AuthService) CurrentUserID(ctx *fiber.Ctx) uint {
  return parseToken(ctx)
}

func NewAuthService(userRepository repositories.UserRepository) AuthService {
  return AuthService{userRepository: userRepository}
}

func HashPassword(password string) string {
  return "hashed-" + password
}

func parseToken(ctx *fiber.Ctx) uint {
	token, ok := ctx.Locals("user").(*jwt.Token)
  if ok != true {
    return 0
  }
	claims := token.Claims.(jwt.MapClaims)
	userID := claims["sub"].(float64)
  return uint(userID)
}

func createToken(userID uint, isSuperuser bool) (string, error) {
  claims := jwt.MapClaims{
    "sub": userID,
		"is_superuser": isSuperuser,
		"exp":   time.Now().Add(time.Hour * 72).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString(secretKey)
  return t, err
}


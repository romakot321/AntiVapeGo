package services

import (
  "gorm.io/gorm"
  models "antivape/db"
  "antivape/schemas"
)

type UserService interface {
  TakeByID(userID uint) schemas.UserSchema
  TakeByName(userName string) schemas.UserSchema
  Find(schema schemas.UserFindSchema) []schemas.UserSchema
  Update(userID uint, schema schemas.UserUpdateSchema)
  Delete(userID uint)
}

type userService struct {
  baseService
  db *gorm.DB
}

func (s userService) modelToSchema(model models.User) schemas.UserSchema {
  return schemas.UserSchema{
    ID: model.ID,
    Name: model.Name,
    IsSuperuser: model.IsSuperuser,
  }
}

func (s userService) TakeByID(userID uint) schemas.UserSchema {
  var model models.User
  s.take(userID, &model, nil)
  return s.modelToSchema(model)
}

func (s userService) TakeByName(userName string) schemas.UserSchema {
  var model models.User
  s.takeByField("name = ?", userName, model, nil)
  return s.modelToSchema(model)
}

func (s userService) Find(schema schemas.UserFindSchema) []schemas.UserSchema {
  var models []models.User
  filters := schemas.SchemaToMap(schema)
  s.find(&models, filters)
  
  var returnSchemas []schemas.UserSchema
  for _, model := range models {
    returnSchemas = append(
      returnSchemas,
      s.modelToSchema(model),
    )
  }
  return returnSchemas
}

func (s userService) Update(userID uint, schema schemas.UserUpdateSchema) {
  m := schemas.SchemaToMap(schema)
  s.update(&models.User{}, userID, m)
}

func (s userService) Delete(userID uint) {
  s.delete(&models.User{}, userID)
}

func NewUserService(db *gorm.DB) UserService {
  return userService{baseService: baseService{db: db}}
}

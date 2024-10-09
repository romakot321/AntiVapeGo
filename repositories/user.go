package repositories

import (
  "log"

  "gorm.io/gorm"
  models "antivape/db"
)

type UserRepository interface {
  TakeByID(userID uint) models.User
  TakeByName(userName string) models.User
  Create(name, password string) models.User
  Find(filters map[string]interface{}) []models.User
  Delete(userID uint)
}

type userRepository struct {
  baseRepository
  db *gorm.DB
}

func (s userRepository) TakeByID(userID uint) models.User {
  var model models.User
  s.take(userID, &model, nil)
  return model
}

func (s userRepository) TakeByName(username string) models.User {
  var model models.User
  s.takeByField("name = ?", username, &model, nil)
  return model
}

func (s userRepository) Create(name, passwordHash string) models.User {
  model := models.User{
    Name: name,
    PasswordHash: passwordHash,
  }
  s.create(&model)
  return model
}

func(s userRepository) Find(filters map[string]interface{}) []models.User {
  var users []models.User
  query := s.db.Select("*")
  roomID, ok := filters["room_id"]
  if ok {
    query = query.Where("guid = (?)", s.db.Table("sensors").Select("sensors.guid").Where("sensors.room_id = ?", roomID.(string)))
  }
  zoneID, ok := filters["zone_id"]
  if ok {
    query = query.Where("guid = (?)", s.db.Table("sensors").Select("sensors.guid").Where("sensors.zone_id = ?", zoneID.(string)))
  }
  result := query.Find(&users)
  if result.Error != nil {
    log.Println("Error find sensordata", result.Error)
  }
  return users
}

func (s userRepository) Update(userID uint, fields map[string]interface{}) {
  s.update(&models.User{}, userID, fields)
}

func (s userRepository) Delete(userID uint) {
  s.delete(&models.User{}, userID)
}

func NewUserRepository(db *gorm.DB) UserRepository {
  return userRepository{baseRepository: baseRepository{db: db}}
}

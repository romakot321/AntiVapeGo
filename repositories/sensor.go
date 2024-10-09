package repositories

import (
  "gorm.io/gorm"
  models "antivape/db"
)

type SensorRepository interface {
  Take(sensorID uint) models.Sensor
  Create(guid, name string, roomID, zoneID, ownerID uint) models.Sensor
  Find(filters map[string]interface{}) []models.Sensor
  Delete(sensorID uint)
}

type sensorRepository struct {
  baseRepository
  db *gorm.DB
}

func (s sensorRepository) Take(sensorID uint) models.Sensor {
  var model models.Sensor
  s.take(sensorID, &model, nil)
  return model
}

func (s sensorRepository) Create(guid, name string, roomID, zoneID, ownerID uint) models.Sensor {
  model := models.Sensor{
    Guid: guid,
    Name: name,
    RoomID: roomID,
    ZoneID: zoneID,
    OwnerID: ownerID,
  }
  s.create(&model)
  return model
}

func (s sensorRepository) Find(filters map[string]interface{}) []models.Sensor {
  var sensors []models.Sensor
  s.find(&sensors, filters)
  return sensors
}

func (s sensorRepository) Update(sensorID uint, fields map[string]interface{}) {
  s.update(&models.Sensor{}, sensorID, fields)
}

func (s sensorRepository) Delete(sensorID uint) {
  s.delete(&models.Sensor{}, sensorID)
}

func NewSensorRepository(db *gorm.DB) SensorRepository {
  return sensorRepository{baseRepository: baseRepository{db: db}}
}

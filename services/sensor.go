package services

import (
  "gorm.io/gorm"
  models "antivape/db"
  "antivape/schemas"
)

type SensorService interface {
  Take(sensorID uint) schemas.SensorSchema
  Create(schema schemas.SensorCreateSchema) schemas.SensorSchema
  Find(schema schemas.SensorFindSchema) []schemas.SensorSchema
  Update(sensorID uint, schema schemas.SensorUpdateSchema)
  Delete(sensorID uint)
  FilterByOwnerID(ownerID uint, sensors ...schemas.SensorSchema) []schemas.SensorSchema
}

type sensorService struct {
  baseService
  db *gorm.DB
}

func (s sensorService) modelToSchema(model models.Sensor) schemas.SensorSchema {
  return schemas.SensorSchema{
    ID: model.ID,
    Name: model.Name,
    Guid: model.Guid,
    RoomID: model.RoomID,
    OwnerID: model.OwnerID,
  }
}

func (s sensorService) Take(sensorID uint) schemas.SensorSchema {
  var model models.Sensor
  s.take(sensorID, &model, nil)
  return s.modelToSchema(model)
}

func (s sensorService) Create(schema schemas.SensorCreateSchema) schemas.SensorSchema {
  var room models.Room
  s.take(schema.RoomID, &room, nil)

  model := models.Sensor{
    Name: schema.Name,
    Guid: schema.Guid,
    RoomID: schema.RoomID,
    ZoneID: room.ZoneID,
    OwnerID: schema.OwnerID,
  }
  s.create(&model)
  return s.modelToSchema(model)
}

func (s sensorService) Find(schema schemas.SensorFindSchema) []schemas.SensorSchema {
  var models []models.Sensor
  filters := schemas.SchemaToMap(schema)
  s.find(&models, filters)
  
  var returnSchemas []schemas.SensorSchema
  for _, model := range models {
    returnSchemas = append(
      returnSchemas,
      s.modelToSchema(model),
    )
  }
  return returnSchemas
}

func (s sensorService) Update(sensorID uint, schema schemas.SensorUpdateSchema) {
  m := schemas.SchemaToMap(schema)
  s.update(&models.Sensor{}, sensorID, m)
}

func (s sensorService) Delete(sensorID uint) {
  s.delete(&models.Sensor{}, sensorID)
}

func (s sensorService) FilterByOwnerID(ownerID uint, sensors ...schemas.SensorSchema) []schemas.SensorSchema {
  var filtered []schemas.SensorSchema
  for _, sensor := range sensors {
    if sensor.OwnerID != ownerID { continue }
    filtered = append(filtered, sensor)
  }
  return filtered
}

func NewSensorService(db *gorm.DB) SensorService {
  return sensorService{baseService: baseService{db: db}}
}

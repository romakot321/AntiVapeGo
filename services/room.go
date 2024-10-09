package services

import (
  "gorm.io/gorm"
  models "antivape/db"
  "antivape/schemas"
  "antivape/repositories"
)

type RoomService interface {
  Take(roomID uint) schemas.RoomSchema
  Create(schema schemas.RoomCreateSchema) schemas.RoomSchema
  Find(schema schemas.RoomFindSchema) []schemas.RoomSchema
  Update(roomID uint, schema schemas.RoomUpdateSchema)
  Delete(roomID uint)
  GetStatistic(roomID uint) schemas.SensorDataRoomSchema
  FilterByOwnerID(ownerID uint, rooms ...schemas.RoomSchema) []schemas.RoomSchema
}

type roomService struct {
  baseService
  db *gorm.DB
  sensorDataRep repositories.SensorDataRepository
}

func (s roomService) modelToSchema(model models.Room) schemas.RoomSchema {
  var sensors []schemas.SensorSchema
  if len(model.Sensors) > 0 {
    sensors = append(
      sensors,
      schemas.SensorSchema{Name: model.Sensors[0].Name, Guid: model.Sensors[0].Guid, RoomID: model.Sensors[0].RoomID},
    )
  }
  return schemas.RoomSchema{
    ID: model.ID,
    Name: model.Name,
    OwnerID: model.OwnerID,
    ZoneID: model.ZoneID,
    Sensors: sensors,
  }
}

func (s roomService) Take(roomID uint) schemas.RoomSchema {
  var model models.Room
  s.take(roomID, &model, "Sensors")
  return s.modelToSchema(model)
}

func (s roomService) Create(schema schemas.RoomCreateSchema) schemas.RoomSchema {
  model := models.Room{
    Name: schema.Name,
    OwnerID: schema.OwnerID,
    ZoneID: schema.ZoneID,
  }
  s.create(&model)
  return s.modelToSchema(model)
}

func (s roomService) Find(schema schemas.RoomFindSchema) []schemas.RoomSchema {
  var models []models.Room
  filters := schemas.SchemaToMap(schema)
  s.find(&models, filters)
  
  var returnSchemas []schemas.RoomSchema
  for _, model := range models {
    returnSchemas = append(
      returnSchemas,
      s.modelToSchema(model),
    )
  }
  return returnSchemas
}

func (s roomService) Update(roomID uint, schema schemas.RoomUpdateSchema) {
  m := schemas.SchemaToMap(schema)
  s.update(&models.Room{}, roomID, m)
}

func (s roomService) Delete(roomID uint) {
  s.delete(&models.Room{}, roomID)
}

func (s roomService) GetStatistic(roomID uint) schemas.SensorDataRoomSchema {
  filters := make(map[string]interface{}, 1)
  filters["room_id"] = roomID
  statistic := s.sensorDataRep.GetStatistic(filters)
  return schemas.SensorDataRoomSchema{Co2: statistic[0].Co2, Tvoc: statistic[0].Tvoc, RoomID: roomID}
}

func (s roomService) FilterByOwnerID(ownerID uint, rooms ...schemas.RoomSchema) []schemas.RoomSchema {
  var filtered []schemas.RoomSchema
  for _, room := range rooms {
    if room.OwnerID != ownerID { continue }
    filtered = append(filtered, room)
  }
  return filtered
}

func NewRoomService(db *gorm.DB, sensorDataRep repositories.SensorDataRepository) RoomService {
  return roomService{baseService: baseService{db: db}, sensorDataRep: sensorDataRep}
}

package services

import (
  "gorm.io/gorm"
  models "antivape/db"
  "antivape/schemas"
  "antivape/repositories"
)

type ZoneService interface {
  Take(zoneID uint) schemas.ZoneSchema
  Create(schema schemas.ZoneCreateSchema) schemas.ZoneSchema
  Find(schema schemas.ZoneFindSchema) []schemas.ZoneSchema
  Update(zoneID uint, schema schemas.ZoneUpdateSchema)
  Delete(zoneID uint)
  FilterByOwnerID(ownerID uint, zones ...schemas.ZoneSchema) []schemas.ZoneSchema
  GetStatistic(zoneID uint) schemas.SensorDataZoneSchema
}

type zoneService struct {
  baseService
  db *gorm.DB
  sensorDataRep repositories.SensorDataRepository
  sensorRep repositories.SensorRepository
}

func (s zoneService) modelToSchema(model models.Zone) schemas.ZoneSchema {
  var rooms []schemas.RoomSchema
  for _, room := range model.Rooms {
    rooms = append(
      rooms,
      schemas.RoomSchema{Name: room.Name, ZoneID: room.ZoneID, ID: room.ID, OwnerID: room.OwnerID},
    )
  }
  return schemas.ZoneSchema{
    ID: model.ID,
    Name: model.Name,
    OwnerID: model.OwnerID,
    Rooms: rooms,
  }
}

func (s zoneService) Take(zoneID uint) schemas.ZoneSchema {
  var model models.Zone
  s.take(zoneID, &model, "Rooms")
  return s.modelToSchema(model)
}

func (s zoneService) Create(schema schemas.ZoneCreateSchema) schemas.ZoneSchema {
  model := models.Zone{
    Name: schema.Name,
    OwnerID: schema.OwnerID,
  }
  s.create(&model)
  return s.modelToSchema(model)
}

func (s zoneService) Find(schema schemas.ZoneFindSchema) []schemas.ZoneSchema {
  var models []models.Zone
  filters := schemas.SchemaToMap(schema)
  s.find(&models, filters)
  
  var returnSchemas []schemas.ZoneSchema
  for _, model := range models {
    returnSchemas = append(
      returnSchemas,
      s.modelToSchema(model),
    )
  }
  return returnSchemas
}

func (s zoneService) Update(zoneID uint, schema schemas.ZoneUpdateSchema) {
  m := schemas.SchemaToMap(schema)
  s.update(&models.Zone{}, zoneID, m)
}

func (s zoneService) Delete(zoneID uint) {
  s.delete(&models.Zone{}, zoneID)
}

func (s zoneService) GetStatistic(zoneID uint) schemas.SensorDataZoneSchema {
  filters := make(map[string]interface{}, 1)
  filters["zone_id"] = zoneID

  statistic := s.sensorDataRep.GetStatistic(filters)

  return schemas.SensorDataZoneSchema{Rooms: statistic, ZoneID: zoneID}
}

func (s zoneService) FilterByOwnerID(ownerID uint, zones ...schemas.ZoneSchema) []schemas.ZoneSchema {
  var filtered []schemas.ZoneSchema
  for _, zone := range zones {
    if zone.OwnerID != ownerID { continue }
    filtered = append(filtered, zone)
  }
  return filtered
}

func NewZoneService(db *gorm.DB, sensorDataRep repositories.SensorDataRepository) ZoneService {
  return zoneService{baseService: baseService{db: db}, sensorDataRep: sensorDataRep}
}

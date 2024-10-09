package services

import (
  "gorm.io/gorm"
  models "antivape/db"
  "antivape/schemas"
)

type ZoneService interface {
  Take(zoneID uint) schemas.ZoneSchema
  Create(schema schemas.ZoneCreateSchema) schemas.ZoneSchema
  Find(schema schemas.ZoneFindSchema) []schemas.ZoneSchema
  Update(zoneID uint, schema schemas.ZoneUpdateSchema)
  Delete(zoneID uint)
  FilterByOwnerID(ownerID uint, zones ...schemas.ZoneSchema) []schemas.ZoneSchema
}

type zoneService struct {
  baseService
  db *gorm.DB
}

func (s zoneService) modelToSchema(model models.Zone) schemas.ZoneSchema {
  var rooms []schemas.RoomSchema
  if len(model.Rooms) > 0 {
    rooms = append(
      rooms,
      schemas.RoomSchema{Name: model.Rooms[0].Name, ZoneID: model.Rooms[0].ZoneID},
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

func (s zoneService) GetStatistics(zoneID uint) {
  var sensors []models.Sensor
  // var sensorData []models.SensorData
  var filters map[string]interface{}
  filters["ZoneID"] = zoneID

  s.find(&sensors, filters)
}

func (s zoneService) FilterByOwnerID(ownerID uint, zones ...schemas.ZoneSchema) []schemas.ZoneSchema {
  var filtered []schemas.ZoneSchema
  for _, zone := range zones {
    if zone.OwnerID != ownerID { continue }
    filtered = append(filtered, zone)
  }
  return filtered
}

func NewZoneService(db *gorm.DB) ZoneService {
  return zoneService{baseService: baseService{db: db}}
}

package schemas

import (
  models "antivape/db"
)

type RoomCreateSchema struct {
  Name string `json:"name" binding:"required"`
  ZoneID uint `json:"zone_id" binding:"required"`
  OwnerID uint `json:"owner_id" binding:"required"`
}

type RoomSchema struct {
  ID uint `json:"id" binding:"required"`
  Name string `json:"name" binding:"required"`
  ZoneID uint `json:"zone_id" binding:"required"`
  OwnerID uint `json:"owner_id" binding:"required"`
  Sensors []SensorSchema `json:"sensors`
}

type RoomUpdateSchema struct {
  Name string `json:"name,omitempty"`
}

type RoomFindSchema struct {
  OwnerID *uint `json:"owner_id,omitempty"`
  ZoneID *uint `json:"zone_id,omitempty"`
}

func (s RoomSchema) ToModel() models.Room {
  sensors := make([]models.Sensor, len(s.Sensors))
  for _, schema := range s.Sensors {
    sensors = append(sensors, schema.ToModel())
  }
  return models.Room{
    Name: s.Name,
    ZoneID: s.ZoneID,
    OwnerID: s.OwnerID,
    Sensors: sensors,
  }
}

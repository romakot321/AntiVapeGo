package schemas

import (
  models "antivape/db"
)

type SensorCreateSchema struct {
  Name string `json:"name" binding:"required"`
  Guid string `json:"guid" binding:"required"`
  RoomID uint `json:"room_id" binding:"required"`
  OwnerID uint `json:"owner_id" binding:"required"`
}

type SensorSchema struct {
  ID uint `json:"id" binding:"required"`
  Name string `json:"name" binding:"required"`
  Guid string `json:"guid" binding:"required"`
  RoomID uint `json:"room_id" binding:"required"`
  OwnerID uint `json:"owner_id" binding:"required"`
}

type SensorUpdateSchema struct {
  Name string `json:"name,omitempty"`
  Guid string `json:"guid,omitempty"`
}

type SensorFindSchema struct {
  RoomID uint `json:"room_id,omitempty"`
  OwnerID *uint `json:"owner_id,omitempty"`
}

func (s SensorSchema) ToModel() models.Sensor {
  return models.Sensor{
    Name: s.Name,
    Guid: s.Guid,
    RoomID: s.RoomID,
    OwnerID: s.OwnerID,
  }
}

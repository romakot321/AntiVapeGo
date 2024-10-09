package db

import (
    "gorm.io/gorm"
)

type ownableModel struct {
}

type Sensor struct {
  gorm.Model
  Name string
  Guid string
  RoomID uint
  ZoneID uint
  OwnerID uint
}

type Room struct {
  gorm.Model
  Name string
  ZoneID uint
  OwnerID uint
  Sensors []Sensor `gorm:"foreignKey:RoomID"`
}

type Zone struct {
  gorm.Model
  Name string
  OwnerID uint
  Rooms []Room `gorm:"foreignKey:ZoneID"`
}

type User struct {
  gorm.Model
  Name string `gorm:"index,unique"`
  PasswordHash string
  IsSuperuser bool `gorm:"default:false"`
  OwnerID uint
}

type SensorData struct {
  gorm.Model
  Guid string `gorm:"index"`
  Co2 int
  Tvoc int
  BatteryCharge int
}

func MigrateModels(db *gorm.DB) {
  db.AutoMigrate(&Sensor{})
  db.AutoMigrate(&Room{})
  db.AutoMigrate(&Zone{})
  db.AutoMigrate(&User{})
  db.AutoMigrate(&SensorData{})
}

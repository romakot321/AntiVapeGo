package schemas

type SensorDataSchema struct {
  Guid string
  Co2 int
  Tvoc int
  BatteryCharge int
}

type SensorDataRoomSchema struct {
  Co2 int
  Tvoc int
  RoomID uint
}

type SensorDataZoneSchema struct {
  Rooms []SensorDataRoomSchema
  ZoneID uint
}

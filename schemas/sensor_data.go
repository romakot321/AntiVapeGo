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

func SensorDataRoomSchemaFromArray(roomID uint, data []SensorDataSchema) SensorDataRoomSchema {
  var co2, tvoc int
  for _, sensorData := range data {
    co2 += sensorData.Co2
    tvoc += sensorData.Tvoc
  }
  co2 /= len(data)
  tvoc /= len(data)
  return SensorDataRoomSchema{Co2: co2, Tvoc: tvoc, RoomID: roomID}
}

func SensorDataZoneSchemaFromArray(zoneID uint, data []SensorDataSchema, sensors []SensorSchema) SensorDataZoneSchema {
  var guidToSensor map[string]SensorSchema
  for _, sensor := range sensors {
    guidToSensor[sensor.Guid] = sensor
  }

  var roomIDToSensorData map[uint][]SensorDataSchema
  for _, sensorData := range data {
    sensor, ok := guidToSensor[sensorData.Guid]
    if !ok {
      continue
    }
    var curr []SensorDataSchema
    curr, ok = roomIDToSensorData[sensor.RoomID]
    roomIDToSensorData[sensor.RoomID] = append(curr, sensorData)
  }

  rooms := make([]SensorDataRoomSchema, len(roomIDToSensorData))
  for roomID, sensorDatas := range roomIDToSensorData {
    rooms = append(rooms, SensorDataRoomSchemaFromArray(roomID, sensorDatas))
  }
  return SensorDataZoneSchema{Rooms: rooms, ZoneID: zoneID}
}

package repositories

import (
  "log"
  "strconv"
  "gorm.io/gorm"
  models "antivape/db"
  "antivape/schemas"
)

const room_statistic_query string = `
SELECT AVG(sensor_data.co2) AS co2, AVG(sensor_data.tvoc) AS tvoc, sensors.room_id
FROM sensor_data JOIN sensors ON sensor_data.guid=sensors.guid
WHERE sensor_data.guid in (select sensors.guid from sensors where sensors.room_id=?) GROUP BY sensors.room_id;
`

const zone_statistic_query string = `
SELECT AVG(sensor_data.co2) AS co2, AVG(sensor_data.tvoc) AS tvoc, sensors.room_id
FROM sensor_data JOIN sensors ON sensor_data.guid=sensors.guid
WHERE sensor_data.guid in (select sensors.guid from sensors where sensors.zone_id=?) GROUP BY sensors.room_id;
`

type dbDataSchema struct {
  Co2 string
  Tvoc string
  RoomID string
}

type SensorDataRepository interface {
  Take(sensorDataID uint) models.SensorData
  Create(guid string, co2, tvoc, batteryCharge int) models.SensorData
  GetStatistic(filters map[string]interface{}) []schemas.SensorDataRoomSchema
  Delete(sensorDataID uint)
}

type sensorDataRepository struct {
  baseRepository
  db *gorm.DB
}

func (s sensorDataRepository) Take(sensorDataID uint) models.SensorData {
  var model models.SensorData
  s.take(sensorDataID, &model, nil)
  return model
}

func (s sensorDataRepository) Create(guid string, co2, tvoc, batteryCharge int) models.SensorData {
  model := models.SensorData{
    Guid: guid,
    Co2: co2,
    Tvoc: tvoc,
    BatteryCharge: batteryCharge,
  }
  s.create(&model)
  return model
}

func (s sensorDataRepository) GetStatistic(filters map[string]interface{}) []schemas.SensorDataRoomSchema {
  statistic := make([]dbDataSchema, 0)
  if roomID, ok := filters["room_id"]; ok {
    var schema dbDataSchema
    s.baseRepository.db.Raw(room_statistic_query, roomID).Take(&schema)
    statistic = append(statistic, schema)
  } else if zoneID, ok := filters["zone_id"]; ok {
    s.baseRepository.db.Raw(zone_statistic_query, zoneID).Find(&statistic)
  }
  resp := make([]schemas.SensorDataRoomSchema, 0, len(statistic))
  for _, schema := range statistic {
    co2, _ := strconv.ParseFloat(schema.Co2, 64)
    tvoc, _ := strconv.ParseFloat(schema.Tvoc, 64)
    roomID, _ := strconv.Atoi(schema.RoomID)
    resp = append(
      resp,
      schemas.SensorDataRoomSchema{Co2: int(co2), Tvoc: int(tvoc), RoomID: uint(roomID)},
    )
  }
  return resp
}

func (s sensorDataRepository) Update(sensorDataID uint, fields map[string]interface{}) {
  s.update(&models.SensorData{}, sensorDataID, fields)
}

func (s sensorDataRepository) Delete(sensorDataID uint) {
  s.delete(&models.SensorData{}, sensorDataID)
}

func NewSensorDataRepository(db *gorm.DB) SensorDataRepository {
  return sensorDataRepository{baseRepository: baseRepository{db: db}}
}

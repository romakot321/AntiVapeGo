package services

import (
  "context"
  "log"
  "time"

  "antivape/schemas"
  models "antivape/db"
  "gorm.io/gorm"
  "github.com/redis/go-redis/v9"
)

type ExternalService interface {
  Store(schema schemas.ExternalSensorDataSchema)
  PopAll() []schemas.ExternalSensorDataSchema
  RunTransferingCycle()
}

type externalService struct {
  baseService
  redisConn *redis.Client
  db *gorm.DB
  ctx context.Context
}

func (s externalService) Store(schema schemas.ExternalSensorDataSchema) {
  key := generateKey(schema.Guid)
  s.redisConn.HSet(s.ctx, key, "guid", schema.Guid).Err()
  s.redisConn.HSet(s.ctx, key, "co2", schema.Co2).Err()
  s.redisConn.HSet(s.ctx, key, "tvoc", schema.Tvoc).Err()
  err := s.redisConn.HSet(s.ctx, key, "batteryCharge", schema.BatteryCharge).Err()
  if err != nil {
    log.Fatal(err)
  }
}

func (s externalService) PopAll() []schemas.ExternalSensorDataSchema {
  var data []schemas.ExternalSensorDataSchema
  var schema schemas.ExternalSensorDataSchema

  iter := s.redisConn.Scan(s.ctx, 0, "*", 0).Iterator()
  for iter.Next(s.ctx) {
    if err := s.redisConn.HGetAll(s.ctx, iter.Val()).Scan(&schema); err != nil {
      log.Fatal(err)
    }
    s.redisConn.Del(s.ctx, iter.Val())
    data = append(data, schema)
  }
  if err := iter.Err(); err != nil {
    log.Fatal(err)
  }

  return data
}

func (s externalService) RunTransferingCycle() {
  for range(time.Tick(time.Second * 3)) {
    data := s.PopAll()
    var dataModels []models.SensorData
    for _, sensorData := range data {
      dataModels = append(
        dataModels,
        models.SensorData{
          Guid: sensorData.Guid,
          Co2: sensorData.Co2,
          Tvoc: sensorData.Tvoc,
          BatteryCharge: sensorData.BatteryCharge,
        },
      )
    }
    s.create(&dataModels)
  }
}

func NewExternalService(redisConn *redis.Client, db *gorm.DB) ExternalService {
  ctx := context.Background()
  return externalService{redisConn: redisConn, ctx: ctx, baseService: baseService{db: db}}
}

func generateKey(guid string) string {
  return time.Now().String() + "-" + guid
}

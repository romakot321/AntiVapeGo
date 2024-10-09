package db

import (
  "gorm.io/gorm"
  "gorm.io/driver/postgres"
)

func InitDatabase(dsn string) (*gorm.DB, error) {
  db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{CreateBatchSize: 1000})
  MigrateModels(db)
  return db, err
}

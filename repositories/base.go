package repositories

import (
  "log"
  "gorm.io/gorm"
)

type baseRepository struct {
  db *gorm.DB
}

func (s baseRepository) take(modelID uint, model interface{}, preload interface{}) {
  query := s.db.Where("id = ?", modelID)
  if preload != nil {
    query = query.Preload(preload.(string))
  }
  result := query.Take(model)
  if result.Error != nil {
    log.Fatal("Error take model: ", result.Error)
  }
}

func (s baseRepository) takeByField(expression string, value interface{}, model interface{}, preload interface{}) {
  query := s.db.Where(expression, value)
  if preload != nil {
    query = query.Preload(preload.(string))
  }
  result := query.Take(model)
  if result.Error != nil {
    log.Fatal("Error take model: ", result.Error)
  }
}

func (s baseRepository) create(model interface{}) {
  result := s.db.Create(model)
  if result.Error != nil {
    log.Fatal("Error create model: ", result.Error)
  }
}

func (s baseRepository) find(models interface{}, filters map[string]interface{}) {
  result := s.db.Where(filters).Find(models)
  if result.Error != nil {
    log.Fatal("Error find models: ", result.Error)
  }
}

func (s baseRepository) update(model interface{}, modelID uint, fields map[string]interface{}) {
  result := s.db.Model(model).Where("id = ?", modelID).Updates(fields)
  if result.Error != nil {
    log.Fatal("Error update model: ", result.Error)
  }
}

func (s baseRepository) delete(model interface{}, modelID uint) {
  result := s.db.Delete(model, modelID)
  if result.Error != nil {
    log.Fatal("Error delete model: ", result.Error)
  }
}

package services

import (
  "log"
  "gorm.io/gorm"
  "gorm.io/gorm/clause"
)

type baseService struct {
  db *gorm.DB
}

func (s baseService) take(modelID uint, model interface{}, preload interface{}) error {
  query := s.db.Where("id = ?", modelID)
  if preload != nil {
    query = query.Preload(clause.Associations)
  }
  result := query.Take(model)
  if result.Error != nil {
    log.Println("Error take model: ", result.Error)
    return result.Error
  }
  return nil
}

func (s baseService) takeByField(expression string, value interface{}, model interface{}, preload interface{}) error {
  query := s.db.Where(expression, value)
  if preload != nil {
    query = query.Preload(clause.Associations)
  }
  result := query.Take(model)
  if result.Error != nil {
    log.Println("Error take model: ", result.Error)
    return result.Error
  }
  return nil
}

func (s baseService) create(model interface{}) error {
  result := s.db.Create(model)
  if result.Error != nil {
    log.Println("Error create model: ", result.Error)
    return result.Error
  }
  return nil
}

func (s baseService) find(models interface{}, filters map[string]interface{}) error {
  result := s.db.Where(filters).Find(models)
  if result.Error != nil {
    log.Println("Error find models: ", result.Error)
    return result.Error
  }
  return nil
}

func (s baseService) update(model interface{}, modelID uint, fields map[string]interface{}) error {
  result := s.db.Model(model).Where("id = ?", modelID).Updates(fields)
  if result.Error != nil {
    log.Println("Error update model: ", result.Error)
    return result.Error
  }
  return nil
}

func (s baseService) delete(model interface{}, modelID uint) error {
  result := s.db.Delete(model, modelID)
  if result.Error != nil {
    log.Println("Error delete model: ", result.Error)
    return result.Error
  }
  return nil
}

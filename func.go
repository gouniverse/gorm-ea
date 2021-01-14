package ea

import (
	"errors"
	"log"

	"gorm.io/gorm"
)

// EntityAttributeFind finds an entity by ID
func EntityAttributeFind(db *gorm.DB, entityID string, attributeKey string) *EntityAttribute {
	entityAttribute := &EntityAttribute{}

	result := db.First(&entityAttribute, "entity_id=? AND attribute_key=?", entityID, attributeKey)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		log.Panic(result.Error)
	}

	return entityAttribute
}

// EntityAttributesUpsert upserts and entity attribute
func EntityAttributesUpsert(db *gorm.DB, entityID string, attributes map[string]string) bool {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false
	}

	for k, v := range attributes {
		entityAttribute := EntityAttributeFind(db, entityID, k)

		if entityAttribute == nil {
			entityAttribute = &EntityAttribute{EntityID: entityID, AttributeKey: k}
			entityAttribute.AttributeValue = v

			dbResult := tx.Create(&entityAttribute)
			if dbResult.Error != nil {
				tx.Rollback()
				return false
			}

		}

		entityAttribute.AttributeValue = v
		dbResult := tx.Save(entityAttribute)
		if dbResult.Error != nil {
			return false
		}
	}

	err := tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return false
	}

	return true

}

// EntityAttributesUpsertJSON upserts and entity attribute
func EntityAttributesUpsertJSON(db *gorm.DB, entityID string, attributes map[string]interface{}) bool {
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return false
	}

	for k, v := range attributes {
		entityAttribute := EntityAttributeFind(db, entityID, k)

		if entityAttribute == nil {
			entityAttribute = &EntityAttribute{EntityID: entityID, AttributeKey: k}
			entityAttribute.SetJSONValue(v)

			dbResult := tx.Create(&entityAttribute)
			if dbResult.Error != nil {
				tx.Rollback()
				return false
			}

		}

		entityAttribute.SetJSONValue(v)
		dbResult := tx.Save(entityAttribute)
		if dbResult.Error != nil {
			return false
		}
	}

	err := tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return false
	}

	return true

}

// EntityAttributeUpsert upserts and entity attribute
func EntityAttributeUpsert(db *gorm.DB, entityID string, attributeKey string, attributeValue interface{}) bool {
	entityAttribute := EntityAttributeFind(db, entityID, attributeKey)

	if entityAttribute == nil {
		entityAttribute = &EntityAttribute{EntityID: entityID, AttributeKey: attributeKey}
		entityAttribute.SetJSONValue(attributeValue)

		dbResult := db.Create(&entityAttribute)
		if dbResult.Error != nil {
			return false
		}

		return true
	}

	entityAttribute.SetJSONValue(attributeValue)
	dbResult := db.Save(entityAttribute)
	if dbResult.Error != nil {
		return false
	}

	return true

}

// EntityCreate creates a new entity
func EntityCreate(db *gorm.DB, entityType string) *Entity {
	entity := &Entity{Type: entityType, Status: EntityStatusActive}

	dbResult := db.Create(&entity)

	if dbResult.Error != nil {
		return nil
	}

	return entity
}

// EntityCreateWithAttributes func
func EntityCreateWithAttributes(db *gorm.DB, entityType string, attributes map[string]string) *Entity {
	// Note the use of tx as the database handle once you are within a transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		return nil
	}

	//return tx.Commit().Error

	entity := &Entity{Type: entityType, Status: EntityStatusActive}

	dbResult := tx.Create(&entity)

	if dbResult.Error != nil {
		tx.Rollback()
		return nil
	}

	//entityAttributes := make([]EntityAttribute, 0)
	for k, v := range attributes {
		ea := EntityAttribute{EntityID: entity.ID, AttributeKey: k} //, AttributeValue: value}
		ea.AttributeValue = v

		dbResult2 := tx.Create(&ea)
		if dbResult2.Error != nil {
			tx.Rollback()
			return nil
		}
	}

	err := tx.Commit().Error

	if err != nil {
		tx.Rollback()
		return nil
	}

	return entity
}

// EntityDelete deletes an entity and all attributes
func EntityDelete(db *gorm.DB, entityID string) bool {
	if entityID == "" {
		return false
	}

	// Note the use of tx as the database handle once you are within a transaction
	tx := db.Begin()

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Error; err != nil {
		log.Println(err)
		return false
	}

	if err := tx.Where("entity_id=?", entityID).Delete(&EntityAttribute{}).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return false
	}

	if err := tx.Where("id=?", entityID).Delete(&Entity{}).Error; err != nil {
		tx.Rollback()
		log.Println(err)
		return false
	}

	err := tx.Commit().Error

	if err == nil {
		return true
	}

	log.Println(err)

	return false
}

// EntityFindByID finds an entity by ID
func EntityFindByID(db *gorm.DB, entityID string) *Entity {
	if entityID == "" {
		return nil
	}

	entity := &Entity{}

	resultEntity := db.First(&entity, "id=?", entityID)

	if resultEntity.Error != nil {
		if errors.Is(resultEntity.Error, gorm.ErrRecordNotFound) {
			return nil
		}
		log.Panic(resultEntity.Error)
	}

	// DEBUG: log.Println(entity)

	return entity
}

// EntityFindByAttribute finds an entity by attribute
func EntityFindByAttribute(db *gorm.DB, entityType string, attributeKey string, attributeValue string) *Entity {
	entityAttribute := &EntityAttribute{}

	subQuery := db.Model(&Entity{}).Select("id").Where("type = ?", entityType)
	result := db.First(&entityAttribute, "entity_id IN (?) AND attribute_key=? AND attribute_value=?", subQuery, attributeKey, attributeValue)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		log.Panic(result.Error)
	}

	// DEBUG: log.Println(entityAttribute)

	entity := &Entity{}

	resultEntity := db.First(&entity, "id=?", entityAttribute.EntityID)

	if errors.Is(resultEntity.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if resultEntity.Error != nil {
		log.Panic(resultEntity.Error)
	}

	// DEBUG: log.Println(entity)

	return entity
}

// EntityListByAttribute finds an entity by attribute
func EntityListByAttribute(db *gorm.DB, entityType string, attributeKey string, attributeValue string) []Entity {
	//entityAttributes := []EntityAttribute{}
	var entityIDs []string

	subQuery := db.Model(&Entity{}).Select("id").Where("type = ?", entityType)
	result := db.Model(&EntityAttribute{}).Select("entity_id").Find(&entityIDs, "entity_id IN (?) AND attribute_key=? AND attribute_value=?", subQuery, attributeKey, attributeValue)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if result.Error != nil {
		log.Panic(result.Error)
	}

	// DEBUG: log.Println(result)

	entities := []Entity{}

	resultEntity := db.Where("id IN (?)", entityIDs).Find(&entities)

	if errors.Is(resultEntity.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	if resultEntity.Error != nil {
		log.Panic(resultEntity.Error)
	}

	// DEBUG: log.Println(entity)

	return entities
}

// EntityList lists entities
func EntityList(db *gorm.DB, entityType string, offset uint64, perPage uint64, search string, orderBy string, sort string) []Entity {
	entityList := []Entity{}
	result := db.Where("type=?", entityType).Order(orderBy + " " + sort).Offset(int(offset)).Limit(int(perPage)).Find(&entityList)

	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return nil
	}

	return entityList
}

// EntityCount counts entities
func EntityCount(db *gorm.DB, entityType string) uint64 {
	var count int64
	db.Model(&Entity{}).Where("type=?", entityType).Count(&count)
	return uint64(count)
	// sqlStr, args, _ := squirrel.Select("COUNT(*) AS count").From(TableArticle).Limit(1).ToSql()

	// entities := Query(sqlStr, args...)

	// count, _ := strconv.ParseUint(entities[0]["count"], 10, 64)

	// return count
}

// func Tree(id string, name string, parentId string) {

// }

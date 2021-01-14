package ea

import (
	"errors"
	"log"
	"time"

	"github.com/gouniverse/uid"
	"gorm.io/gorm"
)

const (
	// EntityStatusActive entity "active" status
	EntityStatusActive = "active"
	// EntityStatusInactive entity "inactive" status
	EntityStatusInactive = "inactive"
)

// Entity type
type Entity struct {
	ID     string `gorm:"type:varchar(40);column:id;primary_key;"`
	Status string `gorm:"type:varchar(10);column:status;"`
	Type   string `gorm:"type:varchar(40);column:type;"`
	//Name        string     `gorm:"type:varchar(255);column:name;DEFAULT NULL;"`
	//Description string     `gorm:"type:longtext;column:description;"`
	CreatedAt time.Time  `gorm:"type:datetime;column:created_at;DEFAULT NULL;"`
	UpdatedAt time.Time  `gorm:"type:datetime;column:updated_at;DEFAULT NULL;"`
	DeletedAt *time.Time `gorm:"type:datetime;column:deleted_at;DEFAULT NULL;"`

	Attributes []EntityAttribute `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

// TableName teh name of the User table
func (Entity) TableName() string {
	return "ea_entity"
}

// BeforeCreate adds UID to model
func (e *Entity) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uid.HumanUid()
	e.ID = uuid
	return nil
}

// GetAttribute the name of the User table
func (e *Entity) GetAttribute(db *gorm.DB, attributeKey string) *EntityAttribute {
	entityAttribute := &EntityAttribute{}

	result := db.First(&entityAttribute, "entity_id=? AND attribute_key=?", e.ID, attributeKey)

	if result.Error != nil {

		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil
		}

		log.Panic(result.Error)
	}

	return entityAttribute
}

// GetAttributeValue the value of the attribute or the default value if it does not exist
func (e *Entity) GetAttributeValue(db *gorm.DB, attributeKey string, defaultValue string) string {
	entityAttribute := e.GetAttribute(db, attributeKey)

	if entityAttribute == nil {
		return defaultValue
	}

	return entityAttribute.AttributeValue
}

// GetAttributeJSONValue the value of the attribute or the default value if it does not exist
func (e *Entity) GetAttributeJSONValue(db *gorm.DB, attributeKey string, defaultValue interface{}) interface{} {
	entityAttribute := e.GetAttribute(db, attributeKey)

	if entityAttribute == nil {
		return defaultValue
	}

	return entityAttribute.GetJSONValue()
}

// UpsertAttributes upserts the attributes
func (e *Entity) UpsertAttributes(db *gorm.DB, attributes map[string]string) bool {
	return EntityAttributesUpsert(db, e.ID, attributes)
}

// UpsertAttributesJSON upserts the attributes
func (e *Entity) UpsertAttributesJSON(db *gorm.DB, attributes map[string]interface{}) bool {
	return EntityAttributesUpsertJSON(db, e.ID, attributes)
}

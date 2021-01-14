package ea

import (
	"encoding/json"
	"time"

	"github.com/gouniverse/uid"
	"gorm.io/gorm"
)

// EntityAttribute type
type EntityAttribute struct {
	ID             string     `gorm:"type:varchar(40);column:id;primary_key;"`
	EntityID       string     `gorm:"type:varchar(40);column:entity_id;"`
	AttributeKey   string     `gorm:"type:varchar(255);column:attribute_key;DEFAULT NULL;"`
	AttributeValue string     `gorm:"type:longtext;column:attribute_value;"`
	CreatedAt      time.Time  `gorm:"type:datetime;column:created_at;DEFAULT NULL;"`
	UpdatedAt      time.Time  `gorm:"type:datetime;column:updated_at;DEFAULT NULL;"`
	DeletedAt      *time.Time `gorm:"type:datetime;column:deleted_at;DEFAULT NULL;"`
}

// TableName teh name of the User table
func (EntityAttribute) TableName() string {
	return "ea_attribute"
}

// BeforeCreate adds UID to model
func (e *EntityAttribute) BeforeCreate(tx *gorm.DB) (err error) {
	uuid := uid.HumanUid()
	e.ID = uuid
	return nil
}

// SetJSONValue serializes the values
func (e *EntityAttribute) SetJSONValue(value interface{}) bool {
	bytes, err := json.Marshal(value)

	if err != nil {
		return false
	}

	e.AttributeValue = string(bytes)

	return true
}

// GetJSONValue serializes the values
func (e *EntityAttribute) GetJSONValue() interface{} {
	var value interface{}
	err := json.Unmarshal([]byte(e.AttributeValue), &value)

	if err != nil {
		panic("JSOB error unmarshaliibg attribute" + err.Error())
	}

	return value
}

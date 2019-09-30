package models

import (
	"time"

	"futuagro.com/pkg/domain/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Variant represents a variant of product or service
type Variant struct {
	ID           primitive.ObjectID      `json:"_id" bson:"_id"`
	Name         string                  `json:"name" bson:"name"`
	LName        string                  `json:"lname" bson:"lname"`
	ItemID       primitive.ObjectID      `json:"itemId" bson:"itemId"`
	Item         *Item                   `json:"item,omitempty" bson:"item"`
	CreatedAt    time.Time               `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time               `json:"updatedAt" bson:"updatedAt"`
	RecordStatus *enums.EnumRecordStatus `json:"recordStatus" bson:"recordStatus"`
}

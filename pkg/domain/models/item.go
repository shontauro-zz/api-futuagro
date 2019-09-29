package models

import (
	"time"

	"futuagro.com/pkg/domain/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Item represents a product or service
type Item struct {
	ID           primitive.ObjectID      `json:"_id" bson:"_id"`
	Name         string                  `json:"name" bson:"name"`
	LName        string                  `json:"lname" bson:"lname"`
	CreatedAt    time.Time               `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time               `json:"updatedAt" bson:"updatedAt"`
	RecordStatus *enums.EnumRecordStatus `json:"recordStatus" bson:"recordStatus"`
}

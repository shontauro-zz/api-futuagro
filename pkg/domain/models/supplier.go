// Package models contains the entities of the domain business.
package models

import (
	"time"

	"futuagro.com/pkg/domain/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Supplier represent the data of a supplier
type Supplier struct {
	ID             primitive.ObjectID      `json:"_id" bson:"_id"`
	Name           string                  `json:"name" bson:"name"`
	Surname        string                  `json:"surname" bson:"surname"`
	DocumentType   string                  `json:"documentType" bson:"documentType"`
	DocumentNumber string                  `json:"documentNumber" bson:"documentNumber"`
	CityID         *primitive.ObjectID     `json:"cityId,omitempty" bson:"cityId"`
	City           *City                   `json:"city,omitempty" bson:"city"`
	Email          string                  `json:"email,omitempty" bson:"email"`
	AddressLine1   string                  `json:"addressLine1,omitempty" bson:"addressLine1"`
	PhoneNumber    string                  `json:"phoneNumber,omitempty" bson:"phoneNumber"`
	Crops          *[]Crop                 `json:"crops,omitempty" bson:"crops"`
	RecordStatus   *enums.EnumRecordStatus `json:"recordStatus" bson:"recordStatus"`
	CreatedAt      time.Time               `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time               `json:"updatedAt" bson:"updatedAt"`
}

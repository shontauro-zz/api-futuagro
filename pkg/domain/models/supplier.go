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
	CityID         string                  `json:"cityId" bson:"cityId"`
	Email          string                  `json:"email,omitempty" bson:"email"`
	AddressLine1   string                  `json:"addressLine1,omitempty" bson:"addressLine1"`
	PhoneNumber    string                  `json:"phoneNumber,omitempty" bson:"phoneNumber"`
	RecordStatus   *enums.EnumRecordStatus `json:"recordStatus" bson:"recordStatus"`
	CreatedAt      time.Time               `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time               `json:"updatedAt" bson:"updatedAt"`
}

//Crop represent the data of a crop
type Crop struct {
	ID           primitive.ObjectID `json:"_id" bson:"_id"`
	CityID       string             `json:"cityId" bson:"cityId"`
	PlantingDate time.Time          `json:"plantingDate" bson:"plantingDate"`
	HarvestDate  time.Time          `json:"harvestDate" bson:"harvestDate"`
	VariantID    *Variant           `json:"variantId" bson:"variantId"`
	CreatedAt    time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time          `json:"updatedAt" bson:"updatedAt"`
}

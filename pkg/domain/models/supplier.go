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
	Crops          []Crop                  `json:"crops,omitempty" bson:"crops"`
	RecordStatus   *enums.EnumRecordStatus `json:"recordStatus" bson:"recordStatus"`
	CreatedAt      time.Time               `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time               `json:"updatedAt" bson:"updatedAt"`
}

//Crop represent the data of a crop
type Crop struct {
	ID             primitive.ObjectID  `json:"_id,omitempty" bson:"_id"`
	CountryStateID string              `json:"countryStateId,omitempty" bson:"countryStateId"`
	CityID         string              `json:"cityId,omitempty" bson:"cityId"`
	PlantingDate   time.Time           `json:"plantingDate,omitempty" bson:"plantingDate"`
	HarvestDate    time.Time           `json:"harvestDate,omitempty" bson:"harvestDate"`
	VariantID      *primitive.ObjectID `json:"variantId,omitempty" bson:"variantId"`
	CreatedAt      time.Time           `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt      time.Time           `json:"updatedAt,omitempty" bson:"updatedAt"`
	Variant        *Variant            `json:"variant,omitempty" bson:"variant"`
}

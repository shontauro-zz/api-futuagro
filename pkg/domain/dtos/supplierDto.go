package dtos

import (
	"time"

	"futuagro.com/pkg/domain/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SupplierDto represents a DTO for a supplier document
type SupplierDto struct {
	Name           string                  `json:"name" bson:"name"`
	Surname        string                  `json:"surname" bson:"surname"`
	DocumentType   string                  `json:"documentType" bson:"documentType"`
	DocumentNumber string                  `json:"documentNumber" bson:"documentNumber"`
	CityID         string                  `json:"cityId" bson:"cityId"`
	Email          string                  `json:"email,omitempty" bson:"email"`
	AddressLine1   string                  `json:"addressLine1,omitempty" bson:"addressLine1"`
	PhoneNumber    string                  `json:"phoneNumber,omitempty" bson:"phoneNumber"`
	RecordStatus   *enums.EnumRecordStatus `json:"recordStatus" bson:"recordStatus"`
}

// CropDto represents a DTO for a crop sub-document
type CropDto struct {
	CountryStateID string              `json:"countryStateId" bson:"countryStateId"`
	CityID         string              `json:"cityId" bson:"cityId"`
	PlantingDate   time.Time           `json:"plantingDate" bson:"plantingDate"`
	HarvestDate    time.Time           `json:"harvestDate" bson:"harvestDate"`
	VariantID      *primitive.ObjectID `json:"variantId" bson:"variantId"`
}

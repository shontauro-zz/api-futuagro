package dtos

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CropDto represents a DTO for a crop sub-document
type CropDto struct {
	CityID       primitive.ObjectID  `json:"cityId" bson:"cityId"`
	PlantingDate time.Time           `json:"plantingDate" bson:"plantingDate"`
	HarvestDate  time.Time           `json:"harvestDate" bson:"harvestDate"`
	VariantID    *primitive.ObjectID `json:"variantId" bson:"variantId"`
	SupplierID   *primitive.ObjectID `json:"supplierId" bson:"supplierId"`
}

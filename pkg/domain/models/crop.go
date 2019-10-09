package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Crop represent the data of a crop
type Crop struct {
	ID           primitive.ObjectID  `json:"_id,omitempty" bson:"_id"`
	CityID       *primitive.ObjectID `json:"cityId,omitempty" bson:"cityId"`
	City         *City               `json:"city,omitempty" bson:"city"`
	PlantingDate time.Time           `json:"plantingDate" bson:"plantingDate"`
	HarvestDate  time.Time           `json:"harvestDate" bson:"harvestDate"`
	VariantID    *primitive.ObjectID `json:"variantId,omitempty" bson:"variantId"`
	Variant      *Variant            `json:"variant,omitempty" bson:"variant"`
	SupplierID   *primitive.ObjectID `json:"supplierId,omitempty" bson:"supplierId"`
	Supplier     *User               `json:"supplier,omitempty" bson:"supplier"`
	CreatedAt    time.Time           `json:"createdAt" bson:"createdAt"`
	UpdatedAt    time.Time           `json:"updatedAt" bson:"updatedAt"`
}

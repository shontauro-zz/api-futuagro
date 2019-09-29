// Package models contains the entities of the domain business.
package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Supplier represent the data of a provider
type Supplier struct {
	ID             primitive.ObjectID `json:"id" bson:"_id"`
	Name           string             `json:"name" bson:"name"`
	DocumentType   string             `json:"documentType" bson:"documentType"`
	DocumentNumber string             `json:"documentNumber" bson:"documentNumber"`
	City           string             `json:"city" bson:"city"`
	Email          string             `json:"email" bson:"email"`
	AddressLine1   string             `json:"addressLine1" bson:"address"`
	PhoneNumber    string             `json:"phoneNumber" bson:"phoneNumber"`
	Products       string             `json:"products" bson:"products"`
	RecordStatus   string             `json:"recordStatus"`
	CreatedAt      time.Time          `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt" bson:"updatedAt"`
}

package dtos

import (
	"time"

	"futuagro.com/pkg/domain/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// UserDto represents a DTO for an user document
type UserDto struct {
	ID             primitive.ObjectID      `json:"_id" bson:"_id"`
	Name           string                  `json:"name" bson:"name"`
	Surname        string                  `json:"surname" bson:"surname"`
	DocumentType   string                  `json:"documentType" bson:"documentType"`
	DocumentNumber string                  `json:"documentNumber" bson:"documentNumber"`
	CityID         *primitive.ObjectID     `json:"cityId,omitempty" bson:"cityId"`
	Email          string                  `json:"email,omitempty" bson:"email"`
	Password       string                  `json:"password,omitempty" bson:"password"`
	AddressLine1   string                  `json:"addressLine1,omitempty" bson:"addressLine1"`
	PhoneNumber    string                  `json:"phoneNumber,omitempty" bson:"phoneNumber"`
	RecordStatus   *enums.EnumRecordStatus `json:"recordStatus" bson:"recordStatus"`
	CreatedAt      time.Time               `json:"createdAt" bson:"createdAt"`
	UpdatedAt      time.Time               `json:"updatedAt" bson:"updatedAt"`
}

package dtos

import (
	"futuagro.com/pkg/domain/enums"
)

// SupplierDto represents a DTO for a supplier object
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

package models

import (
	"futuagro.com/pkg/domain/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//Country represents the data of a country
type Country struct {
	ID           primitive.ObjectID      `json:"id" bson:"_id"`
	CountryName  string                  `json:"countryName" bson:"countryName"`
	CountryCode  string                  `json:"countryCode" bson:"countryCode"`
	States       []CountryState          `json:"states" bson:"states"`
	RecordStatus *enums.EnumRecordStatus `json:"recordStatus" bson:"recordStatus"`
}

// CountryState represent the data of country state
type CountryState struct {
	ID           primitive.ObjectID      `json:"id" bson:"_id"`
	StateName    string                  `json:"stateName" bson:"stateName"`
	RecordStatus *enums.EnumRecordStatus `json:"recordStatus" bson:"recordStatus"`
}

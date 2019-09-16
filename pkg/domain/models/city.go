package models

import (
	"futuagro.com/pkg/domain/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// City represent the data of a city
type City struct {
	ID           primitive.ObjectID      `json:"id" bson:"_id"`
	CityName     string                  `json:"cityName"`
	CountryState CountryState            `json:"countryState"`
	RecordStatus *enums.EnumRecordStatus `json:"recordStatus"`
}

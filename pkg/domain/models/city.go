package models

import (
	"futuagro.com/pkg/domain/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// City represent the data of a city
type City struct {
	ID             primitive.ObjectID      `json:"_id" bson:"_id"`
	CityName       string                  `json:"cityName" bson:"cityName"`
	CountryStateID primitive.ObjectID      `json:"countryStateId" bson:"countryStateId"`
	RecordStatus   *enums.EnumRecordStatus `json:"recordStatus" bson:"recordStatus"`
}

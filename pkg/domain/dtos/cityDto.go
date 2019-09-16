package dtos

import "futuagro.com/pkg/domain/enums"

// CityDto represents a DTO for a city object
type CityDto struct {
	CityName       string                  `json:"cityName"`
	CountryStateID string                  `json:"countryStateId"`
	RecordStatus   *enums.EnumRecordStatus `json:"recordStatus"`
}

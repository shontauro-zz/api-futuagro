package dtos

import "futuagro.com/pkg/domain/enums"

// CityDto represents a DTO for a city object
type CityDto struct {
	CityName     string                  `json:"cityName"`
	RecordStatus *enums.EnumRecordStatus `json:"recordStatus"`
}

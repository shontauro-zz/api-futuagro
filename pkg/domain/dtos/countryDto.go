package dtos

import "futuagro.com/pkg/domain/enums"

// CountryDto represents a DTO for a country object
type CountryDto struct {
	CountryName  string                  `json:"countryName"`
	CountryCode  string                  `json:"countryCode"`
	RecordStatus *enums.EnumRecordStatus `json:"recordStatus"`
}

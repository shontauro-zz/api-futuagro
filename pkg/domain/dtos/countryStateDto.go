package dtos

import (
	"futuagro.com/pkg/domain/enums"
)

// CountryStateDto epresents a DTO for a country state object
type CountryStateDto struct {
	StateName    string                  `json:"stateName"`
	RecordStatus *enums.EnumRecordStatus `json:"recordStatus"`
}

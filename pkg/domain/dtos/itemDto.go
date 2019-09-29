package dtos

import (
	"futuagro.com/pkg/domain/enums"
)

// ItemDto represents a DTO for an Item object
type ItemDto struct {
	Name         string                  `json:"name" bson:"name"`
	RecordStatus *enums.EnumRecordStatus `json:"recordStatus" bson:"recordStatus"`
}

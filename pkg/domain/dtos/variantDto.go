package dtos

import (
	"futuagro.com/pkg/domain/enums"
)

// VariantDto represents a variant of product or service
type VariantDto struct {
	Name         string                  `json:"name" bson:"name"`
	RecordStatus *enums.EnumRecordStatus `json:"recordStatus" bson:"recordStatus"`
}

package enums

import (
	"bytes"
	"encoding/json"

	"github.com/pkg/errors"
)

// EnumRecordStatus represents the status of a database record, moving through active, inactive
type EnumRecordStatus string

const (
	// Active represents a record status enable to be query in the database
	Active EnumRecordStatus = "active"
	// Inactive represents a record status disable, it is a soft delete
	Inactive EnumRecordStatus = "inactive"
)

func (s EnumRecordStatus) String() string {
	return toString[s]
}

var toString = map[EnumRecordStatus]string{
	Active:   "active",
	Inactive: "inactive",
}

var toID = map[string]EnumRecordStatus{
	"active":   Active,
	"inactive": Inactive,
}

// MarshalJSON marshals the enum as a quoted json string
func (s *EnumRecordStatus) MarshalJSON() ([]byte, error) {
	buffer := bytes.NewBufferString(`"`)
	value, ok := toString[*s]
	if !ok {
		return nil, errors.New("Invalid RecordStatus value")
	}
	buffer.WriteString(value)
	buffer.WriteString(`"`)
	return buffer.Bytes(), nil
}

// UnmarshalJSON unmashals a quoted json string to the enum value
func (s *EnumRecordStatus) UnmarshalJSON(b []byte) error {
	var j string
	err := json.Unmarshal(b, &j)
	if err != nil {
		return err
	}
	value, ok := toID[j]
	if !ok {
		return errors.New("Invalid RecordStatus value")
	}
	*s = value
	return nil
}

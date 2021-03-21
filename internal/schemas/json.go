package schemas

import (
	"bytes"
	"database/sql/driver"
	"errors"
)

//JSON custome type json
type JSON []byte

//Value implement gorm interface
func (j JSON) Value() (driver.Value, error) {
	if j.IsNull() {
		return nil, nil
	}
	return string(j), nil
}

//Scan implement gorm interface
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}
	s, ok := value.([]byte)
	if !ok {
		return errors.New("Invalid Scan Source")
	}
	*j = append((*j)[0:0], s...)
	return nil
}

//MarshalJSON custome json marshal
func (j JSON) MarshalJSON() ([]byte, error) {
	if j == nil {
		return []byte("null"), nil
	}
	return j, nil
}

//UnmarshalJSON  custom json unmarshal
func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return errors.New("null point exception")
	}
	*j = append((*j)[0:0], data...)
	return nil
}

//IsNull gorm check is null
func (j JSON) IsNull() bool {
	return len(j) == 0 || string(j) == "null"
}

//Equals gorm check equal
func (j JSON) Equals(j1 JSON) bool {
	return bytes.Equal([]byte(j), []byte(j1))
}

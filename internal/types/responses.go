package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
)

type ResponsesMap map[uuid.UUID]string

// Метод для сериализации
func (r ResponsesMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[uuid.UUID]string(r))
}

// Метод для десериализации
func (r *ResponsesMap) UnmarshalJSON(data []byte) error {
	var aux map[uuid.UUID]string
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	*r = ResponsesMap(aux)
	return nil
}

// Метод для десериализации из базы данных
func (r *ResponsesMap) Scan(value interface{}) error {
	if value == nil {
		*r = ResponsesMap{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("failed to scan ResponsesMap")
	}

	return r.UnmarshalJSON(bytes)
}

// Метод для сериализации в базу данных
func (r ResponsesMap) Value() (driver.Value, error) {
	return json.Marshal(r)
}

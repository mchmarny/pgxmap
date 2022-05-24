package pgxmap

import (
	"database/sql/driver"
	"encoding/json"

	"github.com/pkg/errors"
)

type State[T any] struct {
	Data T `json:"data"`
}

// Value implements the driver.Valuer interface.
func (s State[T]) Value() (driver.Value, error) {
	return json.Marshal(s.Data)
}

// Scan implements the sql.Scanner interface.
func (s *State[T]) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	var data T
	if err := json.Unmarshal(b, &data); err != nil {
		return errors.Wrap(err, "failed to unmarshal data")
	}

	s.Data = data
	return nil
}

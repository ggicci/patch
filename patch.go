package patch

import (
	"encoding/json"
)

// Field is a wrapper which can tell if a field was unmarshalled from the data provided.
// When `Field.Valid` is true, which means `Field.Value` is populated from decoding the raw data.
// Otherwise, no data was provided, i.e. field missing.
type Field[T any] struct {
	Value T
	Valid bool
}

func (f Field[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(f.Value)
}

func (f *Field[T]) UnmarshalJSON(data []byte) error {
	err := json.Unmarshal(data, &f.Value)
	if err == nil {
		f.Valid = true
	}
	return err
}

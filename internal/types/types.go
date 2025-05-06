package types

import (
	"encoding/json"
)

type StrText struct {
	Description string `json:"description"`
	Value       string `json:"value"`
}

func (s *StrText) UnmarshalJSON(data []byte) error {
	var val string
	if err := json.Unmarshal(data, &val); err != nil {
		return err
	}

	s.Value = val
	return nil
}

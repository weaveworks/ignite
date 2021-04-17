package v1alpha4

import (
	"encoding/json"
)

// In this package custom marshal/unmarshal functions are registered

func (s *SSH) MarshalJSON() ([]byte, error) {
	if len(s.PublicKey) != 0 {
		return json.Marshal(s.PublicKey)
	}

	if s.Generate {
		return json.Marshal(true)
	}

	return []byte("{}"), nil
}

func (s *SSH) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err == nil {
		if str == "true" {
			*s = SSH{
				Generate: true,
			}
		} else {
			*s = SSH{
				PublicKey: str,
			}
		}

		return nil
	}

	var boolVar bool
	if err := json.Unmarshal(b, &boolVar); err == nil {
		if boolVar {
			*s = SSH{
				Generate: true,
			}

			return nil
		}
	}

	// The user did not specify this field, just return
	return nil
}

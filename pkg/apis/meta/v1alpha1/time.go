package v1alpha1

import (
	"encoding/json"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
)

type Time struct {
	*metav1.Time
}

var _ fmt.Stringer = Time{}

var _ json.Marshaler = &Time{}
var _ json.Unmarshaler = &Time{}

// The default string for Time is a human readable difference between the Time and the current time
func (t Time) String() string {
	if t.Time == nil {
		return "<unknown>"
	}

	return fmt.Sprintf("%s ago", duration.HumanDuration(time.Now().Sub(t.Time.Time)))
}

// Timestamp returns the current UTC time
func Timestamp() Time {
	return Time{
		&metav1.Time{
			Time: time.Now().UTC(),
		},
	}
}

func newTime() Time {
	return Time{
		&metav1.Time{},
	}
}

func (t Time) MarshalJSON() ([]byte, error) {
	ti := t
	if t.Time == nil {
		// If the embedded metav1.Time is nil,
		// use the current time when marshaling
		ti = Timestamp()
	}

	b, err := ti.MarshalText()
	if err != nil {
		return nil, err
	}

	return json.Marshal(string(b))
}

func (t *Time) UnmarshalJSON(b []byte) error {
	var str string

	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	// If the Time is an empty string,
	// it should not be loaded
	if len(str) > 0 {
		*t = newTime()
		return t.UnmarshalText([]byte(str))
	}

	return nil
}

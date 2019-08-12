package v1alpha1

import (
	"encoding/json"
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
)

type Time struct {
	metav1.Time
}

var _ fmt.Stringer = Time{}

var _ json.Marshaler = &Time{}

// The default string for Time is a human readable difference between the Time and the current time
func (t Time) String() string {
	if t.IsZero() {
		return "<unknown>"
	}

	return duration.HumanDuration(time.Now().Sub(t.Time.Time))
}

// Timestamp returns the current UTC time
func Timestamp() Time {
	return Time{
		metav1.Time{
			Time: time.Now().UTC(),
		},
	}
}

func (t Time) MarshalJSON() (b []byte, err error) {
	if t.Time.IsZero() {
		// If the embedded metav1.Time is zero,
		// use the current time when marshaling
		b, err = Timestamp().Time.MarshalJSON()
	} else {
		b, err = t.Time.MarshalJSON()
	}

	return
}

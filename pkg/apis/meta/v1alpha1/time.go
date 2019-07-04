package v1alpha1

import (
	"fmt"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/duration"
)

type Time struct {
	metav1.Time
}

var _ fmt.Stringer = &Time{}

// The default string for Time is a human readable difference between the Time and the current time
func (t *Time) String() string {
	return fmt.Sprintf("%s ago", duration.HumanDuration(time.Now().Sub(t.Time.Time)))
}

// Timestamp returns the current UTC time
func Timestamp() Time {
	return Time{
		metav1.Time{
			Time: time.Now().UTC(),
		},
	}
}

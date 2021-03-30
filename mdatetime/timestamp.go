package mdatetime

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const TimestampMillisLayout = "2006-01-02T15:04:05Z"

type Mtimestamp time.Time

func (m *Mtimestamp) UnmarshalJSON(b []byte) error {
	var (
		tt  time.Time
		err error
	)
	if InLocalTimeZone {
		tt, err = time.ParseInLocation(TimestampMillisLayout, string(b), time.Local)
	} else {
		tt, err = time.Parse(TimestampMillisLayout, string(b))
	}
	*m = Mtimestamp(tt)
	return err
}

func (ct Mtimestamp) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

func (ct *Mtimestamp) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t.Format(TimestampMillisLayout))
}

// GetBSON customizes the bson serialization for this type
func (t Mtimestamp) GetBSON() (interface{}, error) {
	return t.String(), nil
}

// SetBSON customizes the bson serialization for this type
func (t *Mtimestamp) SetBSON(raw bson.Raw) error {
	return t.UnmarshalJSON(raw)
}

func (t Mtimestamp) GetTime() time.Time {
	return time.Time(t)
}

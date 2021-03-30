package mdatetime

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
)

const MdateLayout = "2006-01-02"

// const MdateLayout = "2006-01-02T15:04:05"

var (
	InLocalTimeZone = true
)

type Mdate time.Time

func (m *Mdate) UnmarshalJSON(b []byte) error {
	var (
		tt  time.Time
		err error
	)
	if InLocalTimeZone {
		tt, err = time.ParseInLocation(MdateLayout, string(b), time.Local)
	} else {
		tt, err = time.Parse(MdateLayout, string(b))
	}
	*m = Mdate(tt)
	return err
}

func (ct Mdate) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

func (ct *Mdate) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t.Format(MdateLayout))
}

// GetBSON customizes the bson serialization for this type
func (t Mdate) GetBSON() (interface{}, error) {
	return t.String(), nil
}

// SetBSON customizes the bson serialization for this type
func (t *Mdate) SetBSON(raw bson.Raw) error {
	return t.UnmarshalJSON(raw)
}

func (t Mdate) GetTime() time.Time {
	return time.Time(t)
}

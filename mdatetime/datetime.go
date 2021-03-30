package mdatetime

import (
	"fmt"
	"time"
)

const MdatetimeLayout = "2006-01-02T15:04:05Z"

type Mdatetime time.Time

func (m *Mdatetime) UnmarshalJSON(b []byte) error {
	var (
		tt  time.Time
		err error
	)
	if InLocalTimeZone {
		tt, err = time.ParseInLocation(MdatetimeLayout, string(b), time.Local)
	} else {
		tt, err = time.Parse(MdatetimeLayout, string(b))
	}
	*m = Mdatetime(tt)
	return err
}

func (ct Mdatetime) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

func (ct *Mdatetime) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t.Format(MdatetimeLayout))
}

func (t Mdatetime) GetTime() time.Time {
	return time.Time(t)
}

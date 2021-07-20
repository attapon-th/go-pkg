package dt

import (
	"fmt"
	"strings"
	"time"
)

type Datetime time.Time

func (m *Datetime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	tt, err := time.ParseInLocation(DATETIME_LAYOUT, s, time.Local)
	*m = Datetime(tt)
	return err
}

func (ct Datetime) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

func (ct *Datetime) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t.Format(DATETIME_LAYOUT))
}

func (t Datetime) GetTime() time.Time {
	return time.Time(t)
}

func (t *Datetime) FromTime(ti time.Time) {
	*t = Datetime(ti)
}

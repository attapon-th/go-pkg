package dt

import (
	"fmt"
	"strings"
	"time"
)

type Timestamp time.Time

func (m *Timestamp) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	tt, err := time.ParseInLocation(TIMESTAMP_LAYOUT, s, time.Local)
	*m = Timestamp(tt)
	return err
}

func (ct Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

func (ct *Timestamp) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t.Format(TIMESTAMP_LAYOUT))
}

func (t Timestamp) GetTime() time.Time {
	return time.Time(t)
}
func (t *Timestamp) FromTime(ti time.Time) {
	*t = Timestamp(ti)
}

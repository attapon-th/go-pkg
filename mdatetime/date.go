package dt

import (
	"fmt"
	"strings"
	"time"
)

// const MdateLayout = "2006-01-02T15:04:05"

type Date time.Time

func (m *Date) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), `"`)
	tt, err := time.ParseInLocation(DATE_LAYOUT, s, time.Local)
	m.FromTime(tt)
	return err
}

func (ct Date) MarshalJSON() ([]byte, error) {
	return []byte(ct.String()), nil
}

func (ct *Date) String() string {
	t := time.Time(*ct)
	return fmt.Sprintf("%q", t.Format(DATE_LAYOUT))
}

func (t Date) GetTime() time.Time {
	return time.Time(t)
}

func (t *Date) FromTime(ti time.Time) {
	*t = Date(ti)
}

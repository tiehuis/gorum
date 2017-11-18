package model

import (
	"database/sql/driver"
	"time"

	"gopkg.in/guregu/null.v3"
)

// sql.NullInt64 doesn't work here as we would expect
type NullInt null.Int

type NullTime struct {
	Time  time.Time
	Valid bool
}

func (nt *NullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

func (nt NullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}

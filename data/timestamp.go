package data

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type Timestamp struct {
	time.Time
}

func NowTimestamp() Timestamp {
	return Timestamp{Time: time.Now()}
}

// Value implements the driver.Valuer interface.
func (ct Timestamp) Value() (driver.Value, error) {
	str := ct.Time.Format(time.RFC3339)
	return str, nil
}

// Scan implements the sql.Scanner interface.
func (ct *Timestamp) Scan(value interface{}) error {
	if value == nil {
		ct.Time = time.Time{}
		return nil
	}

	strValue, ok := value.(string)
	if !ok {
		return fmt.Errorf("expected string value for CustomTime, got %T", value)
	}

	parsedTime, err := time.Parse(time.RFC3339, strValue)
	if err != nil {
		return err
	}

	ct.Time = parsedTime
	return nil
}

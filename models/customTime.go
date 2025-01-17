package models

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type CustomTime struct {
	time.Time
}

const customTimeLayout = "2006-01-02"

func (ct *CustomTime) UnmarshalJSON(b []byte) error {

	str := string(b)
	str = str[1 : len(str)-1]

	parsedTime, err := time.Parse(customTimeLayout, str)

	if err != nil {
		return err
	}

	ct.Time = parsedTime

	return nil
}

func (ct CustomTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + ct.Time.Format(customTimeLayout) + `"`), nil
}

func (ct CustomTime) Value() (driver.Value, error) {
	return ct.Time, nil
}

func (ct *CustomTime) Scan(value interface{}) error {
	if value == nil {
		*ct = CustomTime{}
		return nil
	}

	switch v := value.(type) {
	case time.Time:
		ct.Time = v
	default:
		return fmt.Errorf("cannot convert %v to CustomTime", value)
	}

	return nil
}

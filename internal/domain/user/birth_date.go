package user

import (
	"database/sql/driver"
	"fmt"
	"time"
)

type BirthDate struct {
	time.Time
}

func (b *BirthDate) Scan(value interface{}) error {
	t, ok := value.(time.Time)
	if !ok {
		return fmt.Errorf("failed to scan BirthDate: %v", value)
	}
	b.Time = t
	return nil
}

func (b BirthDate) Value() (driver.Value, error) {
	return b.Time, nil
}

func NewBirthDate(year, month, day int32) (BirthDate, error) {
	birthDate := time.Date(int(year), time.Month(month), int(day), 0, 0, 0, 0, time.Local)

	if birthDate.Year() != int(year) || birthDate.Month() != time.Month(month) || birthDate.Day() != int(day) {
		return BirthDate{}, fmt.Errorf("invalid birth date: %d-%02d-%02d", year, month, day)
	}

	now := time.Now()

	if birthDate.After(now) {
		return BirthDate{}, fmt.Errorf("birth date cannot be in the future")
	}

	if birthDate.Before(now.AddDate(-150, 0, 0)) {
		return BirthDate{}, fmt.Errorf("birth date is too far in the past")
	}

	return BirthDate{birthDate}, nil
}

package user

import (
	"database/sql/driver"
	"testing"
	"time"
)

func TestNewBirthDate(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		year    int32
		month   int32
		day     int32
	}{
		{"success new birth date", true, 2000, 1, 1},
		{"failure birth date is in the future", false, 2500, 1, 1},
		{"failure birth date is too far in the past", false, 1800, 1, 1},
		{"failure invalid birth date", false, 2000, 2, 30},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := NewBirthDate(tt.year, tt.month, tt.day)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestBirthScan(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		success bool
		value   any
	}{
		{"success", true, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
		{"failure invalid type", false, "invalid type"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			birthDate, err := NewBirthDate(2000, 1, 1)
			if err != nil {
				t.Errorf("failed to new birth date: %v", err)
			}

			err = birthDate.Scan(tt.value)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

func TestValue(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		success  bool
		expected driver.Value
	}{
		{"success", true, time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			birthDate, err := NewBirthDate(2000, 1, 1)
			if err != nil {
				t.Errorf("failed to new birth date: %v", err)
			}

			value, err := birthDate.Value()
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}

			if tt.success && value != tt.expected {
				t.Errorf("expected %v, but got %v", tt.expected, value)
			}
		})
	}
}

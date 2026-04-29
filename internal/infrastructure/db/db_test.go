package db

import "testing"

func TestInit(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		success  bool
		host     string
		user     string
		password string
		dbName   string
		port     string
		sslMode  string
	}{
		{"failure unreachable port", false, "127.0.0.1", "user", "password", "dbname", "1", "disable"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_, err := Init(tt.host, tt.user, tt.password, tt.dbName, tt.port, tt.sslMode)
			if tt.success && err != nil {
				t.Errorf("expected no error, but got %v", err)
			}
			if !tt.success && err == nil {
				t.Errorf("expected error, but got nil")
			}
		})
	}
}

package models

import (
	"testing"

	"github.com/hunttraitor/class-notifier/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	//-short flag skips test
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name   string
		userID int
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: 1,
			want:   true,
		},
		{
			name:   "Zero ID",
			userID: 0,
			want:   false,
		},
		{
			name:   "Non-existant ID",
			userID: 2,
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{db}
			exists, err := m.Exists(tt.userID)
			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}

func TestUserModelInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skpping integration tests")
	}

	tests := []struct {
		name      string
		userName  string
		email     string
		password  string
		wantError error
	}{
		{
			name:      "Valid Insert",
			userName:  "Test User",
			email:     "test@gmail.com",
			password:  "pa$$word",
			wantError: nil,
		},
		{
			name:      "Duplicate Insert",
			userName:  "Hunter",
			email:     "hunter@gmail.com",
			password:  "pa$$word",
			wantError: ErrDuplicateEmail,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{db}
			err := m.Insert(tt.userName, tt.email, tt.password)
			assert.Equal(t, err, tt.wantError)
		})
	}
}

func TestUserModelAuthenticate(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration tests")
	}

	tests := []struct {
		name       string
		email      string
		password   string
		wantUserID int
		wantError  error
	}{
		{
			name:       "Valid Authentication",
			email:      "testuserauth@gmail.com",
			password:   "pa$$word",
			wantUserID: 2,
			wantError:  nil,
		},
		{
			name:       "Invalid email",
			email:      "invalidEmail@gmail.com",
			password:   "pa$$word",
			wantUserID: 0,
			wantError:  ErrInvalidCredentials,
		},
		{
			name:       "Invalid Password",
			email:      "testuserauth@gmail.com",
			password:   "wrongPassword",
			wantUserID: 0,
			wantError:  ErrInvalidCredentials,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := UserModel{db}
			m.Insert("test", "testuserauth@gmail.com", "pa$$word")
			userID, err := m.Authenticate(tt.email, tt.password)
			assert.Equal(t, userID, tt.wantUserID)
			assert.Equal(t, err, tt.wantError)
		})
	}
}

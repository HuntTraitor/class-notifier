package models

import (
	"strings"
	"testing"

	"github.com/hunttraitor/class-notifier/internal/assert"
)

func TestClassInsert(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name        string
		classID     int
		className   string
		link        string
		professor   string
		wantClassID int
		wantError   error
	}{
		{
			name:        "Valid Insert",
			classID:     3,
			className:   "Valid Insert Class",
			link:        "Validinsert.com",
			professor:   "Valid Insert Professor",
			wantClassID: 3,
			wantError:   nil,
		},
		{
			name:        "Duplicate Class Insert",
			classID:     1,
			className:   "Duplicate Class",
			link:        "DuplicateClass.com",
			professor:   "Ducplicate Class Professor",
			wantClassID: 0,
			wantError:   ErrDuplicateClass,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := ClassModel{db}

			resultId, err := m.Insert(tt.classID, tt.name, tt.link, tt.professor)
			assert.Equal(t, resultId, tt.wantClassID)
			assert.Equal(t, err, tt.wantError)
		})
	}
}

func TestClassGet(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	tests := []struct {
		name      string
		id        int
		wantClass Class
		wantError error
	}{
		{
			name:      "Valid Get",
			id:        1,
			wantClass: Class{ClassID: 1, Name: "Test Class", Link: "testclass.com", Professor: "Professor Test"},
			wantError: nil,
		},
		{
			name:      "Class Doesnt Exist",
			id:        3,
			wantError: ErrNoRecord,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := ClassModel{db}

			result, err := m.Get(tt.id)
			result.Name = strings.TrimSpace(result.Name)
			result.Link = strings.TrimSpace(result.Link)
			result.Professor = strings.TrimSpace(result.Professor)
			assert.Equal(t, result, tt.wantClass)
			assert.Equal(t, err, tt.wantError)
		})
	}
}

func TestClassList(t *testing.T) {
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	//chat this can't be real
	var (
		ValidClass = func(t *testing.T) []Class {
			var ClassList []Class
			t.Run("Initialization", func(t *testing.T) {
				db := newTestDB(t)
				m := ClassModel{db}
				for i := 1; i <= 2; i++ {
					class, err := m.Get(i)
					if err != nil {
						t.Fatal(err)
					}
					ClassList = append(ClassList, class)
				}
			})
			return ClassList
		}(t)
	)

	tests := []struct {
		name        string
		wantClasses []Class
		wantError   error
	}{
		{
			name:        "Valid ClassList",
			wantClasses: ValidClass,
			wantError:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			db := newTestDB(t)
			m := ClassModel{db}

			ClassList, err := m.Classlist()
			assert.SliceEqual(t, ClassList, tt.wantClasses)
			assert.Equal(t, err, tt.wantError)
		})
	}

}

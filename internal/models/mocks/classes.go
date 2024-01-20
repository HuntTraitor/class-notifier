package mocks

import (
	"github.com/hunttraitor/class-notifier/internal/models"
)

var mockClass = models.Class {
	ClassID: 1,
	Name: "Mock class",
	Link: "www.example.com",
	Professor: "CoolGuys77",
}

type ClassModel struct{}

func (m *ClassModel) Insert(classid int, name string, link string, professor string) (int, error) {
	return 2, nil
}

func (m *ClassModel) Get(ClassID int) (models.Class, error) {
	switch ClassID {
	case 1:
		return mockClass, nil
	default:
		return models.Class{}, models.ErrNoRecord
	}
}

func (m *ClassModel) Classlist() ([]models.Class, error) {
	return []models.Class{mockClass}, nil
}
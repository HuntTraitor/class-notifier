package models

import (
	"database/sql"
)

type Class struct {
	Name string
	Link string
}

type ClassModel struct {
	DB *sql.DB
}

func (m *ClassModel) Insert(name string, link string) (int, error) {
	return 0, nil
}


func (m *ClassModel) Get(name string) (Class, error) {
	return Class{}, nil
}

func (m *ClassModel) List() ([]Class, error) {
	return nil, nil
}
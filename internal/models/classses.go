package models

import (
	"database/sql"
	"errors"
)

type Class struct {
	ClassID   int
	Name      string
	Link      string
	Professor string
}

type ClassModel struct {
	DB *sql.DB
}

func (m *ClassModel) Insert(classid int, name string, link string, professor string) (int, error) {
	stmt := `INSERT INTO classes VALUES ($1, $2, $3, $4) RETURNING classid`

	tx, err := m.DB.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	_, err = tx.Exec(stmt, classid, name, link, professor)
	if err != nil {
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return classid, nil
}

func (m *ClassModel) Get(id int) (Class, error) {

	stmt := `SELECT classid, name, link, professor FROM classes
	WHERE classid = $1`

	row := m.DB.QueryRow(stmt, id)
	var c Class

	err := row.Scan(&c.ClassID, &c.Link, &c.Name, &c.Professor)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Class{}, ErrNoRecord
		} else {
			return Class{}, err
		}
	}
	return c, nil
}

func (m *ClassModel) Classlist() ([]Class, error) {

	stmt := `SELECT classid, name, link, professor FROM classes
	ORDER BY name`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var classes []Class

	for rows.Next() {
		var c Class
		err = rows.Scan(&c.ClassID, &c.Name, &c.Link, &c.Professor)
		if err != nil {
			return nil, err
		}
		classes = append(classes, c)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return classes, nil
}

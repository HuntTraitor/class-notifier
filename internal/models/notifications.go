package models

import (
	"database/sql"
	"time"
)

type Notification struct {
	Classname string
	Email     string
	Expires   time.Time
}

type NotificationModel struct {
	DB *sql.DB
}

func (n *NotificationModel) Insert(email string, classid int, expires int) error {
	stmt := `INSERT INTO notifications VALUES ($1, $2, CURRENT_TIMESTAMP + $3 * INTERVAL '1 day')`

	tx, err := n.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(stmt, email, classid, expires)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (n *NotificationModel) Delete(classname string, email string) error {
	return nil
}

func (n *NotificationModel) List(email string) ([]Notification, error) {
	return nil, nil
}

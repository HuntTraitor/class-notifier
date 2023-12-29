package models

import (
	"database/sql"
	"time"
)

type Notification struct {
	Email     string
	Classname string
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

func (n *NotificationModel) Delete(email string, classid int) error {

	stmt := `DELETE FROM notifications
	WHERE (email, classid) = ($1, $2)`

	tx, err := n.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(stmt, email, classid)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (n *NotificationModel) List(email string) ([]Notification, error) {

	stmt := `SELECT n.email, c.name
	FROM classes c, notifications n
	WHERE c.id = n.classid
	AND n.email = $1`

	rows, err := n.DB.Query(stmt, email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notifications []Notification

	for rows.Next() {
		var n Notification
		err = rows.Scan(&n.Email, &n.Classname)
		if err != nil {
			return nil, err
		}
		notifications = append(notifications, n)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return notifications, nil
}

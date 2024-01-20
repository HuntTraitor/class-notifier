package models

import (
	"database/sql"
	"errors"
	"strings"

	"github.com/lib/pq"
)

type NotificationModelInterface interface {
	Insert(email string, classid int, expires int) error
	Delete(notification int) error
	NotificationList(email string) ([]Notification, error)
}

type Notification struct {
	NotificationID int
	Name           string
	Link           string
}

type NotificationModel struct {
	DB *sql.DB
}

func (n *NotificationModel) Insert(email string, classid int, expires int) error {
	stmt := `INSERT INTO notifications(email, classid, expires) VALUES ($1, $2, CURRENT_TIMESTAMP + $3 * INTERVAL '1 day')`

	tx, err := n.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(stmt, email, classid, expires)
	if err != nil {
		var pqError *pq.Error
		//checks if error is a unique constraint violation
		if errors.As(err, &pqError) {
			if pqError.Code == "23505" && strings.Contains(pqError.Message, "unique_notification") {
				return ErrDuplicateNotification
			}
		}
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (n *NotificationModel) Delete(notificationid int) error {

	stmt := `DELETE FROM notifications
	WHERE notificationid = $1`

	tx, err := n.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(stmt, notificationid)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

func (n *NotificationModel) NotificationList(email string) ([]Notification, error) {

	stmt := `SELECT n.notificationid, c.name, c.link
	FROM classes c, notifications n
	WHERE c.classid = n.classid
	AND n.email = $1`

	rows, err := n.DB.Query(stmt, email)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var notifications []Notification

	for rows.Next() {
		var n Notification
		err = rows.Scan(&n.NotificationID, &n.Name, &n.Link)
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

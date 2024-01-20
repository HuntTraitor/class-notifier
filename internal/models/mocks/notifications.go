package mocks

import (
	"github.com/hunttraitor/class-notifier/internal/models"
)

var mockNotification = models.Notification{
	NotificationID: 1,
	Name:           "Mock Class",
	Link:           "www.example.com",
}

type NotificationModel struct{}

func (m *NotificationModel) Insert(email string, classid int, expires int) error {
	return nil
}

func (m *NotificationModel) Delete(notificationid int) error {
	return nil
}

func (m *NotificationModel) NotificationList(email string) ([]models.Notification, error) {
	return []models.Notification{mockNotification}, nil
}

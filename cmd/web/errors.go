package main

var (
	DuplicateNotificationFlashError = "Class notification already exists!"
)

var flashErrors = map[string]bool{
	DuplicateNotificationFlashError: true,
}

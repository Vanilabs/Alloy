package notifications

type NotificationService interface {
	Send(payload *NotificationPayload) error
}

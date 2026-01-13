package notifications

type NotificationPayload struct {
	To      string
	Subject string
	Body    string
	HTML    string
}

package notifications

type Notification struct {
	Email *Email
}

func NewNotification(email *Email) *Notification {
	return &Notification{
		Email: email,
	}
}

package models

import "time"

const (
	MagicLinkTokenLength = 64

	MagicLinkTokenExpiry = 15 * time.Minute
)

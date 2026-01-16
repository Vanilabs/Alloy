package schemas

type RequestMagicLinkSchema struct {
	Email string `json:"email" validate:"required,email"`
}

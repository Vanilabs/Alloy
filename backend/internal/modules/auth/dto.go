package auth

type InviteUserRequest struct {
	Email string `json:"email"`
	Role  string `json:"role"`
}

type VerifyInvitationRequest struct {
	Token string `json:"token"`
	Email string `json:"email"`
}

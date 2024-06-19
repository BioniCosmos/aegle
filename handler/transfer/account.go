package transfer

type SignUpBody struct {
	Email    string `json:"email"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type SignInBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

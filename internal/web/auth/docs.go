package auth

type LoginRequest struct {
	Email    string `json:"email" example:"example@example.com"`
	Password string `json:"password" example:"passworD1_"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type LoginResponseDoc struct {
	RawResponse
	Payload struct {
		Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9"`
	}
}

type RegistrationRequest struct {
	Name     string `json:"name" example:"kirill"`
	Surname  string `json:"surname" example:"zagrebin"`
	Email    string `json:"email" example:"example@example.com"`
	Password string `json:"password" example:"passworD1_"`
}

type ConfirmEmailRequest struct {
	Email            string `json:"email" example:"example@example.com"`
	VerificationCode string `json:"verification_code" example:"1324"`
}

type ResendCodeRequest struct {
	Email string `json:"email" example:"example@example.com"`
}

type ResetPasswordRequest struct {
	Email string `json:"email" example:"example@example.com"`
}

type VerificationMessageData struct {
	Code int `json:"Code"`
}

type ResetPasswordConfirmRequest struct {
	Email       string `json:"email"`
	NewPassword string `json:"new_password"`
}

type UserResponse struct {
	Id         string  `json:"id" example:"8e3b780c-dfb5-4cd2-ae6b-d83f84a483f9"`
	Name       string  `json:"name" example:"Keril"`
	Surname    string  `json:"surname" example:"Zagrebin"`
	SecondName *string `json:"second_name" example:"Maksimovich"`
	Email      string  `json:"email" example:"example@example.com"`
	TelegramId string  `json:"telegram_id" example:"311441242"`
}

type UserResponseDoc struct {
	RawResponse
	Payload struct {
		UserResponse
	} `json:"payload"`
}

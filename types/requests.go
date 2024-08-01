package types

type LoginUserRequest struct {
	UserEmail    string `json:"user_email" binding:"required,email"`
	UserPassword string `json:"user_password" binding:"required"`
}

type SignupUserRequest struct {
	UserEmail    string `json:"user_email" binding:"required,email"`
	UserPassword string `json:"user_password" binding:"required"`
	UserName     string `json:"user_name" binding:"required"`
}

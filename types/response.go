package types

type BaseHttpResponse struct {
	Status  string      `json:"status"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func GenerateResponse(data interface{}, message string) *BaseHttpResponse {
	return &BaseHttpResponse{
		Status:  "success",
		Data:    data,
		Message: message,
	}
}

type UserDetails struct {
	Id       int64  `json:"user_id"`
	Email    string `json:"user_email"`
	Username string `json:"user_name"`
}

type LoginResponse struct {
	AccessToken string      `json:"access_token"`
	UserDetails UserDetails `json:"user_details"`
}

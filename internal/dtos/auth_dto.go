package dtos

type LoginInput struct {
	Email string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"` 
}

type LoginResponse struct {
	Token string `json:"token"`
	User UserResponse `json:"user"`
}

type UserResponse struct {
	ID uint `json:"id"`	//uint:符号なし整数，DBのIDにマイナスは不要であるため
	Name string `json:"name"`
	Email string `json:"email"`
	IconURL string `json:"icon_url"`
}
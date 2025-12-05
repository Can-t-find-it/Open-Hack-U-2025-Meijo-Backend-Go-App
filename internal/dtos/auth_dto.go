package dtos

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	// User  UserResponse `json:"user"`
}

type UserResponse struct {
	ID      string   `json:"id"` 
	Name    string `json:"name"`
	Email   string `json:"email"`
	IconURL string `json:"icon_url"`
}

type UserSignup struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"` //６文字以上
}

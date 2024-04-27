package form

type LoginInput struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type RegisterInput struct {
	Name     string `json:"name" validate:"required"`
	Surname  string `json:"surname" validate:"required"`
	Username string `json:"username" validate:"required"`
	Email    string `json:"email"  validate:"required"`
	Password string `json:"password" validate:"required"`
}

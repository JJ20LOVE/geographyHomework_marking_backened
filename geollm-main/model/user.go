package model

type SignUpParam struct {
	Name       string `json:"name" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

type LoginParam struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type User struct {
	UserID   int    `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Email    string `db:"email"`
}

type ChangePassParam struct {
	Username string `json:"username" binding:"required"`
	OldPass  string `json:"old_pass" binding:"required"`
	NewPass  string `json:"new_pass" binding:"required"`
}

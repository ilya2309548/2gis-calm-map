package model

type User struct {
	ID    uint   `json:"id" gorm:"primary`
	Name  string `json:"name"`
	Email string `json:"email"`
}

package model

type User struct {
	ID       uint   `json:"id" gorm:"primaryKey"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"-"`    // не отдаём наружу!
	Role     string `json:"role"` // например, "user", "admin"
}

// Organization represents a business entity owned by a user (1:1)
type Organization struct {
	ID               uint     `json:"id" gorm:"primaryKey"`
	OwnerID          uint     `json:"owner_id" gorm:"uniqueIndex"` // один владелец - одна организация
	Owner            User     `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Address          string   `json:"address"`
	Longitude        *float64 `json:"longitude"` // optional
	Latitude         *float64 `json:"latitude"`  // optional
	OrganizationType string   `json:"organization_type"`
}

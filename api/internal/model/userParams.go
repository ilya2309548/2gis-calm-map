package model

type UserParams struct {
	ID            uint `json:"id" gorm:"primaryKey"`
	UserID        uint `json:"user_id" gorm:"uniqueIndex"` // Один к одному
	User          User `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Appearance    bool `json:"appearance"`
	Lighting      bool `json:"lighting"`
	Smell         bool `json:"smell"`
	Temperature   bool `json:"temperature"`
	Tactility     bool `json:"tactility"`
	Signage       bool `json:"signage"`
	Intuitiveness bool `json:"intuitiveness"`
	StaffAttitude bool `json:"staff_attitude"`
	PeopleDensity bool `json:"people_density"`
	SelfService   bool `json:"self_service"`
	Calmness      bool `json:"calmness"`
}

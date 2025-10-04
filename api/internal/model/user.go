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
	ID               uint                `json:"id" gorm:"primaryKey"`
	OwnerID          uint                `json:"owner_id" gorm:"uniqueIndex"` // один владелец - одна организация
	Owner            User                `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Address          string              `json:"address"`
	Longitude        *float64            `json:"longitude"` // optional
	Latitude         *float64            `json:"latitude"`  // optional
	OrganizationType string              `json:"organization_type"`
	Params           *OrganizationParams `json:"params,omitempty" gorm:"foreignKey:OrganizationID;references:ID"`
}

// OrganizationParams aggregates ratings for an organization (1:1)
type OrganizationParams struct {
	ID             uint         `json:"id" gorm:"primaryKey"`
	OrganizationID uint         `json:"organization_id" gorm:"uniqueIndex"`
	Organization   Organization `json:"-" gorm:"constraint:OnDelete:CASCADE"`

	AppearanceAvg   float64 `json:"appearance_avg"`
	AppearanceCount uint    `json:"appearance_count"`
	AppearanceSum   uint    `json:"appearance_sum"`

	LightingAvg   float64 `json:"lighting_avg"`
	LightingCount uint    `json:"lighting_count"`
	LightingSum   uint    `json:"lighting_sum"`

	SmellAvg   float64 `json:"smell_avg"`
	SmellCount uint    `json:"smell_count"`
	SmellSum   uint    `json:"smell_sum"`

	TemperatureAvg   float64 `json:"temperature_avg"`
	TemperatureCount uint    `json:"temperature_count"`
	TemperatureSum   uint    `json:"temperature_sum"`

	TactilityAvg   float64 `json:"tactility_avg"`
	TactilityCount uint    `json:"tactility_count"`
	TactilitySum   uint    `json:"tactility_sum"`

	SignageAvg   float64 `json:"signage_avg"`
	SignageCount uint    `json:"signage_count"`
	SignageSum   uint    `json:"signage_sum"`

	IntuitivenessAvg   float64 `json:"intuitiveness_avg"`
	IntuitivenessCount uint    `json:"intuitiveness_count"`
	IntuitivenessSum   uint    `json:"intuitiveness_sum"`

	StaffAttitudeAvg   float64 `json:"staff_attitude_avg"`
	StaffAttitudeCount uint    `json:"staff_attitude_count"`
	StaffAttitudeSum   uint    `json:"staff_attitude_sum"`

	PeopleDensityAvg   float64 `json:"people_density_avg"`
	PeopleDensityCount uint    `json:"people_density_count"`
	PeopleDensitySum   uint    `json:"people_density_sum"`

	SelfServiceAvg   float64 `json:"self_service_avg"`
	SelfServiceCount uint    `json:"self_service_count"`
	SelfServiceSum   uint    `json:"self_service_sum"`

	CalmnessAvg   float64 `json:"calmness_avg"`
	CalmnessCount uint    `json:"calmness_count"`
	CalmnessSum   uint    `json:"calmness_sum"`
}

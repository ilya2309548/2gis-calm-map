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
	OwnerID          uint                `json:"owner_id"` // ВНИМАНИЕ: для обычного владельца (role=owner) разрешаем только одну организацию логикой приложения; admin может иметь несколько
	Owner            User                `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	Address          string              `json:"address"`
	Longitude        *float64            `json:"longitude"` // optional
	Latitude         *float64            `json:"latitude"`  // optional
	OrganizationType string              `json:"organization_type"`
	Params           *OrganizationParams `json:"params,omitempty" gorm:"foreignKey:OrganizationID;references:ID"`
	MapPath          *string             `json:"map_path"`     // относительный путь к карте (изображение)
	PicturePath      *string             `json:"picture_path"` // относительный путь к общей картинке
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

// OrganizationComment represents a single user comment with optional ratings per parameter.
// Each numeric field (Value) is optional (nil => not provided). For every non-nil value we also can store an optional text comment.
type OrganizationComment struct {
	ID             uint         `json:"id" gorm:"primaryKey"`
	OrganizationID uint         `json:"organization_id" index:"idx_org_comment"`
	Organization   Organization `json:"-" gorm:"constraint:OnDelete:CASCADE"`
	UserID         uint         `json:"user_id"` // author of the comment
	User           User         `json:"-" gorm:"constraint:OnDelete:CASCADE"`

	Text     *string  `json:"text"`    // общий текст комментария (опционально)
	AvgValue *float64 `json:"avg_val"` // средняя по непустым параметрам (вычисляется при создании)

	AppearanceValue      *uint   `json:"appearance_value"`
	AppearanceComment    *string `json:"appearance_comment"`
	LightingValue        *uint   `json:"lighting_value"`
	LightingComment      *string `json:"lighting_comment"`
	SmellValue           *uint   `json:"smell_value"`
	SmellComment         *string `json:"smell_comment"`
	TemperatureValue     *uint   `json:"temperature_value"`
	TemperatureComment   *string `json:"temperature_comment"`
	TactilityValue       *uint   `json:"tactility_value"`
	TactilityComment     *string `json:"tactility_comment"`
	SignageValue         *uint   `json:"signage_value"`
	SignageComment       *string `json:"signage_comment"`
	IntuitivenessValue   *uint   `json:"intuitiveness_value"`
	IntuitivenessComment *string `json:"intuitiveness_comment"`
	StaffAttitudeValue   *uint   `json:"staff_attitude_value"`
	StaffAttitudeComment *string `json:"staff_attitude_comment"`
	PeopleDensityValue   *uint   `json:"people_density_value"`
	PeopleDensityComment *string `json:"people_density_comment"`
	SelfServiceValue     *uint   `json:"self_service_value"`
	SelfServiceComment   *string `json:"self_service_comment"`
	CalmnessValue        *uint   `json:"calmness_value"`
	CalmnessComment      *string `json:"calmness_comment"`
}

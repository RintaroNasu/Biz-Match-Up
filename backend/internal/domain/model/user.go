package model

type User struct {
	ID                 uint   `gorm:"primaryKey;autoIncrement"`
	Email              string `gorm:"unique;not null"`
	Password           string `gorm:"not null"`
	Name               *string
	DesiredJobType     *string
	DesiredLocation    *string
	DesiredCompanySize *string
	CareerAxis1        *string
	CareerAxis2        *string
	SelfPr             *string
	Reasons            []Reason `gorm:"foreignKey:UserID"`
}
type AuthResponse struct {
	Message string `json:"message"`
	User    User   `json:"user"`
	Token   string `json:"token"`
}
type UpdateUserProfileRequest struct {
	Name               string
	DesiredJobType     string
	DesiredLocation    string
	DesiredCompanySize string
	CareerAxis1        string
	CareerAxis2        string
	SelfPr             string
}

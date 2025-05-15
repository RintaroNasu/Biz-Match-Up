package models

type User struct {
	ID                uint           `gorm:"primaryKey;autoIncrement"`
	Email             string         `gorm:"unique;not null"`
	Password          string         `gorm:"not null"`
	Name              *string
	DesiredJobType    *string
	DesiredLocation   *string
	DesiredCompanySize *string
	CareerAxis1       *string
	CareerAxis2       *string
	SelfPr            *string
}

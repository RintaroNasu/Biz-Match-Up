package model

import "time"

type Reason struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Content     string `gorm:"not null"`
	UserID      uint   `gorm:"not null"`
	User        User   `gorm:"foreignKey:UserID"`
	CompanyName string `gorm:"not null"`
	CompanyUrl  string `gorm:"not null"`
	CreatedAt   time.Time
}

type MatchItem struct {
	Axis   string `json:"axis"`
	Score  int    `json:"score"`
	Reason string `json:"reason"`
}

type GenerateReasonsRequest struct {
	MatchResult []MatchItem     `json:"matchResult"`
	Questions   QuestionAnswers `json:"questions"`
}

type QuestionAnswers struct {
	ReasonInterest    string `json:"reasonInterest"`
	AttractiveService string `json:"attractiveService"`
	RelatedExperience string `json:"relatedExperience"`
}

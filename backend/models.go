package main

import (
	"time"

	"gorm.io/datatypes"
)

// Quiz represents a quiz; questions and options are normalized into separate tables.
type Quiz struct {
	ID        string     `json:"id" gorm:"primaryKey;column:id"`
	Title     string     `json:"title" gorm:"column:title"`
	Questions []Question `json:"questions" gorm:"foreignKey:QuizID;constraint:OnDelete:CASCADE"`
	CreatedAt time.Time  `json:"createdAt" gorm:"column:created_at;index"`
}

type Question struct {
	ID       uint     `json:"id" gorm:"primaryKey;autoIncrement"`
	QuizID   string   `json:"-" gorm:"index;column:quiz_id"`
	Text     string   `json:"text" gorm:"column:text"`
	Options  []Option `json:"options" gorm:"foreignKey:QuestionID;constraint:OnDelete:CASCADE"`
	Multiple bool     `json:"multiple" gorm:"column:multiple"`
	// Answers stores the correct answer indices (JSON array of ints)
	Answers datatypes.JSON `json:"answers,omitempty" gorm:"type:jsonb;column:answers"`
}

type Option struct {
	ID         uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	QuestionID uint   `json:"-" gorm:"index;column:question_id"`
	Text       string `json:"text" gorm:"column:text"`
	Ord        int    `json:"ord" gorm:"column:ord"`
}

// User represents an authenticated quiz creator/player.
type User struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	Email        string    `json:"email" gorm:"column:email;uniqueIndex"`
	PasswordHash string    `json:"-" gorm:"column:password_hash"`
	DisplayName  string    `json:"displayName" gorm:"column:display_name"`
	CreatedAt    time.Time `json:"createdAt" gorm:"column:created_at"`
}

// Submission groups responses for a single quiz attempt.
type Submission struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	QuizID    string    `json:"quizId" gorm:"column:quiz_id;index"`
	UserID    *uint     `json:"userId" gorm:"column:user_id;index"`
	Score     *int      `json:"score" gorm:"column:score"`
	CreatedAt time.Time `json:"createdAt" gorm:"column:created_at;index"`
}

// Response stores a user's answer for analytics; Selected can be an array of indices.
type Response struct {
	ID           uint           `json:"id" gorm:"primaryKey"`
	QuizID       string         `json:"quizId" gorm:"index;column:quiz_id"`
	QuestionID   uint           `json:"questionId" gorm:"index;column:question_id"`
	SubmissionID *uint          `json:"submissionId" gorm:"column:submission_id;index"`
	UserID       *uint          `json:"userId,omitempty" gorm:"column:user_id"`
	Selected     datatypes.JSON `json:"selected" gorm:"type:jsonb;column:selected"`
	IsCorrect    *bool          `json:"isCorrect" gorm:"column:is_correct"`
	CreatedAt    time.Time      `json:"createdAt" gorm:"column:created_at;index"`
}

func (Quiz) TableName() string       { return "quizzes" }
func (Question) TableName() string   { return "questions" }
func (Option) TableName() string     { return "options" }
func (User) TableName() string       { return "users" }
func (Submission) TableName() string { return "submissions" }
func (Response) TableName() string   { return "responses" }

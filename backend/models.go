package main

import "time"

// Quiz represents a quiz with simple questions.
type Quiz struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Questions []Question `json:"questions"`
	CreatedAt time.Time  `json:"createdAt"`
}

type Question struct {
	Text    string   `json:"text"`
	Options []string `json:"options"`
	// If multiple answers allowed, this can hold indices.
	Multiple bool  `json:"multiple"`
	Answers  []int `json:"answers"`
}

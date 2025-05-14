package models

type News struct {
	Id         int64   `json:"id"`
	Title      string  `json:"title,omitempty"`
	Content    string  `json:"content,omitempty"`
	Categories []int64 `json:"categories,omitempty"`
}

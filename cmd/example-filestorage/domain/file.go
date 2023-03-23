package domain

import "time"

type File struct {
	Id      *string    `json:"id"`
	Name    *string    `json:"name"`
	Content *string    `json:"content"`
	Owner   *string    `json:"owner"`
	Created *time.Time `json:"created"`
}

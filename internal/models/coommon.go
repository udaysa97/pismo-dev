package models

type Pagination struct {
	Limit     int    `json:"limit"`
	Page      int    `json:"page"`
	Sort      string `json:"sort"`
	Direction string `json:"direction"`
}

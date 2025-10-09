package model

type Class struct {
	ClassID   int    `db:"class_id" json:"class_id"`
	ClassName string `db:"class_name" json:"class_name"`
}

type BaseClass struct {
	ClassName string `json:"class_name" binding:"required"`
}

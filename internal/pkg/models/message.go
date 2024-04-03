package models

import "time"

type Message struct {
	Id       uint32 `gorm:"primaryKey"`
	RoomId   uint32
	SenderId uint32
	Content  string
	Date     time.Time
}

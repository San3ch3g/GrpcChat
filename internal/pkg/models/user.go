package models

type User struct {
	Id       uint32 `gorm:"primaryKey"`
	Nickname string `gorm:"unique"`
	Password string
}

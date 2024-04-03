package models

type Room struct {
	Id   uint32 `gorm:"primaryKey"`
	Name string `gorm:"unique"`
}

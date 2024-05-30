package model

type User struct {
	ID    uint   `gorm:"primaryKey;autoIncrement"`
	Name  string `qorm:"size:255"`
	Email string `gorm:"size:255;unique"`
	Pwd   string `gorm:"size:255"`
}

package model

type Image struct {
	ID    int64  `gorm:"primaryKey"`
	Code  string `gorm:"column:code"`
	Image string `gorm:"column:image"`
}

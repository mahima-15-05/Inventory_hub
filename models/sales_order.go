package models

import "time"

type SalesOrder struct {
	ID         uint `gorm:"primaryKey"`
	ProductID  *uint
	Product    Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Quantity   int
	TotalPrice float64
	OrderDate  time.Time
}

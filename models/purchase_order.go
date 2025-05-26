package models

import "time"

type PurchaseOrder struct {
	ID        uint    `gorm:"primaryKey"`
	ProductID *uint    
	Product   Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Quantity  int
	Supplier  string
	OrderDate time.Time
}

package models

type Product struct {
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"not null"`
	Description   string
	Price         float64
	Quantity      int
}

// In Go, if a field starts with a lowercase letter (like id, name, etc.), it's unexported and invisible to GORM, JSON, or any other reflection-based tools.

// gorm:"..." → tells GORM how to map fields to the DB

// json:"..." → tells encoding/json how to marshal/unmarshal JSON

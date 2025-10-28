package models

type User struct {
	ID       uint      `gorm:"primaryKey"`
	Username string    `gorm:"unique;not null"` // Unique username
	APIKey   string    `gorm:"not null"`        // API key for authentication
	Booleans []Boolean // One-to-many relation
}

type Boolean struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"not null;index:idx_user_boolean,unique"` // Foreign key reference
	Name   string `gorm:"not null;index:idx_user_boolean,unique"` // Name unique per user
	Value  bool   `gorm:"not null"`
}

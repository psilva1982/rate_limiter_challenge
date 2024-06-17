package database

type User struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"size:254" gorm:"uniqueIndex"`
	Password string `gorm:"size:254"`
	APIKey   string `gorm:"size:254"`
}

package persistent

import (
	"github.com/impulseweb3/hypixel-skyblock-auctions-go/internal/hypixel"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB
}

func NewDatabase(db *gorm.DB) *Database {
	return &Database{
		DB: db,
	}
}

func (d *Database) AutoMigrate() error {
	return d.DB.AutoMigrate(&hypixel.Auction{}, &hypixel.EndedAuction{})
}

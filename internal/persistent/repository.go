package persistent

import (
	"github.com/impulseweb3/hypixel-skyblock-auctions-go/internal/hypixel"
	"gorm.io/gorm/clause"
)

type Repository struct {
	database *Database
}

func NewRepository(database *Database) *Repository {
	return &Repository{
		database: database,
	}
}

func (r *Repository) SaveAuctions(auctions []hypixel.Auction) error {
	return r.database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "uuid"}},
		DoNothing: true,
	}).Create(&auctions).Error
}

func (r *Repository) SaveEndedAuctions(endedAuctions []hypixel.EndedAuction) error {
	return r.database.DB.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "auction_id"}},
		DoNothing: true,
	}).Create(&endedAuctions).Error
}

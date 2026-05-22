package persistent

import (
	"time"

	"github.com/impulseweb3/hypixel-skyblock-auctions-go/internal/hypixel"
)

type Repository struct {
	database *Database
}

func NewRepository(database *Database) *Repository {
	return &Repository{
		database: database,
	}
}

func (r *Repository) FindAverage(itemName string) (uint64, error) {
	var average int64

	subquery := r.database.DB.
		Table("ended_auctions e").
		Select("e.price").
		Joins("inner join auctions a on a.uuid = e.auction_id").
		Where("a.item_name = ?", itemName).
		Order("e.timestamp desc").
		Limit(10)

	err := r.database.DB.
		Table("(?) as subquery", subquery).
		Select("coalesce(round(avg(price)), 0)").
		Scan(&average).
		Error

	return uint64(average), err
}

func (r *Repository) FindVolume(itemName string) (uint64, error) {
	var volume int64

	last24Hours := time.
		Now().
		Add(-24 * time.Hour).
		UnixMilli()

	err := r.database.DB.
		Table("ended_auctions e").
		Joins("inner join auctions a on a.uuid = e.auction_id").
		Where("a.item_name = ? and e.timestamp >= ?", itemName, last24Hours).
		Count(&volume).
		Error

	return uint64(volume), err
}

func (r *Repository) SaveAuctions(auctions []hypixel.Auction) error {
	return r.database.DB.Create(&auctions).Error
}

func (r *Repository) SaveEndedAuctions(endedAuctions []hypixel.EndedAuction) error {
	return r.database.DB.Create(&endedAuctions).Error
}

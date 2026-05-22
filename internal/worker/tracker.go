package worker

import (
	"log/slog"
	"time"

	"github.com/impulseweb3/hypixel-skyblock-auctions-go/internal/hypixel"
	"github.com/impulseweb3/hypixel-skyblock-auctions-go/internal/persistent"
)

type Tracker struct {
	repository    *persistent.Repository
	hypixelClient *hypixel.Client
}

func NewTracker(repository *persistent.Repository, hypixelClient *hypixel.Client) *Tracker {
	return &Tracker{
		repository:    repository,
		hypixelClient: hypixelClient,
	}
}

func sleepUntilNextSecond() {
	nextSecond := time.Now().Truncate(time.Second).Add(time.Second)
	time.Sleep(time.Until(nextSecond))
}

func (t *Tracker) Start() {
	auctionsCache := make(map[string]struct{})
	endedAuctionsCache := make(map[string]struct{})

	lastUpdated := uint64(0)
	slog.Info("Tracker started")

	for {
		tempAuctionsCache := make(map[string]struct{})
		tempEndedAuctionsCache := make(map[string]struct{})

		auctionsBatch := make([]hypixel.Auction, 0)
		endedAuctionsBatch := make([]hypixel.EndedAuction, 0)

		auctions, auctionsError := t.hypixelClient.FetchAuctions()
		endedAuctions, endedAuctionsError := t.hypixelClient.FetchEndedAuctions()

		if auctionsError != nil {
			slog.Error(auctionsError.Error())
			continue
		}

		if endedAuctionsError != nil {
			slog.Error(endedAuctionsError.Error())
			continue
		}

		if auctions.LastUpdated == lastUpdated {
			slog.Debug("Auctions not updated")
			sleepUntilNextSecond()
			continue
		}

		if endedAuctions.LastUpdated == lastUpdated {
			slog.Debug("Ended auctions not updated")
			sleepUntilNextSecond()
			continue
		}

		for _, auction := range auctions.Auctions {
			tempAuctionsCache[auction.UUID] = struct{}{}

			if _, exists := auctionsCache[auction.UUID]; !exists {
				auctionsBatch = append(auctionsBatch, auction)
			}
		}

		for _, endedAuction := range endedAuctions.EndedAuctions {
			tempEndedAuctionsCache[endedAuction.AuctionID] = struct{}{}

			if _, exists := endedAuctionsCache[endedAuction.AuctionID]; !exists {
				endedAuctionsBatch = append(endedAuctionsBatch, endedAuction)
			}
		}

		saveAuctionsError := t.repository.SaveAuctions(auctionsBatch)
		saveEndedAuctionsError := t.repository.SaveEndedAuctions(endedAuctionsBatch)

		if saveAuctionsError != nil {
			slog.Error(saveAuctionsError.Error())
			continue
		}

		if saveEndedAuctionsError != nil {
			slog.Error(saveEndedAuctionsError.Error())
			continue
		}

		auctionsCache = tempAuctionsCache
		endedAuctionsCache = tempEndedAuctionsCache

		lastUpdated = (auctions.LastUpdated + endedAuctions.LastUpdated) / 2
		slog.Info("Auctions and ended auctions saved")

		t.track(auctionsBatch)
		slog.Info("Auctions tracked")

		targetTime := time.UnixMilli(int64(lastUpdated)).Add(60 * time.Second)
		sleepDuration := time.Until(targetTime)

		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}
	}
}

func (t *Tracker) track(auctions []hypixel.Auction) {}

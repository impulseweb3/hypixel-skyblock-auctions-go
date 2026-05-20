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
	var lastUpdated uint64

	for {
		auctions, auctionsError := t.hypixelClient.FetchAuctions()
		endedAuctions, endedAuctionsError := t.hypixelClient.FetchEndedAuctions()

		if auctionsError != nil {
			slog.Error(auctionsError.Error())
			sleepUntilNextSecond()
			continue
		}

		if endedAuctionsError != nil {
			slog.Error(endedAuctionsError.Error())
			sleepUntilNextSecond()
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

		saveAuctionsError := t.repository.SaveAuctions(auctions.Auctions)
		saveEndedAuctionsError := t.repository.SaveEndedAuctions(endedAuctions.EndedAuctions)

		if saveAuctionsError != nil {
			slog.Error(saveAuctionsError.Error())
			continue
		}

		if saveEndedAuctionsError != nil {
			slog.Error(saveEndedAuctionsError.Error())
			continue
		}

		lastUpdated = (auctions.LastUpdated + endedAuctions.LastUpdated) / 2
		slog.Info("Auctions updated")

		t.track(auctions.Auctions)
		slog.Info("Auctions tracked")

		targetTime := time.UnixMilli(int64(lastUpdated)).Add(60 * time.Second)
		sleepDuration := time.Until(targetTime)

		if sleepDuration > 0 {
			time.Sleep(sleepDuration)
		}
	}
}

func (t *Tracker) track(auctions []hypixel.Auction) {}

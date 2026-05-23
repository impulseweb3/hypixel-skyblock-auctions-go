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

func sleepUntilNextUpdate(lastUpdated uint64) {
	nextUpdate := time.UnixMilli(int64(lastUpdated)).Add(time.Minute)
	time.Sleep(time.Until(nextUpdate))
}

func (t *Tracker) Start() {
	lastUpdated := uint64(0)
	slog.Info("Tracker started")

	for {
		auctionsResponse, endedAuctionsResponse, err := t.fetchAuctionsResponseAndEndedAuctionsResponse()

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		if auctionsResponse.LastUpdated == lastUpdated || endedAuctionsResponse.LastUpdated == lastUpdated {
			slog.Debug("Auctions or ended auctions not updated")
			sleepUntilNextSecond()
			continue
		}

		auctions, endedAuctions := t.getAuctionsAndEndedAuctions(auctionsResponse, endedAuctionsResponse, lastUpdated)
		err = t.saveAuctionsAndEndedAuctions(auctions, endedAuctions)

		if err != nil {
			slog.Error(err.Error())
			continue
		}

		lastUpdated = auctionsResponse.LastUpdated
		slog.Info("Auctions and ended auctions saved")

		t.trackAuctions(auctions)
		slog.Info("Auctions tracked")

		slog.Info("Sleeping until next update")
		sleepUntilNextUpdate(lastUpdated)
	}
}

func (t *Tracker) fetchAuctionsResponseAndEndedAuctionsResponse() (*hypixel.AuctionsResponse, *hypixel.EndedAuctionsResponse, error) {
	auctionsResponse, err := t.hypixelClient.FetchAuctionsResponse()

	if err != nil {
		return nil, nil, err
	}

	endedAuctionsResponse, err := t.hypixelClient.FetchEndedAuctionsResponse()

	if err != nil {
		return nil, nil, err
	}

	return auctionsResponse, endedAuctionsResponse, nil
}

func (t *Tracker) getAuctionsAndEndedAuctions(auctionsResponse *hypixel.AuctionsResponse, endedAuctionsResponse *hypixel.EndedAuctionsResponse, lastUpdated uint64) ([]hypixel.Auction, []hypixel.EndedAuction) {
	auctions := make([]hypixel.Auction, 0, 1000)

	for _, auction := range auctionsResponse.Auctions {
		if auction.Start >= lastUpdated {
			auctions = append(auctions, auction)
		}
	}

	return auctions, endedAuctionsResponse.EndedAuctions
}

func (t *Tracker) saveAuctionsAndEndedAuctions(auctions []hypixel.Auction, endedAuctions []hypixel.EndedAuction) error {
	err := t.repository.SaveAuctions(auctions)

	if err != nil {
		return err
	}

	err = t.repository.SaveEndedAuctions(endedAuctions)

	if err != nil {
		return err
	}

	return nil
}

func (t *Tracker) trackAuctions(auctions []hypixel.Auction) {}

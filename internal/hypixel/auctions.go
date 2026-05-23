package hypixel

import "encoding/json"

type Auction struct {
	UUID             string `json:"uuid" gorm:"primaryKey"`
	Auctioneer       string `json:"auctioneer"`
	Start            uint64 `json:"start"`
	End              uint64 `json:"end"`
	ItemName         string `json:"item_name" gorm:"index"`
	Extra            string `json:"extra"`
	Category         string `json:"category"`
	Tier             string `json:"tier"`
	StartingBid      uint64 `json:"starting_bid"`
	ItemBytes        string `json:"item_bytes"`
	HighestBidAmount uint64 `json:"highest_bid_amount"`
	Bin              bool   `json:"bin"`
}

type EndedAuction struct {
	AuctionID string `json:"auction_id" gorm:"primaryKey"`
	Seller    string `json:"seller"`
	Buyer     string `json:"buyer"`
	Timestamp uint64 `json:"timestamp" gorm:"index:,sort:desc"`
	Price     uint64 `json:"price"`
	Bin       bool   `json:"bin"`
}

type AuctionsResponse struct {
	LastUpdated uint64    `json:"lastUpdated"`
	Auctions    []Auction `json:"auctions"`
}

type EndedAuctionsResponse struct {
	LastUpdated   uint64         `json:"lastUpdated"`
	EndedAuctions []EndedAuction `json:"auctions"`
}

func (c *Client) getAndDecode(url string, v any) error {
	resp, err := c.httpClient.Get(url)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	err = json.NewDecoder(resp.Body).Decode(v)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) FetchAuctionsResponse() (*AuctionsResponse, error) {
	var auctionsResponse AuctionsResponse

	err := c.getAndDecode("https://api.hypixel.net/v2/skyblock/auctions", &auctionsResponse)

	if err != nil {
		return nil, err
	}

	return &auctionsResponse, nil
}

func (c *Client) FetchEndedAuctionsResponse() (*EndedAuctionsResponse, error) {
	var endedAuctionsResponse EndedAuctionsResponse

	err := c.getAndDecode("https://api.hypixel.net/v2/skyblock/auctions_ended", &endedAuctionsResponse)

	if err != nil {
		return nil, err
	}

	return &endedAuctionsResponse, nil
}

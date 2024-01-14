package shop

import (
	"beerium/pkg/config"
	"beerium/pkg/monitor"
	"fmt"
)

type Collector struct {
	config *config.Config

	shops map[string]shop
}

func NewCollector(c *config.Config, m *monitor.Monitor) *Collector {
	shops := []shop{
		NewBaracek(m.Registry),
		NewManeo(m.Registry),
	}

	shopMap := make(map[string]shop, len(shops))
	for _, s := range shops {
		shopMap[s.getName()] = s
	}

	return &Collector{
		config: c,
		shops:  shopMap,
	}
}

func (c *Collector) GetBeer(shopName, url string) (*Beer, error) {
	shop, found := c.shops[shopName]
	if !found {
		return nil, fmt.Errorf("shopName %s not found", shopName)
	}

	beer, err := shop.getBeer(url)
	if err != nil {
		return nil, fmt.Errorf("could not return beer: %w", err)
	}

	return beer, nil
}

func createHttpClientForShop() {

}

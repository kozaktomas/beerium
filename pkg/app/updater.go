package app

import (
	"beerium/pkg/shop"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"golang.org/x/net/context"
	"sync"
	"time"
)

// RunUpdater is responsible for updating beer prices and stock
// it runs periodically and updates all beers that are older than 2 hours
func (a *App) RunUpdater(ctx context.Context) {
	ticket := time.NewTicker(5 * time.Minute)
	for {
		select {
		case <-ctx.Done():
			a.logger.Info("stopping updater")
			return
		case <-ticket.C:
			a.logger.Info("Full update started")
			a.mtx.Lock()
			a.updateAll()
			a.mtx.Unlock()
			a.logger.Info("Full update finished")
		}
	}
}

// updateAll iterates over all beers and updates the ones that are older than 2 hours
func (a *App) updateAll() {
	wg := sync.WaitGroup{}

	for _, beer := range a.beers {
		for _, s := range beer.Shops {
			if time.Since(s.lastUpdate) > 2*time.Hour {
				wg.Add(1)
				go func(beer *Beer, s *Shop) {
					defer wg.Done()
					a.updateOne(beer, s)
				}(beer, s)
			}
		}
	}

	wg.Wait()
}

func (a *App) updateOne(beer *Beer, s *Shop) {
	b, err := a.collector.GetBeer(s.Name, s.Url)
	if err != nil {
		a.logger.Errorf("could not update beer %s for shop %s", beer.Name, s.Name)
	} else {
		s.Price = b.Price
		s.Stock = b.Stock == shop.StockTypeAvailable
		s.lastUpdate = time.Now()

		a.logger.Infof("beer %s updated in shop %s", beer.Name, s.Name)
		a.monitor.BeerPrice.With(prometheus.Labels{
			"brand":  beer.Brand,
			"name":   beer.Name,
			"type":   beer.Type,
			"degree": fmt.Sprintf("%d", beer.Degree),
			"size":   fmt.Sprintf("%d", beer.Size),
			"shop":   s.Name,
		}).Set(float64(s.Price))
	}
}

package monitor

import (
	"github.com/prometheus/client_golang/prometheus"
)

type Monitor struct {
	Registry  *prometheus.Registry
	BeerPrice *prometheus.GaugeVec
}

func NewMonitor() *Monitor {
	r := prometheus.NewRegistry()
	m := &Monitor{
		Registry: r,
		BeerPrice: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Name: "beer_price_czk",
			Help: "Beer price in CZK",
		}, []string{"brand", "name", "type", "degree", "size", "shop"}),
	}

	r.MustRegister(m.BeerPrice)

	return m
}

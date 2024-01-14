package app

import (
	"beerium/pkg/config"
	"beerium/pkg/monitor"
	"beerium/pkg/shop"
	"github.com/sirupsen/logrus"
	"html/template"
	"path"
	"sync"
	"time"
)

type App struct {
	config    *config.Config
	collector *shop.Collector
	monitor   *monitor.Monitor

	beers []*Beer

	htmlTemplate *template.Template
	logger       *logrus.Logger
	mtx          sync.RWMutex
}

func NewApp(cfg *config.Config, collector *shop.Collector, logger *logrus.Logger, monitor *monitor.Monitor) *App {
	// initialize data structure
	data := make([]*Beer, len(cfg.Beers))
	i := 0
	for _, beer := range cfg.Beers {
		shops := make(map[string]*Shop, len(beer.Shops))
		for shopName, shopUrl := range beer.Shops {
			shops[shopName] = &Shop{
				Name:       shopName,
				Price:      0,
				Stock:      false,
				Url:        shopUrl,
				lastUpdate: time.Time{}.Add(-7 * time.Hour),
			}
		}
		data[i] = &Beer{
			key:    beer.Key,
			Name:   beer.Name,
			Brand:  beer.Brand,
			Type:   beer.Type,
			Degree: beer.Degree,
			Size:   beer.Size,
			Shops:  shops,
		}
		i++
	}

	// initialize html template
	fp := path.Join("pkg", "app", "templates", "homepage.gohtml")
	tmpl, err := template.ParseFiles(fp)
	if err != nil {
		logger.Fatalf("template %s not found: %s", fp, err)
	}

	return &App{
		config:    cfg,
		collector: collector,
		monitor:   monitor,

		beers: data,

		htmlTemplate: tmpl,
		logger:       logger,
		mtx:          sync.RWMutex{},
	}
}

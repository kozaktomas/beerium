package main

import (
	"beerium/pkg/app"
	"beerium/pkg/config"
	"beerium/pkg/monitor"
	"beerium/pkg/shop"
	"github.com/sirupsen/logrus"
	"os"
)

func main() {
	logger := logrus.New()
	logger.SetLevel(logrus.InfoLevel)
	logger.SetFormatter(&logrus.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	if len(os.Args) != 2 {
		logger.Fatalf("usage: %s <config-file>", os.Args[0])
	}

	if _, err := os.Stat(os.Args[1]); os.IsNotExist(err) {
		logger.Fatalf("config file %s does not exist", os.Args[1])
	}

	m := monitor.NewMonitor()
	cfg, err := config.LoadConfig(os.Args[1])
	if err != nil {
		logger.Fatalf("could not load config: %s", err)
	}
	collector := shop.NewCollector(cfg, m)
	beerApi := app.NewApp(cfg, collector, logger, m)
	beerApi.StartServer()
}

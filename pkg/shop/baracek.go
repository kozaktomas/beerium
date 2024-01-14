package shop

import (
	"bytes"
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/travelaudience/go-promhttp"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"strings"
)

type Baracek struct {
	pr *prometheus.Registry
}

func NewBaracek(pr *prometheus.Registry) *Baracek {
	return &Baracek{
		pr: pr,
	}
}

func (d *Baracek) getName() string {
	return "baracek"
}

func (d *Baracek) getBeer(url string) (*Beer, error) {
	promHttpClient := &promhttp.Client{
		Client:     http.DefaultClient,
		Registerer: d.pr,
	}
	httpClient, _ := promHttpClient.ForRecipient(d.getName())

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get response from Baracek: %w", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("count not read data from Baracek: %w", err)
	}

	doc, err := d.parseRoot(b)
	if err != nil {
		return nil, fmt.Errorf("could not parse html from Baracek: %w", err)
	}

	price, err := d.parsePrice(doc)
	if err != nil {
		return nil, fmt.Errorf("could not parse price from Baracek: %w", err)
	}
	stock, err := d.parseStock(doc)
	if err != nil {
		return nil, fmt.Errorf("could not parse stock from Baracek: %w", err)
	}

	return &Beer{
		Price: price,
		Stock: stock,
	}, nil
}

func (d *Baracek) parseRoot(body []byte) (*html.Node, error) {
	return html.Parse(bytes.NewReader(body))
}

func (d *Baracek) parsePrice(node *html.Node) (int, error) {
	xpath := "//span[contains(@class, 'uc-price')]"
	els, err := htmlquery.QueryAll(node, xpath)
	if err != nil {
		return 0, fmt.Errorf("could not parse price: %w", err)
	}
	if len(els) != 1 {
		return 0, fmt.Errorf("could not parse price: %w", err)
	}

	price, err := parsePriceString(els[0].FirstChild.Data)
	if err != nil {
		return 0, fmt.Errorf("could not parse price string: %w", err)
	}

	return price, nil
}

func (d *Baracek) parseStock(node *html.Node) (StockType, error) {
	xpath := "//div[contains(@class, 'field-name-field-skladem')]/div[contains(@class, 'field-items')]/div"
	els, err := htmlquery.QueryAll(node, xpath)
	if err != nil {
		return 0, fmt.Errorf("could not parse stock: %w", err)
	}
	if len(els) != 1 {
		return 0, fmt.Errorf("could not parse stock: %w", err)
	}

	if d.sanitizeString(els[0].FirstChild.Data) == "ano" {
		return StockTypeAvailable, nil
	}

	return StockTypeUnknown, nil
}

func (d *Baracek) sanitizeString(s string) string {
	s = strings.ToLower(s)
	s = strings.TrimSpace(s)
	return s
}

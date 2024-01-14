package shop

import (
	"fmt"
	"github.com/antchfx/htmlquery"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/travelaudience/go-promhttp"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
	"io"
	"net/http"
	"strings"
)

type Maneo struct {
	pr *prometheus.Registry
}

func NewManeo(pr *prometheus.Registry) *Maneo {
	return &Maneo{
		pr: pr,
	}
}

func (d *Maneo) getName() string {
	return "maneo"
}

func (d *Maneo) getBeer(url string) (*Beer, error) {
	promHttpClient := &promhttp.Client{
		Client:     http.DefaultClient,
		Registerer: d.pr,
	}
	httpClient, _ := promHttpClient.ForRecipient(d.getName())

	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("could not get response from Maneo: %w", err)
	}

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("count not read data from Maneo: %w", err)
	}

	body := d.decodeWindows1250(b)
	doc, err := d.parseRoot(body)

	price, err := d.parsePrice(doc)
	if err != nil {
		return nil, fmt.Errorf("could not parse price from Maneo: %w", err)
	}
	stock, err := d.parseStock(doc)
	if err != nil {
		return nil, fmt.Errorf("could not parse stock from Maneo: %w", err)
	}

	return &Beer{
		Price: price,
		Stock: stock,
	}, nil
}

func (d *Maneo) parseRoot(body string) (*html.Node, error) {
	return html.Parse(strings.NewReader(body))
}

func (d *Maneo) parsePrice(node *html.Node) (int, error) {
	xpath := "//div[contains(@class, 'dp-cena') and contains(@class, 'dp-right')]"
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

func (d *Maneo) parseStock(node *html.Node) (StockType, error) {
	xpath := "//span[contains(@class, 'tucne')]"
	els, err := htmlquery.QueryAll(node, xpath)
	if err != nil {
		return StockTypeUnknown, fmt.Errorf("could not parse stock: %w", err)
	}
	if len(els) != 1 {
		return StockTypeUnknown, fmt.Errorf("could not parse stock: %w", err)
	}

	if strings.ToLower(els[0].FirstChild.Data) == "skladem" {
		return StockTypeAvailable, nil
	}

	return StockTypeUnknown, nil
}

// decodeWindows1250 decodes windows-1250 encoded string to utf-8
// because Maneo is using windows-1250 encoding
func (d *Maneo) decodeWindows1250(enc []byte) string {
	dec := charmap.Windows1250.NewDecoder()
	out, _ := dec.Bytes(enc)
	return string(out)
}

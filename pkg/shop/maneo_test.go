package shop

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	body, err := os.ReadFile("./res/maneo.html")
	assert.Nil(t, err)
	m := Maneo{}
	doc, err := m.parseRoot(string(body))
	assert.Nil(t, err)
	assert.NotNil(t, doc)

	price, err := m.parsePrice(doc)
	assert.Nil(t, err)
	assert.Equal(t, 1690, price)

	stock, err := m.parseStock(doc)
	assert.Nil(t, err)
	var expectedStock StockType = StockTypeAvailable
	assert.Equal(t, expectedStock, stock)
}

func TestParsePrice(t *testing.T) {
	cases := []struct {
		in  string
		out int
	}{
		{"100,00 Kč", 100},
		{"100.00 Kč", 100},
		{"1 024,00 Kč", 1024},
		{"1 024,98 Kč", 1024},
	}

	for _, c := range cases {
		v, err := parsePriceString(c.in)
		assert.Nil(t, err)
		assert.Equal(t, c.out, v)
	}
}

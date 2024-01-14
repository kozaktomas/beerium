package shop

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestParseBaracek(t *testing.T) {
	body, err := os.ReadFile("./res/baracek.html")
	assert.Nil(t, err)
	m := Baracek{}
	doc, err := m.parseRoot(body)
	assert.Nil(t, err)
	assert.NotNil(t, doc)

	price, err := m.parsePrice(doc)
	assert.Nil(t, err)
	assert.Equal(t, 1565, price)

	stock, err := m.parseStock(doc)
	assert.Nil(t, err)
	var expectedStock StockType = StockTypeAvailable
	assert.Equal(t, expectedStock, stock)
}

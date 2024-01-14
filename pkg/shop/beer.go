package shop

type StockType int

const (
	StockTypeAvailable = iota
	StockTypeUnknown
)

type Beer struct {
	Price int // price in czk
	Stock StockType
}

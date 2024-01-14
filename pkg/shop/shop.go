package shop

type shop interface {
	getName() string
	getBeer(url string) (*Beer, error)
}

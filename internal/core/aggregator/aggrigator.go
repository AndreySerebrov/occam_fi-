package aggregator

type aggregator struct {
	tickerList []PriceStreamSubscriber
}

func New() *aggregator {
	return &aggregator{}
}

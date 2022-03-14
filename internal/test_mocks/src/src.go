package src

import (
	"occam-fi/internal/core/aggregator/ticker"
)

type Src struct {
	dataChan chan ticker.TickerPrice
	errChan  chan error
}

func New() *Src {
	return &Src{
		dataChan: make(chan ticker.TickerPrice),
		errChan:  make(chan error),
	}
}

func (s *Src) Pub(t ticker.TickerPrice) {
	s.dataChan <- t
}

func (s *Src) SendError(err error) {
	s.errChan <- err
}

func (s *Src) SubscribePriceStream() (chan ticker.TickerPrice, chan error) {
	return s.dataChan, s.errChan
}

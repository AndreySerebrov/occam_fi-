package aggregator

/*
The main logic is here
*/
import (
	"context"
	"log"
	"occam-fi/internal/core/aggregator/price"
	"occam-fi/internal/core/aggregator/publisher"
	"occam-fi/internal/core/aggregator/ticker"
	"time"
)

type errMsg struct {
	Err    error
	Ticker string
}

type aggregator struct {
	closeChan        chan struct{}
	resultChan       chan ticker.TickerPrice
	errorChan        chan errMsg
	pub              publisher.Publisher
	priceAccumulator price.Slice
	currentTimestamp time.Time
}

func New(log log.Logger, period time.Duration, pub publisher.Publisher) *aggregator {
	a := &aggregator{}
	t := time.NewTicker(period)
	a.currentTimestamp = time.Now()
	a.priceAccumulator = price.New()
	a.resultChan = make(chan ticker.TickerPrice)
	a.pub = pub

	go func() {
		for {
			select {
			case value := <-a.resultChan:
				// Requirements: Data from the streams can come with delays
				if value.Time.Before(a.currentTimestamp) {
					log.Printf("delayed data from source: %s", value.Ticker)
					break
				}
				a.priceAccumulator.Append(value.Price)
			case <-a.closeChan:
				log.Print("shut down process")
				return
			case errMsg := <-a.errorChan:
				log.Printf("error from %s : %s ", errMsg.Ticker, errMsg.Err.Error())
				return
			case timestamp := <-t.C:
				a.currentTimestamp = timestamp
				// Requirements: How to combine different prices into the index is up to you.
				// Here is a min price for a period
				fairPrice, err := a.priceAccumulator.Min()
				a.priceAccumulator.Reset()
				if err != nil {
					log.Printf("%s", err.Error())
					break
				}
				a.pub.Publish(float64(fairPrice), a.currentTimestamp)
				a.priceAccumulator = a.priceAccumulator[:0]
			}
		}
	}()
	return a
}

// RegisterSource registers a new data-source
func (a *aggregator) RegisterSource(ctx context.Context, src ticker.PriceStreamSubscriber, srcName string) {
	tickerChan, errorChan := src.SubscribePriceStream()
	go func() {
		select {
		case value := <-tickerChan:
			a.resultChan <- value
		case err := <-errorChan:
			a.errorChan <- errMsg{Err: err, Ticker: srcName}
			// Requirements: Stream can return an error, in that case the channel is closed.
			return
		case <-ctx.Done():
			return
		case <-a.closeChan:
			return
		}
	}()
}

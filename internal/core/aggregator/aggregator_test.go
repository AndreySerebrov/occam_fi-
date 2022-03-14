package aggregator

import (
	"context"
	"fmt"
	"log"
	"occam-fi/internal/core/aggregator/publisher"
	"occam-fi/internal/core/aggregator/ticker"
	source "occam-fi/internal/test_mocks/src"
	"os"
	"sync"

	"testing"
	"time"

	"github.com/golang/mock/gomock"
)

func TestSimple(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pub := publisher.NewMockPublisher(ctrl)

	pub.EXPECT().Publish(2.2, gomock.Any())

	l := log.New(os.Stdout, t.Name(), 0)
	aggregator := New(*l, time.Second*2, pub)

	src1 := source.New()
	src2 := source.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	aggregator.RegisterSource(ctx, src1, "Forex")
	aggregator.RegisterSource(ctx, src2, "Coin base")

	ts := time.Now()
	src1.Pub(ticker.TickerPrice{Ticker: "src1", Time: ts.Add(time.Second), Price: "2.2"})
	src2.Pub(ticker.TickerPrice{Ticker: "src2", Time: ts.Add(time.Second), Price: "2.4"})
	time.Sleep(time.Second * 3)
}

func TestDelayedData(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pub := publisher.NewMockPublisher(ctrl)

	pub.EXPECT().Publish(2.4, gomock.Any())
	l := log.New(os.Stdout, t.Name(), 0)
	aggregator := New(*l, time.Second*2, pub)

	src1 := source.New()
	src2 := source.New()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	aggregator.RegisterSource(ctx, src1, "Forex")
	aggregator.RegisterSource(ctx, src2, "Coin base")

	ts := time.Now()
	src1.Pub(ticker.TickerPrice{Ticker: "src1", Time: ts.Add(-time.Hour), Price: "2.2"})
	src2.Pub(ticker.TickerPrice{Ticker: "src2", Time: ts.Add(time.Second), Price: "2.4"})
	time.Sleep(time.Second * 3)
}

func TestMultiThread(t *testing.T) {

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	pub := publisher.NewMockPublisher(ctrl)

	pub.EXPECT().Publish(1.1, gomock.Any())
	l := log.New(os.Stdout, t.Name(), 0)
	aggregator := New(*l, time.Second*2, pub)
	ts := time.Now()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := sync.WaitGroup{}
	wg.Add(100)
	for i := 0; i < 100; i++ {
		src := source.New()
		aggregator.RegisterSource(ctx, src, "Forex")
		go func(i int, src *source.Src) {
			src.Pub(ticker.TickerPrice{Ticker: ticker.Ticker(fmt.Sprintf("%d", i)), Time: ts, Price: fmt.Sprintf("%d.1", i+1)})
			wg.Done()
		}(i, src)
	}
	wg.Wait()

	time.Sleep(time.Second * 2)
}

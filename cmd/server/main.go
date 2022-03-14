package main

import (
	"context"
	"log"
	"os"
	"time"

	"occam-fi/internal/core/aggregator"
	"occam-fi/internal/core/aggregator/ticker"
	source "occam-fi/internal/test_mocks/src"
	testPub "occam-fi/internal/test_mocks/stdout_pub"
)

func main() {
	l := log.New(os.Stdout, "main", 0)
	publisher := testPub.New()
	agg := aggregator.New(*l, time.Second*2, publisher)
	src := source.New()
	agg.RegisterSource(context.Background(), src, "Coin market")
	src.Pub(ticker.TickerPrice{Time: time.Now(), Price: "1.2"})
	time.Sleep(time.Second * 4)
}

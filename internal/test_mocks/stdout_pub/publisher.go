package stdoutpub

import (
	"fmt"
	"time"
)

type publisher struct {
}

func New() *publisher {
	return &publisher{}
}

func (p *publisher) Publish(price float64, ts time.Time) {
	fmt.Printf("%d, %.2f", ts.Unix(), price)
}

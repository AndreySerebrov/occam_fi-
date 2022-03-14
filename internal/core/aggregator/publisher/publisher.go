package publisher

import (
	"time"
)

/*
	Publisher interface
*/

//go:generate mockgen -destination publisher_mock.go -package publisher -source=publisher.go

type Publisher interface {
	Publish(price float64, timestamp time.Time)
}

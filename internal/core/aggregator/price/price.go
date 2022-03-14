package price

import (
	"strconv"
)

/*
	Package for manipulation with money
*/

type (
	price float64
	Slice []price
)

func New() Slice {
	return make(Slice, 0)
}

//Append new value
func (p *Slice) Append(item string) error {
	pr, err := parse(item)
	if err != nil {
		return err
	}
	*p = append(*p, pr)
	return nil
}

func (p *Slice) Reset() {
	*p = (*p)[:0]
}

// Min get min value
func (p Slice) Min() (price, error) {
	if len(p) == 0 {
		return 0, nil
	}
	var minValue = p[0]
	for _, price := range p {
		if minValue > price {
			minValue = price
		}
	}
	return minValue, nil
}

// Parse converts string to float64
func parse(strPrice string) (price, error) {

	p, err := strconv.ParseFloat(strPrice, 64)
	return price(p), err
}

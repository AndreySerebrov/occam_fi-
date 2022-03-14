package price

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAvg(t *testing.T) {
	accumulator := New()
	accumulator.Append("1.2")
	accumulator.Append("1.6")
	accumulator.Append("1.1")
	min, err := accumulator.Min()
	require.NoError(t, err)
	require.Equal(t, price(1.1), min)

	accumulator.Reset()

	accumulator.Append("2.1")
	accumulator.Append("4.2")
	accumulator.Append("6.3")
	accumulator.Append("3.4")
	min, err = accumulator.Min()
	require.NoError(t, err)
	require.Equal(t, price(2.1), min)
}

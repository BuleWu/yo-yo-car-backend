package utils

// EUR represents EUR euro amount in terms of cents
type EUR int64

// ToEUR converts a float64 to EUR
func ToEUR(f float64) EUR {
	return EUR((f * 100) + 0.5)
}

// Float64 converts an EUR to float64
func (m EUR) Float64() float64 {
	x := float64(m)
	x = x / 100
	return x
}

// Multiply safely multiplies an EUR value by a float64, rounding
// to the nearest cent.
func (m EUR) Multiply(f float64) EUR {
	x := (float64(m) * f) + 0.5
	return EUR(x)
}

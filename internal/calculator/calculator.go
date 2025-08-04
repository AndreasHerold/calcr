package calculator

// Add addiert zwei float64 Zahlen.
func Add(a, b float64) float64 {
	return a + b
}

// Subtract subtrahiert zwei float64 Zahlen.
func Subtract(a, b float64) float64 {
	return a - b
}

// Multiply multipliziert zwei float64 Zahlen.
func Multiply(a, b float64) float64 {
	return a * b
}

// Divide dividiert zwei float64 Zahlen.
func Divide(a, b float64) float64 {
	if b == 0 {
		return 0 // Division durch Null
	}
	return a / b
}
package utils

import (
	"fmt"
	"math"
)

func ValidateFloatPrecision(value float64) error {
	rounded := math.Round(value*100) / 100
	if rounded != value {
		return fmt.Errorf("invalid float precision: %.2f", value)
	}
	return nil
}

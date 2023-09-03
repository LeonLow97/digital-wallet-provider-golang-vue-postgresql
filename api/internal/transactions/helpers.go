package transactions

import (
	"fmt"
	"math"
	"reflect"
)

// Checks if the value passed in is of type float64
func IsFloat64(value interface{}) bool {
	// use reflection to retrieve the type of the value
	valueType := reflect.TypeOf(value)

	return valueType.Kind() == reflect.Float64
}

func ValidateFloatPrecision(value float64) error {
    expected := math.Floor(value*100) / 100
    if expected != value {
        return fmt.Errorf("invalid float precision: %.2f", value)
    }
    return nil
}

package utils

func CurrencyConversion(transferredAmount float64, transferredCurrency, receivedCurrency string) float64 {
	conversionRates := map[string]map[string]float64{
		"USD": {
			"SGD": 1.35,
			"EUR": 0.92,
		},
		"SGD": {
			"USD": 0.74,
			"EUR": 0.68,
		},
		"EUR": {
			"USD": 1.09,
			"SGD": 1.47,
		},
	}

	if transferredCurrency == receivedCurrency {
		return transferredAmount
	}

	if rate, ok := conversionRates[transferredCurrency][receivedCurrency]; ok {
		return transferredAmount * rate
	}

	return 0.0

}

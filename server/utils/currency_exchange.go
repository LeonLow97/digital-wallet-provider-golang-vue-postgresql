package utils

var exchangeRates = map[string]map[string]float64{
	"SGD": {
		"USD": 0.73,
		"AUD": 0.98,
		"MYR": 3.14,
	},
	"USD": {
		"SGD": 1.37,
		"AUD": 1.35,
		"MYR": 4.27,
	},
	"AUD": {
		"SGD": 1.02,
		"USD": 0.74,
		"MYR": 3.21,
	},
	"MYR": {
		"SGD": 0.32,
		"USD": 0.23,
		"AUD": 0.31,
	},
}

var spreads = map[string]map[string]float64{
	"SGD": {
		"USD": 0.005,
		"AUD": 0.004,
		"MYR": 0.007,
	},
	"USD": {
		"SGD": 0.006,
		"AUD": 0.005,
		"MYR": 0.008,
	},
	"AUD": {
		"SGD": 0.003,
		"USD": 0.004,
		"MYR": 0.006,
	},
	"MYR": {
		"SGD": 0.008,
		"USD": 0.007,
		"AUD": 0.009,
	},
}

func CalculateConversionDetails(transferAmount float64, fromCurrency, toCurrency string) (float64, float64) {
	// Calculate pegged rate (because exchange rates are hardcoded)
	peggedRate := exchangeRates[fromCurrency][toCurrency]

	// Calculate conversion amount
	conversionAmount := transferAmount * peggedRate

	// Calculate profit after adding the spread
	profit := conversionAmount * spreads[fromCurrency][toCurrency]

	// Calculate the beneficiary received value
	beneficiaryAmount := conversionAmount - profit

	return profit, beneficiaryAmount
}

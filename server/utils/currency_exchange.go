package utils

var exchangeRates = map[string]map[string]float64{
	"SGD": {
		"USD": 0.76,
		"AUD": 1.03,
		"MYR": 3.29,
	},
	"USD": {
		"SGD": 1.32,
		"AUD": 1.41,
		"MYR": 4.43,
	},
	"AUD": {
		"SGD": 0.97,
		"USD": 0.71,
		"MYR": 3.35,
	},
	"MYR": {
		"SGD": 0.30,
		"USD": 0.22,
		"AUD": 0.30,
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

func CalculateFromAmount(beneficiaryAmount float64, toCurrency, fromCurrency string) float64 {
	// Calculate the conversion amount
	conversionAmount := beneficiaryAmount / (1 - spreads[fromCurrency][toCurrency])

	// Calculate the transfer amount
	transferAmount := conversionAmount / exchangeRates[fromCurrency][toCurrency]

	// Calculate the profit
	// profit := conversionAmount - beneficiaryAmount

	return transferAmount
}

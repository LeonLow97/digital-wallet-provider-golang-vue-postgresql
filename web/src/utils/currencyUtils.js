export function currencyConversion(
  amountTransferred,
  transferredCurrency,
  receivedCurrency
) {
  if (isNaN(amountTransferred)) {
    return 0;
  }

  amountTransferred = parseFloat(amountTransferred);

  // Using objects to store the currency exchange
  const conversionRates = {
    USD: {
      SGD: 1.35,
      EUR: 0.92,
    },
    SGD: {
      USD: 0.74,
      EUR: 0.68,
    },
    EUR: {
      USD: 1.09,
      SGD: 1.47,
    },
  };

  if (amountTransferred < 0) {
    return 0;
  }

  if (transferredCurrency === receivedCurrency) {
    return amountTransferred;
  }

  const conversionRate = conversionRates[transferredCurrency][receivedCurrency];
  return amountTransferred * conversionRate;
}

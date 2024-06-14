export interface Balance {
  id: number;
  balance: number;
  currency: string;
  createdAt: string;
  updatedAt?: string;
}

export interface BalanceHistory {
  amount: number;
  currency: string;
  type: string;
  createdAt: string;
}

export interface GetBalancesResponse {
  balances: Balance[];
}

export interface GetBalanceResponse {
  id: number;
  balance: number;
  currency: string;
  createdAt: string;
}

export interface GetUserBalanceCurrenciesResponse {
  currency: string;
}

export interface GetBalanceHistoryResponse {
  balanceHistory: BalanceHistory[];
}

export interface DepositRequest {
  amount: number;
  currency: string;
}

export interface WithdrawRequest {
  amount: number;
  currency: string;
}

export interface CurrencyExchangeRequest {
  from_amount: number;
  to_currency: string;
}

export interface PreviewExchangeRequest {
  action_type: "amountToSend" | "amountToReceive" | null;
  from_amount: number;
  from_currency: string;
  to_amount: number;
  to_currency: string;
}

export interface PreviewExchangeResponse {
  actionType: string;
  fromAmount: number;
  fromCurrency: string;
  toAmount: number;
  toCurrency: string;
}

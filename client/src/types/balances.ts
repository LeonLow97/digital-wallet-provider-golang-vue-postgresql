export interface Balance {
  id: number;
  balance: number;
  currency: string;
  createdAt: string;
  updatedAt?: string;
}

interface BalanceHistory {
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

export interface DEPOSIT_REQUEST {
  amount: number;
  currency: string;
}

export interface WITHDRAW_REQUEST {
  amount: number;
  currency: string;
}

export interface CURRENCY_EXCHANGE_REQUEST {
  from_amount: number;
  to_currency: string;
}

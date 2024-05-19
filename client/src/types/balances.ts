export interface Balance {
  id: number;
  balance: number;
  currency: string;
  createdAt: string;
  updatedAt: string;
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

export interface GetBalanceHistoryResponse {
  balanceHistory: BalanceHistory[];
}

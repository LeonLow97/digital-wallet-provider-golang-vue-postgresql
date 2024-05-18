export interface Balance {
  id: number;
  balance: number;
  currency: string;
  createdAt: string;
  updatedAt: string;
}

export interface GetBalancesResponse {
  balances: Balance[];
  status: number;
}

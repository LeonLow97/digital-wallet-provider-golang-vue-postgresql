export interface Wallet {
  id: number;
  walletType: string;
  walletTypeID: number;
  userID: number;
  createdAt: string;
  currencyAmount: WalletCurrencyAmount[];
}

export interface WalletCurrencyAmount {
  walletID: number;
  amount: number;
  currency: string;
  createdAt: string;
  updatedAt: string;
}

export interface GetWalletsResponse {
  wallets: Wallet[];
}

export interface GetWalletTypesResponse {
  id: number;
  walletType: string;
}

export interface CreateWalletRequest {
  wallet_type_id: number;
  currency_amount: CurrencyAmount[];
}

export interface CurrencyAmount {
  amount: number;
  currency: string;
}

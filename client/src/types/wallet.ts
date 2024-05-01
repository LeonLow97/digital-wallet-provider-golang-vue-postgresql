export interface Wallet {
  createdAt: string;
  balance: number;
  currency: string;
  type: string;
  typeID: number;
  walletID: number;
}

export interface Wallets {
  wallets: Wallet[];
}

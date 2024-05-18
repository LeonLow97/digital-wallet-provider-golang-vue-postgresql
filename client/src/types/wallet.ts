export interface Wallet {
  createdAt: string;
  balance: number;
  currency: string;
  type: string;
  typeID: number;
  walletID: number;
}

export interface GetWalletsResponse {
  wallets: Wallet[];
}

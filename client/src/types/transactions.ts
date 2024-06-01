export interface Transaction {
  sender_id: number;
  beneficiary_id: number;
  sender_username: string;
  sender_mobile_number: string;
  beneficiary_username: string;
  beneficiary_mobile_number: string;
  source_amount: number;
  source_currency: string;
  destination_amount: number;
  destination_currency: string;
  source_of_transfer: string;
  status: string;
  created_at: string;
}

export interface CreateTransactionRequest {
  sender_wallet_id: number;
  source_currency: string;
  source_amount: number;
  beneficiary_mobile_country_code: string;
  beneficiary_mobile_number: string;
}

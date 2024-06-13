import axios from "axios";

const API_URL = import.meta.env.VITE_APP_API_URL;
import type {
  GENERIC_STATUS_RESPONSE,
  APIResponse,
  PAGINATION,
} from "@/types/generic";
import type {
  CreateTransactionRequest,
  Transaction,
} from "@/types/transactions";

const CREATE_TRANSACTION_URL = `${API_URL}/transaction`;
const GET_TRANSACTIONS_URL = `${API_URL}/transaction/all`;

export const CREATE_TRANSACTION = async (
  body: CreateTransactionRequest,
): Promise<GENERIC_STATUS_RESPONSE> => {
  const { status } = await axios.post(
    CREATE_TRANSACTION_URL,
    JSON.stringify(body),
    { withCredentials: true },
  );

  return { status };
};

export const GET_TRANSACTIONS = async (
  queryParams: PAGINATION,
): Promise<APIResponse<Transaction[]>> => {
  const { data, headers, status } = await axios.get(GET_TRANSACTIONS_URL, {
    params: queryParams,
    withCredentials: true,
  });

  return { data, headers, status };
};

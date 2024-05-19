import axios from "axios";
import type { APIResponse } from "@/types/generic";
import type {
  GetBalanceResponse,
  GetBalancesResponse,
  GetBalanceHistoryResponse,
} from "@/types/balances";

const API_URL = import.meta.env.VITE_APP_API_URL;

const GET_BALANCES_URL = `${API_URL}/balances`;
const GET_BALANCE_URL = `${API_URL}/balances/`;
const GET_BALANCE_HISTORY_URL = `${API_URL}/balances/history/`;

export const GET_BALANCES = async (): Promise<
  APIResponse<GetBalancesResponse>
> => {
  try {
    const { data, status } = await axios.get<GetBalancesResponse>(
      GET_BALANCES_URL,
      {
        withCredentials: true,
      },
    );

    return { data, status };
  } catch (error: unknown) {
    throw error;
  }
};

export const GET_BALANCE = async (
  balanceId: number,
): Promise<APIResponse<GetBalanceResponse>> => {
  try {
    const { data, status } = await axios.get<GetBalanceResponse>(
      GET_BALANCE_URL + balanceId,
      {
        withCredentials: true,
      },
    );

    return { data, status };
  } catch (error: unknown) {
    throw error;
  }
};

export const GET_BALANCE_HISTORY = async (
  balanceId: number,
): Promise<APIResponse<GetBalanceHistoryResponse>> => {
  try {
    const { data, status } = await axios.get<GetBalanceHistoryResponse>(
      GET_BALANCE_HISTORY_URL + balanceId,
      { withCredentials: true },
    );

    return { data, status };
  } catch (error: unknown) {
    throw error;
  }
};

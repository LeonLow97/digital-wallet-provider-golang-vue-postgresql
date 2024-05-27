import axios from "axios";
import type { APIResponse } from "@/types/generic";
import type {
  GetBalanceResponse,
  GetBalancesResponse,
  GetBalanceHistoryResponse,
  GetUserBalanceCurrenciesResponse,
  DEPOSIT_REQUEST,
  WITHDRAW_REQUEST,
} from "@/types/balances";
import type { GENERIC_STATUS_RESPONSE } from "@/types/generic";

const API_URL = import.meta.env.VITE_APP_API_URL;

const GET_BALANCES_URL = `${API_URL}/balances`;
const GET_BALANCE_URL = `${API_URL}/balances/`;
const GET_BALANCE_HISTORY_URL = `${API_URL}/balances/history/`;
const DEPOSIT_URL = `${API_URL}/balances/deposit`;
const WITHDRAW_URL = `${API_URL}/balances/withdraw`;
const GET_USER_BALANCE_CURRENCIES_URL = `${API_URL}/balances/currencies`;

export const GET_USER_BALANCE_CURRENCIES = async (): Promise<
  APIResponse<GetUserBalanceCurrenciesResponse[]>
> => {
  const { data, status } = await axios.get<GetUserBalanceCurrenciesResponse[]>(
    GET_USER_BALANCE_CURRENCIES_URL,
    {
      withCredentials: true,
    },
  );

  return { data, status };
};

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

export const DEPOSIT = async (
  body: DEPOSIT_REQUEST,
): Promise<GENERIC_STATUS_RESPONSE> => {
  try {
    const { status } = await axios.post(DEPOSIT_URL, JSON.stringify(body), {
      withCredentials: true,
    });

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};

export const WITHDRAW = async (
  body: WITHDRAW_REQUEST,
): Promise<GENERIC_STATUS_RESPONSE> => {
  try {
    const { status } = await axios.post(WITHDRAW_URL, JSON.stringify(body), {
      withCredentials: true,
    });

    return { status };
  } catch (error: unknown) {
    throw error;
  }
};

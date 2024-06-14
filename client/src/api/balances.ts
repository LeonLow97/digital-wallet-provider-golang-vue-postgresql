import axios from "axios";
import type { ApiResponse } from "@/types/generic";
import type {
  GetBalanceResponse,
  GetBalancesResponse,
  GetBalanceHistoryResponse,
  GetUserBalanceCurrenciesResponse,
  DepositRequest,
  WithdrawRequest,
  CurrencyExchangeRequest,
  PreviewExchangeRequest,
  PreviewExchangeResponse,
} from "@/types/balances";
import type { HttpStatusResponse } from "@/types/generic";

const API_URL = import.meta.env.VITE_APP_API_URL;

const GET_BALANCES_URL = `${API_URL}/balances`;
const GET_BALANCE_URL = `${API_URL}/balances/`;
const GET_BALANCE_HISTORY_URL = `${API_URL}/balances/history/`;
const DEPOSIT_URL = `${API_URL}/balances/deposit`;
const WITHDRAW_URL = `${API_URL}/balances/withdraw`;
const GET_USER_BALANCE_CURRENCIES_URL = `${API_URL}/balances/currencies`;
const CURRENCY_EXCHANGE_URL = `${API_URL}/balances/currency-exchange`;
const PREVIEW_EXCHANGE_URL = `${API_URL}/balances/preview-exchange`;

export const GET_USER_BALANCE_CURRENCIES = async (): Promise<
  ApiResponse<GetUserBalanceCurrenciesResponse[]>
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
  ApiResponse<GetBalancesResponse>
> => {
  const { data, status } = await axios.get<GetBalancesResponse>(
    GET_BALANCES_URL,
    {
      withCredentials: true,
    },
  );

  return { data, status };
};

export const GET_BALANCE = async (
  balanceId: number,
): Promise<ApiResponse<GetBalanceResponse>> => {
  const { data, status } = await axios.get<GetBalanceResponse>(
    GET_BALANCE_URL + balanceId,
    {
      withCredentials: true,
    },
  );

  return { data, status };
};

export const GET_BALANCE_HISTORY = async (
  balanceId: number,
): Promise<ApiResponse<GetBalanceHistoryResponse>> => {
  const { data, status } = await axios.get<GetBalanceHistoryResponse>(
    GET_BALANCE_HISTORY_URL + balanceId,
    { withCredentials: true },
  );

  return { data, status };
};

export const DEPOSIT = async (
  body: DepositRequest,
): Promise<HttpStatusResponse> => {
  const { status } = await axios.post(DEPOSIT_URL, JSON.stringify(body), {
    withCredentials: true,
  });

  return { status };
};

export const WITHDRAW = async (
  body: WithdrawRequest,
): Promise<HttpStatusResponse> => {
  const { status } = await axios.post(WITHDRAW_URL, JSON.stringify(body), {
    withCredentials: true,
  });

  return { status };
};

export const CURRENCY_EXCHANGE = async (
  body: CurrencyExchangeRequest,
): Promise<HttpStatusResponse> => {
  const { status } = await axios.patch(
    CURRENCY_EXCHANGE_URL,
    JSON.stringify(body),
    { withCredentials: true },
  );

  return { status };
};

export const PREVIEW_EXCHANGE = async (
  body: PreviewExchangeRequest,
): Promise<ApiResponse<PreviewExchangeResponse>> => {
  const { data, status } = await axios.post<PreviewExchangeResponse>(
    PREVIEW_EXCHANGE_URL,
    body,
    {
      withCredentials: true,
    },
  );

  return { data, status };
};
